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
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/l0nax/changelog-go/internal"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const cfgFileName = ".changelog-go.yaml"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "changelog-go",
	Short: "Changelog Management Tool",
	Long: `changelog-go helps you to keep track of your
Changelog (and changes) and its fully compatible with (eg) the Git Flow.

It extends your DevOps Workflow and gives other People
the possibility to read a beautiful formatted CHANGELOG.md
File.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// check if debug flag has been set
		if d, _ := cmd.Flags().GetBool("debug"); d {
			log.SetLevel(logrus.DebugLevel)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.changelog-go.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enables debug output")
}

func initConfig() {
	var err error

	// register all Entry Types
	internal.RegisterAll()

	internal.Cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	internal.GitPath, err = findTopDir(internal.Cwd)
	if err != nil {
		panic(err)
	}

	if internal.GitPath == "" {
		fmt.Println("Unable to find config file! Use `changelog init` to create one.")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

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

		viper.AddConfigPath(internal.GitPath)

		// Search config in home directory with name ".changelog-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".changelog-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to read config: %+v\n", err)
		os.Exit(1)
	}
}

// findTopDir does does in the cwd directory and searches (going against /)
// the directoy where `.changelo-go` does reside.
// An error will be returned if on of the I/O functions does return one.
func findTopDir(cwd string) (string, error) {
	var err error
	var fs os.FileInfo
	var c bool
	var prevPath string // prevPath contains the previous path to prevent an infinite loop

	// check does return true if the p directory does contain our changelog-go config.
	check := func(p string) (bool, error) {
		if fs, err = os.Stat(path.Join(p, cfgFileName)); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}

			return false, err
		}

		return !fs.IsDir(), nil
	}

	c, err = check(cwd)
	if err != nil {
		return "", err
	}

	// the config does reside in the cwd
	if c {
		return cwd, nil
	}

	for prevPath != cwd {
		prevPath = cwd
		cwd = path.Join(cwd, "../")
		c, err = check(cwd)
		if err != nil {
			return "", err
		}

		// the config does reside in the cwd
		if c {
			return cwd, nil
		}
	}

	return "", nil
}
