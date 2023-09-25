package generators

import (
	"bytes"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

// Generator is the top-level generator and calls other generators
type MainGenerator struct {
	Config *config.Config
}

// Execute calls all dependent generators based on configuration
func (g *MainGenerator) Execute() error {
	runner := tasks.NewRunner("generator")

	if g.Config.Models.Enabled {
		runner.Add("database", func() error {
			dg := &DatabaseGenerator{
				Config: g.Config,
				Generator: &Generator{
					Filename: g.Config.Models.Output + "/models.go",
					Buffer:   bytes.NewBufferString(""),
				},
			}
			return dg.Execute()
		})
	}

	if g.Config.Routes.Enabled {
		runner.Add("routes", func() error {
			rg := NewRoutesGenerator(g.Config)
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
