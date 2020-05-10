package database

import (
	"bytes"

	"github.com/dashotv/golem/generators/templates"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type ConnectorGenerator struct {
	*base.Generator
	Config *config.Config
	Models []*ModelGenerator
	Data   *ConnectorGeneratorData
}

type ConnectorGeneratorData struct {
	Package string
	Models  []map[string]string
}

func NewConnectorGenerator(cfg *config.Config, models []*ModelGenerator) *ConnectorGenerator {
	return &ConnectorGenerator{
		Config: cfg,
		Models: models,
		Data: &ConnectorGeneratorData{
			Package: cfg.Models.Package,
			Models:  make([]map[string]string, 0),
		},
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/connector.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (g *ConnectorGenerator) Execute() error {
	err := g.prepare()
	if err != nil {
		return err
	}

	err = templates.New("connector").Execute(g.Buffer, g.Data)
	if err != nil {
		return err
	}

	err = g.Write()
	if err != nil {
		return err
	}

	err = g.Format()
	if err != nil {
		return err
	}

	return nil
}

func (g *ConnectorGenerator) prepare() error {
	def := g.Config.Connections["default"]
	for _, m := range g.Models {
		var c *config.Connection
		if g.Config.Connections[m.Definition.Name] != nil {
			c = g.Config.Connections[m.Definition.Name]
		} else {
			c = &config.Connection{}
		}
		if c.URI == "" {
			c.URI = def.URI
		}
		if c.Database == "" {
			c.Database = def.Database
		}
		if c.Collection == "" {
			c.Collection = m.Definition.Name
		}
		d := map[string]string{
			"Camel":      m.Definition.Camel,
			"Name":       m.Definition.Name,
			"URI":        c.URI,
			"Database":   c.Database,
			"Collection": c.Collection,
		}
		g.Data.Models = append(g.Data.Models, d)
	}
	return nil
}
