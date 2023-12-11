package config

import "github.com/iancoleman/strcase"

var restRoutes = []*Route{
	{
		Name:   "index",
		Path:   "/",
		Method: "GET",
		Params: []*Param{
			{
				Name:  "page",
				Type:  "int",
				Query: true,
			},
			{
				Name:  "limit",
				Type:  "int",
				Query: true,
			},
		},
	},
	{
		Name:   "create",
		Path:   "/",
		Method: "POST",
	},
	{
		Name:   "show",
		Path:   "/:id",
		Method: "GET",
		Params: []*Param{
			{
				Name: "id",
				Type: "string",
			},
		},
	},
	{
		Name:   "update",
		Path:   "/:id",
		Method: "PUT",
		Params: []*Param{
			{
				Name: "id",
				Type: "string",
			},
		},
	},
	{
		Name:   "settings",
		Path:   "/:id",
		Method: "PATCH",
		Params: []*Param{
			{
				Name: "id",
				Type: "string",
			},
		},
	},
	{
		Name:   "delete",
		Path:   "/:id",
		Method: "DELETE",
		Params: []*Param{
			{
				Name: "id",
				Type: "string",
			},
		},
	},
}

// RouteFile corresponds to a Group of routes
type Group struct {
	Name   string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string   `json:"path,omitempty" yaml:"path,omitempty"`
	Rest   bool     `json:"rest,omitempty" yaml:"rest,omitempty"`
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`
}

func (g *Group) CombinedRoutes() []*Route {
	list := []*Route{}
	if g.Rest {
		list = append(list, restRoutes...)
	}
	return append(list, g.Routes...)
}

func (g *Group) Camel() string {
	return strcase.ToCamel(g.Name)
}

func (g *Group) AddRoute(r *Route) {
	g.Routes = append(g.Routes, r)
}

type Route struct {
	Name   string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string   `json:"path,omitempty" yaml:"path,omitempty"`
	Method string   `json:"method,omitempty" yaml:"method,omitempty"`
	Params []*Param `json:"params,omitempty" yaml:"params,omitempty"`
}

func (r *Route) Camel() string {
	return strcase.ToCamel(r.Name)
}

func (r *Route) AddParam(p *Param) {
	r.Params = append(r.Params, p)
}

func (r *Route) Crud() bool {
	if len(r.Params) == 1 && r.Params[0].Name == "id" {
		return true
	}
	return false
}

type Param struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
	Query bool   `json:"query,omitempty" yaml:"query,omitempty"`
}

func (p *Param) TypeCamel() string {
	return strcase.ToCamel(p.Type)
}
