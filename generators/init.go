package generators

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func Init(name, repo string) error {
	var g map[string]string
	cfg := config.DefaultConfig()
	cfg.Name = name
	cfg.Repo = repo
	cfgFile := filepath.Join(name, ".golem", "config.yaml")

	runner := tasks.NewRunner("init")

	core := runner.Group("core")
	core.Add("directory", func() error {
		return tasks.Directory(name)
	})
	core.Add("directory", func() error {
		return tasks.Directory(filepath.Join(name, ".golem"))
	})
	core.Add("golem:create", func() error {
		return tasks.WriteYaml(cfgFile, cfg)
	})
	core.Add("golem:load", func() error {
		err := tasks.ReadYaml(cfgFile, cfg)
		if err != nil {
			return err
		}

		cfg.File = cfgFile
		g = cfg.Data()
		g["Port"] = "3000"

		err = os.Chdir(cfg.Root())
		if err != nil {
			return errors.Wrap(err, "changing directory")
		}

		return nil
	})

	cmd := runner.Group("cmd")
	cmd.Add("cmd", func() error {
		return tasks.Directory("cmd")
	})
	cmd.Add("cmd/root", func() error {
		return tasks.File(filepath.Join("cmd", "root"), filepath.Join("cmd", "root.go"), g)
	})
	cmd.Add("cmd/server", func() error {
		return tasks.File(filepath.Join("cmd", "server"), filepath.Join("cmd", "server.go"), g)
	})

	main := runner.Group("main")
	main.Add("main", func() error {
		return tasks.File("main", "main.go", g)
	})
	main.Add("license", func() error {
		return tasks.File("license", "LICENSE", g)
	})
	main.Add("app", func() error {
		return tasks.Directory(cfg.Output)
	})
	main.Add("makefile", func() error {
		return tasks.File("makefile", "Makefile", g)
	})
	main.Add("dockerfile", func() error {
		return tasks.File("dockerfile", "Dockerfile", g)
	})
	main.Add("docker-compose", func() error {
		return tasks.File("docker-compose", "docker-compose.yml", g)
	})
	main.Add("air", func() error {
		return tasks.File(".air.toml", ".air.toml", g)
	})
	main.Add("drone", func() error {
		return tasks.File(".drone.yml", ".drone.yml", g)
	})
	main.Add("etc", func() error {
		return tasks.Directory("etc")
	})
	main.Add("etc/service", func() error {
		return tasks.File(filepath.Join("etc", "service"), filepath.Join("etc", name+".service"), g)
	})
	main.Add(".env", func() error {
		return tasks.File(".env", ".env", g)
	})
	main.Add(".env.production", func() error {
		return tasks.File(".env.production", ".env.production", g)
	})
	main.Add(".gitignore", func() error {
		return tasks.File(".gitignore", ".gitignore", g)
	})
	main.Add("README", func() error {
		return tasks.File("README", "README.md", g)
	})

	app := runner.Group("app")
	app.Add("app", func() error {
		return tasks.File(filepath.Join("app", "app"), filepath.Join(cfg.Output, "app.go"), g)
	})
	app.Add("config", func() error {
		return tasks.File(filepath.Join("app", "app_config"), filepath.Join(cfg.Output, "config.go"), g)
	})
	app.Add("logger", func() error {
		return tasks.File(filepath.Join("app", "app_logger"), filepath.Join(cfg.Output, "logger.go"), g)
	})
	runner.Add("utils", func() error {
		return tasks.File(filepath.Join("app", "app_utils"), filepath.Join(cfg.Output, "utils.go"), g)
	})

	commands := runner.Group("commands")
	commands.Add("go mod init", func() error {
		return tasks.GoModInit(cfg.Repo)
	})
	commands.Add("git init", func() error {
		return tasks.GitInit()
	})

	err := runner.Run()
	if err != nil {
		return err
	}

	cfg.File = ".golem/config.yaml"
	return App(cfg)
}
