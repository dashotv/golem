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
	Language string
}

func SupportedClients() []string {
	return []string{"go"}
}

func (c *Client) Filename() string {
	switch c.Language {
	case "go":
		return filepath.Join("client", "client.gen.go")
		// case "js":
		// 	return filepath.Join("javascript", "src", "client.js")
		// case "ts":
		// 	return filepath.Join("typescript", "src", "client.ts")
	}
	return ""
}
