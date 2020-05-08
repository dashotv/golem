package base

import (
	"bytes"
	"regexp"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/generators/templates"

	"github.com/dashotv/golem/config"
)

type FileGenerator struct {
	*Generator
	Config *config.Config
	Name   string
	Data   map[string]string
}

func NewFileGenerator(cfg *config.Config, name string, path string, data map[string]string) *FileGenerator {
	return &FileGenerator{
		Config: cfg,
		Name:   name,
		Data:   data,
		Generator: &Generator{
			Filename: path,
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

func (f *FileGenerator) Execute() error {
	err := templates.New(f.Name).Execute(f.Buffer, f.Data)
	if err != nil {
		return errors.Wrap(err, "execute template")
	}

	err = f.Write()
	if err != nil {
		return errors.Wrap(err, "writing template output")
	}

	re := regexp.MustCompile(`.go$`)
	if re.MatchString(f.Filename) {
		err = f.Format()
		if err != nil {
			return errors.Wrap(err, "formatting template output")
		}
	}

	return nil
}
