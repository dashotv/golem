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
	var modelsOutput []string

	clients, err := cfg.Clients()
	if err != nil {
		return fae.Wrap(err, "collecting clients")
	}

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

	runner := tasks.NewRunner("app:clients")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), "client")) // TODO: handle this for other languages
	})

	runner.Add("header", func() error {
		header, err := tasks.Buffer(filepath.Join("client", "models"), data)
		if err != nil {
			return fae.Wrap(err, "models header")
		}
		modelsOutput = append(modelsOutput, header)
		return nil
	})

	for _, c := range clients {
		c := c
		runner.Add("wut", func() error {
			return tasks.File(filepath.Join("client", "client"), filepath.Join(cfg.Root(), c.Filename()), data)
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
				return tasks.File(filepath.Join("client", "group"), filepath.Join(cfg.Root(), "client", g.Name+".gen.go"), groupData)
			})
		}

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
	}

	return runner.Run()
}
