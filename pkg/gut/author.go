package gut

import (
	gitcfg "github.com/tcnksm/go-gitconfig"
)

// GetAuthorName returns the author name configured in the git config
func GetAuthorName() (string, error) {
	name, err := gitcfg.Username()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return name, nil
}
