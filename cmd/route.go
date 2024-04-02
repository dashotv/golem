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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var pathRegex = `(:[a-zA-Z0-9_]+)`
var method string
var methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT"}
var routePath string
var routeResult string
var routeDesc = `... route GROUP NAME [PARAM...]

  GROUP		route group, group must exist
  NAME		route name, must be unique, path defaults to "/NAME"
			custom routes can contain params, e.g. "/:id" or "/something/:name"
  PARAM		parameters: NAME[:TYPE], can be repeated
    TYPE 	name of param, required, must be of base types (string, int, float, bool)
`

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route GROUP NAME [PARAM...]",
	Short: "generate a new route definition",
	Long:  routeDesc,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rx, err := regexp.Compile(pathRegex)
		if err != nil {
			output.FatalTrace("error: %s", err)
		}

		if !slices.Contains(methods, method) {
			output.Fatalf("error: invalid method %s: valid: %s\n", method, strings.Join(methods, ", "))
		}

		group := args[0]
		name := args[1]
		var params []*config.Param

		path := "/" + name
		if routePath != "" {
			path = routePath
		}
		if path[0] != '/' {
			path = "/" + path
		}
		if matches := rx.FindAllStringSubmatch(path, -1); len(matches) > 0 {
			for _, m := range matches {
				params = append(params, &config.Param{Name: m[1][1:], Type: "string"})
			}
		}

		if len(args) > 2 {
			for _, a := range args[2:] {
				p := &config.Param{Type: "string", Query: true}
				s := strings.Split(a, ":")
				p.Name = s[0]
				if len(s) > 1 {
					p.Type = s[1]
				}
				params = append(params, p)
			}
		}

		r := &config.Route{
			Name:   name,
			Path:   path,
			Method: method,
			Params: params,
			Result: routeResult,
		}

		if err := generators.NewRoute(cfg, group, r); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("route"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(routeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	routeCmd.Flags().StringVarP(&method, "method", "m", "GET", "set method for route (GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, TRACE, CONNECT)")
	routeCmd.Flags().StringVarP(&routePath, "path", "p", "", "custom path for route")
	routeCmd.Flags().StringVarP(&routeResult, "result", "r", "", "route result (return type for clients)")
}
