package routes

type RouteGenerator struct {
}
type RouteGeneratorData struct {
	Name string
	Path string
}

type RouteDefinition struct {
	Name   string
	Path   string
	Routes []*RouteDefinition
	Params []*ParamDefinition
}

type ParamDefinition struct {
	Name    string
	Type    string
	Default string
}
