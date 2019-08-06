package internal

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/entry/entries"
)

// define base Pointer to Types
var basePointer = []entry.ChangeEntry{
	entries.Entry_Added{},
}

// RegisterAll() registers all (defined) Changelog Entry Types
func RegisterAll() {
	// --- Initialize base Pointer to Types --- //
	added := entries.Entry_Added{}

	// --- Register --- //
	_ = EntryT.RegisterEntryType(&added)
}