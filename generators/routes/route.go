package routes

import (
	"bytes"
	"strings"

	"github.com/dashotv/golem/generators/app"

	"github.com/dashotv/golem/generators/templates"
	"github.com/dashotv/golem/tasks"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type RouteGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Definition *RouteDefinition
}

type RouteDefinition struct {
	Repo   string
	Name   string             `json:"name" yaml:"name"`
	Path   string             `json:"path" yaml:"path"`
	Method string             `json:"method" yaml:"method"`
	Routes []*RouteDefinition `json:"routes" yaml:"routes"`
	Params []*ParamDefinition `json:"params" yaml:"params"`
}

func (d *RouteDefinition) GetMethod() string {
	if d.Method != "" {
		return strings.Title(d.Method)
	}
	return "Get"
}

type ParamDefinition struct {
	Name    string `json:"name" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	Default string `json:"default" yaml:"default"`
}

func (d *ParamDefinition) GetType() string {
	if d.Type != "" {
		return strings.Title(d.Type)
	}
	return "String"
}

func NewRouteGenerator(cfg *config.Config, d *RouteDefinition) *RouteGenerator {
	d.Repo = cfg.Repo
	return &RouteGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		Generator: &base.Generator{
			Filename: cfg.Routes.Output + "/" + d.Name + "/routes.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (g *RouteGenerator) Execute() error {
	r := tasks.NewTaskRunner("generator:routes:route")

	r.Add("prepare", g.prepare)
	r.Add("make directory", func() error {
		return app.MakeDirectory(g.Config.Routes.Output + "/" + g.Definition.Name)
	})
	r.Add("template", func() error {
		return templates.New("routes").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)

	return r.Run()
}

func (g *RouteGenerator) prepare() error {
	return nil
}
