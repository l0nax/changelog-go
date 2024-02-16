package changelog

const defaultOneRelease = `
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
{{- end -}}
`

// defaultChangelogScheme is the default CHANGELOG.md Scheme
const defaultChangelogScheme = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

` + defaultOneRelease

