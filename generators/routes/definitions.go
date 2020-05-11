package routes

import "strings"

// Definition stores the configuration from the routes file
type Definition struct {
	Name   string
	Repo   string
	Groups map[string]*Group `json:"groups" yaml:"groups"`
	Routes map[string]*Route `json:"routes" yaml:"routes"`
}

// Group stores the configuration of the group
type Group struct {
	Repo   string
	Name   string            `json:"name" yaml:"name"`
	Path   string            `json:"path" yaml:"path"`
	Method string            `json:"method" yaml:"method"`
	Routes map[string]*Route `json:"routes" yaml:"routes"`
}

// Route stores the configuration of the route
type Route struct {
	Repo   string
	Name   string   `json:"name" yaml:"name"`
	Path   string   `json:"path" yaml:"path"`
	Method string   `json:"method" yaml:"method"`
	Params []*Param `json:"params" yaml:"params"`
}

// GetMethod returns the capitalized method of the route or the default get
func (d *Route) GetMethod() string {
	if d.Method != "" {
		return strings.Title(d.Method)
	}
	return "Get"
}

// Param stores the configuration of the parameter
type Param struct {
	Name    string `json:"name" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	Default string `json:"default" yaml:"default"`
}

// GetType returns the capitalized type of the parameter or the default string
func (d *Param) GetType() string {
	if d.Type != "" {
		return strings.Title(d.Type)
	}
	return "String"
}
