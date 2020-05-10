package routes

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config *config.Config
	Routes []*RouteGenerator
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{Config: cfg}
}

func (g *Generator) Execute() error {
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
		r.Add("generate server", g.processServer)
		r.Add("generate routes", g.processRoutes)
	}

	return runner.Run()
}

func (g *Generator) processServer() error {
	sg, err := NewServerGenerator(g.Config)
	if err != nil {
		return err
	}

	return sg.Execute()
}

func (g *Generator) processRoutes() error {
	return nil
}
