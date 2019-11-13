// gut (git utils) implements some Helper Functions around the `git` Library to easier
// interact with the Repository.
package gut

import (
	"os"
	"path"
)

// maxDepth describes how many levels of directories this application will
// search for the Changelog Config.
// This is needed because otherwise the application will search for the Config
// forever.
const maxDepth = 16

// FindGitRoot returns the last found `.git` Directory from the base
// Means if your are in $PATH/cmd/test and the Git Root is $PATH/
// FuncGitRoot will return $PATH/
// It will return an Error if the `.git/` Directory could not be found.
func FindGitRoot(base string) (string, error) {
	var _path = base

	for {
		if _path == "/" {
			return "", os.ErrNotExist
		}

		// check if `$PATH/.git` exists, if not go one Directory back
		if _, err := os.Stat(path.Join(_path, ".git/")); !os.IsNotExist(err) {
			// found
			return _path, nil
		}

		_path = path.Join(_path, "../")
	}
}
