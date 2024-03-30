package cmd

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var clientCmd = &cobra.Command{
	Use:   "client LANGUAGE",
	Short: "Generate a client definition",
	Long:  "Generate a client definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		language := args[0]

		if !lo.Contains(config.SupportedClients(), language) {
			output.FatalTrace("error: %s", errors.Errorf("unsupported client language %s", language))
		}

		client := &config.Client{
			Language: language,
		}

		if err := generators.NewClient(cfg, client); err != nil {
			output.FatalTrace("error: %s", err)
		}
		if err := markdown("client"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(clientCmd)
}
