package cmd

// Copyright © 2019 Robert Gordon <rbg@openrbg.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "simplekv",
	Short: "An exercise in key value store",
	Long:  "An exercise in key value store",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "cfg file (default: $HOME/.simplekv.yaml)")

	RootCmd.PersistentFlags().Bool("debug", false, "Turn on verbose logging")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	RootCmd.PersistentFlags().String("api", ":7800", "API Endpoint")
	viper.BindPFlag("api", RootCmd.PersistentFlags().Lookup("api"))

	// --------------
	//    Airbrake
	// --------------
	// project id
	RootCmd.PersistentFlags().Int64("ab_proj", 0, "Airbrake Project id")
	viper.BindPFlag("ab_proj", RootCmd.PersistentFlags().Lookup("ab_proj"))
	// project key
	RootCmd.PersistentFlags().String("ab_key", "", "Airbrake Project key")
	viper.BindPFlag("ab_key", RootCmd.PersistentFlags().Lookup("ab_key"))
	// project environment
	RootCmd.PersistentFlags().String("ab_env", "dev", "Airbrake Environment")
	viper.BindPFlag("ab_env", RootCmd.PersistentFlags().Lookup("ab_env"))
	// host URL
	RootCmd.PersistentFlags().String("ab_url", "", "Airbrake URL")
	viper.BindPFlag("ab_url", RootCmd.PersistentFlags().Lookup("ab_url"))
	// --------------
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.SetConfigName(".simplekv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
