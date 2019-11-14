// config holds all the required Types and Functions to work with the
// Changelog-Go Configuration File
package config

// V1 represents the Config Version 1.
// Versioning the Config Types is to provide Backward compatibility.
type V1 struct {
	PreRelrease struct {
		// detect pre-releases or not
		Detect           bool `yaml:"detect"`
		DeletePreRelease bool `yaml:"deletePreRelease"` // if true the pre-releases would be deleted on an non pre-release
		FoldPreReleases  bool `yaml:"foldPreReleases"`
	} `yaml:"preRelrease"`

	Entry struct {
		Author rune `yaml:"author"`
	} `yaml:"entry"`

	Changelog struct {
		// Changelog Path where changelog-go will save all the Entries.
		// Will be changed to a relative path after calling Check()
		EntryPath string `yaml:"entryPath"`
		Changelog string `yaml:"changelog"`
	} `yaml:"changelog"`
}

// ===== Internal Types =====
var v1ChangelogFolderStruct = []string{"released", "unreleased"}
