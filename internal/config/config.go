package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"go.l0nax.org/typact"
)

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
}

// C is the loaded configuration.
var C Config

// Load loads the config into C.
func Load(path string) error {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return errors.Wrap(err, "unable to load file from path")
	}

	if err := k.Unmarshal("", &C); err != nil {
		return errors.Wrap(err, "unable to parse configuration file")
	}

	return nil
}
