package generators

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
	"github.com/dashotv/golem/templates"
)

// DatabaseGenerator manages the generation of all database related files
type DatabaseGenerator struct {
	*Generator
	Config *config.Config
	Models []map[string]string
}

// Execute configures and generates all of the database related files
func (g *DatabaseGenerator) Execute() error {
	runner := tasks.NewRunner("generator")
	r := runner.Group("database")

	source := g.Config.Models.Definitions
	if !tasks.PathExists(source) {
		return fmt.Errorf("definitions directory doesn't exist: %s", source)
	}

	dest := g.Config.Models.Output
	if !tasks.PathExists(dest) {
		return fmt.Errorf("output directory doesn't exist: %s", dest)
	}

	if g.Config.Models.Enabled {
		r.Add("header", g.generateHeader)
		r.Add("prepare", g.processModels)
		r.Add("connector", g.generateConnector)
		r.Add("models", g.generateModels)
		r.Add("write", g.Write)
		r.Add("format", g.Format)
		// r.Add("generate connector", func() error {
		// 	cg := NewConnectorGenerator(g.Config, g.Models)
		// 	return cg.Execute()
		// })
		//r.Add("generate schema", func() error {
		//	sg := NewSchemaGenerator(g.Config, g.Models)
		//	return sg.Execute()
		//})
		//r.Add("generate document", func() error {
		//	dg := NewFileGenerator(g.Config, "database_decoder", "models/decoder.go", map[string]string{})
		//	return dg.Execute()
		//})
	}

	return runner.Run()
}

func (g *DatabaseGenerator) generateHeader() error {
	data := map[string]string{
		"Package": g.Config.Models.Package,
	}
	return templates.New("database").Execute(g.Buffer, data)
}

func (g *DatabaseGenerator) generateConnector() error {
	cg := NewConnectorGenerator(g.Config, g.Models)
	return cg.Execute(g.Buffer)
}

func (g *DatabaseGenerator) generateModels() error {
	for _, m := range g.Models {
		file := "database_model"
		if m["Type"] == "struct" {
			file = "database_struct"
		}

		err := templates.New(file).Execute(g.Buffer, m)
		if err != nil {
			return err
		}
	}

	g.Buffer.WriteString("\n")
	return nil
}

// processModels walks all files in the model definitions path and generates model files for each
func (g *DatabaseGenerator) processModels() error {
	err := filepath.Walk(g.Config.Models.Definitions, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return g.addModel(info.Name(), path)
	})
	return err
}

func (g *DatabaseGenerator) addModel(name, path string) error {
	d := &Model{}
	err := ReadYaml(path, d)
	if err != nil {
		return err
	}

	data := map[string]string{
		"Name":  d.Name,
		"Camel": strcase.ToCamel(d.Name),
		"Type":  d.Type,
	}

	s := bytes.NewBufferString("")
	for _, fd := range d.Fields {
		fd.Camel = strcase.ToCamel(fd.Name)
		f := &FieldGenerator{
			Definition: fd,
		}
		err := f.Execute(s)
		if err != nil {
			return err
		}
	}
	data["Fields"] = s.String()

	g.Models = append(g.Models, data)

	return nil
}

