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
	Models []map[string]string
	Data   *ConnectorGeneratorData
}

// ConnectorGeneratorData stores the data for the database connector template
type ConnectorGeneratorData struct {
	Models []*ConnectorGeneratorDataModel
}

type ConnectorGeneratorDataModel struct {
	Name   string
	Type   string
	Camel  string
	Struct bool
	Model  bool
}

// NewConnectorGenerator creates and returns a ConnectorGenerator
func NewConnectorGenerator(cfg *config.Config, models []map[string]string) *ConnectorGenerator {
	return &ConnectorGenerator{
		Config: cfg,
		Models: models,
		Data: &ConnectorGeneratorData{
			Models: []*ConnectorGeneratorDataModel{},
		},
		Generator: &base.Generator{
			Filename: cfg.Models.Output + "/models_connector.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute configures and generates the database connector
func (g *ConnectorGenerator) Execute(buf *bytes.Buffer) error {
	err := g.prepare()
	if err != nil {
		return err
	}

	err = templates.New("database_connector").Execute(buf, g.Data)
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
			Name:   m["Name"],
			Type:   m["Type"],
			Camel:  m["Camel"],
			Struct: m["Type"] == "struct",
			Model:  m["Type"] == "model",
		}
		g.Data.Models = append(g.Data.Models, d)
	}
	return nil
}
