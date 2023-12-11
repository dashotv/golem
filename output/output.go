package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var (
	info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90caf9")).
		Bold(true)
	text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cccccc"))
	success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#66bb6a"))
	warn = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffcc00"))
	err = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Padding(0, 1, 0, 0)
	outputText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))
	errorOutputText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))

	group = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90caf9")).
		Bold(true)
	title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cccccc")).
		Bold(true).
		Padding(0, 1)
	bracket = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#467499"))
)

func Printf(format string, a ...interface{}) {
	fmt.Println(text.Render(fmt.Sprintf(format, a...)))
}

func Success(s string) {
	fmt.Println(success.Render(s))
}

func Successf(format string, a ...interface{}) {
	fmt.Println(success.Render(fmt.Sprintf(format, a...)))
}

func Infof(format string, a ...interface{}) {
	fmt.Println(info.Render(fmt.Sprintf(format, a...)))
}

func Warnf(format string, a ...interface{}) {
	fmt.Println(warn.Render(fmt.Sprintf(format, a...)))
}

func Errorf(format string, a ...interface{}) {
	fmt.Println(err.Render(fmt.Sprintf(format, a...)))
}

func Fatalf(format string, a ...interface{}) {
	Errorf(format, a...)
	os.Exit(1)
}

func FatalTrace(fmt string, err error) {
	Errorf(fmt, err)

	e, ok := errors.Cause(err).(stackTracer)
	if ok {
		st := e.StackTrace()
		for _, l := range st {
			Errorf("   %s", l)
		}
	}

	os.Exit(1)
}

func PrintHeader(name, status string) {
	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		bracket.Render("["),
		group.Render(name),
		bracket.Render("]"),
		title.Render(status),
	)
	fmt.Println(bar)
}

func PrintOutputText(line string, err bool) {
	if err {
		fmt.Println(errorOutputText.Render(line))
		return
	}
	fmt.Println(outputText.Render(line))
}
