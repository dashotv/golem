package tasks

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

type TaskRunner struct {
	Name   string
	Tasks  []*Task
	Groups []*TaskRunner
}

func NewTaskRunner(name string) *TaskRunner {
	return &TaskRunner{Name: name, Tasks: make([]*Task, 0), Groups: make([]*TaskRunner, 0)}
}

type TaskFunction func() error

func (r *TaskRunner) Add(name string, f TaskFunction) {
	r.Tasks = append(r.Tasks, newTask(name, f))
}

func (r *TaskRunner) Group(name string) *TaskRunner {
	n := NewTaskRunner(r.Name + ":" + name)
	r.Groups = append(r.Groups, n)
	return n
}

func (r *TaskRunner) Run() error {
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

func newTask(name string, f TaskFunction) *Task {
	return &Task{Name: name, Function: f}
}

type Task struct {
	Name     string
	Function TaskFunction
}

func (t *Task) Run() error {
	return t.Function()
}
