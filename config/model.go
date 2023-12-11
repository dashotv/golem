package config

import "github.com/iancoleman/strcase"

type Model struct {
	Name    string   `yaml:"name,omitempty"`
	Type    string   `yaml:"type,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
	Fields  []*Field `yaml:"fields,omitempty"`
}

func (m *Model) Model() bool {
	return m.Type == "model"
}

func (m *Model) Struct() bool {
	return m.Type == "struct"
}

func (m *Model) Camel() string {
	return strcase.ToCamel(m.Name)
}

type Field struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	Json string `yaml:"json,omitempty"`
	Bson string `yaml:"bson,omitempty"`
}

func (f *Field) Camel() string {
	return strcase.ToCamel(f.Name)
}

func (f *Field) JsonTag() string {
	if f.Json == "" {
		return f.Name + ",omitempty"
	}
	return f.Json
}

func (f *Field) BsonTag() string {
	if f.Bson == "" {
		return f.Name + ",omitempty"
	}
	return f.Bson
}
