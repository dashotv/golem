package tasks

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// TaskFunction defines the function for tasks
type TaskFunction func() error

// newTask is a simple wrapper for creating an instance of Task
func newTask(name string, f TaskFunction) *Task {
	return &Task{name: name, function: f}
}

// Task store name and function for task
type Task struct {
	name     string
	function TaskFunction
}

func (t *Task) Name() string {
	return t.name
}

// Run executes the task
func (t *Task) Run() error {
	return t.function()
}

// PathExists returns true or false if the path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// NewMakeDirectoryTask is a convenience method for creating a task that creates a directory if it doesn't exist
func NewMakeDirectoryTask(dir string) func() error {
	return func() error {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir")
		}
		return nil
	}
}

// NewExecuteCommandTask is a convenience method for creating a task that executes a command
func NewExecuteCommandTask(name string, arg ...string) func() error {
	return func() error {
		cmd := exec.Command(name, arg...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
}
