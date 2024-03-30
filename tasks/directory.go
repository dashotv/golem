package tasks

import (
	"os"

	"github.com/dashotv/fae"
)

func Directory(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fae.Wrap(err, "mkdir")
	}
	return nil
}
