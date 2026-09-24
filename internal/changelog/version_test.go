package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

// tempProject builds a project in a scratch directory.
//
// released maps a version to its pre-release flag, unreleased lists the change
// type IDs of the pending entries.
func tempProject(t *testing.T, released map[string]bool, unreleased []string) *Project {
	t.Helper()

	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	root := t.TempDir()
	project := NewProject(cfg, root)

	for version, isPre := range released {
		dir := filepath.Join(project.ReleasedDir(), version)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %q: %v", dir, err)
		}

		info := fmt.Sprintf("version = '%s'\ndate = 2024-01-01T00:00:00Z\npre_release = %t\n",
			version, isPre)
		if err := os.WriteFile(filepath.Join(dir, ReleaseInfoFileName), []byte(info), 0o644); err != nil {
			t.Fatalf("write ReleaseInfo: %v", err)
		}
	}

	if len(unreleased) > 0 {
		if err := os.MkdirAll(project.UnreleasedDir(), 0o755); err != nil {
			t.Fatalf("mkdir unreleased: %v", err)
		}
	}

	for i, typeID := range unreleased {
		body := fmt.Sprintf("change_type_id = '%s'\ntitle = 'entry %d'\n", typeID, i)

		path := filepath.Join(project.UnreleasedDir(), fmt.Sprintf("entry-%d", i))
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write entry: %v", err)
		}
	}

	return project
}

func TestNextVersion(t *testing.T) {
	tests := map[string]struct {
		released   map[string]bool
		unreleased []string
		want       string
		wantEffect config.VersionEffect
		wantErr    bool
	}{
		"no releases yet": {
			unreleased: []string{config.DefaultEntryBugFixID},
			want:       "0.1.0",
			wantEffect: config.VersionEffectMinor,
		},
		"major entry bumps the major": {
			released:   map[string]bool{"1.4.0": false},
			unreleased: []string{config.DefaultEntryBugFixID, config.DefaultEntryRemovalID},
			want:       "2.0.0",
			wantEffect: config.VersionEffectMajor,
		},
		"minor entry bumps the minor": {
			released:   map[string]bool{"1.4.0": false},
			unreleased: []string{config.DefaultEntryBugFixID, config.DefaultEntryNewFeatureID},
			want:       "1.5.0",
			wantEffect: config.VersionEffectMinor,
		},
		"only fixes bump the patch": {
			released:   map[string]bool{"1.4.0": false},
			unreleased: []string{config.DefaultEntryBugFixID},
			want:       "1.4.1",
			wantEffect: config.VersionEffectPatch,
		},
		// decision A: the pre-release holds the same entries, so finalizing it
		// is the next version -- bumping would skip 2.0.0 entirely
		"pre-release finalizes to its base version": {
			released:   map[string]bool{"1.4.0": false, "v2.0.0-rc.1": true},
			unreleased: []string{config.DefaultEntryNewFeatureID},
			want:       "2.0.0",
			wantEffect: config.VersionEffectMinor,
		},
		"v-prefixed latest does not panic": {
			released:   map[string]bool{"v1.4.0": false},
			unreleased: []string{config.DefaultEntryBugFixID},
			want:       "1.4.1",
			wantEffect: config.VersionEffectPatch,
		},
		"nothing to release": {
			released: map[string]bool{"1.4.0": false},
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			project := tempProject(t, tc.released, tc.unreleased)

			got, err := project.NextVersion(VersionModeAuto)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %+v", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("NextVersion: %v", err)
			}

			if got.Version != tc.want {
				t.Errorf("Version = %q, want %q", got.Version, tc.want)
			}

			if got.VersionType != tc.wantEffect {
				t.Errorf("VersionType = %v, want %v", got.VersionType, tc.wantEffect)
			}
		})
	}
}

func TestNextVersionRejectsUnknownMode(t *testing.T) {
	project := tempProject(t, nil, []string{config.DefaultEntryBugFixID})

	if _, err := project.NextVersion(VersionFindMode(42)); err == nil {
		t.Fatal("expected an error for an unsupported mode")
	}
}

func TestLatestRelease(t *testing.T) {
	project := tempProject(t, map[string]bool{
		"1.4.0":       false,
		"v2.0.0-rc.1": true,
	}, nil)

	latest, err := project.LatestRelease(false)
	if err != nil {
		t.Fatalf("LatestRelease: %v", err)
	}

	if latest.Info.Version != "v2.0.0-rc.1" {
		t.Errorf("latest = %q, want %q", latest.Info.Version, "v2.0.0-rc.1")
	}

	finalized, err := project.LatestRelease(true)
	if err != nil {
		t.Fatalf("LatestRelease(skipPreReleases): %v", err)
	}

	if finalized.Info.Version != "1.4.0" {
		t.Errorf("latest finalized = %q, want %q", finalized.Info.Version, "1.4.0")
	}
}

func TestLatestReleaseOnEmptyProject(t *testing.T) {
	project := tempProject(t, nil, nil)

	if _, err := project.LatestRelease(false); err == nil {
		t.Fatal("expected an error when nothing has been released")
	}
}

func TestLatestReleaseWithOnlyPreReleases(t *testing.T) {
	project := tempProject(t, map[string]bool{"v0.1.0-rc.1": true}, nil)

	if _, err := project.LatestRelease(true); err == nil {
		t.Fatal("expected an error when every release is a pre-release")
	}
}

func TestParseVersion(t *testing.T) {
	tests := map[string]bool{
		"1.4.0":         true,
		"v1.4.0":        true,
		"v2.0.0-rc.1":   true,
		"2.0.0+build.5": true,
		"not-a-version": false,
		"":              false,
	}

	for version, valid := range tests {
		_, err := ParseVersion(version)
		if valid && err != nil {
			t.Errorf("ParseVersion(%q) = %v, want success", version, err)
		}

		if !valid && err == nil {
			t.Errorf("ParseVersion(%q) succeeded, want an error", version)
		}
	}
}

func TestApplyVersionPrefix(t *testing.T) {
	tests := []struct {
		prefix, version, want string
	}{
		// a stored version keeps whatever spelling it has on disk, so the
		// prefix has to normalize both forms
		{"v", "1.4.0", "v1.4.0"},
		{"v", "v2.0.0-rc.1", "v2.0.0-rc.1"},
		{"", "1.4.0", "1.4.0"},
		{"", "v2.0.0-rc.1", "2.0.0-rc.1"},
		{"release-", "1.4.0", "release-1.4.0"},
	}

	for _, tc := range tests {
		if got := ApplyVersionPrefix(tc.prefix, tc.version); got != tc.want {
			t.Errorf("ApplyVersionPrefix(%q, %q) = %q, want %q",
				tc.prefix, tc.version, got, tc.want)
		}
	}
}
