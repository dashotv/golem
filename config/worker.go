package config

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/tasks"
)

func (c *Config) Workers() (map[string]*Worker, error) {
	dir := c.Path(c.Definitions.Workers)
	workers := make(map[string]*Worker)
	err := c.walk(dir, func(path string) error {
		worker := &Worker{}
		err := tasks.ReadYaml(path, worker)
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("reading worker: %s", path))
		}

		workers[worker.Name] = worker
		return nil
	})
	return workers, err
}

func (c *Config) Queues() (map[string]*Queue, error) {
	dir := c.Path(c.Definitions.Queues)
	queues := make(map[string]*Queue)
	err := c.walk(dir, func(path string) error {
		queue := &Queue{}
		err := tasks.ReadYaml(path, queue)
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("reading queue: %s", path))
		}

		queues[queue.Name] = queue
		return nil
	})
	return queues, err
}

type Worker struct {
	Name     string   `yaml:"name,omitempty"`
	Queue    string   `yaml:"queue,omitempty"`
	Schedule string   `yaml:"schedule,omitempty"`
	Fields   []*Field `yaml:"fields,omitempty"`
}

func (w *Worker) Camel() string {
	return strcase.ToCamel(w.Name)
}
