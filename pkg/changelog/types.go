package changelog

import (
	"gitlab.com/l0nax/changelog-go/pkg/entry"
)

// ReleaseInfo holds all informations about the current release
type ReleaseInfo struct {
	Version      []string // Contains the string slice output of the Regex
	IsPreRelease bool     // True if the release is a pre-release
}

// Release contains the ReleaseInfo and all NEW Changelog entry types.
type Release struct {
	Info    *ReleaseInfo  // All Informations about the current Release
	Entries []entry.Entry // Holds all NEW changelog entry Types
}

// this Type is the default CHANGELOG.md Scheme
const changelogScheme = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).


{{ with .Releases }}
{{ range . }}
## {{ .Version }} ({{ .ReleaseDate }})

{{ with .Entries }}
{{ range . }}
### {{ .ShortTypeName }} ({{ .NumString }})

{{ with .Changes }}
{{ range . }}
- {{ .Title }}
{{ end }}
{{ end }}

{{/* two empty lines so thath the Markdown Parser will put here '<br/>' */}}


{{ end }}
{{/* A empty line to provide a Consistent Layout */}}

{{ end }}
{{ end }}
{{ end }}

`
