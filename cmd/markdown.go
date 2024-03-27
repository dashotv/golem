package cmd

import (
	"embed"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/output"
)

//go:embed *.md
var docs embed.FS

func markdown(path string) error {
	if silent {
		return nil
	}
	data, err := docs.ReadFile(path + ".md")
	if err != nil {
		return errors.Wrap(err, "reading file")
	}
	return output.Markdown(string(data))
}
