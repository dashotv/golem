package database

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/templates"
)

// ModelGenerator is the database model generator
type ModelGenerator struct {
	*base.Generator
	Config     *config.Config
	Name       string
	Path       string
	Definition *Model
	data       map[string]string
}

// NewModelGenerator creates and returns an instance of ModelGenerator
func NewModelGenerator(cfg *config.Config, name, path string) (*ModelGenerator, error) {
	d := &Model{}
	err := base.ReadYaml(path, d)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(name, ".")
	d.Package = cfg.Models.Package
	m := &ModelGenerator{
		Config:     cfg,
		Name:       name,
		Path:       path,
		Definition: d,
		data:       make(map[string]string),
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/models_" + parts[0] + ".go",
			Buffer:   bytes.NewBufferString(""),
		},
	}

	return m, nil
}

// Execute generates the model file from the template
func (m *ModelGenerator) Execute() error {
	err := m.prepare()
	if err != nil {
		return err
	}

	file := "database_model"
	if m.data["Type"] == "struct" {
		file = "database_struct"
	}

	err = templates.New(file).Execute(m.Buffer, m.data)
	if err != nil {
		return err
	}

	err = m.Write()
	if err != nil {
		return err
	}

	err = m.Format()
	if err != nil {
		return err
	}

	return nil
}

// prepare configures the data for the template
func (m *ModelGenerator) prepare() error {
	m.Definition.Camel = strcase.ToCamel(m.Definition.Name)
	m.data = map[string]string{
		"Package": m.Config.Models.Package,
		"Name":    m.Definition.Camel,
		"Type":    m.Definition.Type,
	}

	err := m.prepareFields()
	if err != nil {
		return err
	}

	return nil
}

// prepareFields configures the field data for the template
func (m *ModelGenerator) prepareFields() error {
	s := bytes.NewBufferString("")
	for _, fd := range m.Definition.Fields {
		fd.Camel = strcase.ToCamel(fd.Name)
		f := &FieldGenerator{
			Definition: fd,
		}
		err := f.Execute(s)
		if err != nil {
			return err
		}
	}
	m.data["Fields"] = s.String()
	return nil
}