// Model holds the data from the YAML model
type Model struct {
	Package string   `yaml:"package,omitempty"`
	Camel   string   `yaml:"camel,omitempty"`
	Name    string   `yaml:"name,omitempty"`
	Type    string   `yaml:"type,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
	Fields  []*Field `yaml:"fields,omitempty"`
}

// Field holds the data from the YAML field
type Field struct {
	Name   string   `yaml:"name,omitempty"`
	Camel  string   `yaml:"camel,omitempty"`
	Type   string   `yaml:"type,omitempty"`
	Json   string   `yaml:"json,omitempty"`
	Bson   string   `yaml:"bson,omitempty"`
	Tags   string   `yaml:"tags,omitempty"`
	Fields []*Field `yaml:"fields,omitempty"`
}

// ConnectorGenerator manages generation of the database connector
type ConnectorGenerator struct {
	*Generator
	Config *config.Config
	Models []map[string]string
	Data   *ConnectorGeneratorData
}

// ConnectorGeneratorData stores the data for the database connector template
type ConnectorGeneratorData struct {
	Models []*ConnectorGeneratorDataModel
}

type ConnectorGeneratorDataModel struct {
	Name   string
	Type   string
	Camel  string
	Struct bool
	Model  bool
}

// NewConnectorGenerator creates and returns a ConnectorGenerator
func NewConnectorGenerator(cfg *config.Config, models []map[string]string) *ConnectorGenerator {
	return &ConnectorGenerator{
		Config: cfg,
		Models: models,
		Data: &ConnectorGeneratorData{
			Models: []*ConnectorGeneratorDataModel{},
		},
		Generator: &Generator{
			Filename: cfg.Models.Output + "/models_connector.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute configures and generates the database connector
func (g *ConnectorGenerator) Execute(buf *bytes.Buffer) error {
	err := g.prepare()
	if err != nil {
		return err
	}

	err = templates.New("database_connector").Execute(buf, g.Data)
	if err != nil {
		return err
	}

	return nil
}

// Prepare the template configuration data
func (g *ConnectorGenerator) prepare() error {
	for _, m := range g.Models {
		//d := map[string]string{
		//	"Camel": m.Definition.Camel,
		//	"Name":  m.Definition.Name,
		//}
		d := &ConnectorGeneratorDataModel{
			Name:   m["Name"],
			Type:   m["Type"],
			Camel:  m["Camel"],
			Struct: m["Type"] == "struct",
			Model:  m["Type"] == "model",
		}
		g.Data.Models = append(g.Data.Models, d)
	}
	return nil
}

// FieldGenerator is the generator of model fields
type FieldGenerator struct {
	Definition *Field
	data       map[string]string
}

// Execute generates the field using the template
func (f *FieldGenerator) Execute(s *bytes.Buffer) error {
	err := f.Prepare()
	if err != nil {
		return err
	}

	if f.Definition.Type == "struct" || f.Definition.Type == "[]struct" || f.Definition.Type == "[]*struct" {
		err := f.executeStruct(s)
		if err != nil {
			return err
		}
	} else {
		err := f.executeField(s)
		if err != nil {
			return err
		}
	}

	return nil
}

// Prepare configures the data for the template
func (f *FieldGenerator) Prepare() error {
	f.data = map[string]string{
		"Name": strcase.ToCamel(f.Definition.Name),
		"Type": f.Definition.Type,
	}

	err := f.prepareTags()
	if err != nil {
		return err
	}

	return nil
}

// executeStruct generates the field if it is a struct type
func (f *FieldGenerator) executeStruct(s *bytes.Buffer) error {
	t, err := template.New("simplefield").Parse(`    {{.Name}} {{.Type}} {
		{{.Fields}}
  } {{.Tags}}` + "\n")
	if err != nil {
		return err
	}

	err = f.prepareStructFields()
	if err != nil {
		return err
	}

	err = t.Execute(s, f.data)
	if err != nil {
		return err
	}

	return nil
}

// executeField generates the field if it is not a struct type
func (f *FieldGenerator) executeField(s *bytes.Buffer) error {
	t, err := template.New("simplefield").Parse(`    {{.Name}} {{.Type}} {{.Tags}}` + "\n")
	if err != nil {
		return err
	}

	err = t.Execute(s, f.data)
	if err != nil {
		return err
	}

	return nil
}

// prepareStructFields configures the data for struct field template
func (f *FieldGenerator) prepareStructFields() error {
	s := bytes.NewBufferString("")
	for _, fd := range f.Definition.Fields {
		f := &FieldGenerator{Definition: fd}
		err := f.Execute(s)
		if err != nil {
			return err
		}
	}
	f.data["Fields"] = s.String()
	return nil
}

// prepareTags configures the tag data for the field template
func (f *FieldGenerator) prepareTags() error {
	json := f.Definition.Name
	bson := f.Definition.Name
	if f.Definition.Json != "" {
		json = f.Definition.Json
	}
	if f.Definition.Bson != "" {
		bson = f.Definition.Bson
	}

	f.data["Tags"] = fmt.Sprintf("`json:\"%s\" bson:\"%s\"`", json, bson)
	return nil
}
