package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/tasks"
)

func (c *Config) Groups() (map[string]*Group, error) {
	dir := c.Path(c.Definitions.Routes)
	groups := make(map[string]*Group)
	err := c.walk(dir, func(path string) error {
		group := &Group{}
		err := tasks.ReadYaml(path, group)
		if err != nil {
			return fae.Wrap(err, fmt.Sprintf("reading group: %s", path))
		}

		groups[group.Name] = group
		return nil
	})
	return groups, err
}
func restRoutes(model string) []*Route {
	modelList := ""
	if model != "" {
		modelList = "[]" + model
	}
	return []*Route{
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
			Result: modelList,
		},
		{
			Name:   "create",
			Path:   "/",
			Method: "POST",
			Params: []*Param{
				{
					Name: "subject",
					Type: model,
					Bind: true,
				},
			},
			Result: model,
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
			Result: model,
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
				{
					Name: "subject",
					Type: model,
					Bind: true,
				},
			},
			Result: model,
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
				{
					Name: "setting",
					Type: "*Setting",
					Bind: true,
				},
			},
			Result: model,
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
			Result: model,
		},
	}
}

// RouteFile corresponds to a Group of routes
type Group struct {
	Name   string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string   `json:"path,omitempty" yaml:"path,omitempty"`
	Rest   bool     `json:"rest,omitempty" yaml:"rest,omitempty"`
	Model  string   `json:"model,omitempty" yaml:"model,omitempty"`
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`
}

func (g *Group) CombinedRoutes() []*Route {
	list := []*Route{}
	if g.Rest {
		list = append(list, restRoutes(g.Model)...)
	}
	return append(list, g.Routes...)
}

func (g *Group) Camel() string {
	return strcase.ToCamel(g.Name)
}

func (g *Group) AddRoute(r *Route) {
	g.Routes = append(g.Routes, r)
}

func (g *Group) TypescriptType() string {
	if g.Model == "" {
		return ""
	}

	t := TypescriptType(g.Model)
	if t == "any" {
		return ""
	}

	return t
}
func (g *Group) GetModels() string {
	list := []string{g.TypescriptType()}
	for _, r := range g.Routes {
		t := r.TypescriptResult()
		if t != "" {
			if t[len(t)-2:] == "[]" {
				t = t[:len(t)-2]
			}
			list = append(list, t)
		}
	}
	return strings.Join(list, ", ")
}

type Route struct {
	Name   string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string   `json:"path,omitempty" yaml:"path,omitempty"`
	Method string   `json:"method,omitempty" yaml:"method,omitempty"`
	Params []*Param `json:"params,omitempty" yaml:"params,omitempty"`
	Result string   `json:"result,omitempty" yaml:"result,omitempty"`
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

func (r *Route) Index() bool {
	return r.Name == "index"
}
func (r *Route) Create() bool {
	return r.Name == "create"
}
func (r *Route) Show() bool {
	return r.Name == "show"
}
func (r *Route) Update() bool {
	return r.Name == "update"
}
func (r *Route) Settings() bool {
	return r.Name == "settings"
}
func (r *Route) Delete() bool {
	return r.Name == "delete"
}
func (r *Route) HasModel() bool {
	for _, p := range r.Params {
		if p.Bind {
			return true
		}
	}
	return false
}
func (r *Route) GetModel() *Param {
	for _, p := range r.Params {
		if p.Bind {
			return p
		}
	}
	return nil
}
func (r *Route) HasParams() bool {
	return len(r.Params) > 0
}
func (r *Route) QueryParams() []*Param {
	list := []*Param{}
	for _, p := range r.Params {
		if p.Query {
			list = append(list, p)
		}
	}
	return list
}
func (r *Route) PathParams() []*Param {
	list := []*Param{}
	for _, p := range r.Params {
		if !p.Query && !p.Bind {
			list = append(list, p)
		}
	}
	return list
}
func (r *Route) ClientMethod() string {
	return strcase.ToCamel(r.Method)
}
func (r *Route) ClientPath() string {
	return convertPathParams(r.Path)
}
func (r *Route) TypescriptResult() string {
	return TypescriptType(r.Result)
}
func (r *Route) TypescriptMethod() string {
	return strings.ToLower(r.Method)
}
func (r *Route) TypescriptPath() string {
	return convertPathParamsTypescript(r.Path)
}
func (r *Route) QueryString() string {
	list := []string{}
	for _, p := range r.QueryParams() {
		list = append(list, fmt.Sprintf("%s=${params.%s}", p.Name, p.Name))
	}
	return strings.Join(list, "&")
}

var pathParam = regexp.MustCompile(`:(\w+)`)

func convertPathParams(path string) string {
	return pathParam.ReplaceAllString(path, "{$1}")
}

func convertPathParamsTypescript(path string) string {
	return pathParam.ReplaceAllString(path, "${params.$1}")
}

type Param struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
	Query bool   `json:"query,omitempty" yaml:"query,omitempty"`
	Bind  bool   `json:"bind,omitempty" yaml:"bind,omitempty"`
}

func (p *Param) Camel() string {
	return strcase.ToCamel(p.Name)
}

func (p *Param) TypeCamel() string {
	return strcase.ToCamel(p.Type)
}

func (p *Param) TypescriptType() string {
	return TypescriptType(p.Type)
}

func TypescriptType(t string) string {
	if t == "" {
		return "any"
	}

	if t == "time.Time" || t == "primitive.ObjectID" {
		return "string"
	}

	if t[0] == '*' {
		t = t[1:]
	}

	if strings.HasPrefix(t, "[]") {
		return TypescriptType(t[2:]) + "[]"
	}

	switch t {
	case "string":
		return "string"
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	}

	return t
}
