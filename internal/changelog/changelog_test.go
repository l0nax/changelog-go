package changelog

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// testProject loads the fixture project in testdata/project.
func testProject(t *testing.T, mutate func(*config.Config)) *Project {
	t.Helper()

	root := filepath.Join("testdata", "project")

	cfg, err := config.Load(filepath.Join(root, ".changelog-go.toml"))
	if err != nil {
		t.Fatalf("load fixture config: %v", err)
	}

	if mutate != nil {
		mutate(&cfg)
	}

	return NewProject(cfg, root)
}

// assertGolden compares got against testdata/<name>, rewriting it under -update.
func assertGolden(t *testing.T, name string, got []byte) {
	t.Helper()

	path := filepath.Join("testdata", name)

	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %q: %v", path, err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %q (run go test -update to create it): %v", path, err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("%s does not match the golden file.\n--- got ---\n%s\n--- want ---\n%s",
			name, got, want)
	}
}

func TestRenderGolden(t *testing.T) {
	project := testProject(t, nil)

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	got, err := cl.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	assertGolden(t, "CHANGELOG.golden.md", got)
}

func TestRenderIsStable(t *testing.T) {
	project := testProject(t, nil)

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	first, err := cl.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	// re-rendering the same changelog must not shuffle anything around
	second, err := cl.Render()
	if err != nil {
		t.Fatalf("re-render: %v", err)
	}

	if !bytes.Equal(first, second) {
		t.Error("rendering the same changelog twice produced different output")
	}
}

func TestVersionPrefixRendering(t *testing.T) {
	tests := map[string]struct {
		prefix string
		want   []string
	}{
		"v": {
			prefix: "v",
			// 1.4.0 is stored bare, v2.0.0-rc.1 with a prefix; both render prefixed
			want: []string{"## v1.0.0 ", "## v1.1.0 ", "## v1.1.0-rc.1 ", "## v2.0.0-rc.1 "},
		},
		"empty": {
			prefix: "",
			want:   []string{"## 1.0.0 ", "## 1.1.0 ", "## 1.1.0-rc.1 ", "## 2.0.0-rc.1 "},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			project := testProject(t, func(c *config.Config) {
				c.VersionPrefix = tc.prefix
			})

			cl, err := project.ParseReleased()
			if err != nil {
				t.Fatalf("parse released: %v", err)
			}

			got, err := cl.Render()
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			for _, want := range tc.want {
				if !strings.Contains(string(got), want) {
					t.Errorf("missing heading %q in:\n%s", want, got)
				}
			}
		})
	}
}

func TestSupersededPreReleaseIsCollapsed(t *testing.T) {
	project := testProject(t, nil)

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	collapsed := map[string]bool{}
	for _, release := range cl.Releases {
		collapsed[release.Info.Version] = release.Collapse
	}

	// 1.1.0 supersedes it
	if !collapsed["v1.1.0-rc.1"] {
		t.Error("superseded pre-release v1.1.0-rc.1 should be collapsed")
	}

	// nothing supersedes it yet, so hiding it would hide the newest changes
	if collapsed["v2.0.0-rc.1"] {
		t.Error("unsuperseded pre-release v2.0.0-rc.1 must stay expanded")
	}

	if collapsed["1.1.0"] {
		t.Error("a final release must never be collapsed")
	}
}

func TestFoldPreReleasesDisabled(t *testing.T) {
	project := testProject(t, func(c *config.Config) {
		c.PreRelease.FoldPreReleases = false
	})

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	for _, release := range cl.Releases {
		if release.Collapse {
			t.Errorf("release %q collapsed although fold_pre_releases is off", release.Info.Version)
		}
	}
}

