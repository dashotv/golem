package database

import (
	"bytes"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Field is the generator of model fields
type FieldGenerator struct {
	Definition *Field
	data       map[string]string
}

// FieldDefinition holds the data from the YAML field
type Field struct {
	Name   string
	Camel  string
	Type   string
	Json   string
	Bson   string
	Tags   string
	Fields []*Field
}

func (f *FieldGenerator) Execute(s *bytes.Buffer) error {
	err := f.Prepare()
	if err != nil {
		return err
	}

	if f.Definition.Type == "struct" || f.Definition.Type == "[]struct" {
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

func (f *FieldGenerator) prepareTags() error {
	s := bytes.NewBufferString("")

	t, err := template.New("fieldtags").Parse("`json:\"{{.json}}\" bson:\"{{.bson}}\"`")
	if err != nil {
		return err
	}

	data := map[string]string{
		"json": f.Definition.Name,
		"bson": f.Definition.Name,
	}
	if f.Definition.Json != "" {
		data["json"] = f.Definition.Json
	}
	if f.Definition.Bson != "" {
		data["bson"] = f.Definition.Bson
	}

	err = t.Execute(s, data)
	if err != nil {
		return err
	}

	f.data["Tags"] = s.String()
	return nil
}
