package changelog

import (
	"bytes"
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/blang/semver/v4"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"go.l0nax.org/typact"
	"gopkg.in/yaml.v3"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

const (
	// ReleasedDir is the name of the directory holding the released
	// versions.
	ReleasedDir = "released"
	// UnreleasedDir is the name of the directory holding the unreleased
	// entries.
	UnreleasedDir = "unreleased"
	// ReleaseInfoFileName is the name of the "ReleaseInfo" filename.
	ReleaseInfoFileName = "ReleaseInfo"
)

// Changelog represents the final CHANGELOG file.
type Changelog struct {
	Releases []Release
}

// SaveToFile generates a new changelog and saves it to path.
func (c *Changelog) SaveToFile(path string) error {
	c.sortReleaseEntries()

	// TODO: Allow overriding the default
	tmpl, err := template.New("changelog-tmpl").
		Funcs(template.FuncMap{
			"formatTime": formatTime,
		}).
		Parse(defaultChangelogScheme)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	out.Grow(1024 * 1024) // 1 MB

	err = tmpl.Execute(&out, c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(out.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Changelog) sortReleaseEntries() error {
	// first we sort all releases
	slices.SortStableFunc(c.Releases, func(a, b Release) int {
		aVer, _ := semver.Make(a.Info.Version)
		bVer, _ := semver.Make(b.Info.Version)

		return bVer.Compare(aVer)
	})

	// now we sort all changelog entries
	for i, release := range c.Releases {
		grouped := groupBy(release.Entries, func(item Entry) string {
			return item.ChangeTypeID
		})

		groupedEntries := make([]GrouppedEntries, 0, len(grouped))

		for typeID, entries := range grouped {
			changeType, ok := resolveChangeTypeID(typeID)
			if !ok {
				panic(fmt.Sprintf("Unable to resolve already parsed and validated changelog type ID %q", typeID))
			}

			// sort the entries
			slices.SortStableFunc(entries, func(a, b Entry) int {
				return cmp.Compare(a.Title, b.Title)
			})

			for _, ent := range entries {
				err := ent.LoadChangeType()
				if err != nil {
					return errors.Wrapf(err, "unknown change type with ID %q at %v", ent.ChangeTypeID, ent.EntryPath)
				}
			}

			groupedEntries = append(groupedEntries, GrouppedEntries{
				ChangeType: changeType,
				Entries:    entries,
			})
		}

		// now we can sort all the entries
		slices.SortStableFunc(groupedEntries, func(a, b GrouppedEntries) int {
			return cmp.Compare(a.ChangeType.ID, b.ChangeType.ID)
		})

		release.GrouppedEntries = groupedEntries
		c.Releases[i] = release
	}

	return nil
}

// LoadUnreleasedEntries loads all unreleased entries.
func LoadUnreleasedEntries() ([]Entry, error) {
	dir := filepath.Join(config.C.ChangelogDir, UnreleasedDir)

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		} else if strings.HasPrefix(entry.Name(), ".") {
			// hidden files are ignored
			continue
		}

		slog.Debug("Processing unreleased changelog entry in directory",
			slog.String("root_dir", dir), slog.String("entry_name", entry.Name()))

		entry, err := parseChangelogEntry(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func ParseReleased() (*Changelog, error) {
	dir := filepath.Join(config.C.ChangelogDir, ReleasedDir)

	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	cl := new(Changelog)

	for _, entry := range dirs {
		if !entry.IsDir() {
			// we are only interested in directories
			continue
		} else if strings.HasPrefix(entry.Name(), ".") {
			// hidden files are ignored
			continue
		}

		slog.Debug("Processing entry in directory",
			slog.String("root_dir", dir), slog.String("entry_name", entry.Name()))

		rel, err := parseReleaseDirectory(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		cl.Releases = append(cl.Releases, *rel)
	}

	return cl, nil
}

func parseReleaseDirectory(path string) (*Release, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	rel := new(Release)

	for _, entry := range dirs {
		if entry.Name() == ReleaseInfoFileName || entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(path, entry.Name())

		slog.Debug("Parsing change entry", slog.String("entry_path", entryPath))

		change, err := parseChangelogEntry(entryPath)
		if err != nil {
			return nil, err
		}

		rel.Entries = append(rel.Entries, change)
	}

	releaseInfoPath := filepath.Join(path, ReleaseInfoFileName)

	raw, err := os.ReadFile(releaseInfoPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read release info file %q", releaseInfoPath)
	}

	err = toml.Unmarshal(raw, &rel.Info)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse release info file %q", releaseInfoPath)
	}

	// TODO: Validate whether it is a PreRelease and set Collapse accordingly

	return rel, nil
}

// parseChangelogEntry parses a changelog [Entry] at the given path.
func parseChangelogEntry(path string) (Entry, error) {
	slog.Debug("Parsing change entry", slog.String("entry_path", path))

	raw, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, errors.Wrapf(err, "unable to read %q entry", path)
	}

	var change Entry

	err = yaml.Unmarshal(raw, &change)
	if err != nil {
		return Entry{}, errors.Wrapf(err, "unable to parse %q entry", path)
	}

	change.EntryPath = typact.Some(path)

	return change, nil
}
