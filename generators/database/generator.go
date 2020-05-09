package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config *config.Config
	Models []*ModelGenerator
}

func (g *Generator) Execute() error {
	runner := tasks.NewTaskRunner()
	r := runner.Group("database")

	source := g.Config.Models.Definitions
	if !exists(source) {
		return fmt.Errorf("definitions directory doesn't exist: %s", source)
	}

	dest := g.Config.Models.Output
	if !exists(dest) {
		return fmt.Errorf("output directory doesn't exist: %s", dest)
	}

	if g.Config.Models.Enabled {
		r.Add("generate models", g.processModels)
		r.Add("generate schema", g.processSchema)
		r.Add("generate connector", g.processConnector)
		r.Add("generate document", g.processDocument)
		// TODO: decide on using document and null decoder
	}

	return runner.Run()
}

func (g *Generator) processModels() error {
	err := filepath.Walk(g.Config.Models.Definitions, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		m, err := NewModelGenerator(g.Config, info.Name(), path)
		if err != nil {
			return err
		}
		g.Models = append(g.Models, m)

		return m.Execute()
	})
	return err
}

func (g *Generator) processSchema() error {
	sg, err := NewSchemaGenerator(g.Config, g.Models)
	if err != nil {
		return err
	}

	return sg.Execute()
}

func (g *Generator) processConnector() error {
	cg := NewConnectorGenerator(g.Config, g.Models)
	return cg.Execute()
}

func (g *Generator) processDocument() error {
	dg := base.NewFileGenerator(g.Config, "document", "models/document.go", map[string]string{})
	return dg.Execute()
}

func exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
