package database

// Model holds the data from the YAML model
type Model struct {
	Package string
	Camel   string
	Name    string
	Type    string
	Imports []string
	Fields  []*Field
}

// Field holds the data from the YAML field
type Field struct {
	Name   string
	Camel  string
	Type   string
	Json   string
	Bson   string
	Tags   string
	Fields []*Field
}
