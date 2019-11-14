package entry

// ChangeEntry is the basic Changelog Entry Type. This Type is used by other
// Entry Types.
type ChangeEntry interface {
	// Returns a new (initialized) Changelog Entry Struct
	NewChangeType() *ChangeEntry

	// GetShortTypeName() returns the short Type Name
	// For Example: "Added"
	GetShortTypeName() string

	// GetTypeDescription() returns the Description of the Changelog Entry
	// Type which will be shown to the User.
	// For Example: "New Feature"
	GetTypeDescription() string

	// Returns the Type ID
	// !!Please choose the Type ID to be unique!!
	GetTypeID() int

	// Returns a List of all Entries
	// of this Changelog Type
	GetListEntries() []Entry

	// Add a new Entry to the List of Changes
	AddEntry(ent Entry)
}

type Entry struct {
	// ChangeTitle represents the Changelog Entry
	// Title/ Description
	ChangeTitle string `yaml:"title"`

	// Type represents Type of the Changelog
	Type   *ChangeEntry `yaml:"-"`
	TypeID int          `yaml:"type"`

	// Author Contains the Name (+ Mail) of
	// the Author from the Changelog Entry
	// If the Author Name (or mail) could not
	// be found from the Git Config, then this
	// String will be contain nothing
	Author string `yaml:"author,omitempty"`
}
