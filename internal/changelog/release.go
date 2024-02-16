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

// Release represents a single release.
type Release struct {
	// Info holds all the meta informations about the release.
	Info ReleaseInfo
	// Entries holds all the change entries.
	Entries []Entry
	// Collapse defines whether the release should be collapsed
	// in the generated file or not.
	//
	// It is set to true based on the project configuration.
	Collapse bool
}
