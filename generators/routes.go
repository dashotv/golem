package generators

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/tasks"
	"github.com/dashotv/golem/templates"
)

// RoutesGenerator manages the generation of all routes related files
type RoutesGenerator struct {
	Config     *config.Config
	Definition *Definition
	Groups     []*GroupGenerator
}

// Group stores the configuration of the group
type Group struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string            `json:"path,omitempty" yaml:"path,omitempty"`
	Method string            `json:"method,omitempty" yaml:"method,omitempty"`
	Routes map[string]*Route `json:"routes,omitempty" yaml:"routes,omitempty"`
	Rest   bool              `json:"rest,omitempty" yaml:"rest,omitempty"`
}

func (g *Group) Camel() string {
	return strcase.ToCamel(g.Name)
}

// Route stores the configuration of the route
type Route struct {
	Name   string   `json:"name,omitempty" yaml:"name,omitempty"`
	Path   string   `json:"path,omitempty" yaml:"path,omitempty"`
	Method string   `json:"method,omitempty" yaml:"method,omitempty"`
	Params []*Param `json:"params,omitempty" yaml:"params,omitempty"`
}

func (r *Route) Camel() string {
	return strcase.ToCamel(r.Name)
}

// Definition stores the configuration from the routes file
type Definition struct {
	Repo   string            `json:"repo,omitempty" yaml:"repo,omitempty"`
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`
	Groups map[string]*Group `json:"groups,omitempty" yaml:"groups,omitempty"`
	Routes map[string]*Route `json:"routes,omitempty" yaml:"routes,omitempty"`
}

func (d *Definition) FindGroup(name string, rest bool) *Group {
	if rest {
		group := &Group{
			Path: "/" + name,
			Rest: true,
		}
		d.Groups = map[string]*Group{name: group}
		return group
	}

	if d.Groups != nil {
		for _, g := range d.Groups {
			if name == g.Name {
				return g
			}
		}
	}

	if d.Groups == nil {
		d.Groups = make(map[string]*Group)
	}

	if d.Groups[name] == nil {
		d.Groups[name] = &Group{
			Name:   name,
			Path:   "/" + name,
			Routes: make(map[string]*Route),
		}
	}

	return d.Groups[name]
}

// FindRoute finds a route corresponding to the path name passed as argument
func (d *Definition) FindRoute(path string) *Route {
	parts := strings.Split(path, "/")
	group := d.FindGroup(parts[1], false)
	if group != nil && group.Rest {
		return nil
	}

	route := "/" + strings.Join(parts[2:], "/")

	if group != nil && group.Routes != nil {
		for _, r := range group.Routes {
			if r.Path == route {
				return r
			}
		}
	}

	if len(parts) == 2 {
		if group.Routes == nil {
			group.Routes = make(map[string]*Route)
		}
		if group.Routes["index"] == nil {
			group.Routes["index"] = &Route{Name: "index", Path: "/"}
		}
		return group.Routes["index"]
	}

	if group.Routes == nil {
		group.Routes = make(map[string]*Route)
	}
	if group.Routes[parts[2]] == nil {
		group.Routes[parts[2]] = &Route{Name: parts[2], Path: route}
	}
	return group.Routes[route]
}

// NewRoutesGenerator creates and returns an instance of Generator
func NewRoutesGenerator(cfg *config.Config) *RoutesGenerator {
	return &RoutesGenerator{
		Config:     cfg,
		Groups:     make([]*GroupGenerator, 0),
		Definition: &Definition{},
	}
}

// Execute generates the files from the templates
func (g *RoutesGenerator) Execute() error {
	if err := g.prepare(); err != nil {
		return err
	}
	runner := tasks.NewRunner("generator")
	r := runner.Group("routes")

	if g.Config.Routes.Enabled {
		r.Add("server routes", func() error {
			sg := NewServerRoutesGenerator(g.Config, g.Definition)
			return sg.Execute()
		})
	}

	return runner.Run()
}

// prepare configures the data for the templates
func (g *RoutesGenerator) prepare() error {
	if err := ReadYaml(g.Config.Routes.Definition, g.Definition); err != nil {
		return err
	}

	for name, group := range g.Definition.Groups {
		if group.Method == "" {
			group.Method = "GET"
		}
		if group.Path[0] != '/' {
			group.Path = "/" + group.Path
		}
		if group.Rest {
			group.Routes = map[string]*Route{
				"index": {
					Method: "GET",
					Name:   name + strings.Title("index"),
					Path:   "/",
				},
				"create": {
					Method: "POST",
					Name:   name + strings.Title("create"),
					Path:   "/",
				},
				"show": {
					Method: "GET",
					Name:   name + strings.Title("show"),
					Path:   "/:id",
					Params: []*Param{
						{Name: "id", Type: "string"},
					},
				},
				"update": {
					Method: "PUT",
					Name:   name + strings.Title("update"),
					Path:   "/:id",
					Params: []*Param{
						{Name: "id", Type: "string"},
					},
				},
				"delete": {
					Method: "DELETE",
					Name:   name + strings.Title("delete"),
					Path:   "/:id",
					Params: []*Param{
						{Name: "id", Type: "string"},
					},
				},
			}
		}
		for rn, route := range group.Routes {
			if route.Path == "" {
				route.Path = "/"
			}
			route.Name = name + strings.Title(rn)
		}
	}
	return nil
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
	Default string `json:"default,omitempty" yaml:"default,omitempty"`
	Query   bool   `json:"query,omitempty" yaml:"query,omitempty"`
}

// GetType returns the capitalized type of the parameter or the default string
func (d *Param) GetType() string {
	if d.Type != "" {
		return strings.Title(d.Type)
	}
	return "String"
}

// GroupGenerator manages the creation of group routes
type GroupGenerator struct {
	*Generator
	Config     *config.Config
	Output     string
	Name       string
	Definition *Group
}

// NewGroupGenerator configures and returns an instance of GroupGenerator
func NewGroupGenerator(cfg *config.Config, name string, d *Group) *GroupGenerator {
	return &GroupGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Name:       name,
		Definition: d,
		Generator: &Generator{
			Filename: cfg.Routes.Output + "/routs_" + name + ".go",
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute manages creation of group routes files with the template
func (g *GroupGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:" + g.Name)
	r.Add("prepare", g.prepare)
	return r.Run()
}

// prepare configures the data for the template
func (g *GroupGenerator) prepare() error {
	for n, r := range g.Definition.Routes {
		r.Name = n
	}
	return nil
}

// ServerGenerator manages the creation of the server file
type ServerGenerator struct {
	*Generator
	Config     *config.Config
	Output     string
	Definition *Definition
}

// NewServerGenerator creates and returns an instance of ServerGenerator
func NewServerGenerator(cfg *config.Config, d *Definition) *ServerGenerator {
	return &ServerGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		Generator: &Generator{
			Filename: cfg.Path(cfg.Routes.Output, "server.go"),
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute creates the server file from the template
func (g *ServerGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:server")

	r.Add("prepare", g.prepare)
	r.Add("template", func() error {
		return templates.New("app/server").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)
	r.Add("server routes", func() error {
		srg := NewServerRoutesGenerator(g.Config, g.Definition)
		return srg.Execute()
	})

	return r.Run()
}

// prepare the configuration for the template
func (g *ServerGenerator) prepare() error {
	return nil
}

// ServerRoutesGenerator manages the creation of the server routes file
type ServerRoutesGenerator struct {
	*Generator
	Config     *config.Config
	Output     string
	Definition *Definition
}

// NewServerRoutesGenerator creates and returns an instance of ServerRoutesGenerator
func NewServerRoutesGenerator(cfg *config.Config, d *Definition) *ServerRoutesGenerator {
	return &ServerRoutesGenerator{
		Config:     cfg,
		Output:     cfg.Routes.Output,
		Definition: d,
		Generator: &Generator{
			Filename: cfg.Path(cfg.Routes.Output, "routes.go"),
			Buffer:   bytes.NewBufferString(""),
		},
	}
}

// Execute creates the server routes file from the template
func (g *ServerRoutesGenerator) Execute() error {
	r := tasks.NewRunner("generator:routes:server routes")

	r.Add("prepare", g.prepare)
	r.Add("template", func() error {
		return templates.New("app/routes").Execute(g.Buffer, g.Definition)
	})
	r.Add("write", g.Write)
	r.Add("format", g.Format)

	return r.Run()
}

// prepare the configuration for the template
func (g *ServerRoutesGenerator) prepare() error {
	return nil
}
