/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/output"
)

// readmeCmd represents the readme command
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "view README.md file",
	Long:  "view README.md file",
	Run: func(cmd *cobra.Command, args []string) {
		err := output.MarkdownFile("README.md")
		if err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(readmeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readmeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readmeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
