package config

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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
	Connections map[string]*Connection
}

type Connection struct {
	URI        string
	Database   string `json:"database"`
	Collection string `json:"collection"`
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
	if err := c.validateDefaultConnection(); err != nil {
		return err
	}
	// TODO: add more validations?
	return nil
}

func (c *Config) validateDefaultConnection() error {
	if len(c.Connections) == 0 {
		return errors.New("you must specify a default connection")
	}

	var def *Connection
	for n, c := range c.Connections {
		if n == "default" || n == "Default" {
			def = c
			break
		}
	}

	if def == nil {
		return errors.New("no 'default' found in connections list")
	}
	if def.Database == "" {
		return errors.New("default connection must specify database")
	}
	if def.URI == "" {
		return errors.New("default connection must specify URI")
	}

	return nil
}
