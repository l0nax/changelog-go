package tools

import (
	"os"
)

// CheckDirExist checks if a directory exists or not and returns 'true' or
// 'false'.
func CheckDirExist(_path string) bool {
	if _, err := os.Stat(_path); os.IsNotExist(err) {
		return false
	}

	return true
}
