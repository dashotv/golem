package generators

import (
	"path/filepath"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func EnablePlugin(cfg *config.Config, name string) error {
	return pluginEnable(cfg, name)
}

func DisablePlugin(cfg *config.Config, name string) error {
	return pluginDisable(cfg, name)
}

func pluginEnable(cfg *config.Config, name string) error {
	if cfg.Enabled(name) {
		return nil
	}
	runner := tasks.NewRunner("plugin:enable")
	runner.Add("config", func() error {
		return cfg.Enable(name)
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(".golem", "config.yaml"), cfg)
	})
	env := runner.Group("env")
	env.Add("dev", func() error {
		return tasks.AppendFile(".env_partial_"+name, ".env", cfg)
	})
	env.Add("prod", func() error {
		return tasks.AppendFile(".env_partial_"+name, ".env.production", cfg)
	})
	return runner.Run()
}

func pluginDisable(cfg *config.Config, name string) error {
	if !cfg.Enabled(name) {
		return nil
	}
	runner := tasks.NewRunner("plugin:disable")
	runner.Add("disable", func() error {
		return cfg.Disable(name)
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(".golem", "config.yaml"), cfg)
	})
	return runner.Run()
}
