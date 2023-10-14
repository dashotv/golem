package generators

import (
	"bytes"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

// AppGenerator manages generating an application
type AppGenerator struct {
	*Generator
	Config *config.Config
	Name   string
	Repo   string
}

// NewAppGenerator returns a new instance of Generator
func NewAppGenerator(cfg *config.Config, name, repo string) *AppGenerator {
	return &AppGenerator{Config: cfg, Name: name, Repo: repo}
}

// Execute processes all of the configurations and generates an application
func (g *AppGenerator) Execute() error {
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
		err := ReadYaml(g.Name+"/.golem/.golem.yaml", g.Config)
		if err != nil {
			return err
		}
		g.Config.File = g.Name + "/.golem/.golem.yaml"
		return nil
	})
	runner.Add("create config file", func() error {
		return writeAppConfig(g.Name, g.Name)
	})
	runner.Add("create main", func() error {
		d := map[string]string{"Repo": g.Repo}
		fg := NewFileGenerator(g.Config, "main", g.Name+"/main.go", d)
		return fg.Execute()
	})
	runner.Add("create application license", func() error {
		d := map[string]string{}
		fg := NewFileGenerator(g.Config, "license", g.Name+"/LICENSE", d)
		return fg.Execute()
	})
	runner.Add("make application directory", tasks.NewMakeDirectoryTask(g.Name+"/app"))
	runner.Add("make application app", func() error {
		data := map[string]string{"Repo": g.Config.Repo}
		fg := NewFileGenerator(g.Config, "app/app", g.Name+"/app/app.go", data)
		return fg.Execute()
	})
	runner.Add("make config directory", tasks.NewMakeDirectoryTask(g.Name+"/app"))
	runner.Add("create application config", func() error {
		d := map[string]string{}
		fg := NewFileGenerator(g.Config, "app/config", g.Name+"/app/config.go", d)
		return fg.Execute()
	})
	runner.Add("create cron config", func() error {
		d := map[string]string{}
		fg := NewFileGenerator(g.Config, "app/cron", g.Name+"/app/cron.go", d)
		return fg.Execute()
	})
	runner.Add("create cache config", func() error {
		d := map[string]string{}
		fg := NewFileGenerator(g.Config, "app/cache", g.Name+"/app/cache.go", d)
		return fg.Execute()
	})
	runner.Add("make command directory", tasks.NewMakeDirectoryTask(g.Name+"/cmd"))
	runner.Add("create application root command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		fg := NewFileGenerator(g.Config, "cmd/root", g.Name+"/cmd/root.go", d)
		return fg.Execute()
	})
	runner.Add("create application server command", func() error {
		d := map[string]string{"Name": g.Name, "Repo": g.Repo}
		fg := NewFileGenerator(g.Config, "cmd/server", g.Name+"/cmd/server.go", d)
		return fg.Execute()
	})
	runner.Add("create makefile", func() error {
		fg := NewFileGenerator(g.Config, "makefile", g.Name+"/Makefile", map[string]string{"Name": g.Name})
		return fg.Execute()
	})
	runner.Add("create dockerfile", func() error {
		fg := NewFileGenerator(g.Config, "dockerfile", g.Name+"/Dockerfile", map[string]string{"Name": g.Name})
		return fg.Execute()
	})
	runner.Add("create docker-compose file", func() error {
		fg := NewFileGenerator(g.Config, "docker-compose", g.Name+"/docker-compose.yml", map[string]string{"Name": g.Name, "Port": "3000"})
		return fg.Execute()
	})
	runner.Add("create air config file", func() error {
		fg := NewFileGenerator(g.Config, ".air.toml", g.Name+"/.air.toml", map[string]string{"Name": g.Name})
		return fg.Execute()
	})
	runner.Add("create drone file", func() error {
		fg := NewFileGenerator(g.Config, ".drone.yml", g.Name+"/.drone.yml", map[string]string{"Name": g.Name})
		return fg.Execute()
	})
	runner.Add("make etc directory", tasks.NewMakeDirectoryTask(g.Name+"/etc"))
	runner.Add("create service config", func() error {
		fg := NewFileGenerator(g.Config, "etc/service", g.Name+"/etc/"+g.Name+".service", map[string]string{"Name": g.Name})
		return fg.Execute()
	})
	runner.Add("create production config file", func() error {
		return writeAppConfig(g.Name, g.Name+"/etc")
	})
	runner.Add("create server main", func() error {
		if tasks.PathExists(g.Name + "/app/server.go") {
			return nil
		}
		d := &Definition{Name: g.Name, Repo: g.Repo}
		sg := NewServerGenerator(g.Config, d)
		return sg.Execute()
	})

	err := runner.Run()
	if err != nil {
		return err
	}

	return nil
}

// ModelDefinitionGenerator manages the generation of model definitions
type ModelDefinitionGenerator struct {
	*Generator
	Config     *config.Config
	Name       string
	Type       string
	Fields     []string
	Definition *Model
}

