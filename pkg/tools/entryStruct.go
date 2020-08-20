package tools

import (
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"

	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
)

func YAMLMarshal(in entry.Entry) ([]byte, error) {
	// we need this Function because we have to Convert the 'type' Field
	in.TypeID = (*in.Type).GetTypeID()
	return yaml.Marshal(in)
}

func YAMLUnmarshal(in []byte, out *entry.Entry) error {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		return err
	}

	// after unmarshaling the input data into the struct, we can now search
	// for the type of the ChangeEntry
	typeSearch := &internal.SEntryType{
		TypeID: (*out).TypeID,
	}

	// search for the Change entry type and call the 'Add' method
	(*out).Type, err = internal.EntryT.SearchEntryType(typeSearch)
	if err != nil {
		return err
	}

	// add the unmarshalled type to our internal change list to keep track
	// of him
	(*(*out).Type).AddEntry(out)

	log.Debugf("Add new Changelog Entry to list of changes: %# v\n", pretty.Formatter(out))

	return nil
}
