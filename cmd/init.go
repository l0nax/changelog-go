/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"path"

	"gitlab.com/l0nax/changelog-go/internal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the default config file",
	Run: func(cmd *cobra.Command, args []string) {
		// check if Config file already exists
		if _, err := os.Stat(path.Join(internal.GitPath, ".changelog-go.yaml")); !os.IsNotExist(err) {
			if !force {
				log.Fatalln("There's already a config file. To overwrite run with the flag '--force'.")
			}
		}

		// create new config file
		file, err := os.Create(path.Join(internal.GitPath, ".changelog-go.yaml"))
		if err != nil {
			log.Fatal(err)
		}

		// write default config
		//			       /example/.changelog-go.yaml
		data, _ := internal.FS.String("/example/.changelog-go.yaml")
		file.WriteString(data)

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.Flags().BoolVar(&force, "force", false, "force writting a new config file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
