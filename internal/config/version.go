package config

// Version represents the config version.
type Version string

func (v Version) IsValid() bool {
	return v == "2"
}
