package routes

import (
	"strings"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	Config     *config.Config
	Definition *Definition
	Groups     []*GroupGenerator
}

type Definition struct {
	Name   string
	Repo   string
	Groups map[string]*GroupDefinition `json:"groups" yaml:"groups"`
	Routes map[string]*RouteDefinition `json:"routes" yaml:"routes"`
}

func (d *Definition) FindRoute(name string) *RouteDefinition {
	parts := strings.Split(name, "/")
	if len(parts) > 2 {
		return nil
	}
	if len(parts) == 1 {
		if d.Routes[parts[0]] == nil {
			return nil
		}
		return d.Routes[parts[0]]
	}
	if d.Groups[parts[0]] == nil {
		return nil
	}
	if d.Groups[parts[0]].Routes[parts[1]] == nil {
		return nil
	}
	return d.Groups[parts[0]].Routes[parts[1]]
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		Config:     cfg,
		Groups:     make([]*GroupGenerator, 0),
		Definition: &Definition{},
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
	runner := tasks.NewRunner("generator")
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
		for name, group := range g.Definition.Groups {
			r.Add("generate groups "+name, func() error {
				rg := NewGroupGenerator(g.Config, name, group)
				g.Groups = append(g.Groups, rg)
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
