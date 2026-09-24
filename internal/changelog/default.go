package changelog

// numChangesDefine renders the "(3 changes)" suffix of a group heading.
const numChangesDefine = `{{ define "numChanges" }}
{{- $totalChanges := len . }}
{{- if gt $totalChanges 1 }}changes
{{- else }}change{{end -}}
{{- end -}}
`

// defaultOneRelease renders every release.
//
// The blank lines are load-bearing: two of them separate consecutive releases
// so that Markdown renders a break between them, the one after "</summary>" is
// what makes the lists inside a collapsed pre-release render at all on GitLab
// and GitHub, and the one before an entry body is what keeps the body in its
// list item instead of running into the title.
const defaultOneRelease = `
{{- range .Releases }}
{{ if $.ShowHeading }}## {{ .DisplayVersion }}{{ with .DisplayDate }} ({{ . }}){{ end }}

{{ end }}{{ if .Collapse }}<details>
<summary>This is a Pre-Release, Click to see details.</summary>

{{ end }}{{ range .GrouppedEntries }}### {{ .ChangeType.GroupTitle }} ({{ len .Entries }} {{ template "numChanges" .Entries }})
{{ range .Entries }}- {{ .Title }}
{{ if .Body.IsSome }}
{{ indentBody .Body.UnwrapOrZero }}
{{ end }}{{ end }}
{{ end }}{{ if .Collapse }}</details>

{{ end }}{{ end -}}
`

// defaultChangelogScheme is the default CHANGELOG.md Scheme
const defaultChangelogScheme = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
` + numChangesDefine + defaultOneRelease

// defaultReleaseScheme renders releases on their own, without the preamble.
//
// This is what "changelog show" prints, i.e. the body a forge publishes as the
// release notes.
const defaultReleaseScheme = numChangesDefine + defaultOneRelease
