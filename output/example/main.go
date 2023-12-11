package main

import (
	"embed"
	"fmt"

	"github.com/dashotv/golem/output"
)

//go:embed *.md
var docs embed.FS

func main() {
	output.Printf("This is a test")
	output.Successf("This is success")
	output.Infof("This is info")
	output.Warnf("This is a warning")
	output.Errorf("This is an error")

	fmt.Printf("\n")
	output.PrintHeader("main:header", "this is a header")

	fmt.Printf("\n")
	output.PrintOutputText("This is output", false)
	output.PrintOutputText("This is output", false)
	output.PrintOutputText("This is output", false)

	fmt.Printf("\n")
	output.PrintOutputText("This is error output", true)
	output.PrintOutputText("This is error output", true)
	output.PrintOutputText("This is error output", true)

	fmt.Printf("\nMarkdown:\n")
	data, _ := docs.ReadFile("test.md")
	output.MarkdownBytes(data)
}
