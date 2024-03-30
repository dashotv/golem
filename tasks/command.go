package tasks

import (
	"os/exec"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/output"
)

func Command(name, cmd string, args ...string) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return fae.Wrap(err, "finding binary")
	}

	c := exec.Command(path, args...)

	out, err := c.CombinedOutput()
	if err != nil {
		output.Errorf("Error running command: %s", name)
		output.Errorf(string(out))
		return fae.Wrap(err, "running command")
	}

	return nil
}

func GoFmt() error {
	return Command("go fmt ./...", "go", "fmt", "./...")
}
func GoImports() error {
	return Command("goimports", "go", "run", "--", "golang.org/x/tools/cmd/goimports@latest", "-w", ".")
}
func GoModInit(repo string) error {
	return Command("go mod init", "go", "mod", "init", repo)
}
func GoModTidy() error {
	return Command("go mod tidy", "go", "mod", "tidy")
}
func GitInit() error {
	return Command("git init", "git", "init", ".")
}
