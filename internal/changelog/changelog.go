package changelog

import (
	"bytes"
	"os"
	"text/template"
)

// Changelog represents the final CHANGELOG file.
type Changelog struct {
	Releases []Release
}

// SaveToFile generates a new changelog and saves it to path.
func (c Changelog) SaveToFile(path string) error {
	// TODO: Allow overriding the default
	tmpl, err := template.New("changelog-tmpl").Parse(defaultChangelogScheme)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	out.Grow(1024 * 1024) // 1 MB

	err = tmpl.Execute(&out, c)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(out.Bytes())
	if err != nil {
	  return err
	}

	return nil
}
