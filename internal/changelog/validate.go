package changelog

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/v2/internal/tomlx"
)

// Problem is one thing wrong with the project.
type Problem struct {
	// Path is the offending file, relative to the project root.
	Path string
	// Message says what is wrong with it.
	Message string
}

func (p Problem) String() string {
	return p.Path + ": " + p.Message
}

// Validate parses the whole project and reports everything wrong with it.
//
// It keeps going after the first problem, because the point is to hand back a
// list a person can work through rather than one error at a time. Unlike the
// render path, an entry naming a change type that is not configured is a
// problem here rather than something to fall back from.
func (p *Project) Validate() ([]Problem, error) {
	var problems []Problem

	released, err := p.validateReleased()
	if err != nil {
		return nil, err
	}

	problems = append(problems, released...)

	unreleased, err := p.validateEntries(p.UnreleasedDir())
	if err != nil {
		return nil, err
	}

	problems = append(problems, unreleased...)

	sort.Slice(problems, func(i, j int) bool {
		if problems[i].Path != problems[j].Path {
			return problems[i].Path < problems[j].Path
		}

		return problems[i].Message < problems[j].Message
	})

	return problems, nil
}

func (p *Project) validateReleased() ([]Problem, error) {
	dirs, err := os.ReadDir(p.ReleasedDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "unable to read directory %q", p.ReleasedDir())
	}

	var problems []Problem

	for _, dir := range dirs {
		if !dir.IsDir() || strings.HasPrefix(dir.Name(), ".") {
			continue
		}

		path := filepath.Join(p.ReleasedDir(), dir.Name())

		problems = append(problems, p.validateReleaseInfo(path, dir.Name())...)

		entries, err := p.validateEntries(path)
		if err != nil {
			return nil, err
		}

		problems = append(problems, entries...)
	}

	return problems, nil
}

func (p *Project) validateReleaseInfo(dir, name string) []Problem {
	path := filepath.Join(dir, ReleaseInfoFileName)

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Problem{{
				Path:    p.displayPath(dir),
				Message: "release directory has no " + ReleaseInfoFileName + " file",
			}}
		}

		return []Problem{{Path: p.displayPath(path), Message: err.Error()}}
	}

	var info ReleaseInfo
	if err := tomlx.Unmarshal(raw, &info); err != nil {
		return []Problem{{Path: p.displayPath(path), Message: tomlx.Explain(err)}}
	}

	if info.Version == "" {
		return []Problem{{Path: p.displayPath(path), Message: "no version recorded"}}
	}

	version, err := ParseVersion(info.Version)
	if err != nil {
		return []Problem{{Path: p.displayPath(path), Message: err.Error()}}
	}

	var problems []Problem

	// The directory name is what "release" and "show" address a version by, so
	// a directory that disagrees with its own metadata is a trap.
	if dirVersion, err := ParseVersion(name); err != nil || !dirVersion.EQ(version) {
		problems = append(problems, Problem{
			Path: p.displayPath(path),
			Message: "version " + info.Version +
				" does not match the directory name " + name,
		})
	}

	if info.ReleaseDate.IsZero() {
		problems = append(problems, Problem{
			Path:    p.displayPath(path),
			Message: "no release date recorded",
		})
	}

	return problems
}

// validateEntries checks every changelog entry in dir.
func (p *Project) validateEntries(dir string) ([]Problem, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "unable to read directory %q", dir)
	}

	var problems []Problem

	for _, file := range files {
		if file.IsDir() || file.Name() == ReleaseInfoFileName || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		path := filepath.Join(dir, file.Name())

		raw, err := os.ReadFile(path)
		if err != nil {
			problems = append(problems, Problem{Path: p.displayPath(path), Message: err.Error()})

			continue
		}

		var entry Entry
		if err := tomlx.Unmarshal(raw, &entry); err != nil {
			problems = append(problems, Problem{Path: p.displayPath(path), Message: tomlx.Explain(err)})

			continue
		}

		if entry.ChangeTypeID == "" {
			problems = append(problems, Problem{Path: p.displayPath(path), Message: "no change_type_id set"})
		} else if _, ok := p.cfg.ResolveType(entry.ChangeTypeID); !ok {
			problems = append(problems, Problem{
				Path: p.displayPath(path),
				Message: "unknown change type " + entry.ChangeTypeID +
					`; configure it or mark the old type "hidden = true" instead of removing it`,
			})
		}

		if strings.TrimSpace(entry.Title) == "" {
			problems = append(problems, Problem{Path: p.displayPath(path), Message: "no title set"})
		}
	}

	return problems, nil
}

// displayPath renders path relative to the project root, so that the report
// reads the same wherever it was run from.
func (p *Project) displayPath(path string) string {
	rel, err := filepath.Rel(p.rootDir, path)
	if err != nil {
		return path
	}

	return filepath.ToSlash(rel)
}
