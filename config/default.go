package config

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Version = "2.0"
	cfg.Name = "test"
	cfg.Repo = "github.com/dashotv/test"
	cfg.Package = "app"
	cfg.Output = "app"

	cfg.Plugins.Models = false
	cfg.Plugins.Routes = false
	cfg.Plugins.Workers = false
	cfg.Plugins.Cache = false
	cfg.Plugins.Events = false

	cfg.Definitions.Models = ".golem/models"
	cfg.Definitions.Routes = ".golem/routes"
	cfg.Definitions.Events = ".golem/events"
	cfg.Definitions.Workers = ".golem/workers"

	return cfg
}
