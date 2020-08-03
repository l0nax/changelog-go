package entries

import (
	"github.com/kr/pretty"

	"gitlab.com/l0nax/changelog-go/pkg/entry"

	log "github.com/sirupsen/logrus"
)

/// ++++++] Added [++++++

//// register Type
//internal.EntryT.RegisterEntryType

//noinspection ALL
type Entry_Added struct {
	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Added) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Added",
		typeDescription: "New Feature",
		typeID:          0,
	}

	return &ret
}

func (e *Entry_Added) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Added) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Added) GetTypeID() int { return e.typeID }

func (e *Entry_Added) GetListEntries() []*entry.Entry {
	log.Debugf("Returning: %# v\n", pretty.Formatter(e.entries))
	return e.entries
}

func (e *Entry_Added) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Fixed [++++++
//noinspection ALL
type Entry_Fixed struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Fixed) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Fixed",
		typeDescription: "Bug Fix",
		typeID:          1,
	}

	return &ret
}

func (e *Entry_Fixed) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Fixed) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Fixed) GetTypeID() int { return e.typeID }

func (e *Entry_Fixed) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Fixed) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Changed [++++++
//noinspection ALL
type Entry_Changed struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Changed) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Changed",
		typeDescription: "Feature change",
		typeID:          2,
	}

	return &ret
}

func (e *Entry_Changed) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Changed) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Changed) GetTypeID() int { return e.typeID }

func (e *Entry_Changed) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Changed) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Deprecated [++++++
//noinspection ALL
type Entry_Deprecated struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Deprecated) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Deprecated",
		typeDescription: "New deprecation",
		typeID:          3,
	}

	return &ret
}

func (e *Entry_Deprecated) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Deprecated) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Deprecated) GetTypeID() int { return e.typeID }

func (e *Entry_Deprecated) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Deprecated) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Removed [++++++
//noinspection ALL
type Entry_Removed struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Removed) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Removed",
		typeDescription: "Feature removal",
		typeID:          4,
	}

	return &ret
}

func (e *Entry_Removed) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Removed) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Removed) GetTypeID() int { return e.typeID }

func (e *Entry_Removed) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Removed) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Security [++++++
//noinspection ALL
type Entry_Security struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Security) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Security",
		typeDescription: "Security fix",
		typeID:          5,
	}

	return &ret
}

func (e *Entry_Security) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Security) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Security) GetTypeID() int { return e.typeID }

func (e *Entry_Security) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Security) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}

/// ++++++] Other [++++++
//noinspection ALL
type Entry_Other struct {

	// internal Struct Data
	entries []*entry.Entry

	// internal Struct Information
	shortTypeName   string
	typeDescription string
	typeID          int
}

func (e *Entry_Other) NewChangeType() *entry.ChangeEntry {
	var ret entry.ChangeEntry = &Entry_Added{
		// -- internal struct data --
		shortTypeName:   "Other",
		typeDescription: "Other",
		typeID:          6,
	}

	return &ret
}

func (e *Entry_Other) GetShortTypeName() string { return e.shortTypeName }

func (e *Entry_Other) GetTypeDescription() string { return e.typeDescription }

func (e *Entry_Other) GetTypeID() int { return e.typeID }

func (e *Entry_Other) GetListEntries() []*entry.Entry { return e.entries }

func (e *Entry_Other) AddEntry(ent *entry.Entry) {
	if len(e.entries) == 0 || e.entries == nil {
		e.entries = []*entry.Entry{}
	}

	// add Entry to Entries Slice
	e.entries = append(e.entries, ent)
}
