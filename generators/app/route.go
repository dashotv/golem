package app

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/generators/routes"
	"github.com/dashotv/golem/generators/templates"
	"github.com/dashotv/golem/tasks"
)

// RouteDefinitionGenerator manages the generation of and updates to routes definition
type RouteDefinitionGenerator struct {
	*base.Generator
	Config     *config.Config
	Name       string
	Params     []string
	Rest       bool
	Definition *routes.Definition
}

// NewRouteDefinitionGenerator creates and returns an instance of RouteDefinitionGenerator
func NewRouteDefinitionGenerator(cfg *config.Config, name string, rest bool, params ...string) *RouteDefinitionGenerator {
	if name[0] != '/' {
		name = "/" + name
	}
	return &RouteDefinitionGenerator{
		Config:     cfg,
		Name:       name,
		Params:     params,
		Definition: &routes.Definition{},
		Rest:       rest,
		Generator: &base.Generator{
			Filename: ".golem/routes.yaml",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute generates or updates the routes definition
func (g *RouteDefinitionGenerator) Execute() error {
	if !tasks.PathExists(".golem") {
		return errors.New(".golem directory does not exist, run from inside app directory")
	}
	if tasks.PathExists(".golem/routes.yaml") {
		if err := base.ReadYaml(".golem/routes.yaml", g.Definition); err != nil {
			return err
		}
	}

	g.Definition.Name = g.Config.Name
	g.Definition.Repo = g.Config.Repo

	g.prepare()

	runner := tasks.NewRunner("generator:route")
	runner.Add("generate route definition", func() error {
		err := templates.New("app_routes_yaml").Execute(g.Buffer, g.Definition)
		if err != nil {
			return err
		}

		return g.Write()
	})

	return runner.Run()
}

// prepare the data for generation
func (g *RouteDefinitionGenerator) prepare() {
	rd := g.Definition.FindRoute(g.Name)
	if rd != nil {
		g.prepareModify(rd)
	} else {
		g.prepareCreate()
	}
}

// prepareModify updates an existing routes definition
func (g *RouteDefinitionGenerator) prepareModify(rd *routes.Route) {

}

// prepareCreate creates a new routes definition
func (g *RouteDefinitionGenerator) prepareCreate() {
	parts := strings.Split(g.Name, "/")
	route := parts[1]
	path := "index"
	if len(parts) > 2 {
		path = parts[2]
	}

	// ensure groups exists
	if g.Definition.Groups == nil {
		g.Definition.Groups = make(map[string]*routes.Group)
	}

	// ensure group exists
	var gd *routes.Group
	if g.Definition.Groups[route] != nil {
		gd = g.Definition.Groups[route]
		gd.Path = "/" + route
	} else {
		gd = &routes.Group{}
		if g.Definition.Groups == nil {
			g.Definition.Groups = make(map[string]*routes.Group)
		}
		g.Definition.Groups[route] = gd
		gd.Path = route
	}

	if g.Rest {
		gd.Rest = true
		return
	}

	// ensure route exists
	var rd *routes.Route
	if gd.Routes[path] != nil {
		rd = gd.Routes[path]
		rd.Path = "/" + path
	} else {
		rd = &routes.Route{}
		if gd.Routes == nil {
			gd.Routes = make(map[string]*routes.Route)
		}
		gd.Routes[path] = rd
	}

	for _, f := range g.Params {
		f := strings.Split(f, ":")
		n := f[0]
		t := "string"
		if len(f) > 1 {
			t = f[1]
		}
		rd.Params = append(rd.Params, &routes.Param{Name: n, Type: t})
	}
}

// routeExists determines if a route path already exists
//func (g *RouteDefinitionGenerator) routeExists(name string, routes []*routes.Route) bool {
//	if g.Definition.FindRoute(name) != nil {
//		return true
//	}
//	return false
//}
