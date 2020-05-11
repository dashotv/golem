package app

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

type Generator struct {
	*base.Generator
	Config *config.Config
	Name   string
	Repo   string
}

func NewAppGenerator(cfg *config.Config, name, repo string) *Generator {
	return &Generator{Config: cfg, Name: name, Repo: repo}
}

func (g *Generator) Execute() error {
	runner := tasks.NewRunner("generator:app")
	runner.Add("make app directory", tasks.NewMakeDirectoryTask(g.Name))
	runner.Add("create default config", func() error {
		cfg := config.DefaultConfig()
		cfg.Name = g.Name
		cfg.Repo = g.Repo
		return writeConfig(g.Name+"/.golem", cfg)
	})
	runner.Add("create application config", func() error {
		return writeAppConfig(g.Name)
	})
	runner.Add("create application main", func() error {
		d := map[string]string{"Repo": g.Repo}
		g := base.NewFileGenerator(g.Config, "app_main", g.Name+"/main.go", d)
		return g.Execute()
	})
	runner.Add("create application license", func() error {
		d := map[string]string{}
		g := base.NewFileGenerator(g.Config, "app_license", g.Name+"/LICENSE", d)
		return g.Execute()
	})
	runner.Add("make config directory", tasks.NewMakeDirectoryTask(g.Name+"/config"))
	runner.Add("create application config", func() error {
		d := map[string]string{}
		g := base.NewFileGenerator(g.Config, "app_config_config", g.Name+"/config/config.go", d)
		return g.Execute()
	})
	runner.Add("make command directory", tasks.NewMakeDirectoryTask(g.Name+"/cmd"))
	runner.Add("create application root command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		g := base.NewFileGenerator(g.Config, "app_cmd_root", g.Name+"/cmd/root.go", d)
		return g.Execute()
	})
	runner.Add("create application server command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		g := base.NewFileGenerator(g.Config, "app_cmd_server", g.Name+"/cmd/server.go", d)
		return g.Execute()
	})
	runner.Add("make server directory", tasks.NewMakeDirectoryTask(g.Name+"/server"))
	runner.Add("make server directory keep file", func() error {
		g := base.NewFileGenerator(g.Config, "keep", g.Name+"/server/.keep", map[string]string{})
		return g.Execute()
	})
	//runner.Add("create server main", func() error {
	//	d := map[string]string{"Name": g.Name, "Repo": g.Repo}
	//	g := base.NewFileGenerator(g.Config, "app_server_main", g.Name+"/server/main.go", d)
	//	return g.Execute()
	//})
	runner.Add("make models directory", tasks.NewMakeDirectoryTask(g.Name+"/models"))
	runner.Add("make models directory keep file", func() error {
		g := base.NewFileGenerator(g.Config, "keep", g.Name+"/models/.keep", map[string]string{})
		return g.Execute()
	})

	err := runner.Run()
	if err != nil {
		return err
	}

	// TODO: generate example model definition => .golem/models/example.yaml
	// TODO: generate example route definition => .golem/routes.yaml
	// TODO: generate example job definition   => .golem/jobs.yaml

	// TODO: run golem generate

	return nil
}
