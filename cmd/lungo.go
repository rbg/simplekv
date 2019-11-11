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
	"github.com/apex/log"
	"github.com/rbg/simplekv/api"
	"github.com/rbg/simplekv/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// memCmd represents the mem command
var lungoCmd = &cobra.Command{
	Use:   "lungo",
	Short: "Simple kv using lungo MongoDB look-a-like",
	Long:  "Simple kv using lungo MongoDB look-a-like",
	Run:   lungoKV,
}

func init() {
	RootCmd.AddCommand(lungoCmd)

	lungoCmd.Flags().String("db", "lungo", "database name")
	viper.BindPFlag("db", lungoCmd.Flags().Lookup("db"))

	lungoCmd.Flags().String("col", "kv", "Collection Name")
	viper.BindPFlag("col", lungoCmd.Flags().Lookup("col"))
}

func lungoKV(cmd *cobra.Command, args []string) {

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if be := store.NewLungo(viper.GetString("db"), viper.GetString("col")); be != nil {
		if kvapi := api.New(be, viper.GetString("api")); kvapi != nil {
			kvapi.Run()
		}
	}
}
