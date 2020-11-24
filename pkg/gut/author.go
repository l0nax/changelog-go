package gut

import (
	gitcfg "github.com/tcnksm/go-gitconfig"
)

// GetAuthorName returns the author name configured in the git config
//
// Deprecated: The gut package is marked as deprecated since we will drop the git dependency in the future.
func GetAuthorName() (string, error) {
	name, err := gitcfg.Username()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return name, nil
}
