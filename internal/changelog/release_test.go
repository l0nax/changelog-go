package changelog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

// pendingEntries writes count unreleased entries and returns them.
func pendingEntries(t *testing.T, project *Project, count int) []Entry {
	t.Helper()

	if err := os.MkdirAll(project.UnreleasedDir(), 0o755); err != nil {
		t.Fatalf("mkdir unreleased: %v", err)
	}

	entries := make([]Entry, 0, count)

	for i := range count {
		entry := Entry{
			ChangeTypeID: config.DefaultEntryBugFixID,
			Title:        "Fix something",
		}

		path := filepath.Join(project.UnreleasedDir(), "entry-"+string(rune('a'+i)))
		if err := entry.SaveToFile(path); err != nil {
			t.Fatalf("save entry: %v", err)
		}

		entry.EntryPath = typact.Some(path)
		entries = append(entries, entry)
	}

	return entries
}

func newRelease(version string, isPre bool, entries []Entry) Release {
	return Release{
		Info: ReleaseInfo{
			Version:      version,
			ReleaseDate:  time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC),
			IsPreRelease: isPre,
		},
		Entries: entries,
	}
}

func TestCreateReleaseRejectsExistingDirectory(t *testing.T) {
	project := tempProject(t, map[string]bool{"1.0.0": false}, nil)

	err := project.CreateRelease(newRelease("1.0.0", false, nil))
	if err == nil {
		t.Fatal("expected an error when the release directory already exists")
	}
}

func TestCreateReleaseMovesFragments(t *testing.T) {
	project := tempProject(t, nil, nil)
	entries := pendingEntries(t, project, 2)

	if err := project.CreateRelease(newRelease("1.0.0", false, entries)); err != nil {
		t.Fatalf("CreateRelease: %v", err)
	}

	// a final release consumes the fragments
	left, err := project.LoadUnreleasedEntries()
	if err != nil {
		t.Fatalf("load unreleased: %v", err)
	}

	if len(left) != 0 {
		t.Errorf("expected the unreleased entries to be consumed, %d left", len(left))
	}

	released, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	if len(released.Releases) != 1 || len(released.Releases[0].Entries) != 2 {
		t.Fatalf("expected one release holding two entries, got %+v", released.Releases)
	}
}

func TestCreateReleaseKeepsFragmentsForAPreRelease(t *testing.T) {
	project := tempProject(t, nil, nil)
	entries := pendingEntries(t, project, 2)

	if err := project.CreateRelease(newRelease("1.0.0-rc.1", true, entries)); err != nil {
		t.Fatalf("CreateRelease: %v", err)
	}

	// a pre-release keeps them, which is what lets the final release of the
	// same base version carry the same entries
	left, err := project.LoadUnreleasedEntries()
	if err != nil {
		t.Fatalf("load unreleased: %v", err)
	}

	if len(left) != 2 {
		t.Errorf("expected the two unreleased entries to survive, got %d", len(left))
	}
}

func TestRemoveSupersededPreReleases(t *testing.T) {
	project := tempProject(t, map[string]bool{
		"1.0.0":       false,
		"v1.1.0-rc.1": true,
		"v1.1.0-rc.2": true,
		"v1.2.0-rc.1": true,
	}, nil)

	if err := project.RemoveSupersededPreReleases("1.1.0"); err != nil {
		t.Fatalf("RemoveSupersededPreReleases: %v", err)
	}

	for _, version := range []string{"v1.1.0-rc.1", "v1.1.0-rc.2"} {
		if _, err := os.Stat(filepath.Join(project.ReleasedDir(), version)); !os.IsNotExist(err) {
			t.Errorf("superseded pre-release %q was not removed", version)
		}
	}

	// a pre-release of a different base version is none of its business
	for _, version := range []string{"1.0.0", "v1.2.0-rc.1"} {
		if _, err := os.Stat(filepath.Join(project.ReleasedDir(), version)); err != nil {
			t.Errorf("release %q should have been left alone: %v", version, err)
		}
	}
}

func TestReleaseInfoIsStoredAtSecondPrecision(t *testing.T) {
	project := tempProject(t, nil, nil)

	// what cmd/release.go writes
	stamp := time.Now().Truncate(time.Second)

	release := newRelease("1.0.0", false, nil)
	release.Info.ReleaseDate = stamp

	if err := project.CreateRelease(release); err != nil {
		t.Fatalf("CreateRelease: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(project.ReleasedDir(), "1.0.0", ReleaseInfoFileName))
	if err != nil {
		t.Fatalf("read ReleaseInfo: %v", err)
	}

	if stamp.Nanosecond() != 0 {
		t.Fatalf("truncation did not take: %v", stamp)
	}

	released, err := project.ParseReleased()
	if err != nil {
		t.Fatalf("parse released: %v", err)
	}

	if !released.Releases[0].Info.ReleaseDate.Equal(stamp) {
		t.Errorf("date round-trip changed the value: got %v want %v\n%s",
			released.Releases[0].Info.ReleaseDate, stamp, raw)
	}
}

func TestVersionEffectOfARelease(t *testing.T) {
	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	typeOf := func(id string) config.ChangeType {
		ct, ok := cfg.ResolveType(id)
		if !ok {
			t.Fatalf("type %q missing from the default config", id)
		}

		return ct
	}

	release := Release{
		Entries: []Entry{
			{ChangeType: typeOf(config.DefaultEntryBugFixID)},
			{ChangeType: typeOf(config.DefaultEntryRemovalID)},
			{ChangeType: typeOf(config.DefaultEntryNewFeatureID)},
		},
	}

	// the strongest effect wins
	if got := release.VersionEffect(); got != config.VersionEffectMajor {
		t.Errorf("VersionEffect = %v, want major", got)
	}

	// an entry whose type never resolved must not drag the effect to zero
	release.Entries = append(release.Entries, Entry{})

	if got := release.VersionEffect(); got != config.VersionEffectMajor {
		t.Errorf("VersionEffect with an unresolved entry = %v, want major", got)
	}
}

func TestParseRejectsUnknownKeysInAFragment(t *testing.T) {
	project := tempProject(t, nil, nil)

	if err := os.MkdirAll(project.UnreleasedDir(), 0o755); err != nil {
		t.Fatalf("mkdir unreleased: %v", err)
	}

	// "change_type" is a plausible typo for "change_type_id"
	body := "change_type = 'bug_fix'\ntitle = 'Fix something'\n"
	if err := os.WriteFile(filepath.Join(project.UnreleasedDir(), "typo"), []byte(body), 0o644); err != nil {
		t.Fatalf("write entry: %v", err)
	}

	_, err := project.LoadUnreleasedEntries()
	if err == nil {
		t.Fatal("expected a typo in a fragment to be reported")
	}
}
