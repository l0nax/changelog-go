package changelog

// defaultOneRelease is the template rendering every release.
//
// The blank lines are load-bearing: two separate consecutive releases, and the
// one after "</summary>" is required for the lists inside a collapsed
// pre-release to render on GitLab and GitHub.
const defaultOneRelease = `
{{- range .Releases }}
## {{ .DisplayVersion }} ({{ formatTime .Info.ReleaseDate }})

{{ if .Collapse }}<details>
<summary>This is a Pre-Release, Click to see details.</summary>

{{ end }}{{ range .GrouppedEntries }}### {{ .ChangeType.GroupTitle }} ({{ len .Entries }} {{ template "numChanges" .Entries }})
{{ range .Entries }}- {{ .Title }}
{{ end }}
{{ end }}{{ if .Collapse }}</details>

{{ end }}{{ end -}}
`

// defaultChangelogScheme is the template of the whole changelog file.
const defaultChangelogScheme = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
{{ define "numChanges" }}
{{- $totalChanges := len . }}
{{- if gt $totalChanges 1 }}changes
{{- else }}change{{end -}}
{{- end -}}
` + defaultOneRelease
