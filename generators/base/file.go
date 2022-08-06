package base

import (
	"bytes"
	"regexp"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/templates"
)

// FileGenerator manages generic or simple file generation
type FileGenerator struct {
	*Generator
	Config *config.Config
	Name   string
	Data   map[string]string
}

// NewFileGenerator creates and returns a FileGenerator
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

// Execute configures and generates a file
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
