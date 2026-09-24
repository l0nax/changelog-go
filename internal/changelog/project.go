package changelog

import (
	"path/filepath"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// Project ties a loaded configuration to the directory it was loaded from.
// Every path in the configuration is interpreted relative to that directory.
type Project struct {
	cfg     config.Config
	rootDir string
}

// NewProject returns the project described by cfg, loaded from a config file
// in rootDir.
func NewProject(cfg config.Config, rootDir string) *Project {
	return &Project{
		cfg:     cfg,
		rootDir: rootDir,
	}
}

// Config returns the project configuration.
func (p *Project) Config() config.Config {
	return p.cfg
}

// RootDir returns the directory holding the config file.
func (p *Project) RootDir() string {
	return p.rootDir
}

// resolve returns path relative to the project root. An absolute path is
// returned unchanged.
func (p *Project) resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(p.rootDir, path)
}

// ChangelogDir returns the directory holding the changelog files.
func (p *Project) ChangelogDir() string {
	return p.resolve(p.cfg.ChangelogDir)
}

// ReleasedDir returns the directory holding the released versions.
func (p *Project) ReleasedDir() string {
	return filepath.Join(p.ChangelogDir(), ReleasedDir)
}

// UnreleasedDir returns the directory holding the unreleased entries.
func (p *Project) UnreleasedDir() string {
	return filepath.Join(p.ChangelogDir(), UnreleasedDir)
}

// OutputPath returns the file the changelog is rendered to.
func (p *Project) OutputPath() string {
	return p.resolve(p.cfg.OutputPath.UnwrapOr(config.DefaultOutputPath))
}

// DisplayVersion returns version with the configured version prefix applied.
func (p *Project) DisplayVersion(version string) string {
	return ApplyVersionPrefix(p.cfg.VersionPrefix, version)
}

// ReleaseChangelog returns a changelog holding only the named release.
//
// The version is matched semantically, so the caller may pass it with or
// without the configured prefix -- which is what makes
// "changelog show $(changelog latest)" work.
func (p *Project) ReleaseChangelog(version string) (*Changelog, error) {
	want, err := ParseVersion(version)
	if err != nil {
		return nil, err
	}

	released, err := p.ParseReleased()
	if err != nil {
		return nil, err
	}

	for i := range released.Releases {
		got, err := ParseVersion(released.Releases[i].Info.Version)
		if err != nil {
			return nil, err
		}

		if !got.EQ(want) {
			continue
		}

		return &Changelog{
			VersionPrefix: p.cfg.VersionPrefix,
			Releases:      []Release{released.Releases[i]},
		}, nil
	}

	return nil, &VersionNotFoundError{Version: version}
}

// UnreleasedChangelog returns a changelog holding the pending entries as a
// single, versionless release.
func (p *Project) UnreleasedChangelog() (*Changelog, error) {
	entries, err := p.LoadUnreleasedEntries()
	if err != nil {
		return nil, err
	}

	return &Changelog{
		VersionPrefix: p.cfg.VersionPrefix,
		Releases:      []Release{{Entries: entries}},
	}, nil
}

// VersionNotFoundError reports that a version has not been released.
type VersionNotFoundError struct {
	Version string
}

func (e *VersionNotFoundError) Error() string {
	return "no release found for version " + e.Version
}
