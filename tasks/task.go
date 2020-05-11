package tasks

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type TaskFunction func() error

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

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewMakeDirectoryTask(dir string) func() error {
	return func() error {
		if !Exists(dir) {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return errors.Wrap(err, "mkdir")
			}
		}
		return nil
	}
}

func NewExecuteCommandTask(name string, arg ...string) func() error {
	return func() error {
		cmd := exec.Command(name, arg...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
}
