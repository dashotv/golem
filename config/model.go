package config

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type Model struct {
	Name          string                 `yaml:"name,omitempty"`
	Type          string                 `yaml:"type,omitempty"`
	QueryDefaults map[string]interface{} `yaml:"query_defaults,omitempty"`
	Imports       []string               `yaml:"imports,omitempty"`
	Fields        []*Field               `yaml:"fields,omitempty"`
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

func (m *Model) QueryDefaultsString() string {
	out := []string{}
	if len(m.QueryDefaults) == 0 {
		return ""
	}
	for k, v := range m.QueryDefaults {
		switch val := v.(type) {
		case string:
			out = append(out, fmt.Sprintf("{\"%s\": \"%s\"}", k, val))
		default:
			out = append(out, fmt.Sprintf("{\"%s\": %v}", k, val))
		}
	}
	return "[]bson.M{" + strings.Join(out, ",") + "}"
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
