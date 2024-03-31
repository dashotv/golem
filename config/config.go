package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/fae"
)

type Config struct {
	File    string `yaml:"-"`
	Version string `yaml:"version,omitempty"`
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
		Clients bool `yaml:"clients"`
	} `yaml:"plugins"`
	Definitions struct {
		Models  string `yaml:"models,omitempty"`
		Routes  string `yaml:"routes,omitempty"`
		Events  string `yaml:"events,omitempty"`
		Workers string `yaml:"workers,omitempty"`
		Queues  string `yaml:"queues,omitempty"`
		Clients string `yaml:"clients,omitempty"`
	} `yaml:"definitions,omitempty"`
}

func (c *Config) Data() map[string]string {
	return map[string]string{
		"Name":     c.Name,
		"Camel":    strcase.ToCamel(c.Name),
		"Repo":     c.Repo,
		"Package":  c.Package,
		"Output":   c.Output,
		"Models":   fmt.Sprintf("%t", c.Plugins.Models),
		"Routes":   fmt.Sprintf("%t", c.Plugins.Routes),
		"Cache":    fmt.Sprintf("%t", c.Plugins.Cache),
		"Events":   fmt.Sprintf("%t", c.Plugins.Events),
		"Workers":  fmt.Sprintf("%t", c.Plugins.Workers),
		"Clients":  fmt.Sprintf("%t", c.Plugins.Clients),
		"NameHash": c.NameHash(),
	}
}

func (c *Config) Enable(name string) error {
	switch name {
	case "models", "Models":
		c.Plugins.Models = true
		if c.Definitions.Models == "" {
			c.Definitions.Models = "models"
		}
	case "routes", "Routes":
		c.Plugins.Routes = true
		if c.Definitions.Routes == "" {
			c.Definitions.Routes = "routes"
		}
	case "cache", "Cache":
		c.Plugins.Cache = true
	case "events", "Events":
		c.Plugins.Events = true
		if c.Definitions.Events == "" {
			c.Definitions.Events = "events"
		}
	case "workers", "Workers":
		c.Plugins.Workers = true
		if c.Definitions.Workers == "" {
			c.Definitions.Workers = "workers"
		}
	case "clients", "Clients":
		c.Plugins.Clients = true
		if c.Definitions.Clients == "" {
			c.Definitions.Clients = "clients"
		}
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
	case "clients", "Clients":
		return c.Plugins.Clients
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
	case "clients", "Clients":
		c.Plugins.Clients = false
	default:
		return fmt.Errorf("unknown plugin: %s", name)
	}

	return nil
}

func (c *Config) walk(dir string, fn func(yaml string) error) error {
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("walking %s: %s", dir, path))
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			return fn(path)
		}

		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		return err
	}
	return nil
}

func (c *Config) NameHash() string {
	sum := md5.Sum([]byte(c.Name))
	hash := hex.EncodeToString(sum[:])
	return c.Name + "-" + hash
}

func (c *Config) Root() string {
	abs, err := filepath.Abs(c.File)
	if err != nil {
		return ""
	}
	return strings.Replace(abs, "/.golem/config.yaml", "", 1)
}

func (c *Config) Join(names ...string) string {
	list := []string{c.Root(), c.Output}
	list = append(list, names...)
	return filepath.Join(list...)
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
