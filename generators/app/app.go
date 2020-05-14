package app

import (
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators/base"
	"github.com/dashotv/golem/tasks"
)

// Generator manages generating an application
type AppGenerator struct {
	*base.Generator
	Config *config.Config
	Name   string
	Repo   string
}

// NewGenerator returns a new instance of Generator
func NewAppGenerator(cfg *config.Config, name, repo string) *AppGenerator {
	return &AppGenerator{Config: cfg, Name: name, Repo: repo}
}

// Execute processes all of the configurations and generates an application
func (g *AppGenerator) Execute() error {
	runner := tasks.NewRunner("generator:app")

	err := runner.Run()
	if err != nil {
		return err
	}

	return nil
}
