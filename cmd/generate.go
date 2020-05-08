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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/dashotv/golem/generators"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate mongo database models based on MDM library",
	Long:  "generate mongo database models based on MDM library",
	Run: func(cmd *cobra.Command, args []string) {
		source := cfg.Models.Definitions
		if !exists(source) {
			logrus.Fatalf("definitions directory doesn't exist: %s", source)
		}

		dest := cfg.Models.Output
		if !exists(dest) {
			logrus.Fatalf("output directory doesn't exist: %s", dest)
		}

		g := &generators.Generator{Config: cfg}
		err := g.Process()
		if err != nil {
			logrus.Fatalf("error occured during process: %s", err)
		}
	},
}

func exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
