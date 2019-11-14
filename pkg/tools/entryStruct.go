package tools

import (
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gopkg.in/yaml.v2"
)

func YAMLMarshal(in entry.Entry) ([]byte, error) {
	// we need this Function because we have to Convert the 'type' Field
	in.TypeID = (*in.Type).GetTypeID()
	return yaml.Marshal(in)
}

func YAMLUnmarshal(in []byte, out *entry.Entry) error {
	typeSearch := internal.SEntryType{
		TypeID: (*out).TypeID,
	}

	var err error
	// simply search for the Change Entry Type
	(*out).Type, err = internal.EntryT.SearchEntryType(&typeSearch)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(in, out)
}
