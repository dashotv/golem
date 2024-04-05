package generators

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewClient(cfg *config.Config, client *config.Client) error {
	runner := tasks.NewRunner("client")
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "clients")
	})
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Clients))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Clients, client.Language+".yaml"), client)
	})

	return runner.Run()
}

func Clients(cfg *config.Config) error {
	clients, err := cfg.Clients()
	if err != nil {
		return fae.Wrap(err, "collecting clients")
	}

	runner := tasks.NewRunner("app:clients")
	for lang, c := range clients {
		switch lang {
		case "go":
			runner.Add("client:"+lang, func() error {
				return clientGolang(cfg, c)
			})
		case "typescript":
			runner.Add("client:"+lang, func() error {
				return clientTypescript(cfg, c)
			})
		}
	}

	return runner.Run()
}

func clientGolang(cfg *config.Config, client *config.Client) error {
	var modelsOutput []string

	groups, err := cfg.Groups()
	if err != nil {
		return fae.Wrap(err, "collecting groups")
	}

	data := struct {
		Config map[string]string
		Groups map[string]*config.Group
	}{
		Config: cfg.Data(),
		Groups: groups,
	}

	// collect models for connector registration
	models, err := cfg.Models()
	if err != nil {
		return fae.Wrap(err, "collecting models")
	}

	keys := make([]string, 0, len(models))
	for k := range models {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	runner := tasks.NewRunner("client:golang")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), filepath.Dir(client.Output())))
	})
	runner.Add("client", func() error {
		return tasks.File(filepath.Join("client", "go", "client"), filepath.Join(cfg.Root(), client.Output(), "client.gen.go"), data)
	})
	for _, g := range groups {
		g := g
		groupData := struct {
			Config map[string]string
			Group  *config.Group
		}{
			Config: cfg.Data(),
			Group:  g,
		}
		runner.Add("route:"+g.Name, func() error {
			return tasks.File(filepath.Join("client", "go", "group"), filepath.Join(cfg.Root(), client.Output(), g.Name+".gen.go"), groupData)
		})
	}

	runner.Add("header", func() error {
		header, err := tasks.Buffer(filepath.Join("client", "go", "models"), data)
		if err != nil {
			return fae.Wrap(err, "models header")
		}
		modelsOutput = append(modelsOutput, header)
		return nil
	})

	for _, k := range keys {
		k := k
		v := models[k]
		runner.Add("model:"+k, func() error {
			t := "app/partial_model"
			if v.Type == "struct" {
				t = "app/partial_struct"
			}
			buf, err := tasks.Buffer(t, v)
			if err != nil {
				return fae.Wrap(err, fmt.Sprintf("model buffer: %s", k))
			}
			modelsOutput = append(modelsOutput, buf)
			return nil
		})
		if v.Type == "struct" {
			continue
		}
	}

	runner.Add("save", func() error {
		return tasks.RawFile(filepath.Join("client", "models.gen.go"), strings.Join(modelsOutput, "\n"))
	})

	return runner.Run()
}

func clientTypescript(cfg *config.Config, client *config.Client) error {
	groups, err := cfg.Groups()
	if err != nil {
		return fae.Wrap(err, "collecting groups")
	}

	groupData := struct {
		Config map[string]string
		Groups map[string]*config.Group
	}{
		Config: cfg.Data(),
		Groups: groups,
	}

	runner := tasks.NewRunner("typescript")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), filepath.Dir(client.Output())))
	})
	runner.Add("client", func() error {
		return tasks.File(filepath.Join("client", "typescript", "client"), filepath.Join(cfg.Root(), client.Output(), "client.gen.ts"), groupData)
	})
	runner.Add("index", func() error {
		return tasks.File(filepath.Join("client", "typescript", "index"), filepath.Join(cfg.Root(), client.Output(), "index.ts"), groupData)
	})
	for _, g := range groups {
		g := g
		groupData := struct {
			Config map[string]string
			Group  *config.Group
		}{
			Config: cfg.Data(),
			Group:  g,
		}
		runner.Add("route:"+g.Name, func() error {
			return tasks.File(filepath.Join("client", "typescript", "group"), filepath.Join(cfg.Root(), client.Output(), g.Name+".gen.ts"), groupData)
		})
	}

	// collect models for connector registration
	models, err := cfg.Models()
	if err != nil {
		return fae.Wrap(err, "collecting models")
	}

	keys := make([]string, 0, len(models))
	for k := range models {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	modelData := struct {
		Config map[string]string
		Models map[string]*config.Model
	}{
		Config: cfg.Data(),
		Models: models,
	}

	runner.Add("models", func() error {
		return tasks.File(filepath.Join("client", "typescript", "models"), filepath.Join(cfg.Root(), client.Output(), "models.gen.ts"), modelData)
	})

	list := []string{}
	for _, g := range groups {
		list = append(list, g.TypescriptPackages()...)
	}
	for _, m := range models {
		list = append(list, m.TypescriptImports()...)
	}

	for _, p := range list {
		p := p
		packageData := map[string]string{"Package": p}
		runner.Add("package:index", func() error {
			return tasks.FileDoesntExist(filepath.Join("client", "typescript", "package_index"), filepath.Join(cfg.Root(), client.Output(), p, "index.ts"), packageData)
		})
		runner.Add("package:"+p, func() error {
			return tasks.FileDoesntExist(filepath.Join("client", "typescript", "package"), filepath.Join(cfg.Root(), client.Output(), p, p+".ts"), packageData)
		})
	}

	runner.Add("prettier", func() error {
		return tasks.Prettier(filepath.Join(cfg.Root(), client.Output()))
	})

	return runner.Run()
}
