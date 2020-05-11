package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

// Generator manages the generation of all database related files
type Generator struct {
	Config *config.Config
	Models []*ModelGenerator
}

// Execute configures and generates all of the database related files
func (g *Generator) Execute() error {
	runner := tasks.NewRunner("generator")
	r := runner.Group("database")

	source := g.Config.Models.Definitions
	if !tasks.PathExists(source) {
		return fmt.Errorf("definitions directory doesn't exist: %s", source)
	}

	dest := g.Config.Models.Output
	if !tasks.PathExists(dest) {
		return fmt.Errorf("output directory doesn't exist: %s", dest)
	}

	if g.Config.Models.Enabled {
		r.Add("generate models", g.processModels)
		r.Add("generate schema", func() error {
			sg := NewSchemaGenerator(g.Config, g.Models)
			return sg.Execute()
		})
		r.Add("generate connector", func() error {
			cg := NewConnectorGenerator(g.Config, g.Models)
			return cg.Execute()
		})
		r.Add("generate document", func() error {
			dg := base.NewFileGenerator(g.Config, "database_document", "models/document.go", map[string]string{})
			return dg.Execute()
		})
	}

	return runner.Run()
}

// processModels walks all files in the model definitions path and generates model files for each
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
