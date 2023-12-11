package config

import "github.com/iancoleman/strcase"

type Event struct {
	Name            string   `yaml:"name,omitempty"`
	Channel         string   `yaml:"channel,omitempty"`
	Receiver        bool     `yaml:"receiver"`
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
