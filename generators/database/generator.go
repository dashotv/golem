package database

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
	"github.com/dashotv/golem/templates"
)

// Generator manages the generation of all database related files
type Generator struct {
	*base.Generator
	Config *config.Config
	Models []map[string]string
}

// Execute configures and generates all of the database related files
func (g *Generator) Execute() error {
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
		//	dg := base.NewFileGenerator(g.Config, "database_decoder", "models/decoder.go", map[string]string{})
		//	return dg.Execute()
		//})
	}

	return runner.Run()
}

func (g *Generator) generateHeader() error {
	data := map[string]string{
		"Package": g.Config.Models.Package,
	}
	return templates.New("database").Execute(g.Buffer, data)
}

func (g *Generator) generateConnector() error {
	cg := NewConnectorGenerator(g.Config, g.Models)
	return cg.Execute(g.Buffer)
}

func (g *Generator) generateModels() error {
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
func (g *Generator) processModels() error {
	err := filepath.Walk(g.Config.Models.Definitions, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return g.addModel(info.Name(), path)
	})
	return err
}

func (g *Generator) addModel(name, path string) error {
	d := &Model{}
	err := base.ReadYaml(path, d)
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
