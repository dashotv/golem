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
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dashotv/golem/config"
)

var cfg = &config.Config{}
var cfgFile string
var colors bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golem [options] [command]",
	Short: "generate mongo database models based on MDM library",
	Long:  "generate mongo database models based on MDM library",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.golem.yaml)")
	rootCmd.PersistentFlags().BoolVar(&colors, "colors", true, "use color output")
	au = aurora.NewAurora(colors)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".golem" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("./.golem")
		viper.AddConfigPath(home)
		viper.SetConfigName(".golem")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("unable to read config: %s\n", err)
		return
	}

	if err := viper.Unmarshal(cfg); err != nil {
		logrus.Fatalf("failed to unmarshal configuration file: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("failed to validate config: %s", err)
	}

	cfg.File = cfgFile
}

// func executeCommand(name string, arg ...string) error {
// 	cmd := exec.Command(name, arg...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
//
// 	return cmd.Run()
// }
//
// func checkBinaries() error {
// 	bins := []string{
// 		"cobra",
// 	}
//
// 	for _, b := range bins {
// 		_, err := exec.LookPath(b)
// 		if err != nil {
// 			return errors.Wrap(err, "couldn't find binary in path: "+b)
// 		}
// 	}
//
// 	return nil
// }
