package templates

import "embed"

//go:embed *
var content embed.FS

func ReadDir(dir string) ([]string, error) {
	files, err := content.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}
