package generators

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewWorker(cfg *config.Config, worker *config.Worker) error {
	runner := tasks.NewRunner("worker")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Workers))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Workers, worker.Name+".yaml"), worker)
	})
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "workers")
	})
	runner.Add("hook", func() error {
		d := struct {
			Package string
			Worker  *config.Worker
		}{
			Package: cfg.Package,
			Worker:  worker,
		}
		return tasks.FileDoesntExist(filepath.Join("app", "workers"), filepath.Join("app", "workers_"+worker.Name+".go"), d)
	})

	return runner.Run()
}

func NewQueue(cfg *config.Config, queue *config.Queue) error {
	runner := tasks.NewRunner("queue")
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "workers")
	})
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Queues))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Queues, queue.Name+".yaml"), queue)
	})

	return runner.Run()
}

func Workers(cfg *config.Config) error {
	// collect workers
	workers, err := cfg.Workers()
	if err != nil {
		return errors.Wrap(err, "collecting models")
	}
	// collect queues
	queues, err := cfg.Queues()
	if err != nil {
		return errors.Wrap(err, "collecting models")
	}

	data := struct {
		Config  *config.Config
		Workers map[string]*config.Worker
		Queues  map[string]*config.Queue
	}{
		Config:  cfg,
		Workers: workers,
		Queues:  queues,
	}

	runner := tasks.NewRunner("workers")
	runner.Add("file", func() error {
		return tasks.File(filepath.Join("app", "app_workers"), filepath.Join("app", "workers.gen.go"), data)
	})

	return runner.Run()
}
