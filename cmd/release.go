/*
Copyright © 2019 Emanuel Bennici <eb@fabmation.de>

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
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/l0nax/changelog-go/pkg/changelog"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <version>",
	Short: "This Command will generate the CHANGELOG.md file.",
	Long: `The CHANGELOG.md will be generated and – depending on
the configuration and flags – the Entry files will be moved to the 'released'
Folder.`,
	Example: `  release 1.0.0`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newRelease := changelog.Release{}
		newRelease.Info = &changelog.ReleaseInfo{}

		r := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
		newRelease.Info.Version = r.FindStringSubmatch(args[0])

		// this Variable describes if the current release is a Pre-Release
		var isPreRelrease bool

		// check if Version is a pre-release
		if viper.GetBool("preRelease.detect") {
			// check if Version is a pre-release
			for i, group := range r.SubexpNames() {
				if group == "prerelease" {
					if len(newRelease.Info.Version[i]) != 0 {
						isPreRelrease = true
					}

					break
				}
			}
		}

		log.Debugf("Releasing version '%#v'\n", newRelease.Info.Version)

		fIsPreRelease, err := cmd.Flags().GetBool("pre-release")
		if err != nil {
			log.Fatal(err)
		}

		if isPreRelrease || fIsPreRelease {
			newRelease.Info.IsPreRelease = true
			log.Debugln("Current release is a Pre-Release!")
		}

		newRelease.Info.ReleaseDate, err = cmd.Flags().GetString("release-date")
		if err != nil {
			log.Fatal(err)
		}

		// generate CHANGELOG.md
		changelog.GenerateChangelog(&newRelease)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringP("release-date", "d", time.Now().Format("2006-01-02"), "set the release date, defaults to today")
	releaseCmd.Flags().BoolP("pre-release", "p", false, `Mark this Release as pre-release in the CHANGELOG.md.
The Output and how the Application reacts depends on your Configuration.`)
}
