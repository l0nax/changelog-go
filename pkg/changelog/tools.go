package changelog

// This File contains all basic Functions

import (
	"github.com/spf13/viper"
	"os"
	"path"
)

// CheckDir will check if the changelog directory exists
// if not it will create and initialize it.
func CheckDir() {
	// get Changelog data directory
	dataDir := viper.GetString("changelog.entryPath")
	if len(dataDir) == 0 {
		panic("'changelog.entryPath' must be specified!")
	}

	// create Directory structure if it does not exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(dataDir, "unreleased"), 0755)
		if err != nil {
			panic(err)
		}

		err = os.Mkdir(path.Join(dataDir, "released"), 0755)
		if err != nil {
			panic(err)
		}
	}
}
