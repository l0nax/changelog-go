package entries

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
)

type BasicEntry struct {
	entries []*entry.Entry

	shortTypeName string
	typeDesc      string
	typeID        int
}

// NewBasicEntry returns a basic initialized entry based on the input parameters.
func NewBasicEntry(short, desc string, id int) *BasicEntry {
	return &BasicEntry{
		shortTypeName: short,
		typeDesc:      desc,
		typeID:        id,
	}
}

func (b *BasicEntry) AddEntry(ent *entry.Entry)      { b.entries = append(b.entries, ent) }
func (b *BasicEntry) GetListEntries() []*entry.Entry { return b.entries }
func (b *BasicEntry) GetTypeDescription() string     { return b.typeDesc }
func (b *BasicEntry) GetTypeID() int                 { return b.typeID }
func (b *BasicEntry) GetShortTypeName() string       { return b.shortTypeName }
func (b *BasicEntry) NewChangeType() *entry.ChangeEntry {
	// TODO(l0nax): Why TF are we returning a pointer to a slice to a pointer of a struct???
	var ret entry.ChangeEntry = &BasicEntry{
		shortTypeName: b.shortTypeName,
		typeDesc:      b.typeDesc,
		typeID:        b.typeID,
	}

	return &ret
}
