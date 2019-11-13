package internal

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/entry/entries"
)

//// define base Pointer to Types
//var basePointer []entry.ChangeEntry

// RegisterAll() registers all (defined) Changelog Entry Types
func RegisterAll() {
	// --- Initialize base Pointer to Types --- //
	var added entry.ChangeEntry = &entries.Entry_Added{}
	var changed entry.ChangeEntry = &entries.Entry_Changed{}
	var deprecated entry.ChangeEntry = &entries.Entry_Deprecated{}
	var removed entry.ChangeEntry = &entries.Entry_Removed{}
	var security entry.ChangeEntry = &entries.Entry_Security{}
	var other entry.ChangeEntry = &entries.Entry_Other{}

	// --- Register --- //
	_ = EntryT.RegisterEntryType(added.NewChangeType())
	_ = EntryT.RegisterEntryType(changed.NewChangeType())
	_ = EntryT.RegisterEntryType(deprecated.NewChangeType())
	_ = EntryT.RegisterEntryType(removed.NewChangeType())
	_ = EntryT.RegisterEntryType(security.NewChangeType())
	_ = EntryT.RegisterEntryType(other.NewChangeType())
}
