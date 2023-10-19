package config

import (
	"path/filepath"
	"strings"
)

type Config struct {
	File   string `yaml:"-"`
	Name   string `yaml:"name,omitempty"`
	Repo   string `yaml:"repo,omitempty"`
	Models struct {
		Enabled     bool   `yaml:"enabled,omitempty"`
		Package     string `yaml:"package,omitempty"`
		Output      string `yaml:"output,omitempty"`
		Definitions string `yaml:"definitions,omitempty"`
	} `yaml:"models,omitempty"`
	Routes struct {
		Enabled    bool   `yaml:"enabled,omitempty"`
		Name       string `yaml:"name,omitempty"`
		Definition string `yaml:"definition,omitempty"`
		Output     string `yaml:"output,omitempty"`
		Repo       string `yaml:"repo,omitempty"`
	} `yaml:"routes,omitempty"`
	Jobs struct {
		Enabled     bool   `yaml:"enabled,omitempty"`
		Package     string `yaml:"package,omitempty"`
		Definitions string `yaml:"definitions,omitempty"`
		Output      string `yaml:"output,omitempty"`
	} `yaml:"jobs,omitempty"`
}

func (c *Config) Root() string {
	abs, err := filepath.Abs(c.File)
	if err != nil {
		return ""
	}

	return strings.Replace(abs, "/.golem/.golem.yaml", "", 1)
}

func (c *Config) Path(arg ...string) string {
	list := []string{c.Root()}
	list = append(list, arg...)
	return strings.Join(list, "/")
}

func (c *Config) Validate() error {
	// TODO: add validations?
	return nil
}
