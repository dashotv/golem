package config

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Name = "test"
	cfg.Repo = "github.com/dashotv/test"

	cfg.Models.Enabled = true
	cfg.Models.Package = "models"
	cfg.Models.Output = "./models"
	cfg.Models.Definitions = "./.golem/models"

	cfg.Routes.Enabled = false
	cfg.Routes.Output = "./server"
	cfg.Routes.Definition = "./.golem/routes.yaml"
	cfg.Routes.Name = "Blarg"
	cfg.Routes.Repo = "github.com/dashotv/blarg"

	cfg.Jobs.Enabled = false
	cfg.Jobs.Package = "jobs"
	cfg.Jobs.Output = "./jobs"
	cfg.Jobs.Definitions = "./.golem/jobs.yaml"

	cfg.Connections = make(map[string]*Connection)
	cfg.Connections["default"] = &Connection{
		URI:      "mongodb://localhost:27017",
		Database: "database",
	}

	return cfg
}
