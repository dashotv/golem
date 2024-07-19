package config

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/tasks"
)

func (c *Config) Events() (map[string]*Event, bool, error) {
	dir := c.Path(c.Definitions.Events)
	events := make(map[string]*Event)
	var hasReceiver bool
	err := c.walk(dir, func(path string) error {
		event := &Event{}
		err := tasks.ReadYaml(path, event)
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("reading event: %s", path))
		}

		if !hasReceiver && event.Receiver {
			hasReceiver = true
		}

		events[event.Name] = event
		return nil
	})
	return events, hasReceiver, err
}

type Event struct {
	Name            string   `yaml:"name,omitempty"`
	Channel         string   `yaml:"channel,omitempty"`
	Receiver        bool     `yaml:"receiver"`
	Concurrency     int      `yaml:"concurrency,omitempty"`
	ExistingPayload string   `yaml:"existing_payload,omitempty"`
	ProxyTo         string   `yaml:"proxy_to,omitempty"`
	ProxyType       string   `yaml:"proxy_type,omitempty"`
	Fields          []*Field `yaml:"fields,omitempty"` // create type
}

func (e *Event) Sender() bool {
	return !e.Receiver
}

func (e *Event) Camel() string {
	return strcase.ToCamel(e.Name)
}

func (e *Event) Create() bool {
	return e.ExistingPayload == ""
}

func (e *Event) Payload() string {
	if e.ExistingPayload != "" {
		return e.ExistingPayload
	}

	return "Event" + e.Camel()
}
