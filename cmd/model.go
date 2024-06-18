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
	"strings"

	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var modelStruct bool
var modelDesc = `Generate a new model definition

  NAME 		name, must be unique
  FIELD 	field name:type, can be repeated
`

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model NAME [FIELD...]",
	Short: "generate a new model definition",
	Long:  modelDesc,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fields := args[1:]

		m := &config.Model{
			Name:    name,
			Type:    "model",
			Indexes: []string{"created_at", "updated_at"},
		}
		if modelStruct {
			m.Type = "struct"
		}

		for _, f := range fields {
			s := strings.Split(f, ":")
			n := s[0]
			t := "string"
			if len(s) > 1 {
				t = s[1]
			}
			m.Fields = append(m.Fields, &config.Field{Name: n, Type: t, Json: "", Bson: ""})
		}

		if err := generators.NewModel(cfg, m); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("model"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	modelCmd.Flags().BoolVarP(&modelStruct, "struct", "s", false, "create model type struct")
}
