package generators

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/templates"
)

// Generator is the base generate for all other generators
type Generator struct {
	Filename string
	Buffer   *bytes.Buffer
}

// Write the output of the generator to a file
func (g *Generator) Write() error {
	logrus.Debugf("Model Output:\n\n" + g.Buffer.String())
	return ioutil.WriteFile(g.Filename, g.Buffer.Bytes(), 0644)
}

// Format the file using goimports
func (g *Generator) Format() error {
	cmd := exec.Command("goimports", "-w", g.Filename)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// ReadYaml reads a yaml file into a structure
func ReadYaml(path string, object interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, object)
	if err != nil {
		return err
	}
	return nil
}

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
