package config

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Name = "test"
	cfg.Repo = "github.com/dashotv/test"

	cfg.Models.Enabled = true
	cfg.Models.Package = "models"
	cfg.Models.Output = "models"
	cfg.Models.Definitions = ".golem/models"

	cfg.Routes.Enabled = true
	cfg.Routes.Repo = "github.com/dashotv/blarg"
	cfg.Routes.Output = "app"
	cfg.Routes.Definition = ".golem/routes.yaml"

	cfg.Jobs.Enabled = false
	cfg.Jobs.Package = "jobs"
	cfg.Jobs.Output = "jobs"
	cfg.Jobs.Definitions = ".golem/jobs.yaml"

	return cfg
}
