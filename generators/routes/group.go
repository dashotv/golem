package routes

import (
	"bytes"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/generators/templates"
	"github.com/dashotv/golem/tasks"
)

// GroupGenerator manages the creation of group routes
type GroupGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Name       string
	Definition *Group
}

// NewGroupGenerator configures and returns an instance of GroupGenerator
func NewGroupGenerator(cfg *config.Config, name string, d *Group) *GroupGenerator {
	return &GroupGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Name:       name,
		Definition: d,
		Generator: &base.Generator{
			Filename: cfg.Routes.Output + "/" + name + "/routes.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute manages creation of group routes files with the template
func (g *GroupGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:" + g.Name)
	dir := g.Config.Routes.Output + "/" + g.Definition.Name
	r.Add("prepare", g.prepare)
	r.Add("make directory: "+dir, tasks.NewMakeDirectoryTask(dir))
	if g.Definition.Rest == true {
		r.Add("template", func() error {
			return templates.New("routes_group_rest").Execute(g.Buffer, g.Definition)
		})
	} else {
		r.Add("template", func() error {
			return templates.New("routes_group").Execute(g.Buffer, g.Definition)
		})
	}
	r.Add("write: "+g.Filename, g.Write)
	r.Add("format: "+g.Filename, g.Format)

	return r.Run()
}

// prepare configures the data for the template
func (g *GroupGenerator) prepare() error {
	for n, r := range g.Definition.Routes {
		r.Name = n
	}
	return nil
}
