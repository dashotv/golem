package config

type Config struct {
	Models struct {
		Enabled     bool
		Package     string
		Output      string
		Definitions string
	}
	Routes struct {
		Enabled    bool
		Name       string
		Definition string
		Output     string
		Repo       string
	}
	Jobs struct {
		Enabled     bool
		Package     string
		Definitions string
		Output      string
	}
}
