package generators

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewModel(cfg *config.Config, m *config.Model) error {
	runner := tasks.NewRunner("model")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Models))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Models, m.Name+".yaml"), m)
	})
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "models")
	})
	if m.Type == "model" {
		runner.Add("hook", func() error {
			d := struct {
				Package string
				Model   *config.Model
			}{
				Package: cfg.Package,
				Model:   m,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "models"), cfg.Join("models_"+m.Name+".go"), d)
		})
	}

	return runner.Run()
}

func Models(cfg *config.Config) error {
	data := cfg.Data()
	var output []string

	runner := tasks.NewRunner("app:models")
	runner.Add("header", func() error {
		header, err := tasks.Buffer(filepath.Join("app", "app_models"), data)
		if err != nil {
			return fae.Wrap(err, "models header")
		}
		output = append(output, header)
		return nil
	})

	// collect models for connector registration
	models, err := cfg.Models()
	if err != nil {
		return fae.Wrap(err, "collecting models")
	}

	modelData := struct {
		Config *config.Config
		Models map[string]*config.Model
	}{
		Config: cfg,
		Models: models,
	}

	runner.Add("connector", func() error {
		buf, err := tasks.Buffer(filepath.Join("app", "partial_connector"), modelData)
		if err != nil {
			return fae.Wrap(err, "models connector buffer")
		}
		output = append(output, buf)
		return nil
	})

	keys := make([]string, 0, len(models))
	for k := range models {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		k := k
		v := models[k]
		runner.Add("model:"+k, func() error {
			t := "app/partial_model"
			if v.Type == "struct" {
				t = "app/partial_struct"
			}
			buf, err := tasks.Buffer(t, v)
			if err != nil {
				return fae.Wrap(err, fmt.Sprintf("model buffer: %s", k))
			}
			output = append(output, buf)
			return nil
		})
		if v.Type == "struct" {
			continue
		}
		runner.Add("hook:"+k, func() error {
			d := struct {
				Package string
				Model   *config.Model
			}{
				Package: cfg.Package,
				Model:   v,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "models"), cfg.Join("models_"+k+".go"), d)
		})
	}

	runner.Add("save", func() error {
		return tasks.RawFile(cfg.Join("models.gen.go"), strings.Join(output, "\n"))
	})

	return runner.Run()
}