// The nested lists only render as Markdown on GitLab and GitHub if a blank line
// follows the </summary> tag.
func TestCollapsedBlockHasBlankLineAfterSummary(t *testing.T) {
	project := testProject(t, nil)

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	got, err := cl.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	lines := strings.Split(string(got), "\n")

	var seen int

	for i, line := range lines {
		if !strings.Contains(line, "</summary>") {
			continue
		}

		seen++

		if i+1 >= len(lines) || strings.TrimSpace(lines[i+1]) != "" {
			t.Errorf("no blank line after </summary> at line %d:\n%s",
				i+1, strings.Join(lines[i:min(i+4, len(lines))], "\n"))
		}
	}

	if seen == 0 {
		t.Fatal("fixture rendered no collapsed pre-release")
	}
}

// Removing a type from the config must not brick the tool: past releases still
// parse and render, under the raw ID of the retired type.
func TestRetiredChangeTypeStillRenders(t *testing.T) {
	project := testProject(t, nil)

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	got, err := cl.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(got), "Something from a retired type") {
		t.Errorf("entry with a retired type was dropped from the output:\n%s", got)
	}

	if !strings.Contains(string(got), "### retired_type ") {
		t.Errorf("retired type did not render under its raw ID:\n%s", got)
	}
}

func TestSortByReleaseRejectsMalformedVersion(t *testing.T) {
	cl := &Changelog{
		Releases: []Release{
			{Info: ReleaseInfo{Version: "1.0.0"}},
			{Info: ReleaseInfo{Version: "not-a-version"}},
		},
	}

	// must report the problem rather than panicking or sorting it as 0.0.0
	if err := cl.SortByRelease(); err == nil {
		t.Fatal("expected an error for a malformed version")
	}
}

func TestSortByReleaseOrder(t *testing.T) {
	cl := &Changelog{
		Releases: []Release{
			{Info: ReleaseInfo{Version: "1.4.0"}},
			{Info: ReleaseInfo{Version: "2.0.0"}},
			{Info: ReleaseInfo{Version: "v2.0.0-rc.1"}},
			{Info: ReleaseInfo{Version: "v1.10.0"}},
			{Info: ReleaseInfo{Version: "2.0.0+build.5"}},
		},
	}

	if err := cl.SortByRelease(); err != nil {
		t.Fatalf("sort: %v", err)
	}

	got := make([]string, len(cl.Releases))
	for i, release := range cl.Releases {
		got[i] = release.Info.Version
	}

	// descending; build metadata is ignored for precedence, so 2.0.0 and
	// 2.0.0+build.5 keep their relative input order
	want := []string{"2.0.0", "2.0.0+build.5", "v2.0.0-rc.1", "v1.10.0", "1.4.0"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("wrong order:\ngot  %v\nwant %v", got, want)
		}
	}
}

func TestMissingDirectoriesAreTolerated(t *testing.T) {
	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	// a project that has been configured but never used
	project := NewProject(cfg, t.TempDir())

	cl, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released on an empty project: %v", err)
	}

	if len(cl.Releases) != 0 {
		t.Errorf("expected no releases, got %d", len(cl.Releases))
	}

	entries, err := project.LoadUnreleasedEntries()
	if err != nil {
		t.Fatalf("load unreleased on an empty project: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestPathsResolveAgainstTheConfigDir(t *testing.T) {
	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	project := NewProject(cfg, filepath.Join("some", "project"))

	if got, want := project.ChangelogDir(), filepath.Join("some", "project", ".changelogs"); got != want {
		t.Errorf("ChangelogDir = %q, want %q", got, want)
	}

	if got, want := project.OutputPath(), filepath.Join("some", "project", "CHANGELOG.md"); got != want {
		t.Errorf("OutputPath = %q, want %q", got, want)
	}
}

func TestAbsolutePathsAreLeftAlone(t *testing.T) {
	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	abs := filepath.Join(t.TempDir(), "elsewhere")
	cfg.ChangelogDir = abs

	project := NewProject(cfg, filepath.Join("some", "project"))

	if got := project.ChangelogDir(); got != abs {
		t.Errorf("ChangelogDir = %q, want %q", got, abs)
	}
}
