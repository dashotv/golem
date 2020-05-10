package base

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Generator struct {
	Filename string
	Buffer   *bytes.Buffer
}

func (g *Generator) Write() error {
	logrus.Debugf("Model Output:\n\n" + g.Buffer.String())
	return ioutil.WriteFile(g.Filename, g.Buffer.Bytes(), 0644)
}

func (g *Generator) Format() error {
	cmd := exec.Command("goimports", "-w", g.Filename)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

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
