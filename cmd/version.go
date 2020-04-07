/*
Copyright © 2020 Emanuel Bennici <benniciemanuel78@gmail.com>

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

	"github.com/spf13/cobra"
	"gitlab.com/l0nax/changelog-go/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the Version information about this binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-20v: %s\n", "Version", version.Version)
		fmt.Printf("%-20v: %s\n", "Build", version.BuildTime)
		fmt.Printf("%-20v: %s\n", "Git Hash", version.Hash)
		fmt.Printf("%-20v: %s\n", "Build Environment", version.Environment)

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
