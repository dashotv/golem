package generators

import (
	"path/filepath"
	"strings"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
	"github.com/dashotv/golem/templates"
)

func App(cfg *config.Config) error {
	g := cfg.Data()

	runner := tasks.NewRunner("app")

	if cfg.Plugins.Cache {
		runner.Add("cache", func() error {
			f := filepath.Join("app", "app_cache")
			return tasks.File(f, f+".go", g)
		})
	}
	if cfg.Plugins.Routes {
		runner.Add("routes", func() error {
			return Routes(cfg)
		})
	}
	if cfg.Plugins.Models {
		runner.Add("models", func() error {
			return Models(cfg)
		})
	}
	if cfg.Plugins.Events {
		runner.Add("events", func() error {
			return Events(cfg)
		})
	}
	if cfg.Plugins.Workers {
		runner.Add("workers", func() error {
			return Workers(cfg)
		})
	}
	if cfg.Plugins.Clients {
		runner.Add("clients", func() error {
			return Clients(cfg)
		})
	}

	runner.Group("app").Add("modify", func() error {
		return tasks.Modify(filepath.Join("app", "app.go"), g)
	})
	runner.Group("config").Add("modify", func() error {
		return tasks.Modify(filepath.Join("app", "app_config.go"), g)
	})
	hooks := runner.Group("hooks")
	hooks.Add("directory", func() error {
		return tasks.Directory("hooks")
	})
	files, err := templates.ReadDir("hooks")
	if err != nil {
		return err
	}
	for _, f := range files {
		hooks.Add(f, func() error {
			n := strings.TrimSuffix(f, ".tgo")
			d := filepath.Join("hooks", n)
			return tasks.FileDoesntExist(d, d, g)
		})
	}

	commands := runner.Group("commands")
	commands.Add("go mod tidy", func() error {
		return tasks.GoModTidy()
	})
	// TODO: this breaks the echo import in app/app.go on first run
	// and I don't know why
	commands.Add("goimports", func() error {
		return tasks.GoImports()
	})

	return runner.Run()
}
