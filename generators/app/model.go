package app

import (
	"bytes"
	"strings"

	"github.com/dashotv/golem/tasks"

	"github.com/iancoleman/strcase"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/generators/templates"

	"github.com/dashotv/golem/generators/database"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type ModelDefinitionGenerator struct {
	*base.Generator
	Config     *config.Config
	Name       string
	Fields     []string
	Definition *database.Model
}

func NewModelDefinitionGenerator(cfg *config.Config, name string, fields ...string) *ModelDefinitionGenerator {
	return &ModelDefinitionGenerator{
		Config:     cfg,
		Name:       name,
		Fields:     fields,
		Definition: &database.Model{},
		Generator: &base.Generator{
			Filename: ".golem/models/" + name + ".yaml",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (g *ModelDefinitionGenerator) Execute() error {
	if !exists(".golem") {
		return errors.New(".golem directory does not exist, run from inside app directory")
	}
	if exists(g.Filename) {
		return errors.New("model definition already exists: " + g.Filename)
	}

	g.prepare()

	runner := tasks.NewTaskRunner("generator:model")
	runner.Add("ensure models directory exists", func() error {
		return makeDirectory(".golem/models")
	})
	runner.Add("generate model definition", func() error {
		err := templates.New("app_model_yaml").Execute(g.Buffer, g.Definition)
		if err != nil {
			return err
		}

		return g.Write()
	})

	return runner.Run()
}

func (g *ModelDefinitionGenerator) prepare() {
	g.Definition.Name = g.Name
	g.Definition.Camel = strcase.ToCamel(g.Name)
	g.Definition.Type = "model"
	g.Definition.Fields = make([]*database.Field, 0)

	for _, f := range g.Fields {
		f := strings.Split(f, ":")
		n := f[0]
		t := f[1]
		if t == "" {
			t = "string"
		}
		g.Definition.Fields = append(g.Definition.Fields, &database.Field{Name: n, Type: t})
	}
}
