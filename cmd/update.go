/*
Copyright © 2020 Francesco Emanuel Bennici <benniciemanuel78@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/l0nax/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"gitlab.com/l0nax/changelog-go/pkg/version"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Manually run the auto-updater",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Checking for new available versions...")

		var updater = &selfupdate.Updater{
			CurrentVersion: version.Version,
			ApiURL:         "https://l0nax.gitlab.io/",
			BinURL:         "https://l0nax.gitlab.io/",
			DiffURL:        "https://l0nax.gitlab.io/",
			Dir:            "update/",
			CmdName:        "changelog-go",

			// force the Updater to run the update check
			ForceCheck: true,
		}

		if updater == nil {
			log.Fatalln("Could not create 'updater' object (nil)")
		}

		// run Updater
		if err := updater.BackgroundRun(); err != nil {
			log.Fatalf("Error while updating: %s\n", err)
		} else {
			log.Infoln("Updated successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
