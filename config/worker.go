package config

import "github.com/iancoleman/strcase"

type Worker struct {
	Name     string   `yaml:"name,omitempty"`
	Queue    string   `yaml:"queue,omitempty"`
	Schedule string   `yaml:"schedule,omitempty"`
	Fields   []*Field `yaml:"fields,omitempty"`
}

func (w *Worker) Camel() string {
	return strcase.ToCamel(w.Name)
}
