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

// Execute manages creation of group routes files with the template
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

// prepare configures the data for the template
func (g *GroupGenerator) prepare() error {
	return nil
}
