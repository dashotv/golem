package output

import (
	"embed"
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/pkg/errors"
)

//go:embed style.json
var files embed.FS

func Markdown(in string) error {
	return glamourText(in)
}

func MarkdownFile(path string) error {
	return glamourFile(path)
}

func MarkdownBytes(data []byte) error {
	return glamourBytes(data)
}

func glamourRenderer() (*glamour.TermRenderer, error) {
	data, err := files.ReadFile("style.json")
	if err != nil {
		return nil, errors.Wrap(err, "reading styles")
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(data),
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating renderer")
	}

	return r, nil
}

func glamourText(in string) error {
	r, err := glamourRenderer()
	if err != nil {
		return errors.Wrap(err, "getting renderer")
	}

	out, err := r.Render(in)
	if err != nil {
		return errors.Wrap(err, "rendering markdown")
	}

	fmt.Print(out)
	return nil
}

func glamourBytes(data []byte) error {
	r, err := glamourRenderer()
	if err != nil {
		return errors.Wrap(err, "getting renderer")
	}

	out, err := r.RenderBytes(data)
	if err != nil {
		return errors.Wrap(err, "rendering markdown")
	}

	fmt.Print(string(out))
	return nil
}

func glamourFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "reading file")
	}

	return glamourBytes(data)
}
