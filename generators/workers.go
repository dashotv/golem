package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewWorker(cfg *config.Config, worker *config.Worker) error {
	runner := tasks.NewRunner("worker")
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "workers")
	})
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Workers))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Workers, worker.Name+".yaml"), worker)
	})

	return runner.Run()
}

func Workers(cfg *config.Config) error {
	dir := filepath.Join(cfg.Root(), cfg.Definitions.Workers)

	workers := make(map[string]*config.Worker)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			worker := &config.Worker{}
			err := tasks.ReadYaml(path, worker)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading worker: %s", path))
			}

			workers[worker.Name] = worker
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "walking routes")
	}

	data := struct {
		Config  *config.Config
		Workers map[string]*config.Worker
	}{
		Config:  cfg,
		Workers: workers,
	}

	runner := tasks.NewRunner("workers")
	runner.Add("file", func() error {
		return tasks.File(filepath.Join("app", "app_workers"), filepath.Join("app", "app_workers.go"), data)
	})

	keys := make([]string, 0, len(workers))
	for k := range workers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		k := k
		v := workers[k]
		runner.Add("hook:"+k, func() error {
			d := struct {
				Package string
				Worker  *config.Worker
			}{
				Package: cfg.Package,
				Worker:  v,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "workers"), filepath.Join("app", "workers_"+k+".go"), d)
		})
	}
	return runner.Run()
}
