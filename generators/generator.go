package generators

import (
	"os"
	"path/filepath"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/database"
	"github.com/dashotv/golem/generators/routes"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config *config.Config
	Models []*database.ModelGenerator
}

func (g *Generator) Process() error {
	runner := tasks.NewTaskRunner()

	if g.Config.Models.Enabled {
		r := runner.Group("database")
		r.Add("generate models", g.processModels)
		r.Add("generate schema", g.processSchema)
		r.Add("generate connector", g.processConnector)
		// TODO: decide on using document and null decoder
	}

	if g.Config.Routes.Enabled {
		r := runner.Group("routes")
		r.Add("generate routes", g.processRoutes)
		// TODO: Add Route generator
	}

	if g.Config.Jobs.Enabled {
		r := runner.Group("jobs")
		r.Add("generate jobs", g.processJobs)
		// TODO: Add Job generator
	}

	return runner.Run()
}

func (g *Generator) processJobs() error {
	return nil
}

func (g *Generator) processRoutes() error {
	rg, err := routes.NewGenerator(g.Config)
	if err != nil {
		return err
	}

	return rg.Execute()
}

func (g *Generator) processModels() error {
	err := filepath.Walk(g.Config.Models.Definitions, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		m, err := database.NewModelGenerator(g.Config, info.Name(), path)
		if err != nil {
			return err
		}
		g.Models = append(g.Models, m)

		return m.Execute()
	})
	return err
}

func (g *Generator) processSchema() error {
	sg, err := database.NewSchemaGenerator(g.Config, g.Models)
	if err != nil {
		return err
	}

	return sg.Execute()
}

func (g *Generator) processConnector() error {
	cg := database.NewConnectorGenerator(g.Config, g.Models)
	return cg.Execute()
}
