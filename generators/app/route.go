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

type RouteDefinitionGenerator struct {
	*base.Generator
	Config     *config.Config
	Name       string
	Params     []string
	Crud       bool
	Definition *routes.Definition
}

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

func (g *RouteDefinitionGenerator) Execute() error {
	if !tasks.Exists(".golem") {
		return errors.New(".golem directory does not exist, run from inside app directory")
	}
	if tasks.Exists(".golem/routes.yaml") {
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

func (g *RouteDefinitionGenerator) prepare() {
	rd := g.Definition.FindRoute(g.Name)
	if rd != nil {
		g.prepareModify(rd)
	} else {
		g.prepareCreate()
	}
}

func (g *RouteDefinitionGenerator) prepareModify(rd *routes.RouteDefinition) {

}

func (g *RouteDefinitionGenerator) prepareCreate() {
	parts := strings.Split(g.Name, "/")

	// ensure groups exists
	if g.Definition.Groups == nil {
		g.Definition.Groups = make(map[string]*routes.GroupDefinition)
	}

	// ensure group exists
	var gd *routes.GroupDefinition
	if g.Definition.Groups[parts[1]] != nil {
		gd = g.Definition.Groups[parts[1]]
	} else {
		gd = &routes.GroupDefinition{}
		if g.Definition.Groups == nil {
			g.Definition.Groups = make(map[string]*routes.GroupDefinition)
		}
		g.Definition.Groups[parts[1]] = gd
	}
	gd.Path = "/" + parts[1]

	// ensure route exists
	var rd *routes.RouteDefinition
	if gd.Routes[parts[2]] != nil {
		rd = gd.Routes[parts[2]]
	} else {
		rd = &routes.RouteDefinition{}
		if gd.Routes == nil {
			gd.Routes = make(map[string]*routes.RouteDefinition)
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
			rd.Params = append(rd.Params, &routes.ParamDefinition{Name: n, Type: t})
		}
	}
}

func (g *RouteDefinitionGenerator) routeExists(name string, routes []*routes.RouteDefinition) bool {
	if g.Definition.FindRoute(name) != nil {
		return true
	}
	return false
}
