package entry

// ++++++] Added [++++++
//noinspection ALL
type Entry_Added struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Added) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Added) GetListEntries() []Entry { return e.entries }

func (e *Entry_Added) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}

// ++++++] Fixed [++++++
//noinspection ALL
type Entry_Fixed struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Fixed) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Fixed) GetListEntries() []Entry { return e.entries }

func (e *Entry_Fixed) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}

// ++++++] Changed [++++++
//noinspection ALL
type Entry_Changed struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Changed) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Changed) GetListEntries() []Entry { return e.entries }

func (e *Entry_Changed) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}

// ++++++] Deprecated [++++++
//noinspection ALL
type Entry_Deprecated struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Deprecated) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Deprecated) GetListEntries() []Entry { return e.entries }

func (e *Entry_Deprecated) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}

// ++++++] Removed [++++++
//noinspection ALL
type Entry_Removed struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Removed) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Removed) GetListEntries() []Entry { return e.entries }

func (e *Entry_Removed) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}

// ++++++] Security [++++++
//noinspection ALL
type Entry_Security struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Security) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Security) GetListEntries() []Entry { return e.entries }

func (e *Entry_Security) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}


// ++++++] Other [++++++
//noinspection ALL
type Entry_Other struct {

	// internal Struct Data
	entries			[]Entry

	// internal Struct Information
	shortTypeName 	string
	typeDescription string
	typeID			int
}

func (e *Entry_Other) newChangeType() *ChangeEntry {
	var ret ChangeEntry = &Entry_Added{
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

func (e *Entry_Other) GetListEntries() []Entry { return e.entries }

func (e *Entry_Other) AddEntry(entry Entry) {
	if len(e.entries) == 0 || e.entries == nil { e.entries = []Entry{} }

	// add Entry to Entries Slice
	e.entries = append(e.entries, entry)
}