package changelog

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// ReleaseInfo is the metadata of a single release.
type ReleaseInfo struct {
	Version     string    `toml:"version"`
	ReleaseDate time.Time `toml:"date"`
	// IsPreRelease reports whether the release is a pre-release.
	IsPreRelease bool `toml:"pre_release"`
}

// SaveToFile writes r to path.
func (r ReleaseInfo) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	raw, err := toml.Marshal(r)
	if err != nil {
		return err
	}

	_, err = file.Write(raw)

	return err
}

// GrouppedEntries are the entries of one release that share a change type.
type GrouppedEntries struct {
	// ChangeType is the type the entries share.
	ChangeType config.ChangeType

	Entries []Entry
}

// Release is a single release and its entries.
type Release struct {
	// Info is the metadata of the release.
	Info ReleaseInfo

	// Entries are the change entries of the release.
	Entries []Entry

	// GrouppedEntries holds the entries grouped by change type and sorted.
	//
	// NOTE: [Changelog.Render] fills this field.
	GrouppedEntries []GrouppedEntries

	// DisplayVersion is the version with the configured version prefix
	// applied.
	//
	// NOTE: [Changelog.Render] fills this field.
	DisplayVersion string

	// DisplayDate is the release date as it should be rendered. It is empty
	// for the pending set, which has no date yet.
	//
	// This means that this field needs to be manually filled!
	DisplayDate string

	// Collapse renders the release inside a collapsed "<details>" block.
	Collapse bool
}

// VersionEffect returns the strongest version effect across the entries of the
// release.
func (r Release) VersionEffect() config.VersionEffect {
	return versionEffect(r.Entries)
}

// CreateRelease creates the directory of release r and moves the changelog
// entries into it. It returns an error if the release already exists.
func (p *Project) CreateRelease(r Release) error {
	releaseDir := filepath.Join(p.ReleasedDir(), r.Info.Version)

	// prevent overwriting if the release already exists
	if _, err := os.Stat(releaseDir); err == nil {
		return errors.Errorf("a release directory at %q already exists", releaseDir)
	}

	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return errors.Wrapf(err, "unable to create release directory at %q", releaseDir)
	}

	releaseInfoPath := filepath.Join(releaseDir, ReleaseInfoFileName)

	if err := r.Info.SaveToFile(releaseInfoPath); err != nil {
		return err
	}

	// Copied rather than renamed: a rename across filesystems fails and
	// would need a fallback.
	for i := range r.Entries {
		entry := &r.Entries[i]

		srcPath, ok := entry.EntryPath.Deconstruct()
		if !ok {
			return errors.Errorf("missing entry path for entry %+v", entry)
		}

		dstPath := filepath.Join(releaseDir, filepath.Base(srcPath))

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}

		// a pre-release keeps its fragments for the final release
		if r.Info.IsPreRelease {
			continue
		}

		if err := os.Remove(srcPath); err != nil {
			return errors.Wrapf(err, "unable to remove old entry file at %q", srcPath)
		}
	}

	return nil
}

// RemoveSupersededPreReleases deletes the release directory of every
// pre-release that version supersedes.
func (p *Project) RemoveSupersededPreReleases(version string) error {
	base, err := ParseVersion(version)
	if err != nil {
		return err
	}

	released, err := p.ParseReleased()
	if err != nil {
		return err
	}

	for i := range released.Releases {
		release := &released.Releases[i]

		if !release.Info.IsPreRelease {
			continue
		}

		other, err := ParseVersion(release.Info.Version)
		if err != nil {
			return err
		}

		if !BaseVersion(other).EQ(BaseVersion(base)) {
			continue
		}

		dir := filepath.Join(p.ReleasedDir(), release.Info.Version)

		slog.Info("Removing superseded pre-release",
			slog.String("version", release.Info.Version), slog.String("path", dir))

		if err := os.RemoveAll(dir); err != nil {
			return errors.Wrapf(err, "unable to remove pre-release directory %q", dir)
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrapf(err, "unable to open file %q", srcPath)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "unable to open file %q", dstPath)
	}

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()

		return errors.Wrapf(err, "unable to copy file from %q to %q", srcPath, dstPath)
	}

	// closed explicitly: a failure to flush would truncate the copy
	if err := dst.Close(); err != nil {
		return errors.Wrapf(err, "unable to close file %q", dstPath)
	}

	return nil
}
