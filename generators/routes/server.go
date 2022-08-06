package routes

import (
	"bytes"
	"github.com/dashotv/golem/templates"

	"github.com/dashotv/golem/tasks"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

// ServerGenerator manages the creation of the server file
type ServerGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Definition *Definition
}

// NewServerGenerator creates and returns an instance of ServerGenerator
func NewServerGenerator(cfg *config.Config, d *Definition) *ServerGenerator {
	return &ServerGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		Generator: &base.Generator{
			Filename: cfg.Path(cfg.Routes.Output, "server.go"),
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute creates the server file from the template
func (g *ServerGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:server")

	r.Add("directory", tasks.NewMakeDirectoryTask("server"))
	r.Add("prepare", g.prepare)
	r.Add("template", func() error {
		return templates.New("routes_server").Execute(g.Buffer, g.Definition)
	})
	r.Add("write "+g.Filename, g.Write)
	r.Add("format", g.Format)
	r.Add("server routes", func() error {
		srg := NewServerRoutesGenerator(g.Config, g.Definition)
		return srg.Execute()
	})

	return r.Run()
}

// prepare the configuration for the template
func (g *ServerGenerator) prepare() error {
	return nil
}

// ServerRoutesGenerator manages the creation of the server routes file
type ServerRoutesGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Definition *Definition
}

// NewServerRoutesGenerator creates and returns an instance of ServerRoutesGenerator
func NewServerRoutesGenerator(cfg *config.Config, d *Definition) *ServerRoutesGenerator {
	return &ServerRoutesGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		Generator: &base.Generator{
			Filename: cfg.Path(cfg.Routes.Output, "routes.go"),
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute creates the server routes file from the template
func (g *ServerRoutesGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:server routes")

	r.Add("prepare", g.prepare)
	r.Add("template", func() error {
		return templates.New("routes_routes").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)

	return r.Run()
}

// prepare the configuration for the template
func (g *ServerRoutesGenerator) prepare() error {
	return nil
}
