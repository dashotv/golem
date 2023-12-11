/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable NAME",
	Short: "enable a plugin for golem application",
	Long:  "enable a plugin for golem application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := generators.EnablePlugin(cfg, name); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("enable"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	pluginCmd.AddCommand(enableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
