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
	var added entry.ChangeEntry = entries.NewBasicEntry("Added", "New Feature", 0)
	var fixed entry.ChangeEntry = entries.NewBasicEntry("Fixed", "Bug Fixed", 1)
	var changed entry.ChangeEntry = entries.NewBasicEntry("Changed", "Feature change", 2)
	var deprecated entry.ChangeEntry = entries.NewBasicEntry("Deprecated", "New deprecation", 3)
	var removed entry.ChangeEntry = entries.NewBasicEntry("Removed", "Feature removal", 4)
	var security entry.ChangeEntry = entries.NewBasicEntry("Security", "Security fix", 5)
	var other entry.ChangeEntry = entries.NewBasicEntry("Other", "Other", 6)

	// --- Register --- //
	_ = EntryT.RegisterEntryType(added.NewChangeType())
	_ = EntryT.RegisterEntryType(fixed.NewChangeType())
	_ = EntryT.RegisterEntryType(changed.NewChangeType())
	_ = EntryT.RegisterEntryType(deprecated.NewChangeType())
	_ = EntryT.RegisterEntryType(removed.NewChangeType())
	_ = EntryT.RegisterEntryType(security.NewChangeType())
	_ = EntryT.RegisterEntryType(other.NewChangeType())
}
