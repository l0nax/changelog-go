// Package config holds the project configuration.
package config

import (
	stderrors "errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/v2/internal/tomlx"
)

// VersionEffect is the effect a [ChangeType] has on the next version.
type VersionEffect uint8

// IsValid reports whether v is one of the known version effects.
func (v VersionEffect) IsValid() bool {
	return v == VersionEffectMajor || v == VersionEffectMinor || v == VersionEffectPatch
}

const (
	VersionEffectMajor VersionEffect = iota + 1
	VersionEffectMinor
	VersionEffectPatch
)

// String returns the name v has in the config file.
//
// NOTE: String panics if v is not valid.
func (v VersionEffect) String() string {
	switch v {
	case VersionEffectMajor:
		return "major"
	case VersionEffectMinor:
		return "minor"
	case VersionEffectPatch:
		return "patch"
	default:
		panic(fmt.Sprintf("unknown version effect %d", uint8(v)))
	}
}

// MarshalText returns the name v has in the config file.
func (v VersionEffect) MarshalText() ([]byte, error) {
	if !v.IsValid() {
		return nil, errors.Errorf("unknown version effect %d", uint8(v))
	}

	return []byte(v.String()), nil
}

// UnmarshalText sets v to the effect named by b.
func (v *VersionEffect) UnmarshalText(b []byte) error {
	switch string(b) {
	case "major":
		*v = VersionEffectMajor
	case "minor":
		*v = VersionEffectMinor
	case "patch":
		*v = VersionEffectPatch

	default:
		return errors.Errorf("unknown version effect %q", b)
	}

	return nil
}

// ChangeType is one kind of change a changelog entry can record.
type ChangeType struct {
	// ID is the unique identifier of the type.
	ID string `toml:"id"`
	// Title is the title of the type.
	Title string `toml:"title"`
	// Description is the optional description of the type.
	Description typact.Option[string] `toml:"description,omitzero"`
	// GroupTitle is the heading the entries of this type are grouped under.
	GroupTitle string `toml:"group_title"`

	// Effect is the version effect the type has on the next version.
	Effect VersionEffect `toml:"effect"`

	// Hidden hides the type in the selection input. A retired type that past
	// entries still refer to belongs here.
	Hidden bool `toml:"hidden"`
}

// Validate returns an error if c is missing an ID or carries an unknown
// version effect.
func (c ChangeType) Validate() error {
	if c.ID == "" {
		return errors.New("missing type ID")
	}

	if !c.Effect.IsValid() {
		return errors.New("unknown version effect")
	}

	return nil
}

// PreRelease configures how pre-releases are detected and rendered.
type PreRelease struct {
	// Detect enables deriving the pre-release flag from the released version.
	Detect bool `toml:"detect"`

	// DeletePreRelease removes the directories of superseded pre-releases
	// when a non pre-release of the same base version is released.
	//
	// It takes precedence over [PreRelease.FoldPreReleases].
	DeletePreRelease bool `toml:"delete_pre_release"`

	// FoldPreReleases renders a superseded pre-release inside a collapsed
	// "<details>" block.
	FoldPreReleases bool `toml:"fold_pre_releases"`
}

// Check configures "changelog check".
type Check struct {
	// BaseBranch is the ref a branch is compared against when looking for a
	// new changelog entry.
	//
	// Empty means the base branch is detected from the repository, which is
	// what keeps the shipped default config portable.
	BaseBranch string `toml:"base_branch"`
}

// Config is the project configuration.
type Config struct {
	// Version holds the config version.
	Version Version `toml:"version"`

	// ChangelogDir is the path to the directory holding all changelog
	// files, relative to the directory of the config file.
	ChangelogDir string `toml:"changelog_dir"`

	// OutputPath is the file the changelog is written to, relative to the
	// directory of the config file.
	//
	// Defaults to [DefaultOutputPath].
	OutputPath typact.Option[string] `toml:"output_path,omitzero"`

	// VersionPrefix is prepended to every version the tool renders.
	// An empty string renders bare versions.
	VersionPrefix string `toml:"version_prefix"`

	PreRelease PreRelease `toml:"pre_release"`

	Check Check `toml:"check"`

	// Entry configures a single changelog entry.
	Entry struct {
		// Types are the available change types.
		Types []ChangeType `toml:"types"`
	} `toml:"entry"`
}

// Validate returns an error if c defines a change type ID twice or holds an
// invalid change type.
func (c Config) Validate() error {
	knownTypes := make(map[string]struct{}, len(c.Entry.Types))

	for _, typ := range c.Entry.Types {
		_, ok := knownTypes[typ.ID]
		if ok {
			return fmt.Errorf("changelog type with ID %q already defined", typ.ID)
		}

		if err := typ.Validate(); err != nil {
			return errors.Wrapf(err, "error while validating type %q", typ.ID)
		}

		knownTypes[typ.ID] = struct{}{}
	}

	return nil
}

