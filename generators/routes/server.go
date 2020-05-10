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
	Definition *ServerDefinition
	data       map[string]string
}

//type ServerGeneratorData struct {
//	Name   string
//	Camel  string
//	Repo   string
//	Routes []*RouteGeneratorData
//}

type ServerDefinition struct {
	Package string
	Name    string
	Repo    string
	Routes  []*RouteDefinition
}

func NewServerGenerator(cfg *config.Config) (*ServerGenerator, error) {
	d := &ServerDefinition{Name: cfg.Name, Package: "server", Repo: cfg.Repo}
	g := &ServerGenerator{
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
