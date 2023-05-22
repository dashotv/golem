package routes

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// Definition stores the configuration from the routes file
type Definition struct {
	Repo   string            `json:"repo" yaml:"repo"`
	Name   string            `json:"name" yaml:"name"`
	Groups map[string]*Group `json:"groups" yaml:"groups"`
	Routes map[string]*Route `json:"routes" yaml:"routes"`
}

// Group stores the configuration of the group
type Group struct {
	Repo   string            `json:"repo" yaml:"repo"`
	Name   string            `json:"name" yaml:"name"`
	Path   string            `json:"path" yaml:"path"`
	Method string            `json:"method" yaml:"method"`
	Routes map[string]*Route `json:"routes" yaml:"routes"`
	Rest   bool              `json:"rest" yaml:"rest"`
}

func (g *Group) Camel() string {
	return strcase.ToCamel(g.Name)
}

// Route stores the configuration of the route
type Route struct {
	Repo   string   `json:"repo" yaml:"repo"`
	Name   string   `json:"name" yaml:"name"`
	Path   string   `json:"path" yaml:"path"`
	Method string   `json:"method" yaml:"method"`
	Params []*Param `json:"params" yaml:"params"`
}

func (r *Route) Camel() string {
	return strcase.ToCamel(r.Name)
}

// GetMethod returns the capitalized method of the route or the default get
func (r *Route) GetMethod() string {
	if r.Method != "" {
		return strcase.ToScreamingSnake(r.Method)
	}
	return "GET"
}

// Param stores the configuration of the parameter
type Param struct {
	Name    string `json:"name" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	Default string `json:"default" yaml:"default"`
	Query   bool   `json:"query" yaml:"query"`
}

// GetType returns the capitalized type of the parameter or the default string
func (d *Param) GetType() string {
	if d.Type != "" {
		return strings.Title(d.Type)
	}
	return "String"
}
