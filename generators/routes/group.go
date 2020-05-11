package routes

import (
	"bytes"
	"strings"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/generators/templates"
	"github.com/dashotv/golem/tasks"
)

type GroupGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Name       string
	Definition *GroupDefinition
}

type GroupDefinition struct {
	Repo   string
	Name   string                      `json:"name" yaml:"name"`
	Path   string                      `json:"path" yaml:"path"`
	Method string                      `json:"method" yaml:"method"`
	Routes map[string]*RouteDefinition `json:"routes" yaml:"routes"`
}

type RouteDefinition struct {
	Repo   string
	Name   string             `json:"name" yaml:"name"`
	Path   string             `json:"path" yaml:"path"`
	Method string             `json:"method" yaml:"method"`
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

func NewGroupGenerator(cfg *config.Config, name string, d *GroupDefinition) *GroupGenerator {
	d.Name = name
	d.Repo = cfg.Repo
	return &GroupGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Name:       name,
		Definition: d,
		Generator: &base.Generator{
			Filename: cfg.Routes.Output + "/" + d.Name + "/routes.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (g *GroupGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:group")

	r.Add("prepare", g.prepare)
	r.Add("make directory", tasks.NewMakeDirectoryTask(g.Config.Routes.Output+"/"+g.Definition.Name))
	r.Add("template", func() error {
		return templates.New("routes_group").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)

	return r.Run()
}

func (g *GroupGenerator) prepare() error {
	return nil
}
