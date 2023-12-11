package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewGroup(cfg *config.Config, group *config.Group) error {
	runner := tasks.NewRunner("group")
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "routes")
	})
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Routes))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Routes, group.Name+".yaml"), group)
	})

	return runner.Run()
}

func NewRoute(cfg *config.Config, group string, route *config.Route) error {
	g := &config.Group{}
	err := tasks.ReadYaml(filepath.Join(cfg.Root(), cfg.Definitions.Routes, group+".yaml"), g)
	if err != nil {
		return errors.Wrap(err, "reading group")
	}

	g.AddRoute(route)

	runner := tasks.NewRunner("route")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Routes))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Routes, group+".yaml"), g)
	})

	return runner.Run()
}

func Routes(cfg *config.Config) error {
	dir := filepath.Join(cfg.Root(), cfg.Definitions.Routes)
	data := cfg.Data()
	var output []string

	runner := tasks.NewRunner("app:routes")
	runner.Add("header", func() error {
		header, err := tasks.Buffer(filepath.Join("app", "app_routes"), data)
		if err != nil {
			return errors.Wrap(err, "routes header")
		}
		output = append(output, header)
		return nil
	})

	// collect groups for route registration
	groups := make(map[string]*config.Group)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			group := &config.Group{}
			err := tasks.ReadYaml(path, group)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading group: %s", path))
			}

			groups[group.Name] = group
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "walking routes")
	}

	routes := struct {
		Config *config.Config
		Groups map[string]*config.Group
	}{
		Config: cfg,
		Groups: groups,
	}

	runner.Add("registration", func() error {
		buf, err := tasks.Buffer(filepath.Join("app", "partial_routes"), routes)
		if err != nil {
			return errors.Wrap(err, "routes registration buffer")
		}
		output = append(output, buf)
		return nil
	})

	for k, v := range groups {
		k := k
		v := v
		runner.Add("group:"+k, func() error {
			buf, err := tasks.Buffer(filepath.Join("app", "partial_route"), v)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("routes buffer: %s", k))
			}
			output = append(output, buf)
			return nil
		})
		runner.Add("hook:"+k, func() error {
			d := struct {
				Package string
				Group   *config.Group
			}{
				Package: cfg.Package,
				Group:   v,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "routes"), filepath.Join("app", "routes_"+k+".go"), d)
		})
	}

	runner.Add("save", func() error {
		return tasks.RawFile(filepath.Join("app", "app_routes.go"), strings.Join(output, "\n"))
	})

	return runner.Run()
}
