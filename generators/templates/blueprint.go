package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

const DIR = "generators/templates"

type Blueprint struct {
	Name     string
	Template *template.Template
}

func New(name string) *Blueprint {
	return &Blueprint{
		Name: name,
	}
}

func (b *Blueprint) readTemplate() (string, error) {
	filename := fmt.Sprintf("%s/%s.tgo", DIR, b.Name)

	text, err := Asset(filename)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func (b *Blueprint) Execute(buf *bytes.Buffer, data interface{}) error {
	text, err := b.readTemplate()
	if err != nil {
		return err
	}

	b.Template, err = template.New(b.Name).Parse(text)
	if err != nil {
		return err
	}

	return b.Template.Execute(buf, data)
}
