/*
Copyright © 2019 Emanuel Bennici <benniciemanuel78@gmail.com>

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
package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"

	"gitlab.com/l0nax/changelog-go/cmd"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/internal/log"
	"gitlab.com/l0nax/changelog-go/pkg/version"
)

func main() {
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(time.Second * 2)
	defer sentry.Recover()

	// check if there are any updates available, if so print info message
	// has a timeout of 500ms to not slow down the changelog application
	// if its offline.
	if internal.CheckUpdate(version.Version,
		"https://update.l0nax.org/changelog-go") {
		fmt.Printf("[UPDATE] There is a new version available, run" +
			" 'changelog update' or 'snap refresh changelog' " +
			"to update to the latest version.\n\n")
	}

	log.Log.Debug("changelog-go. Version [%s] Commit [%s]", version.Version, version.Hash)

	// NOTE: We are currently disabling the AUTOMATIC update, because we
	//	 have now the 'update' subcommand. So we do not force clients
	//	 to update to the newest version. See Issue #2
	// // configure and start self-updater
	// var updater = &selfupdate.Updater{
	//         CurrentVersion: version.Version,
	//         ApiURL:         "https://l0nax.gitlab.io/",
	//         BinURL:         "https://l0nax.gitlab.io/",
	//         DiffURL:        "https://l0nax.gitlab.io/",
	//         Dir:            "update/",
	//         CmdName:        "changelog-go",
	// }

	// if updater != nil {
	//         go updater.BackgroundRun()
	// }

	cmd.Execute()
}
