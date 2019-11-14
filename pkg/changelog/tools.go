package changelog

// This File contains all basic Functions

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/tools"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
		err = os.MkdirAll(path.Join(dataDir, "unreleased"), 0755)
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(path.Join(dataDir, "released"), 0755)
		if err != nil {
			panic(err)
		}
	}
}

// GetEntries returns a List of all Entries found in the Changelog-Data directory.
func GetEntries(r *Release) error {
	// contains all found unreleased Files
	var files map[string][]byte

	// initialize r.Info if its not already initialized
	if r.Info == nil {
		r.Info = &ReleaseInfo{}
	}

	// First check if everything is ok
	CheckDir()

	changelogPath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	unreleasedPath := path.Join(changelogPath, "unreleased")

	err := filepath.Walk(unreleasedPath, func(path string, info os.FileInfo, err error) error {
		// read file content
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		files[path] = data
		return nil
	})

	if err != nil {
		log.Panic(err)
		return err
	}

	for _, fileData := range files {
		// unmarshal Filecontent into Struct
		changeEntry := entry.Entry{}
		err = tools.YAMLUnmarshal(fileData, &changeEntry)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// append new Filecontent to Release Data
		(*r).Entries = append(r.Entries, changeEntry)
	}

	return nil
}

// MoveEntries will move all unreleased Entries to 'released/{{ .Version }}'
func MoveEntries(version string) error {
	changelogPath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	unreleasedPath := path.Join(changelogPath, "unreleased")
	releasedPath := path.Join(changelogPath, "released")
	verPath := path.Join(releasedPath, version)

	// check if a Directory already exists
	if _, err := os.Stat(verPath); !os.IsNotExist(err) {
		return os.ErrExist
	}

	// create new release directory
	err := os.Mkdir(verPath, 0755)
	if err != nil {
		return err
	}

	// move all Files to the new release directory
	err = filepath.Walk(unreleasedPath, func(_path string, info os.FileInfo, err error) error {
		// move File
		return os.Rename(_path, path.Join(verPath, info.Name()))
	})

	if err != nil {
		return err
	}

	log.Debug("Moved all unreleased Files")

	return nil
}
