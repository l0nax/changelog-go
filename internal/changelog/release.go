package changelog

import "time"

// ReleaseInfo represents a single release with all its
// meta informations.
type ReleaseInfo struct {
	Version     string    `koanf:"version"`
	ReleaseDate time.Time `koanf:"date"`
	// IsPreRelease defines whether the release is a pre-release
	// or not.
	IsPreRelease bool `koanf:"prerelease"`
}
