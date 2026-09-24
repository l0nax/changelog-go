// Package tomlx wraps the TOML decoder with the settings used across the
// project.
package tomlx

import (
	"bytes"
	"errors"

	"github.com/pelletier/go-toml/v2"
)

// Unmarshal decodes data into v and returns an error for any key that has no
// counterpart in v.
func Unmarshal(data []byte, v any) error {
	dec := toml.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}

// Explain returns the message of err, including the offending lines when err
// is a [toml.StrictMissingError].
func Explain(err error) string {
	var strict *toml.StrictMissingError
	if errors.As(err, &strict) {
		return strict.String()
	}

	return err.Error()
}

// RejectedKeys returns the keys that err reports as unknown, or nil if err is
// not a strictness error.
func RejectedKeys(err error) []string {
	var strict *toml.StrictMissingError
	if !errors.As(err, &strict) {
		return nil
	}

	keys := make([]string, 0, len(strict.Errors))

	for i := range strict.Errors {
		key := strict.Errors[i].Key()
		if len(key) == 0 {
			continue
		}

		keys = append(keys, key[len(key)-1])
	}

	return keys
}
