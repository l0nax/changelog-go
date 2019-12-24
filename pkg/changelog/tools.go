package changelog

// This File contains all basic Functions

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/tools"
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
	var files = make(map[string][]byte)

	// initialize r.Info if its not already initialized
	if r.Info == nil {
		r.Info = &ReleaseInfo{}
	}

	// First check if everything is ok
	CheckDir()

	unreleasedPath := path.Join(internal.GitPath,
		viper.GetString("changelog.entryPath"), "unreleased")

	files, err := ReadEntryFiles(unreleasedPath)
	if err != nil {
		log.Panic(err)
	}

	// parse file data into struct
	r.Entries, err = ParseFiles(files)
	if err != nil {
		return err
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
	err := os.MkdirAll(verPath, 0755)
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

// ReadEntryFiles reads all changelog entry files form a given path
func ReadEntryFiles(filesPath string) (map[string][]byte, error) {
	var files = make(map[string][]byte)

	err := filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
		// check if path is File
		if info.IsDir() {
			return nil
		}

		log.Debugf("Reading '%s'\n", path)

		// read file content
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		files[path] = data
		return nil
	})

	return files, err
}

// ParseFiles parses all given files into an Entries struct
// @param files is a map of 'FILE_PATH' => 'FILE_CONTENT'
func ParseFiles(files map[string][]byte) ([]entry.Entry, error) {
	var entries = []entry.Entry{}

	for _, file := range files {
		// unmarshal Filecontent into Struct
		changeEntry := entry.Entry{}
		err := tools.YAMLUnmarshal(file, &changeEntry)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		// append new Filecontent to Release Data
		entries = append(entries, changeEntry)
	}

	return entries, nil
}