// NewModelDefinitionGenerator creates and returns an instance of ModelDefinitionGenerator
func NewModelDefinitionGenerator(cfg *config.Config, name string, fields ...string) *ModelDefinitionGenerator {
	return &ModelDefinitionGenerator{
		Config:     cfg,
		Name:       name,
		Type:       "model",
		Fields:     fields,
		Definition: &Model{},
		Generator: &Generator{
			Filename: ".golem/models/" + name + ".yaml",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute generates the model definition
func (g *ModelDefinitionGenerator) Execute() error {
	if !tasks.PathExists(".golem") {
		return errors.New(".golem directory does not exist, run from inside app directory")
	}
	if tasks.PathExists(g.Filename) {
		return errors.New("model definition already exists: " + g.Filename)
	}

	g.prepare()

	runner := tasks.NewRunner("generator:model")
	runner.Add("ensure models directory exists", tasks.NewMakeDirectoryTask(".golem/models"))
	runner.Add("generate model definition", func() error {
		data, err := yaml.Marshal(g.Definition)
		if err != nil {
			return err
		}
		return os.WriteFile(g.Filename, data, 0644)
	})

	return runner.Run()
}

// prepare the definition and data
func (g *ModelDefinitionGenerator) prepare() {
	g.Definition.Name = g.Name
	g.Definition.Type = g.Type
	g.Definition.Fields = make([]*Field, 0)

	for _, f := range g.Fields {
		f := strings.Split(f, ":")
		n := f[0]
		t := "string"
		if len(f) > 1 {
			t = f[1]
		}

		g.Definition.Fields = append(g.Definition.Fields, &Field{Name: n, Type: t})
	}
}

// RouteDefinitionGenerator manages the generation of and updates to routes definition
type RouteDefinitionGenerator struct {
	*Generator
	Config     *config.Config
	Name       string
	Path       string
	Params     []string
	Rest       bool
	Definition *Definition
}

// NewRouteDefinitionGenerator creates and returns an instance of RouteDefinitionGenerator
func NewRouteDefinitionGenerator(cfg *config.Config, path string, rest bool, params ...string) *RouteDefinitionGenerator {
	if path[0] != '/' {
		path = "/" + path
	}
	parts := strings.Split(path, "/")
	name := parts[1]

	return &RouteDefinitionGenerator{
		Config:     cfg,
		Name:       name,
		Path:       path,
		Params:     params,
		Definition: &Definition{},
		Rest:       rest,
		Generator: &Generator{
			Filename: ".golem/routes.yaml",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute generates or updates the routes definition
func (g *RouteDefinitionGenerator) Execute() error {
	if !tasks.PathExists(".golem") {
		return errors.New(".golem directory does not exist, run from inside app directory")
	}
	if tasks.PathExists(".golem/routes.yaml") {
		if err := ReadYaml(".golem/routes.yaml", g.Definition); err != nil {
			return err
		}
	}

	g.Definition.Name = g.Config.Name
	g.Definition.Repo = g.Config.Repo

	g.prepare()

	runner := tasks.NewRunner("generator:route")
	runner.Add("generate route definition", func() error {
		data, err := yaml.Marshal(g.Definition)
		if err != nil {
			return err
		}

		return os.WriteFile(".golem/routes.yaml", data, 0644)
	})

	return runner.Run()
}

// prepare the data for generation
func (g *RouteDefinitionGenerator) prepare() {
	group := g.Definition.FindGroup(g.Name, g.Rest)
	if g.Rest {
		return
	}

	rd := g.Definition.FindRoute(g.Path)

	if g.Rest {
		group.Rest = true
		return
	}

	if len(g.Params) == 0 {
		return
	}

	for _, f := range g.Params {
		f := strings.Split(f, ":")
		n := f[0]
		t := "string"
		if len(f) > 1 {
			t = f[1]
		}
		if rd.Params == nil {
			rd.Params = make([]*Param, 0)
		}
		rd.Params = append(rd.Params, &Param{Name: n, Type: t})
	}
}

func writeConfig(dir string, cfg *config.Config) error {
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

	err = os.WriteFile(dir+"/.golem.yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}

type defaultAppConfig struct {
	Mode        string
	Port        int
	Connections map[string]*Connection
}

type Connection struct {
	URI        string
	Database   string `json:"database,omitempty" yaml:"database,omitempty"`
	Collection string `json:"collection,omitempty" yaml:"collection,omitempty"`
}

func writeAppConfig(name string, dir string) error {
	cfg := &defaultAppConfig{Mode: "dev", Port: 3000}

	cfg.Connections = make(map[string]*Connection)
	cfg.Connections["default"] = &Connection{
		URI:      "mongodb://localhost:27017",
		Database: "database",
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "could not marshal config")
	}

	err = os.WriteFile(dir+"/."+name+".yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}
