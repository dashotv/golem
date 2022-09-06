package routes

import (
	"strings"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

// Generator manages the generation of all routes related files
type Generator struct {
	Config     *config.Config
	Definition *Definition
	Groups     []*GroupGenerator
}

// FindRoute finds a route corresponding to the path name passed as argument
func (d *Definition) FindRoute(name string) *Route {
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

// NewGenerator creates and returns an instance of Generator
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		Config:     cfg,
		Groups:     make([]*GroupGenerator, 0),
		Definition: &Definition{},
	}
}

// Execute generates the files from the templates
func (g *Generator) Execute() error {
	if err := g.prepare(); err != nil {
		return err
	}
	runner := tasks.NewRunner("generator")
	r := runner.Group("routes")

	if g.Config.Routes.Enabled {
		r.Add("server routes", func() error {
			sg := NewServerRoutesGenerator(g.Config, g.Definition)
			return sg.Execute()
		})
	}

	return runner.Run()
}

// prepare configures the data for the templates
func (g *Generator) prepare() error {
	if err := base.ReadYaml(g.Config.Routes.Definition, g.Definition); err != nil {
		return err
	}

	for name, group := range g.Definition.Groups {
		if group.Method == "" {
			group.Method = "GET"
		}
		if group.Path[0] != '/' {
			group.Path = "/" + group.Path
		}
		if group.Rest {
			group.Routes = map[string]*Route{
				"index": &Route{
					Method: "GET",
					Name:   name + strings.Title("index"),
				},
				"create": &Route{
					Method: "POST",
					Name:   name + strings.Title("create"),
				},
				"show": &Route{
					Method: "GET",
					Name:   name + strings.Title("show"),
				},
				"update": &Route{
					Method: "PUT",
					Name:   name + strings.Title("update"),
				},
				"delete": &Route{
					Method: "DELETE",
					Name:   name + strings.Title("delete"),
				},
			}
		}
		for rn, route := range group.Routes {
			if route.Path == "" {
				route.Path = "/"
			}
			route.Name = name + strings.Title(rn)
		}
	}
	return nil
}
