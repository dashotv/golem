package database

import (
	"bytes"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/templates"
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
	Models  []*ConnectorGeneratorDataModel
	Repo    string
}

type ConnectorGeneratorDataModel struct {
	Name   string
	Type   string
	Camel  string
	Struct bool
	Model  bool
}

// NewConnectorGenerator creates and returns a ConnectorGenerator
func NewConnectorGenerator(cfg *config.Config, models []*ModelGenerator) *ConnectorGenerator {
	return &ConnectorGenerator{
		Config: cfg,
		Models: models,
		Data: &ConnectorGeneratorData{
			Package: cfg.Models.Package,
			Models:  make([]*ConnectorGeneratorDataModel, 0),
			Repo:    cfg.Repo,
		},
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/models_connector.go",
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
		//d := map[string]string{
		//	"Camel": m.Definition.Camel,
		//	"Name":  m.Definition.Name,
		//}
		d := &ConnectorGeneratorDataModel{
			Name:   m.Definition.Name,
			Type:   m.Definition.Type,
			Camel:  m.Definition.Camel,
			Struct: m.Definition.Type == "struct",
			Model:  m.Definition.Type == "model",
		}
		g.Data.Models = append(g.Data.Models, d)
	}
	return nil
}
