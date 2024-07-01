package tasks

import (
	"os/exec"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/output"
)

func CommandWithCwd(name, cwd, cmd string, args ...string) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return fae.Wrap(err, "finding binary")
	}

	c := exec.Command(path, args...)
	if cwd != "" {
		c.Dir = cwd
	}

	out, err := c.CombinedOutput()
	if err != nil {
		output.Errorf("Error running command: %s", name)
		output.Errorf(string(out))
		return fae.Wrap(err, "running command")
	}

	return nil
}
func Command(name string, cmd string, args ...string) error {
	return CommandWithCwd(name, "", cmd, args...)
}

func GoFmt() error {
	return Command("go fmt ./...", "go", "fmt", "./...")
}
func GoGet(pkg string) error {
	return Command("go get", "go", "get", "-u", pkg)
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
func Prettier(dir string) error {
	return CommandWithCwd("prettier", dir, "npx", "prettier", "--write", ".")
}
