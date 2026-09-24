package changelog

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// Entry is a single change entry.
type Entry struct {
	// ChangeTypeID is the ID identifying the type of the change.
	ChangeTypeID string `toml:"change_type_id"`

	// ChangeType is the resolved type named by ChangeTypeID, available to
	// the CHANGELOG template.
	ChangeType config.ChangeType `toml:"-"`

	// Title describes the change and is what the changelog lists.
	Title string `toml:"title"`

	// Author is the author of the change, if recorded.
	Author typact.Option[string] `toml:"author,omitzero"`

	// EntryPath is the path where the changelog entry is stored.
	EntryPath typact.Option[string] `toml:"-"`
}

// SaveToFile writes e to path.
func (e Entry) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	raw, err := toml.Marshal(e)
	if err != nil {
		return err
	}

	_, err = file.Write(raw)

	return err
}

// synthesizeChangeType returns a stand-in for a change type ID that the
// configuration no longer defines. It renders under its raw ID and carries the
// weakest version effect.
func synthesizeChangeType(id string) config.ChangeType {
	return config.ChangeType{
		ID:         id,
		Title:      id,
		GroupTitle: id,
		Effect:     config.VersionEffectPatch,
		Hidden:     true,
	}
}
