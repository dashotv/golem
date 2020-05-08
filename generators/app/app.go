package app

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

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
	runner := tasks.NewTaskRunner()
	runner.Add("make app directory", func() error {
		return makeDirectory(g.Name)
	})
	runner.Add("create default config", func() error {
		return writeDefaultConfig(g.Name + "/.golem")
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
	runner.Add("make config directory", func() error {
		return makeDirectory(g.Name + "/config")
	})
	runner.Add("create application config", func() error {
		d := map[string]string{}
		g := base.NewFileGenerator(g.Config, "app_config_config", g.Name+"/config/config.go", d)
		return g.Execute()
	})
	runner.Add("make command directory", func() error {
		return makeDirectory(g.Name + "/cmd")
	})
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
	runner.Add("make server directory", func() error {
		return makeDirectory(g.Name + "/server")
	})
	runner.Add("create server main", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		g := base.NewFileGenerator(g.Config, "app_server_main", g.Name+"/server/main.go", d)
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

func makeDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir")
		}
	}
	return nil
}

func executeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func writeDefaultConfig(dir string) error {
	cfg := defaultConfig()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir")
		}
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "could not marshal config")
	}

	err = ioutil.WriteFile(dir+"/.golem.yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}

type defaultAppConfig struct {
	Mode string
	Port int
}

func writeAppConfig(name string) error {
	cfg := &defaultAppConfig{Mode: "dev", Port: 3000}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "could not marshal config")
	}

	err = ioutil.WriteFile(name+"/."+name+".yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}

func defaultConfig() *config.Config {
	cfg := &config.Config{}
	cfg.Models.Enabled = true
	cfg.Models.Package = "models"
	cfg.Models.Output = "./models"
	cfg.Models.Definitions = "./.golem/models"

	cfg.Routes.Enabled = true
	cfg.Routes.Output = "./server"
	cfg.Routes.Definition = "./.golem/routes.yaml"
	cfg.Routes.Name = "Blarg"
	cfg.Routes.Repo = "github.com/dashotv/blarg"

	cfg.Jobs.Enabled = true
	cfg.Jobs.Package = "jobs"
	cfg.Jobs.Output = "./jobs"
	cfg.Jobs.Definitions = "./.golem/jobs.yaml"

	return cfg
}
