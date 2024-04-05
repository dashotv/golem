package config

import (
	"fmt"
	"path/filepath"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/tasks"
)

func (c *Config) Clients() (map[string]*Client, error) {
	dir := c.Path(c.Definitions.Clients)
	clients := make(map[string]*Client)
	err := c.walk(dir, func(path string) error {
		client := &Client{}
		err := tasks.ReadYaml(path, client)
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("reading client: %s", path))
		}

		clients[client.Language] = client
		return nil
	})
	return clients, err
}

type Client struct {
	Language    string `yaml:"language"`
	Destination string `yaml:"destination,omitempty"`
}

func SupportedClients() []string {
	return []string{"go", "typescript"}
}

func (c *Client) Output() string {
	if c.Destination != "" {
		return c.Destination
	}
	switch c.Language {
	case "go":
		return filepath.Join("client")
	// case "js":
	// 	return filepath.Join("javascript", "src", "client.js")
	case "typescript":
		return filepath.Join("ui", "src", "client")
	}
	return ""
}
