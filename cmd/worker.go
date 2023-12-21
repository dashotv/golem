/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var workerSchedule string
var workerQueue string
var workerDesc = `... worker NAME [FIELD...]

  NAME		worker name, must be unique
  FIELD		field NAME[:TYPE], can be repeated
`

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "generate a new worker definition",
	Long:  workerDesc,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		n := strcase.ToSnake(args[0])
		w := &config.Worker{Name: n, Schedule: workerSchedule}

		if workerQueue != "" {
			queues, err := cfg.Queues()
			if err != nil {
				output.FatalTrace("error: %s", err)
			}
			if _, ok := queues[workerQueue]; !ok {
				output.Fatalf("error: queue %s does not exist", workerQueue)
			}
			w.Queue = workerQueue
		}

		if len(args) > 1 {
			for _, a := range args[1:] {
				s := strings.Split(a, ":")
				f := &config.Field{Name: s[0], Type: "string"}
				if len(s) > 1 {
					f.Type = s[1]
				}
				w.Fields = append(w.Fields, f)
			}
		}

		if err := generators.NewWorker(cfg, w); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("worker"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(workerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	workerCmd.Flags().StringVarP(&workerSchedule, "schedule", "s", "", "schedule for worker (cron format with seconds)")
	workerCmd.Flags().StringVarP(&workerQueue, "queue", "q", "", "queue for worker, must already exist")
}
