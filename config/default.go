package config

import "path/filepath"

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Version = "2.0"
	cfg.Name = "test"
	cfg.Repo = "github.com/dashotv/test"
	cfg.Package = "app"
	cfg.Output = filepath.Join("internal", "app")

	cfg.Plugins.Models = false
	cfg.Plugins.Routes = false
	cfg.Plugins.Workers = false
	cfg.Plugins.Cache = false
	cfg.Plugins.Events = false
	cfg.Plugins.Clients = false

	cfg.Definitions.Models = ".golem/models"
	cfg.Definitions.Routes = ".golem/routes"
	cfg.Definitions.Events = ".golem/events"
	cfg.Definitions.Workers = ".golem/workers"
	cfg.Definitions.Queues = ".golem/queues"
	cfg.Definitions.Clients = ".golem/clients"

	return cfg
}
