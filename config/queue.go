package config

type Queue struct {
	Name        string `yaml:"name,omitempty"`
	Concurrency int    `yaml:"concurrency,omitempty"`
	Buffer      int    `yaml:"buffer,omitempty"`
	Interval    int    `yaml:"interval,omitempty"`
}
