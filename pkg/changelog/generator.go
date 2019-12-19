package changelog

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/tools"
	"strings"
)

// GenerateChangelog will generate a new CHANGELOG.md
func GenerateChangelog(r *Release) {
	// load all unreleased Files
	err := GetEntries(r)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Implement Feature (and Config variable) to allow the User
	// to regenerate the complete CHANGELOG.md.

	// initialize and add the current Relase to `r.Releases`
	r.Releases[0] = TplRelease{
		Version: strings.Join(r.Info.Version[:], ""),
		Colapse: false,
		Entries: []TplEntries{},
	}

	// get Entries by Change Type
	for _, cType := range internal.EntryT.ListAvailableTypes() {
		changeEntries := (*cType).GetListEntries()
		if len(changeEntries) == 0 {
			// skip this Change entry type because it has no changes
			continue
		}

		// change exist
		newEntry := TplEntries{
			ShortTypeName: (*cType).GetShortTypeName(),
		}

		// write 'n changes' instead of '1 change' if more than one
		if len(changeEntries) > 1 {
			newEntry.NumString = string(len(changeEntries)) + " changes"
		} else {
			newEntry.NumString = "1 change"
		}

		// add change-entry to the list of changes
		for iType, tChange := range changeEntries {
			newEntry.Changes[iType] = tChange
		}

		// append the new Change entry to the list of changes
		r.Releases[0].Entries = append(r.Releases[0].Entries, newEntry)
	}

	// debug: print the Releases
}
