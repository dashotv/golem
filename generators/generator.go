package generators

import (
	"bytes"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/generators/database"
	"github.com/dashotv/golem/generators/routes"
	"github.com/dashotv/golem/tasks"
)

// Generator is the top-level generator and calls other generators
type Generator struct {
	Config *config.Config
}

// Execute calls all dependent generators based on configuration
func (g *Generator) Execute() error {
	runner := tasks.NewRunner("generator")

	if g.Config.Models.Enabled {
		runner.Add("database", func() error {
			dg := database.Generator{
				Config: g.Config,
				Generator: &base.Generator{
					Filename: g.Config.Models.Output + "/models.go",
					Buffer:   bytes.NewBufferString(""),
				},
			}
			return dg.Execute()
		})
	}

	if g.Config.Routes.Enabled {
		runner.Add("routes", func() error {
			rg := routes.NewGenerator(g.Config)
			return rg.Execute()
		})
	}

	if g.Config.Jobs.Enabled {
		runner.Add("jobs", func() error {
			return nil
		})
		// TODO: Add Job generator
	}

	return runner.Run()
}
