package database

import (
	"bytes"

	"github.com/dashotv/golem/generators/templates"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

// ConnectorGenerator manages generation of the database connector
type ConnectorGenerator struct {
	*base.Generator
	Config *config.Config
	Models []*ModelGenerator
	Data   *ConnectorGeneratorData
}

// ConnectorGeneratorData stores the data for the database connector template
type ConnectorGeneratorData struct {
	Package string
	Models  []map[string]string
	Repo    string
}

// NewConnectorGenerator creates and returns a ConnectorGenerator
func NewConnectorGenerator(cfg *config.Config, models []*ModelGenerator) *ConnectorGenerator {
	return &ConnectorGenerator{
		Config: cfg,
		Models: models,
		Data: &ConnectorGeneratorData{
			Package: cfg.Models.Package,
			Models:  make([]map[string]string, 0),
			Repo:    cfg.Repo,
		},
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/connector.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute configures and generates the database connector
func (g *ConnectorGenerator) Execute() error {
	err := g.prepare()
	if err != nil {
		return err
	}

	err = templates.New("database_connector").Execute(g.Buffer, g.Data)
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

// Prepare the template configuration data
func (g *ConnectorGenerator) prepare() error {
	for _, m := range g.Models {
		d := map[string]string{
			"Camel": m.Definition.Camel,
			"Name":  m.Definition.Name,
		}
		g.Data.Models = append(g.Data.Models, d)
	}
	return nil
}
