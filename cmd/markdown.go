package cmd

import (
	"embed"

	"github.com/dashotv/fae"
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
		return fae.Wrap(err, "reading file")
	}
	return output.Markdown(string(data))
}
