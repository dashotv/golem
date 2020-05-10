package routes

import (
	"bytes"

	"github.com/dashotv/golem/generators/templates"
	"github.com/dashotv/golem/tasks"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
)

type ServerGenerator struct {
	*base.Generator
	Config     *config.Config
	Output     string
	Definition *Definition
	data       map[string]string
}

//type ServerDefinition struct {
//	Package string             `json:"package" yaml:"package"`
//	Name    string             `json:"name" yaml:"name"`
//	Repo    string             `json:"repo" yaml:"repo"`
//	Routes  []*RouteDefinition `json:"routes" yaml:"routes"`
//}

func NewServerGenerator(cfg *config.Config, d *Definition) *ServerGenerator {
	return &ServerGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		data:       make(map[string]string),
		Generator: &base.Generator{
			Filename: cfg.Routes.Output + "/server.go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (g *ServerGenerator) Execute() error {
	r := tasks.NewTaskRunner("generator:routes:server")

	r.Add("prepare", g.prepare)
	r.Add("template", func() error {
		return templates.New("app_server_main").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)

	return r.Run()
}

func (g *ServerGenerator) prepare() error {
	//g.data["Name"] = g.Config.Routes.Name
	//g.data["Camel"] = strcase.ToCamel(g.Config.Routes.Name)
	//g.data["Repo"] = g.Config.Routes.Repo
	//err := g.prepareRoutes()
	//if err != nil {
	//	return err
	//}
	return nil
}

//func (g *ServerGenerator) prepareRoutes() error {
//
//}

//
//func getTemplate(name string) string {
//	filename := fmt.Sprintf("generators/routes/templates/%s.tgo", name)
//	data, err := Asset(filename)
//	if err != nil {
//		panic(err)
//	}
//	return string(data)
//}
