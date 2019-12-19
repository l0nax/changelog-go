package tools

import (
	// "fmt"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

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
	typeSearch := &internal.SEntryType{
		TypeID: (*out).TypeID,
	}

	log.Debugf("new SEntryType: %# v\n", pretty.Formatter(typeSearch))
	log.Debugf("Out Struct: %# v\n", pretty.Formatter(out))

	// var err error
	// // search for the Change entry type and call the 'Add' method
	// (*out).Type, err = internal.EntryT.SearchEntryType(typeSearch)
	// if err != nil {
	//         return err
	// }
	// fmt.Printf("%#v\n", *out)
	// (*out.Type).AddEntry(out)

	return yaml.Unmarshal(in, out)
}
