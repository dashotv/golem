/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var queueConcurrency int
var queueBuffer int
var queueInterval int

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue NAME",
	Short: "generate a new queue definition",
	Long:  "generate a new queue definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		n := strcase.ToSnake(args[0])
		q := &config.Queue{Name: n, Concurrency: queueConcurrency, Buffer: queueBuffer, Interval: queueInterval}
		if err := generators.NewQueue(cfg, q); err != nil {
			output.FatalTrace("error: %s", err)
		}
		if err := markdown("queue"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(queueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	queueCmd.Flags().IntVarP(&queueConcurrency, "concurrency", "c", 0, "concurrency for queue")
	queueCmd.Flags().IntVarP(&queueBuffer, "buffer", "b", 0, "buffer size for queue")
	queueCmd.Flags().IntVarP(&queueInterval, "interval", "i", 0, "polling interval for queue")
}
