package tasks

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

type Runner struct {
	Name   string
	Tasks  []*Task
	Groups []*Runner
}

func NewRunner(name string) *Runner {
	return &Runner{Name: name, Tasks: make([]*Task, 0), Groups: make([]*Runner, 0)}
}

func (r *Runner) Add(name string, f TaskFunction) {
	r.Tasks = append(r.Tasks, newTask(name, f))
}

func (r *Runner) Group(name string) *Runner {
	n := NewRunner(r.Name + ":" + name)
	r.Groups = append(r.Groups, n)
	return n
}

func (r *Runner) Run() error {
	for _, t := range r.Tasks {
		fmt.Printf("* %s %s\n", aurora.Cyan(r.Name).Bold(), aurora.White(t.Name))
		if err := t.Function(); err != nil {
			return err
		}
	}

	for _, r := range r.Groups {
		//fmt.Printf("%s: ", aurora.Magenta(r.Name))
		if err := r.Run(); err != nil {
			return err
		}
	}
	return nil
}
