package changelog

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

const (
	// ReleasedDir is the name of the directory holding the released
	// versions.
	ReleasedDir = "released"
	// ReleaseInfoFileName is the name of the "ReleaseInfo" filename.
	ReleaseInfoFileName = "ReleaseInfo"
)

// Changelog represents the final CHANGELOG file.
type Changelog struct {
	Releases []Release
}

// SaveToFile generates a new changelog and saves it to path.
func (c Changelog) SaveToFile(path string) error {
	// TODO: Allow overriding the default
	tmpl, err := template.New("changelog-tmpl").Parse(defaultChangelogScheme)
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

	spew.Dump(cl)

	panic("TODO")
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

		k := koanf.New(".")
		entryPath := filepath.Join(path, entry.Name())

		slog.Debug("Parsing change entry", slog.String("entry_path", entryPath))

		err = k.Load(file.Provider(entryPath), yaml.Parser())
		if err != nil {
			return nil, err
		}

		var change Entry

		if err = k.Unmarshal("", &change); err != nil {
			return nil, err
		}

		// load all relevant information
		if err = change.LoadChangeType(); err != nil {
			return nil, err
		}

		rel.Entries = append(rel.Entries, change)
	}

	k := koanf.New(".")
	if err = k.Load(file.Provider(filepath.Join(path, ReleaseInfoFileName)), yaml.Parser()); err != nil {
		return nil, err
	}

	if err = k.Unmarshal("", &rel.Info); err != nil {
		return nil, err
	}

	// TODO: Validate whether it is a PreRelease and set Collapse accordingly

	return rel, nil
}
