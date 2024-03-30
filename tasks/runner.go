package tasks

import (
	"os"

	"github.com/dashotv/golem/output"
)

var debug = os.Getenv("GOLEM_DEBUG") == "true"

type Runnable interface {
	Name() string
	Run() error
}

// Runner collects tasks to run
type Runner struct {
	name string
	list []Runnable
}

// NewRunner returns a new instance of Runner
func NewRunner(name string) *Runner {
	return &Runner{name, []Runnable{}}
}

func (r *Runner) Name() string {
	return r.name
}

// Add a new task to runner
func (r *Runner) Add(name string, f TaskFunction) {
	r.list = append(r.list, newTask(name, f))
}

// Add a new group to runner and return the group instance
func (r *Runner) Group(name string) *Runner {
	n := NewRunner(r.Name() + "." + name)
	r.Add(name, n.Run)
	return n
}

// Run all tasks and groups
func (r *Runner) Run() error {
	for _, t := range r.list {
		if debug {
			output.PrintHeader(r.Name(), t.Name())
		}
		if err := t.Run(); err != nil {
			return err
		}
	}

	return nil
}

func cwd() {
	cwd, err := os.Getwd()
	if err != nil {
		output.Errorf("Error getting current working directory: %s", err)
		return
	}
	output.Errorf("Current working directory: %s", cwd)
}
