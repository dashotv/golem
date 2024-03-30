package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var rest bool
var groupPath string
var groupModel string
var groupDesc = `... group NAME

  NAME	group name, must be unique, path defaults to "/NAME"
`

// groupCmd represents the route command
var groupCmd = &cobra.Command{
	Use:   "group NAME",
	Short: "generate a new route definition",
	Long:  groupDesc,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := "/" + name
		if groupPath != "" {
			path = groupPath
		}

		if path[0] != '/' {
			path = "/" + path
		}

		if rest && groupModel == "" {
			output.Fatalf("error: model is required for REST-style routes")
		}

		g := &config.Group{
			Name:  name,
			Path:  path,
			Rest:  rest,
			Model: groupModel,
		}

		if err := generators.NewGroup(cfg, g); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("group"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(groupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	groupCmd.Flags().BoolVarP(&rest, "rest", "r", false, "generate REST-style routes (Create/Retrieve/Update/Destroy)")
	groupCmd.Flags().StringVarP(&groupPath, "path", "p", "", "custom path for group")
	groupCmd.Flags().StringVarP(&groupModel, "model", "m", "", "model for REST-style routes")
}
