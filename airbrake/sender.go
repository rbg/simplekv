package airbrake

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
	goab "github.com/airbrake/gobrake/v4"
	"github.com/apex/log"
	"github.com/spf13/viper"
)

var (
	// GoBrake is our interface to airbrake.
	GoBrake *goab.Notifier
)

// Init should make happy notifier juju.
func Init() {
	GoBrake = goab.NewNotifierWithOptions(
		&goab.NotifierOptions{
			ProjectId:   viper.GetInt64("ab_proj"),
			ProjectKey:  viper.GetString("ab_key"),
			Environment: viper.GetString("ab_env"),
			Host:        viper.GetString("ab_url"),
		},
	)
	log.Infof("Here is the GoBrake struct: %+#v", GoBrake)
}
