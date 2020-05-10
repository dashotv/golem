package routes

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config     *config.Config
	Definition *ServerDefinition
	Routes     []*RouteGenerator
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		Config:     cfg,
		Routes:     make([]*RouteGenerator, 0),
		Definition: &ServerDefinition{},
		//	Name:    cfg.Name,
		//	Package: "server",
		//	Repo:    cfg.Repo,
		//},
	}
}

func (g *Generator) Execute() error {
	if err := g.prepare(); err != nil {
		return err
	}
	runner := tasks.NewTaskRunner("generator")
	r := runner.Group("routes")

	//source := g.Config.Models.Definitions
	//if !exists(source) {
	//	return fmt.Errorf("definitions directory doesn't exist: %s", source)
	//}
	//
	//dest := g.Config.Models.Output
	//if !exists(dest) {
	//	return fmt.Errorf("output directory doesn't exist: %s", dest)
	//}

	if g.Config.Routes.Enabled {
		r.Add("generate server", func() error {
			sg := NewServerGenerator(g.Config, g.Definition)
			return sg.Execute()
		})
		for _, route := range g.Definition.Routes {
			r.Add("generate routes "+route.Name, func() error {
				rg := NewRouteGenerator(g.Config, route)
				g.Routes = append(g.Routes, rg)
				return rg.Execute()
			})
		}
	}

	return runner.Run()
}

func (g *Generator) prepare() error {
	if err := base.ReadYaml(g.Config.Routes.Definition, g.Definition); err != nil {
		return err
	}
	return nil
}
