package changelog

import "go.l0nax.org/typact"

// Entry is a single change entry.
type Entry struct {
	// ChangeTypeID is the ID identifying the type of the change.
	ChangeTypeID string `koanf:"change_type_id"`
	// Title is the title describing the change, which will be used in
	// the resulting changelog.
	Title  string `koanf:"title"`
	// Author is the author, if defined
	Author typact.Option[string] `koanf:"author"`
}
