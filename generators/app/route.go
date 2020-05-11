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
	Crud       bool
	Definition *routes.Definition
}

// NewRouteDefinitionGenerator creates and returns an instance of RouteDefinitionGenerator
func NewRouteDefinitionGenerator(cfg *config.Config, name string, crud bool, params ...string) *RouteDefinitionGenerator {
	return &RouteDefinitionGenerator{
		Config:     cfg,
		Name:       name,
		Params:     params,
		Definition: &routes.Definition{},
		Crud:       crud,
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

	// ensure groups exists
	if g.Definition.Groups == nil {
		g.Definition.Groups = make(map[string]*routes.Group)
	}

	// ensure group exists
	var gd *routes.Group
	if g.Definition.Groups[parts[1]] != nil {
		gd = g.Definition.Groups[parts[1]]
	} else {
		gd = &routes.Group{}
		if g.Definition.Groups == nil {
			g.Definition.Groups = make(map[string]*routes.Group)
		}
		g.Definition.Groups[parts[1]] = gd
	}
	gd.Path = "/" + parts[1]

	// ensure route exists
	var rd *routes.Route
	if gd.Routes[parts[2]] != nil {
		rd = gd.Routes[parts[2]]
	} else {
		rd = &routes.Route{}
		if gd.Routes == nil {
			gd.Routes = make(map[string]*routes.Route)
		}
		gd.Routes[parts[2]] = rd
	}
	rd.Path = "/" + parts[2]

	if g.Crud {
		// add crud routes
	} else {
		for _, f := range g.Params {
			f := strings.Split(f, ":")
			n := f[0]
			t := f[1]
			if t == "" {
				t = "string"
			}
			rd.Params = append(rd.Params, &routes.Param{Name: n, Type: t})
		}
	}
}

// routeExists determines if a route path already exists
//func (g *RouteDefinitionGenerator) routeExists(name string, routes []*routes.Route) bool {
//	if g.Definition.FindRoute(name) != nil {
//		return true
//	}
//	return false
//}
