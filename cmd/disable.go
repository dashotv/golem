/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable NAME",
	Short: "disable a plugin for golem application",
	Long:  "disable a plugin for golem application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := generators.DisablePlugin(cfg, name); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("disable"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	pluginCmd.AddCommand(disableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
