package tasks

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/dashotv/fae"
)

// ReadYaml reads a yaml file into a structure
func ReadYaml(path string, object interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, object)
	if err != nil {
		return err
	}
	return nil
}

// WriteYaml writes a yaml file from a structure
func WriteYaml(path string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return fae.Wrap(err, "marshal yaml")
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return fae.Wrap(err, "writing yaml")
	}

	return nil
}
