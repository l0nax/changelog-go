package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
)

// Config holds the config structure.
type Config struct {
	PreRelease struct {
		// detect pre-releases or not
		Detect           bool `koanf:"detect"`
		DeletePreRelease bool `koanf:"deletePreRelease"` // if true the pre-releases would be deleted on an non pre-release
		FoldPreReleases  bool `koanf:"foldPreReleases"`
	} `koanf:"preRelease"`

	Entry struct {
		Author bool `koanf:"author"`
	} `koanf:"entry"`

	Changelog struct {
		// Changelog Path where changelog-go will save all the Entries.
		// Will be changed to a relative path after calling Check()
		EntryPath    string `koanf:"entryPath"`
		Changelog    string `koanf:"changelog"`
		CustomScheme bool   `koanf:"customScheme"`
	} `koanf:"changelog"`
}

// C holds the loaded config.
var C Config

// Load will load the config at the provided path.
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
