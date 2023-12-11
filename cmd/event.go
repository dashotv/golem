/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dashotv/golem/config"
	"github.com/dashotv/golem/generators"
	"github.com/dashotv/golem/output"
)

var eventChannel string
var eventProxy string
var eventPayload string
var eventReceiver bool
var eventDesc = `Generate a new event and channel definition

  NAME 		name, must be unique
  FIELD 	field name:type, can be repeated
`

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event NAME [FIELD...]",
	Short: "Add event/channel to your project",
	Long:  eventDesc,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fields := args[1:]

		e := &config.Event{
			Name:            name,
			Channel:         cfg.Name + "." + name,
			Receiver:        eventReceiver || eventProxy != "",
			ExistingPayload: eventPayload,
			ProxyTo:         eventProxy,
		}
		if eventChannel != "" {
			e.Channel = eventChannel
		}

		if len(fields) > 0 {
			for _, f := range fields {
				s := strings.Split(f, ":")
				n := s[0]
				t := "string"
				if len(s) > 1 {
					t = s[1]
				}
				e.Fields = append(e.Fields, &config.Field{Name: n, Type: t, Json: "", Bson: ""})
			}
		}

		if e.ProxyTo != "" {
			events, _, err := generators.WalkEvents(filepath.Join(cfg.Root(), cfg.Definitions.Events))
			if err != nil {
				output.FatalTrace("error: %s", err)
			}
			for _, v := range events {
				if v.Name == e.ProxyTo {
					e.ProxyType = v.Payload()
					break
				}
			}
			if e.ProxyType == "" {
				output.FatalTrace("error: %s", err)
			}
		}

		if err := generators.NewEvent(cfg, e); err != nil {
			output.FatalTrace("error: %s", err)
		}

		if err := markdown("event"); err != nil {
			output.FatalTrace("error: %s", err)
		}
	},
}

func init() {
	addCmd.AddCommand(eventCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	eventCmd.Flags().BoolVarP(&eventReceiver, "receiver", "r", false, "channel is a receiver")
	eventCmd.Flags().StringVarP(&eventPayload, "payload", "p", "", "channel will use this payload instead of creating one")
	eventCmd.Flags().StringVarP(&eventProxy, "proxy", "t", "", "proxy channel to another event (should match existing event name), implies receiver")
	eventCmd.Flags().StringVarP(&eventChannel, "channel", "c", "", "override channel default")
}
