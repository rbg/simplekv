// Copyright © 2019 Robert Gordon <rbg@openrbg.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/apex/log"
	"github.com/rbg/simplekv/api"
	"github.com/rbg/simplekv/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Start up a redis backed simplekv",
	Long:  "Start up a redis backed simplekv",
	Run:   redisKV,
}

func init() {
	RootCmd.AddCommand(redisCmd)

	redisCmd.Flags().String("rw", "localhost:6379", "Read/Write endpoint")
	viper.BindPFlag("rw", redisCmd.Flags().Lookup("rw"))

	redisCmd.Flags().String("ro", "", "Read only endpoint")
	viper.BindPFlag("ro", redisCmd.Flags().Lookup("ro"))
}

func redisKV(cmd *cobra.Command, args []string) {

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if be := store.NewRedis(
		viper.GetString("rw"),
		viper.GetString("ro")); be != nil {
		if kvapi := api.New(be, viper.GetString("api")); kvapi != nil {
			kvapi.Run()
		}
	}
}
