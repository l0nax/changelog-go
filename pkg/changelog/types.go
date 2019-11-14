package changelog

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
)

// ReleaseInfo holds all informations about the current release
type ReleaseInfo struct {
	Version      []string // Contains the string slice output of the Regex
	IsPreRelease bool     // True if the release is a pre-release
}

// Release contains the ReleaseInfo and all NEW Changelog entry types.
type Release struct {
	Info    *ReleaseInfo  // All Informations about the current Release
	Entries []entry.Entry // Holds all NEW changelog entry Types
}
