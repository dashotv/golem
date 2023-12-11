package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/output"
	"github.com/dashotv/golem/tasks"
)

// routeCmd represents the route command
var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "list routes",
	Long:  "list routes",
	Run: func(cmd *cobra.Command, args []string) {
		if !cfg.Plugins.Routes {
			output.Warnf("error: routes plugin is not enabled")
			return
		}

		dir := filepath.Join(cfg.Root(), cfg.Definitions.Routes)
		if !tasks.PathExists(dir) {
			output.Fatalf("error: routes directory does not exist: %s", dir)
		}

		groups := make(map[string]*config.Group)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".yaml") {
				group := &config.Group{}
				err := tasks.ReadYaml(path, group)
				if err != nil {
					output.Printf("reading group: %s", path)
					os.Exit(1)
				}

				groups[group.Name] = group
			}

			return nil
		})
		if err != nil {
			output.FatalTrace("error: walking groups: %s", err)
		}

		for _, g := range groups {
			output.Infof("%s", g.Camel())
			for _, r := range g.CombinedRoutes() {
				output.Printf("  %10s %-25.25s %s (%s)", r.Method, g.Path+r.Path, g.Camel()+r.Camel(), r.Name+g.Camel()+"Handler")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(routesCmd)
}
