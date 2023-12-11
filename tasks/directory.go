package tasks

import (
	"os"

	"github.com/pkg/errors"
)

func Directory(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Wrap(err, "mkdir")
	}
	return nil
}
