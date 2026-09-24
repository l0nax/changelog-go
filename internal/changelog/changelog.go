package changelog

import (
	"bytes"
	"cmp"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/config"
	"gitlab.com/l0nax/changelog-go/internal/tomlx"
)

const (
	// ReleasedDir is the name of the directory holding the released
	// versions.
	ReleasedDir = "released"
	// UnreleasedDir is the name of the directory holding the unreleased
	// entries.
	UnreleasedDir = "unreleased"
	// ReleaseInfoFileName is the name of the file holding a release's
	// metadata.
	ReleaseInfoFileName = "ReleaseInfo"
)

// Changelog is the set of releases rendered into the changelog file.
type Changelog struct {
	// VersionPrefix is prepended to every rendered version.
	VersionPrefix string

	Releases []Release
}

// Render returns the rendered changelog.
func (c *Changelog) Render() ([]byte, error) {
	if err := c.prepare(); err != nil {
		return nil, err
	}

	// TODO: Allow overriding the default
	tmpl, err := template.New("changelog-tmpl").
		Funcs(template.FuncMap{
			"formatTime": formatTime,
		}).
		Parse(defaultChangelogScheme)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	out.Grow(1024 * 1024) // 1 MB

	if err := tmpl.Execute(&out, c); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// SaveToFile renders the changelog and writes it to path.
func (c *Changelog) SaveToFile(path string) error {
	out, err := c.Render()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(out)

	return err
}

// prepare sorts the releases, fills in their display version and groups their
// entries by change type.
func (c *Changelog) prepare() error {
	if err := c.SortByRelease(); err != nil {
		return err
	}

	for i := range c.Releases {
		release := &c.Releases[i]

		release.DisplayVersion = ApplyVersionPrefix(c.VersionPrefix, release.Info.Version)
		release.GrouppedEntries = groupEntries(release.Entries)
	}

	return nil
}

// SortByRelease sorts the releases by their version in descending order.
func (c *Changelog) SortByRelease() error {
	// Parsed up front because a comparison function cannot report an error.
	versions := make(map[string]semver.Version, len(c.Releases))

	for _, release := range c.Releases {
		version, err := ParseVersion(release.Info.Version)
		if err != nil {
			return err
		}

		versions[release.Info.Version] = version
	}

	slices.SortStableFunc(c.Releases, func(a, b Release) int {
		return versions[b.Info.Version].Compare(versions[a.Info.Version])
	})

	return nil
}

// groupEntries groups entries by change type and sorts both the groups and the
// entries within them.
func groupEntries(entries []Entry) []GrouppedEntries {
	grouped := groupBy(entries, func(item Entry) string {
		return item.ChangeTypeID
	})

	groupedEntries := make([]GrouppedEntries, 0, len(grouped))

	for _, entries := range grouped {
		slices.SortStableFunc(entries, func(a, b Entry) int {
			return cmp.Compare(a.Title, b.Title)
		})

		groupedEntries = append(groupedEntries, GrouppedEntries{
			// every entry in the group carries the same resolved type
			ChangeType: entries[0].ChangeType,
			Entries:    entries,
		})
	}

	slices.SortStableFunc(groupedEntries, func(a, b GrouppedEntries) int {
		return cmp.Compare(a.ChangeType.ID, b.ChangeType.ID)
	})

	return groupedEntries
}

// LoadUnreleasedEntries returns all unreleased entries. A missing directory
// yields no entries.
func (p *Project) LoadUnreleasedEntries() ([]Entry, error) {
	dir := p.UnreleasedDir()

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "unable to read directory %q", dir)
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

		entry, err := p.parseChangelogEntry(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// ParseReleased returns the changelog of all released versions. A missing
// directory yields an empty changelog.
func (p *Project) ParseReleased() (*Changelog, error) {
	dir := p.ReleasedDir()

	cl := &Changelog{VersionPrefix: p.cfg.VersionPrefix}

	dirs, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return cl, nil
		}

		return nil, errors.Wrapf(err, "unable to read directory %q", dir)
	}

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

		rel, err := p.parseReleaseDirectory(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		cl.Releases = append(cl.Releases, *rel)
	}

	if err := p.markSupersededPreReleases(cl); err != nil {
		return nil, err
	}

	return cl, nil
}

// markSupersededPreReleases sets Collapse on every pre-release that a final
// release of the same base version supersedes.
func (p *Project) markSupersededPreReleases(cl *Changelog) error {
	if !p.cfg.PreRelease.FoldPreReleases {
		return nil
	}

	finalized := make(map[string]struct{}, len(cl.Releases))

	for _, release := range cl.Releases {
		if release.Info.IsPreRelease {
			continue
		}

		version, err := ParseVersion(release.Info.Version)
		if err != nil {
			return err
		}

		finalized[BaseVersion(version).String()] = struct{}{}
	}

	for i := range cl.Releases {
		release := &cl.Releases[i]
		if !release.Info.IsPreRelease {
			continue
		}

		version, err := ParseVersion(release.Info.Version)
		if err != nil {
			return err
		}

		_, release.Collapse = finalized[BaseVersion(version).String()]
	}

	return nil
}

func (p *Project) parseReleaseDirectory(path string) (*Release, error) {
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

		change, err := p.parseChangelogEntry(entryPath)
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

	if err := tomlx.Unmarshal(raw, &rel.Info); err != nil {
		return nil, errors.Wrapf(errors.New(tomlx.Explain(err)),
			"unable to parse release info file %q", releaseInfoPath)
	}

	return rel, nil
}

// parseChangelogEntry returns the changelog [Entry] stored at path.
func (p *Project) parseChangelogEntry(path string) (Entry, error) {
	slog.Debug("Parsing change entry", slog.String("entry_path", path))

	raw, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, errors.Wrapf(err, "unable to read %q entry", path)
	}

	var change Entry

	if err := tomlx.Unmarshal(raw, &change); err != nil {
		return Entry{}, errors.Wrapf(errors.New(tomlx.Explain(err)), "unable to parse %q entry", path)
	}

	change.ChangeType = p.changeTypeOf(change.ChangeTypeID, path)
	change.EntryPath = typact.Some(path)

	return change, nil
}

// changeTypeOf returns the configured change type named by id, or a
// synthesized one when the configuration no longer defines it.
func (p *Project) changeTypeOf(id, path string) config.ChangeType {
	if changeType, ok := p.cfg.ResolveType(id); ok {
		return changeType
	}

	slog.Warn("Changelog entry uses a change type that is no longer configured",
		slog.String("change_type_id", id), slog.String("entry_path", path),
		slog.String("hint", `set "hidden = true" on the type instead of removing it`))

	return synthesizeChangeType(id)
}
