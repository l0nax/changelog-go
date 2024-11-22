package config

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"go.l0nax.org/typact"
)

// ChangeType is a single change type.
type ChangeType struct {
	// ID is the unique identifier of the type.
	ID string `koanf:"id"`
	// Title is the title of the type.
	Title string `koanf:"title"`
	// Description is the optional description of the type.
	Description typact.Option[string] `koanf:"description"`
	// GroupTitle is the title which is used in the CHANGELOG.md
	GroupTitle string `koanf:"group_title"`

	// Hidden hides the type in the selection input.
	// This can be used if the type has been deprecated and should
	// not be used anymore.
	Hidden bool `koanf:"hidden"`
}

// Config is the configuration structure.
type Config struct {
	// Version holds the config version.
	Version Version `koanf:"version"`

	// ChangelogDir is the relative path to the config file
	// where all the changelog files are stored.
	ChangelogDir string `koanf:"changelog_dir"`

	// OutputPath is the path to the file where the resulting
	// file should be stored.
	//
	// Defaults to "CHANGELOG.md".
	OutputPath typact.Option[string] `koanf:"output_path"`

	PreRelease struct {
		// detect pre-releases or not
		Detect           bool `koanf:"detect"`
		DeletePreRelease bool `koanf:"deletePreRelease"` // if true the pre-releases would be deleted on an non pre-release
		FoldPreReleases  bool `koanf:"foldPreReleases"`
	} `koanf:"preRelease"`

	// Entry configures a single changelog entry.
	Entry struct {
		// Types are the different change types which are vailable.
		Types []ChangeType `koanf:"types"`
	} `koanf:"entry"`
}

func (c Config) Validate() error {
	knownTypes := make(map[string]struct{}, len(c.Entry.Types))
	for _, typ := range c.Entry.Types {
		_, ok := knownTypes[typ.ID]
		if ok {
			return fmt.Errorf("changelog type with ID %q already defined", typ.ID)
		}

		knownTypes[typ.ID] = struct{}{}
	}

	return nil
}

// C is the loaded configuration.
var C Config

// Load loads the config into C.
func Load(path string) error {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return errors.Wrap(err, "unable to load file from path")
	}

	if err := k.Unmarshal("", &C); err != nil {
		return errors.Wrap(err, "unable to parse configuration file")
	}

	return C.Validate()
}

const DefaultOutputPath = "CHANGELOG.md"

// The default entry type IDs.
const (
	DefaultEntryNewFeatureID    = "new_feat"
	DefaultEntryBugFixID        = "bug_fix"
	DefaultEntryOtherID         = "other"
	DefaultEntryFeatureChangeID = "feat_change"
	DefaultEntryDeprecateID     = "deprecate"
	DefaultEntryRemovalID       = "rem_feat"
	DefaultEntrySecurityID      = "security"
)

const defaultConfig = `
changelog_dir = '.changelogs'
output_path = 'null'
version = '2'

[entry]
  [[entry.types]]
  id = 'new_feat'
  group_title = 'Added'
  title = 'New Feature'

  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'
  title = 'Bug Fixed'

  [[entry.types]]
  id = 'feat_change'
  group_title = 'Changed'
  title = 'Feature change'

  [[entry.types]]
  id = 'deprecate'
  group_title = 'Deprecated'
  title = 'Deprecation'

  [[entry.types]]
  id = 'rem_feat'
  group_title = 'Removed'
  title = 'Feature removal'

  [[entry.types]]
  id = 'security'
  group_title = 'Security'
  title = 'Security fix'

  [[entry.types]]
  id = 'other'
  group_title = 'Other'
  title = 'Other'

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
