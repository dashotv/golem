package app

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/generators/routes"
	"github.com/dashotv/golem/tasks"
)

// Generator manages generating an application
type Generator struct {
	*base.Generator
	Config *config.Config
	Name   string
	Repo   string
}

// NewGenerator returns a new instance of Generator
func NewGenerator(cfg *config.Config, name, repo string) *Generator {
	return &Generator{Config: cfg, Name: name, Repo: repo}
}

// Execute processes all of the configurations and generates an application
func (g *Generator) Execute() error {
	runner := tasks.NewRunner("generator:app")
	runner.Add("make app directory", tasks.NewMakeDirectoryTask(g.Name))
	runner.Add("create default config", func() error {
		cfg := config.DefaultConfig()
		cfg.Name = g.Name
		cfg.Repo = g.Repo
		cfg.Routes.Repo = g.Repo
		return writeConfig(g.Name+"/.golem", cfg)
	})
	runner.Add("load default config", func() error {
		err := base.ReadYaml(g.Name+"/.golem/.golem.yaml", g.Config)
		if err != nil {
			return err
		}
		g.Config.File = g.Name + "/.golem/.golem.yaml"
		return nil
	})
	runner.Add("create application config", func() error {
		return writeAppConfig(g.Name)
	})
	runner.Add("create application main", func() error {
		d := map[string]string{"Repo": g.Repo}
		fg := base.NewFileGenerator(g.Config, "app_main", g.Name+"/main.go", d)
		return fg.Execute()
	})
	runner.Add("create application license", func() error {
		d := map[string]string{}
		fg := base.NewFileGenerator(g.Config, "app_license", g.Name+"/LICENSE", d)
		return fg.Execute()
	})
	runner.Add("make config directory", tasks.NewMakeDirectoryTask(g.Name+"/config"))
	runner.Add("create application config", func() error {
		d := map[string]string{}
		fg := base.NewFileGenerator(g.Config, "app_config_config", g.Name+"/config/config.go", d)
		return fg.Execute()
	})
	runner.Add("make command directory", tasks.NewMakeDirectoryTask(g.Name+"/cmd"))
	runner.Add("create application root command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		fg := base.NewFileGenerator(g.Config, "app_cmd_root", g.Name+"/cmd/root.go", d)
		return fg.Execute()
	})
	runner.Add("create application server command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		fg := base.NewFileGenerator(g.Config, "app_cmd_server", g.Name+"/cmd/server.go", d)
		return fg.Execute()
	})
	runner.Add("make server directory", tasks.NewMakeDirectoryTask(g.Name+"/server"))
	runner.Add("make server directory keep file", func() error {
		fg := base.NewFileGenerator(g.Config, "keep", g.Name+"/server/.keep", map[string]string{})
		return fg.Execute()
	})
	runner.Add("create server main", func() error {
		if tasks.PathExists(g.Name + "/server/server.go") {
			return nil
		}
		d := &routes.Definition{Name: g.Name, Repo: g.Repo}
		sg := routes.NewServerGenerator(g.Config, d)
		return sg.Execute()
	})
	runner.Add("make models directory", tasks.NewMakeDirectoryTask(g.Name+"/models"))
	runner.Add("make models directory keep file", func() error {
		fg := base.NewFileGenerator(g.Config, "keep", g.Name+"/models/.keep", map[string]string{})
		return fg.Execute()
	})
	runner.Add("make application directory", tasks.NewMakeDirectoryTask(g.Name+"/application"))
	runner.Add("make application app", func() error {
		data := map[string]string{"Repo": g.Config.Repo}
		fg := base.NewFileGenerator(g.Config, "app_application", g.Name+"/application/app.go", data)
		return fg.Execute()
	})

	err := runner.Run()
	if err != nil {
		return err
	}

	return nil
}
