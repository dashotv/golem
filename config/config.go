package config

import (
	"path/filepath"
	"strings"
)

type Config struct {
	File   string
	Name   string
	Repo   string
	Models struct {
		Enabled     bool
		Package     string
		Output      string
		Definitions string
	}
	Routes struct {
		Enabled    bool
		Name       string
		Definition string
		Output     string
		Repo       string
	}
	Jobs struct {
		Enabled     bool
		Package     string
		Definitions string
		Output      string
	}
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
