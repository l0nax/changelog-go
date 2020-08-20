package gut

import (
	"os"
	"path"
)

// maxDepth describes how many levels of directories this application will
// search for the '.git' Directory.
// This is needed because otherwise the application will search for the Git Root
// forever.
const maxDepth = 16

// FindGitRoot returns the last found `.git` Directory from the base
// Means if your are in $PATH/cmd/test and the Git Root is $PATH/
// FuncGitRoot will return $PATH/
// It will return an Error if the `.git/` Directory could not be found.
func FindGitRoot(base string) (string, error) {
	var i int = 0
	_path := base

	for {
		// if we are at the Root Directory stop searching
		if i >= maxDepth {
			return "", os.ErrNotExist
		}

		// check if `$PATH/.git` exists, if not go one Directory back
		if _, err := os.Stat(path.Join(_path, ".git/")); !os.IsNotExist(err) {
			// found
			return _path, nil
		}

		_path = path.Join(_path, "../")
		i++
	}
}
