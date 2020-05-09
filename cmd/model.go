/*
Copyright © 2020 Shawn Catanzarite <me@shawncatz.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/generators/app"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model <name> [field...]",
	Short: "generate a new model definition",
	Long:  "generate a new model definition",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fields := args[1:]

		g := app.NewModelDefinitionGenerator(cfg, name, fields...)
		if err := g.Execute(); err != nil {
			logrus.Fatalf("error generating new model definition: %s", err)
		}
	},
}

func init() {
	newCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
