package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/dashotv/golem/tasks"
)

type Config struct {
	File    string `yaml:"-"`
	Name    string `yaml:"name,omitempty"`
	Repo    string `yaml:"repo,omitempty"`
	Package string `yaml:"package,omitempty"`
	Output  string `yaml:"output,omitempty"`
	Plugins struct {
		Models  bool `yaml:"models"`
		Routes  bool `yaml:"routes"`
		Workers bool `yaml:"workers"`
		Cache   bool `yaml:"cache"`
		Events  bool `yaml:"events"`
	} `yaml:"plugins"`
	Definitions struct {
		Models  string `yaml:"models,omitempty"`
		Routes  string `yaml:"routes,omitempty"`
		Events  string `yaml:"events,omitempty"`
		Workers string `yaml:"workers,omitempty"`
	} `yaml:"definitions,omitempty"`
}

func (c *Config) Data() map[string]string {
	return map[string]string{
		"Name":    c.Name,
		"Camel":   strcase.ToCamel(c.Name),
		"Repo":    c.Repo,
		"Package": c.Package,
		"Output":  c.Output,
		"Models":  fmt.Sprintf("%t", c.Plugins.Models),
		"Routes":  fmt.Sprintf("%t", c.Plugins.Routes),
		"Cache":   fmt.Sprintf("%t", c.Plugins.Cache),
		"Events":  fmt.Sprintf("%t", c.Plugins.Events),
		"Workers": fmt.Sprintf("%t", c.Plugins.Workers),
	}
}

func (c *Config) Enable(name string) error {
	switch name {
	case "models", "Models":
		c.Plugins.Models = true
	case "routes", "Routes":
		c.Plugins.Routes = true
	case "cache", "Cache":
		c.Plugins.Cache = true
	case "events", "Events":
		c.Plugins.Events = true
	case "workers", "Workers":
		c.Plugins.Workers = true
	default:
		return fmt.Errorf("unknown plugin: %s", name)
	}

	return nil
}

func (c *Config) Enabled(name string) bool {
	switch name {
	case "models", "Models":
		return c.Plugins.Models
	case "routes", "Routes":
		return c.Plugins.Routes
	case "cache", "Cache":
		return c.Plugins.Cache
	case "events", "Events":
		return c.Plugins.Events
	case "workers", "Workers":
		return c.Plugins.Workers
	default:
		return false
	}
}

func (c *Config) Disable(name string) error {
	switch name {
	case "models", "Models":
		c.Plugins.Models = false
	case "routes", "Routes":
		c.Plugins.Routes = false
	case "cache", "Cache":
		c.Plugins.Cache = false
	case "events", "Events":
		c.Plugins.Events = false
	case "workers", "Workers":
		c.Plugins.Workers = false
	default:
		return fmt.Errorf("unknown plugin: %s", name)
	}

	return nil
}

func (c *Config) Models() (map[string]*Model, error) {
	dir := c.Path(c.Definitions.Models)
	models := make(map[string]*Model)
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("walking models: %s", path))
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			model := &Model{}
			err := tasks.ReadYaml(path, model)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading model: %s", path))
			}

			models[model.Name] = model
		}

		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}
	return models, nil
}

func (c *Config) Events() (map[string]*Event, bool, error) {
	dir := c.Path(c.Definitions.Events)
	events := make(map[string]*Event)
	var hasReceiver bool
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("walking event: %s", path))
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			event := &Event{}
			err := tasks.ReadYaml(path, event)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading event: %s", path))
			}

			if !hasReceiver && event.Receiver {
				hasReceiver = true
			}

			events[event.Name] = event
		}

		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		return nil, false, err
	}
	return events, hasReceiver, nil
}

func (c *Config) Root() string {
	abs, err := filepath.Abs(c.File)
	if err != nil {
		return ""
	}
	return strings.Replace(abs, "/.golem/config.yaml", "", 1)
}

func (c *Config) Path(arg ...string) string {
	list := []string{c.Root()}
	list = append(list, arg...)
	return filepath.Join(list...)
}

func (c *Config) Validate() error {
	// TODO: add validations?
	return nil
}
