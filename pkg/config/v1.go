package config

import (
	"os"
	"path"

	"github.com/spf13/viper"

	"gitlab.com/l0nax/changelog-go/internal"
)

// Check will test if all required Paths and files are represent
// and if not it will TRY to create them. If an Error occurs while
// creating those Files/ Dirs Check() will return the Error.
func (c *V1) Check() error {
	// try to Unmarshal the Config
	err := viper.Unmarshal(c)
	if err != nil {
		return err
	}

	c.Changelog.EntryPath = path.Join(path.Join(internal.GitPath, "../"), c.Changelog.EntryPath)

	// check if Changelog Directory exists
	if _, err := os.Stat(c.Changelog.EntryPath); os.IsNotExist(err) {
		// try to create the Directory-Structure
		for _, entry := range v1ChangelogFolderStruct {
			err = os.MkdirAll(path.Join(c.Changelog.EntryPath, entry), 0o766)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
