package changelog

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
)

// ReleaseInfo holds all informations about the current release
type ReleaseInfo struct {
	Version      []string `yaml:"-"`                  // Contains the string slice output of the Regex
	IsPreRelease bool     `mapstructure:"prerelease"` // True if the release is a pre-release
	ReleaseDate  string   `mapstructure:"date"`       // Release Date in ISO foramt (eg. "2019-12-24")
}

// Release contains the ReleaseInfo and all Changelog entry types.
type Release struct {
	Info    *ReleaseInfo  // All Informations about the current Release
	Entries []entry.Entry // Holds all NEW changelog entry Types

	// == Fields used in the Changelog Template ==
	Releases []TplRelease // Holds all Releases
}

// TplRelease contains all Data of a Release
type TplRelease struct {
	Info     *ReleaseInfo // All Informations about the this Release
	Version  string       // Release Version
	Collapse bool         // This Field indicates if this Release should be collapsed
	Entries  []TplEntries // Contains a list with all available Change-Entries
}

// TplEntries is like a normal Entry, but it contains only Changes of a specific
// Change Type.
type TplEntries struct {
	ShortTypeName string         // Short Type Name, eg. "Added"
	NumString     string         // Contains the 'Number of Changes' String, eg: "1 change" OR "5 changes"
	Changes       []*entry.Entry // Contains the raw Change-Entry struct
}

// this Type is the default CHANGELOG.md Scheme
const changelogScheme = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

{{ with .Releases -}}
{{ range . }}
## {{ .Version }} ({{ .Info.ReleaseDate }})

{{- /* Collapse if PreRelease */ -}}
{{- if .Collapse -}}<details>{{- end -}}
{{- if .Collapse -}}<summary>This is a Pre-Release, Click to see details.</summary>{{- end }}

{{ range .Entries -}}
### {{ .ShortTypeName }} ({{ .NumString }})

{{- range .Changes }}
- {{ .ChangeTitle -}}
{{ end }}

{{ end }}
{{- if .Collapse -}}</details>{{- end -}}
{{- /* two empty lines so that the Markdown Parser will put here '<br/>' */ -}}


{{- /* A empty line to provide a Consistent Layout */ -}}

{{- end -}}
{{- end -}}`
