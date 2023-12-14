package generators

import (
	"fmt"
	"os"
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

	return runner.Run()
}

func Events(cfg *config.Config) error {
	dir := filepath.Join(cfg.Root(), cfg.Definitions.Events)
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

	events, hasReceiver, err := WalkEvents(dir)
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
		runner.Add("hook:"+k, func() error {
			d := struct {
				Package string
				Event   *config.Event
			}{
				Package: cfg.Package,
				Event:   v,
			}
			return tasks.FileDoesntExist(filepath.Join("app", "events"), filepath.Join("app", "events_"+k+".go"), d)
		})
	}

	runner.Add("save", func() error {
		return tasks.RawFile(filepath.Join("app", "app_events.go"), strings.Join(output, "\n"))
	})

	return runner.Run()
}

func WalkEvents(dir string) (map[string]*config.Event, bool, error) {
	events := make(map[string]*config.Event)
	var hasReceiver bool
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("walking event: %s", path))
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			event := &config.Event{}
			err := tasks.ReadYaml(path, event)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading event: %s", path))
			}

			if !hasReceiver && event.Receiver {
				hasReceiver = true
			}

			events[event.Name] = event
		}

		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		return nil, false, err
	}
	return events, hasReceiver, nil
}
