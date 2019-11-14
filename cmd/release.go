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
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <version>",
	Short: "This Command will generate the CHANGELOG.md file.",
	Long: `The CHANGELOG.md will be generated and – depending on
the configuration and flags – the Entry Fiels will be moved to the 'released'
Folder.`,
	Example: `  release v1.0.0`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().BoolP("pre-relrease", "p", false, `Mark this Release as pre-release in the CHANGELOG.md.
The Output and how the Application reacts depends on your Configuration.`)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