// ResolveType returns the configured change type with the given ID.
func (c Config) ResolveType(id string) (ChangeType, bool) {
	for _, typ := range c.Entry.Types {
		if typ.ID == id {
			return typ, true
		}
	}

	return ChangeType{}, false
}

// Default returns the configuration with every defaultable key filled in.
//
// [Load] decodes the config file on top of this value. A key absent from the
// file keeps its default; a key present in the file wins, including when it is
// set to the zero value.
func Default() Config {
	var c Config

	c.ChangelogDir = DefaultChangelogDir
	c.VersionPrefix = DefaultVersionPrefix
	c.PreRelease.Detect = true
	c.PreRelease.FoldPreReleases = true

	return c
}

// DefaultProjectConfig returns the parsed configuration a new project starts
// with.
func DefaultProjectConfig() (Config, error) {
	cfg := Default()

	if err := tomlx.Unmarshal([]byte(defaultConfig), &cfg); err != nil {
		return Config{}, errors.Wrap(err, "unable to parse the built-in default config")
	}

	return cfg, cfg.Validate()
}

// Marshal renders c as a config file.
func Marshal(c Config) ([]byte, error) {
	return toml.Marshal(c)
}

// legacyKeys maps the camelCase keys shipped by v2.0.0-rc.1 to their
// snake_case replacement.
var legacyKeys = map[string]string{
	"preRelease":       "pre_release",
	"deletePreRelease": "delete_pre_release",
	"foldPreReleases":  "fold_pre_releases",
}

// Load reads and validates the config file at path.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "unable to read config file")
	}

	cfg := Default()

	if err := tomlx.Unmarshal(data, &cfg); err != nil {
		return Config{}, errors.Wrapf(explainConfigError(err), "unable to parse config file %q", path)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, errors.Wrapf(err, "invalid config file %q", path)
	}

	return cfg, nil
}

// explainConfigError returns err with the replacement named for every rejected
// key that has a known legacy spelling.
func explainConfigError(err error) error {
	hints := make([]string, 0, len(legacyKeys))

	for _, key := range tomlx.RejectedKeys(err) {
		replacement, ok := legacyKeys[key]
		if !ok {
			continue
		}

		hints = append(hints, fmt.Sprintf("  %q has been renamed to %q", key, replacement))
	}

	if len(hints) == 0 {
		return stderrors.New(tomlx.Explain(err))
	}

	return fmt.Errorf("%s\n\nthe following keys changed in v2.0.0:\n%s",
		tomlx.Explain(err), strings.Join(hints, "\n"))
}

const (
	// DefaultOutputPath is the file the changelog is rendered to.
	DefaultOutputPath = "CHANGELOG.md"
	// DefaultChangelogDir is the directory holding the changelog files.
	DefaultChangelogDir = ".changelogs"
	// DefaultVersionPrefix is prepended to every rendered version.
	DefaultVersionPrefix = "v"
)

// The default entry type IDs.
const (
	DefaultEntryNewFeatureID    = "new_feat"
	DefaultEntryBugFixID        = "bug_fix"
	DefaultEntryFeatureChangeID = "feat_change"
	DefaultEntryDeprecateID     = "deprecate"
	DefaultEntryRemovalID       = "rem_feat"
	DefaultEntrySecurityID      = "security"
	DefaultEntryOtherID         = "other"
)

const defaultConfig = `
changelog_dir = '.changelogs'
output_path = 'CHANGELOG.md'
version = '2'
version_prefix = 'v'

[entry]
  [[entry.types]]
  id = 'new_feat'
  group_title = 'Added'
  title = 'New Feature'
  effect = 'minor'

  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'
  title = 'Bug Fixed'
  effect = 'patch'

  [[entry.types]]
  id = 'feat_change'
  group_title = 'Changed'
  title = 'Feature change'
  effect = 'minor'

  [[entry.types]]
  id = 'deprecate'
  group_title = 'Deprecated'
  title = 'Deprecation'
  effect = 'minor'

  [[entry.types]]
  id = 'rem_feat'
  group_title = 'Removed'
  title = 'Feature removal'
  effect = 'major'

  [[entry.types]]
  id = 'security'
  group_title = 'Security'
  title = 'Security fix'
  effect = 'patch'

[pre_release]
delete_pre_release = false
detect = true
fold_pre_releases = true

[check]
# The ref "changelog check" compares a branch against. Leave it empty to detect
# the repository's default branch.
base_branch = ''
`

// DefaultConfigFile returns the contents written by "changelog init".
func DefaultConfigFile() string {
	return defaultConfig
}

// CreateDefault writes the default config file to path. force replaces an
// existing file.
func CreateDefault(path string, force bool) error {
	if force {
		_ = os.Remove(path)
	}

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.WriteString(defaultConfig)
	if err != nil {
		return err
	}

	return nil
}
