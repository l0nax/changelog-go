package changelog

import "sort"

// SortEntries does sort the `Entries` by their `ShortTypeName` field.
func (t *TplRelease) SortEntries() {
	sort.SliceStable(t.Entries, func(i, j int) bool {
		return t.Entries[i].ShortTypeName < t.Entries[j].ShortTypeName
	})
}
