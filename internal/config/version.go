package config

// Version represents the config version.
type Version string

// IsValid reports whether v is a config version this release understands.
func (v Version) IsValid() bool {
	return v == "2"
}
