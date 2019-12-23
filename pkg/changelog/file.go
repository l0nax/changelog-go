package changelog

import (
	"errors"
	"os"
	"path"

	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
)

// createReleaseDir creates the needed directory structure for all Changelog
// releases
func createReleaseDir(version string) error {

	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released")
	// unreleased := path.Join(basePath, "unreleased")

	verPath := path.Join(releasedPath, version)

	// NOTE: Should we add an 'force' flag to delete the existing directory?
	if _, err := os.Stat(verPath); !os.IsNotExist(err) {
		// error since this path already exists
		return errors.New("Release Path already exists.")
	}

	return nil
}
