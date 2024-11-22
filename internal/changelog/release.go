package changelog

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// ReleaseInfo represents a single release with all its
// meta informations.
type ReleaseInfo struct {
	Version     string    `toml:"version"`
	ReleaseDate time.Time `toml:"date"`
	// IsPreRelease defines whether the release is a pre-release
	// or not.
	IsPreRelease bool `toml:"pre_release"`
}

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

type GrouppedEntries struct {
	// ChangeType holds the change type information.
	// It can be used in the CHANGELOG template.
	ChangeType config.ChangeType

	Entries []Entry
}

// Release represents a single release.
type Release struct {
	// Info holds all the meta informations about the release.
	Info ReleaseInfo

	// Entries holds all the change entries.
	Entries []Entry

	// GrouppedEntries is a meta field containing a
	// groupped and sorted list of all entries
	// from the Entries field.
	//
	// This means that this field needs to be manually filled!
	GrouppedEntries []GrouppedEntries

	// Collapse defines whether the release should be collapsed
	// in the generated file or not.
	//
	// It is set to true based on the project configuration.
	Collapse bool
}

// Create creates the release, i.e. it does not exist yet
// and moves all changelog entries to the new release directory.
func (r Release) Create() error {
	releaseDir := filepath.Join(config.C.ChangelogDir, ReleasedDir, r.Info.Version)

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

	// NOTE: we always copy the files instead of renaming them.
	//       The performance degredation shouldn't be much of a problem.
	//       This saves us from having to implement a fallback (osRename -> cannot move because of FS)
	//       If it ever becomes a problem, we can implement it.
	for _, entry := range r.Entries {
		srcPath, ok := entry.EntryPath.Deconstruct()
		if !ok {
			return errors.Errorf("missing entry path for entry %+v", entry)
		}

		dstPath := filepath.Join(releaseDir, filepath.Base(srcPath))

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}

		// we will not delete changlog entry files if it is a pre-release
		if r.Info.IsPreRelease {
			continue
		}

		if err := os.Remove(srcPath); err != nil {
			return errors.Wrapf(err, "unable to remove old entry file at %q", srcPath)
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "unable to open file %q", dstPath)
	}
	defer dst.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrapf(err, "unable to open file %q", srcPath)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return errors.Wrapf(err, "unable to copy file from %q to %q", srcPath, dstPath)
	}

	return nil
}
