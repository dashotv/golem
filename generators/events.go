package generators

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
)

func NewEvent(cfg *config.Config, event *config.Event) error {
	runner := tasks.NewRunner("events")
	runner.Add("directory", func() error {
		return tasks.Directory(filepath.Join(cfg.Root(), cfg.Definitions.Events))
	})
	runner.Add("file", func() error {
		return tasks.WriteYaml(filepath.Join(cfg.Root(), cfg.Definitions.Events, event.Name+".yaml"), event)
	})
	runner.Add("plugin:enable", func() error {
		return pluginEnable(cfg, "events")
	})
	if event.Receiver {
		runner.Add("hook", func() error {
			d := struct {
				Package string
				Event   *config.Event
			}{
				Package: cfg.Package,
				Event:   event,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "events"), cfg.Join("events_"+event.Name+".go"), d)
		})
	}

	return runner.Run()
}

func Events(cfg *config.Config) error {
	data := cfg.Data()
	var output []string

	runner := tasks.NewRunner("app:events")
	runner.Add("header", func() error {
		header, err := tasks.Buffer(filepath.Join("app", "app_events"), data)
		if err != nil {
			return errors.Wrap(err, "events header")
		}
		output = append(output, header)
		return nil
	})

	events, hasReceiver, err := cfg.Events()
	if err != nil {
		return errors.Wrap(err, "walking events")
	}

	runner.Add("registration", func() error {
		combined := struct {
			Config      *config.Config
			Events      map[string]*config.Event
			HasReceiver bool
		}{
			Config:      cfg,
			Events:      events,
			HasReceiver: hasReceiver,
		}
		buf, err := tasks.Buffer(filepath.Join("app", "partial_events"), combined)
		if err != nil {
			return errors.Wrap(err, "events registration")
		}
		output = append(output, buf)
		return nil
	})

	keys := make([]string, 0, len(events))
	for k := range events {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		k := k
		v := events[k]
		runner.Add("event:"+k, func() error {
			buf, err := tasks.Buffer(filepath.Join("app", "partial_event"), v)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("events buffer: %s", k))
			}
			output = append(output, buf)
			return nil
		})
		if !v.Receiver {
			continue
		}
	}

	runner.Add("save", func() error {
		return tasks.RawFile(cfg.Join("events.gen.go"), strings.Join(output, "\n"))
	})

	return runner.Run()
}
