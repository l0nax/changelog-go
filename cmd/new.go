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
	"bufio"
	"fmt"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/changelog"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/gut"
	"os"
	"strconv"
	"strings"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new Changelog-Entry",
	Long: `"new" creates a new Changelog-Entry so you can easily commit
your entry.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check if first Argument is set
		title := args[0]

		if len(title) == 0 {
			log.Fatalln("The Title of your Change MUST be specified!")
		}

		// check if needed Directories exists
		changelog.CheckDir()

		var author string
		var err error

		// get Author if enabled in Config
		if viper.GetBool("entry.author") {
			log.Debug("'entry.author' is enabled, so the Author will be grabbed and used.")

			// get the Author
			author, err = gut.GetAuthorName()
			if err != nil {
				os.Exit(1)
			}
		}

		// create Entry
		changelogEntry := entry.Entry{
			ChangeTitle: title,
			Author:      author,
		}

		// get all available Entry Types
		for i, entryType := range internal.EntryT.ListAvailableTypes() {
			fmt.Printf("[%d] %- 20s (%s)\n", i,
				(*entryType).GetTypeDescription(),
				(*entryType).GetShortTypeName())
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")

		// ask User about Entry Type
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		// convert input to Number
		var choosenType int

		if choosenType, err = strconv.Atoi(text); err != nil {
			log.Fatalln("Please enter only VALID numbers (error while converting)!")
		}

		// check if Number does exists in our Entry Types List
		if choosenType >= len(internal.EntryT.ListAvailableTypes()) {
			log.Debugln(len(internal.EntryT.ListAvailableTypes()))
			log.Debugf("%# v\n", pretty.Formatter(internal.EntryT.ListAvailableTypes()))
			log.Fatalln("Please choose a VALID number!")
		}

		// create Entry
		changelogEntry.Type = internal.EntryT.ListAvailableTypes()[choosenType]
		changelog.AddEntry(changelogEntry)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
