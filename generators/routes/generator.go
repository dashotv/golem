package routes

import (
	"bytes"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type Generator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Definition *Definition
	data       map[string]string
}

type Definition struct {
	Package string
	Routes  []*RouteDefinition
}

type RouteDefinition struct {
	Path   string
	Routes []*RouteDefinition
	Params []*ParamDefinition
}

type ParamDefinition struct {
	Name    string
	Type    string
	Default string
}

func NewGenerator(cfg *config.Config) (*Generator, error) {
	d := &Definition{}
	g := &Generator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		data:       make(map[string]string),
		Generator: &base.Generator{
			Filename: cfg.Routes.Output + "/server.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}

	err := base.ReadYaml(cfg.Routes.Definition, d)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Generator) Execute() error {
	err := g.prepare()
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) prepare() error {
	g.data["Name"] = g.Config.Routes.Name
	g.data["Camel"] = strcase.ToCamel(g.Config.Routes.Name)
	g.data["Repo"] = g.Config.Routes.Repo
	//g.data["Routes"] = g.Definition
	return nil
}

//
//func getTemplate(name string) string {
//	filename := fmt.Sprintf("generators/routes/templates/%s.tgo", name)
//	data, err := Asset(filename)
//	if err != nil {
//		panic(err)
//	}
//	return string(data)
//}
