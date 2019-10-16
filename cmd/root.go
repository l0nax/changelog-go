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
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/l0nax/changelog-go/internal"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"gitlab.com/l0nax/changelog-go/pkg/gut"
	"gopkg.in/src-d/go-git.v4"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "changelog-go",
	Short: "Changelog Management Tool",
	Long: `changelog-go  helps you to keep track of you
Changelog (and changes) and its fully compatible with (eg) the Git Flow.

It extends your DevOps Workflow and gives other People
the possibility to read a beautiful formatted CHANGELOG.md
File.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	__path, _ := gut.FindGitRoot("/home/l0nax/.go/src/github.com/l0nax/kubenab/cmd/kubenab")
	log.Printf("--> %s", __path)

	_, err := git.PlainOpen(__path)
	if err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.changelog-go.yaml)")
}

func initConfig() {
	// register all Entry Types
	internal.RegisterAll()

	// TODO: Change this to scan for a Git Directory

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

		// Search config in home directory with name ".changelog-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".changelog-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
