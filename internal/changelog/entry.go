package changelog

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// Entry is a single change entry.
type Entry struct {
	// ChangeTypeID is the ID identifying the type of the change.
	ChangeTypeID string `koanf:"change_type_id"`

	// ChangeType holds the change type information.
	// It can be used in the CHANGELOG template.
	ChangeType config.ChangeType `koanf:"-"`

	// Title is the title describing the change, which will be used in
	// the resulting changelog.
	Title  string `koanf:"title"`
	// Author is the author, if defined
	Author typact.Option[string] `koanf:"author"`
}

// SaveToFile saves e to the given path.
func (e Entry) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	k := koanf.New(".")

	if err = k.Load(structs.Provider(e, "koanf"), nil); err != nil {
		return err
	}

	raw, err := k.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	_, err = file.Write(raw)
	if err != nil {
		return err
	}

	return nil
}
