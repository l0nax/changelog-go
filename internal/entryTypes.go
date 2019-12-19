package internal

import (
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/errors"
)

var EntryT EntryTypes = EntryTypes{}

/// =========> Code and Functions <=========

// SEntryType is only used to search Entry Types
type SEntryType struct {
	ShortTypeName   string
	TypeDescription string // The Description is only to document the Type better
	TypeID          int    // This Type ID must be unique
}

// EntryTypes contains all available Entry Types
type EntryTypes struct {
	// internal struct Data
	entryTypes []*entry.ChangeEntry
}

// ListAvailableTypes returns a list of all available Entry Types
func (e *EntryTypes) ListAvailableTypes() []*entry.ChangeEntry {
	return e.entryTypes
}

func (e *EntryTypes) RegisterEntryType(et *entry.ChangeEntry) error {
	// check if a Entry Type already exists
	for _, v := range e.entryTypes {
		// Short Type Name
		if (*v).GetShortTypeName() == (*et).GetShortTypeName() {
			return errors.NewEntryTypeExistsError()
		}

		// Type Description
		if (*v).GetTypeDescription() == (*et).GetTypeDescription() {
			return errors.NewEntryTypeExistsError()
		}

		// Type ID
		if (*v).GetTypeID() == (*et).GetTypeID() {
			return errors.NewEntryTypeExistsError()
		}
	}

	// register Changelog Entry Type
	e.entryTypes = append(e.entryTypes, et)

	return nil
}

func (e *EntryTypes) DeRegisterEntryType(et *entry.ChangeEntry) error {
	// find Entry Type
	for i, _ := range e.entryTypes {
		// Short Type Name
		if ((*e.entryTypes[i]).GetShortTypeName() == (*et).GetShortTypeName()) &&
			((*e.entryTypes[i]).GetTypeDescription() == (*et).GetTypeDescription()) &&
			((*e.entryTypes[i]).GetTypeID() == (*et).GetTypeID()) {
			e.entryTypes = append(e.entryTypes[:i], e.entryTypes[i+1:]...)
			return nil
		}
	}

	// TODO: Return custom Error
	return nil
}

func (e *EntryTypes) SearchEntryType(se *SEntryType) (*entry.ChangeEntry, error) {
	// find Entry Type
	for i, _ := range e.entryTypes {
		log.Debugf("Searching (%# v) => Entry (%# v)\n",
			pretty.Formatter(*e.entryTypes[i]), pretty.Formatter(se))

		// short Type Name
		if ((*e.entryTypes[i]).GetShortTypeName() == se.ShortTypeName) &&
			((*e.entryTypes[i]).GetTypeDescription() == se.TypeDescription) &&
			((*e.entryTypes[i]).GetTypeID() == se.TypeID) {
			return e.entryTypes[i], nil
		}
	}

	// TODO: Add custom Error
	return nil, nil
}
