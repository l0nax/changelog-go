package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

// the v1 config file, comments and all
const legacyConfigFile = `# PLEASE DO NOT CHANGE THIS VALUE!
version: "1"

preRelease:
  detect: true
  deletePreRelease: true
  foldPreReleases: true

entry:
  author: true

changelog:
  entryPath: ".changes"
  changelog: "docs/CHANGELOG.md"
  customScheme: false
`

func writeFile(t *testing.T, path, body string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %q: %v", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}

// legacyProject lays out a v1 project and returns its root.
func legacyProject(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	writeFile(t, filepath.Join(root, LegacyConfigFileName), legacyConfigFile)

	writeFile(t, filepath.Join(root, ".changes", "released", "1.0.0", "ReleaseInfo"),
		"releasedate: \"2021-08-02\"\nisprerelease: false\n")
	writeFile(t, filepath.Join(root, ".changes", "released", "1.0.0", "develop-aaaa"),
		"author: emanuel\ntitle: Add the widget\ntype: 0\n")
	// type 6 is v1's "other", which has no counterpart in the stock v2 types
	writeFile(t, filepath.Join(root, ".changes", "released", "1.0.0", "develop-bbbb"),
		"author: emanuel\ntitle: Something uncategorized\ntype: 6\n")

	writeFile(t, filepath.Join(root, ".changes", "unreleased", "develop-cccc"),
		"author: emanuel\ntitle: Fix the sprocket\ntype: 1\n")

	return root
}

func TestMigrateCarriesOverV1Settings(t *testing.T) {
	root := legacyProject(t)

	if err := migrateProject(root, false); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg, err := config.Load(filepath.Join(root, ConfigFileName))
	if err != nil {
		t.Fatalf("the migrated config does not load: %v", err)
	}

	if cfg.ChangelogDir != ".changes" {
		t.Errorf("ChangelogDir = %q, want %q -- the v1 entryPath was dropped",
			cfg.ChangelogDir, ".changes")
	}

	if got := cfg.OutputPath.UnwrapOr(""); got != "docs/CHANGELOG.md" {
		t.Errorf("OutputPath = %q, want %q -- the v1 changelog path was dropped",
			got, "docs/CHANGELOG.md")
	}

	if !cfg.PreRelease.DeletePreRelease {
		t.Error("deletePreRelease = true was dropped")
	}

	if !cfg.PreRelease.FoldPreReleases {
		t.Error("foldPreReleases = true was dropped")
	}

	if !cfg.Version.IsValid() {
		t.Errorf("config version = %q, want a v2 one", cfg.Version)
	}

	// v1's numeric type 6 needs somewhere to land
	if _, ok := cfg.ResolveType(config.DefaultEntryOtherID); !ok {
		t.Error(`the "other" type was not added, so v1 type 6 entries have no group`)
	}

	if _, err := os.Stat(filepath.Join(root, LegacyConfigFileName)); !os.IsNotExist(err) {
		t.Error("the v1 config file was left behind")
	}
}

func TestMigrateTranslatesFragments(t *testing.T) {
	root := legacyProject(t)

	if err := migrateProject(root, false); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg, err := config.Load(filepath.Join(root, ConfigFileName))
	if err != nil {
		t.Fatalf("load migrated config: %v", err)
	}

	project := changelog.NewProject(cfg, root)

	released, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse migrated releases: %v", err)
	}

	if len(released.Releases) != 1 {
		t.Fatalf("expected one release, got %d", len(released.Releases))
	}

	release := released.Releases[0]
	if release.Info.Version != "1.0.0" {
		t.Errorf("version = %q, want %q", release.Info.Version, "1.0.0")
	}

	byTitle := map[string]string{}
	for _, entry := range release.Entries {
		byTitle[entry.Title] = entry.ChangeTypeID
	}

	if got := byTitle["Add the widget"]; got != config.DefaultEntryNewFeatureID {
		t.Errorf("v1 type 0 became %q, want %q", got, config.DefaultEntryNewFeatureID)
	}

	if got := byTitle["Something uncategorized"]; got != config.DefaultEntryOtherID {
		t.Errorf("v1 type 6 became %q, want %q", got, config.DefaultEntryOtherID)
	}

	unreleased, err := project.LoadUnreleasedEntries()
	if err != nil {
		t.Fatalf("parse migrated unreleased entries: %v", err)
	}

	if len(unreleased) != 1 || unreleased[0].ChangeTypeID != config.DefaultEntryBugFixID {
		t.Errorf("unreleased entries did not migrate: %+v", unreleased)
	}

	if got := unreleased[0].Author.UnwrapOr(""); got != "emanuel" {
		t.Errorf("author = %q, want %q", got, "emanuel")
	}

	// every entry must render, which is what the lost "other" type broke
	out, err := released.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	for _, want := range []string{"Add the widget", "Something uncategorized", "### Other "} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in the rendered changelog:\n%s", want, out)
		}
	}
}

func TestMigrateRefusesAnAlreadyMigratedProject(t *testing.T) {
	root := legacyProject(t)

	if err := migrateProject(root, false); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := migrateProject(root, false); err == nil {
		t.Fatal("expected the second migration to be refused")
	}
}

// A project that has no v1 config still migrates, it just keeps the defaults.
func TestMigrateWithoutALegacyConfig(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, ".changelogs", "unreleased", "develop-aaaa"),
		"title: Add the widget\ntype: 0\n")

	if err := migrateProject(root, false); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg, err := config.Load(filepath.Join(root, ConfigFileName))
	if err != nil {
		t.Fatalf("the migrated config does not load: %v", err)
	}

	if cfg.ChangelogDir != config.DefaultChangelogDir {
		t.Errorf("ChangelogDir = %q, want the default %q",
			cfg.ChangelogDir, config.DefaultChangelogDir)
	}
}
