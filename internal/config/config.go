package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/fabmation-gmbh/toml"
	"go.l0nax.org/typact"
)

// VersionEffect is the effect a [ChangeType] has on the next version.
type VersionEffect uint8

func (v VersionEffect) IsValid() bool {
	return v == VersionEffectMajor || v == VersionEffectMinor || v == VersionEffectPatch
}

const (
	VersionEffectMajor VersionEffect = iota + 1
	VersionEffectMinor
	VersionEffectPatch
)

func (v VersionEffect) String() string {
	switch v {
	case VersionEffectMajor:
		return "major"
	case VersionEffectMinor:
		return "minor"
	case VersionEffectPatch:
		return "patch"
	default:
		panic(fmt.Sprintf("unknown version effect %d", v))
	}
}

func (v *VersionEffect) UnmarshalText(b []byte) error {
	switch string(b) {
	case "major":
		*v = VersionEffectMajor
	case "minor":
		*v = VersionEffectMinor
	case "patch":
		*v = VersionEffectPatch

	default:
		return errors.New("unknown version effect")
	}

	return nil
}

// ChangeType is a single change type.
type ChangeType struct {
	// ID is the unique identifier of the type.
	ID string `toml:"id"`
	// Title is the title of the type.
	Title string `toml:"title"`
	// Description is the optional description of the type.
	Description typact.Option[string] `toml:"description"`
	// GroupTitle is the title which is used in the CHANGELOG.md
	GroupTitle string `toml:"group_title"`

	// Effect is the version effect the type has on the next version.
	Effect VersionEffect `toml:"effect"`

	// Hidden hides the type in the selection input.
	// This can be used if the type has been deprecated and should
	// not be used anymore.
	Hidden bool `toml:"hidden"`
}

func (c ChangeType) Validate() error {
	if !c.Effect.IsValid() {
		return errors.New("unknown version effect")
	}

	return nil
}

// Config is the configuration structure.
type Config struct {
	// Version holds the config version.
	Version Version `toml:"version"`

	// ChangelogDir is the relative path to the config file
	// where all the changelog files are stored.
	ChangelogDir string `toml:"changelog_dir"`

	// OutputPath is the path to the file where the resulting
	// file should be stored.
	//
	// Defaults to "CHANGELOG.md".
	OutputPath typact.Option[string] `toml:"output_path"`

	PreRelease struct {
		// detect pre-releases or not
		Detect           bool `toml:"detect"`
		DeletePreRelease bool `toml:"deletePreRelease"` // if true the pre-releases would be deleted on an non pre-release
		FoldPreReleases  bool `toml:"foldPreReleases"`
	} `toml:"preRelease"`

	// Entry configures a single changelog entry.
	Entry struct {
		// Types are the different change types which are vailable.
		Types []ChangeType `toml:"types"`
	} `toml:"entry"`
}

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

// C is the loaded configuration.
var C Config

// Load loads the config into C.
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "unable to read config file")
	}

	if err := toml.Unmarshal(data, &C); err != nil {
		return errors.Wrap(err, "unable to parse config file")
	}

	return C.Validate()
}

const DefaultOutputPath = "CHANGELOG.md"

// The default entry type IDs.
const (
	DefaultEntryNewFeatureID    = "new_feat"
	DefaultEntryBugFixID        = "bug_fix"
	DefaultEntryFeatureChangeID = "feat_change"
	DefaultEntryDeprecateID     = "deprecate"
	DefaultEntryRemovalID       = "rem_feat"
	DefaultEntrySecurityID      = "security"
)

const defaultConfig = `
changelog_dir = '.changelogs'
output_path = 'CHANGELOG.md'
version = '2'

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

[preRelease]
deletePreRelease = false
detect = true
foldPreReleases = false
`

// CreateDefault creates a file at path with the contents
// of [Default].
func CreateDefault(path string, force bool) error {
	if force {
		_ = os.Remove(path)
	}

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write([]byte(defaultConfig))
	if err != nil {
		return err
	}

	return nil
}
