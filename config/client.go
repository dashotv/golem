package config

import "path/filepath"

type Client struct {
	Language string
}

func SupportedClients() []string {
	return []string{"go"}
}

func (c *Client) Filename() string {
	switch c.Language {
	case "go":
		return filepath.Join("client", "client.go")
		// case "js":
		// 	return filepath.Join("javascript", "src", "client.js")
		// case "ts":
		// 	return filepath.Join("typescript", "src", "client.ts")
	}
	return ""
}
