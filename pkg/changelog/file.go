package changelog

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
)

// createReleaseDir creates the needed directory structure for all Changelog
// releases
func createReleaseDir(version string) error {
	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released", version)
	// unreleased := path.Join(basePath, "unreleased")

	log.Debugf("New release path: '%s'\n", releasedPath)

	// NOTE: Should we add an 'force' flag to delete the existing directory?
	if _, err := os.Stat(releasedPath); !os.IsNotExist(err) {
		// error because this release already exists
		return errors.New("Release Path already exists.")
	}

	err := os.MkdirAll(releasedPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// moveToReleaseFolder moves all unreleased changelog-entry files to the specifc
// release directory.
func moveToReleaseFolder(version string) error {
	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released", version)
	unreleased := path.Join(basePath, "unreleased")

	// get all files from the unrealed folder
	err := filepath.Walk(unreleased, func(_path string, info os.FileInfo, err error) error {
		// check if type is file, if not just skip
		if info.IsDir() {
			return nil
		}

		// since it's a file, we can now move it to the new directory.
		return os.Rename(_path, path.Join(releasedPath, info.Name()))
	})

	if err != nil {
		return errors.Wrap(err, "Error while moving changelog-entry to release directory")
	}

	return nil
}
