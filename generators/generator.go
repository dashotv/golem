package generators

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/database"
	"github.com/dashotv/golem/generators/routes"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config *config.Config
}

func (g *Generator) Execute() error {
	runner := tasks.NewTaskRunner()

	if g.Config.Models.Enabled {
		runner.Add("database", func() error {
			dg := database.Generator{Config: g.Config}
			return dg.Execute()
		})
	}

	if g.Config.Routes.Enabled {
		runner.Add("routes", func() error {
			rg, err := routes.NewGenerator(g.Config)
			if err != nil {
				return err
			}

			return rg.Execute()
		})
		// TODO: Add Route generator
	}

	if g.Config.Jobs.Enabled {
		runner.Add("jobs", func() error {
			return nil
		})
		// TODO: Add Job generator
	}

	return runner.Run()
}
