package app

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/dashotv/golem/config"
)

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// TODO: move this
func MakeDirectory(dir string) error {
	if !exists(dir) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir")
		}
	}
	return nil
}

func executeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func writeDefaultConfig(dir string) error {
	cfg := config.DefaultConfig()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "mkdir")
		}
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "could not marshal config")
	}

	err = ioutil.WriteFile(dir+"/.golem.yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}

type defaultAppConfig struct {
	Mode string
	Port int
}

func writeAppConfig(name string) error {
	cfg := &defaultAppConfig{Mode: "dev", Port: 3000}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "could not marshal config")
	}

	err = ioutil.WriteFile(name+"/."+name+".yaml", b, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	return nil
}
