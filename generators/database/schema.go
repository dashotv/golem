package database

import (
	"bytes"

	"github.com/dashotv/golem/generators/templates"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type SchemaGenerator struct {
	*base.Generator
	Config *config.Config
	Models []*ModelGenerator
	Data   *SchemaGeneratorData
}

type SchemaGeneratorData struct {
	Package string
	Models  []*Model
}

func NewSchemaGenerator(cfg *config.Config, models []*ModelGenerator) (*SchemaGenerator, error) {
	g := &SchemaGenerator{
		Config: cfg,
		Models: models,
		Data: &SchemaGeneratorData{
			Package: cfg.Models.Package,
			Models:  make([]*Model, 0),
		},
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/schema.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}

	return g, nil
}

func (g *SchemaGenerator) Execute() error {
	err := g.prepare()
	if err != nil {
		return err
	}

	err = templates.New("schema").Execute(g.Buffer, g.Data)
	if err != nil {
		return err
	}

	err = g.Write()
	if err != nil {
		return err
	}

	return nil
}

func (g *SchemaGenerator) prepare() error {
	for _, m := range g.Models {
		g.Data.Models = append(g.Data.Models, m.Definition)
	}
	return nil
}
