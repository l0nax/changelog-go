package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/v2/internal/tomlx"
)

// writeConfig writes body to a scratch config file and returns its path.
func writeConfig(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".changelog-go.toml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return path
}

func TestLoadDefaults(t *testing.T) {
	// only the keys without a sensible zero value
	path := writeConfig(t, "version = '2'\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.ChangelogDir != DefaultChangelogDir {
		t.Errorf("ChangelogDir = %q, want %q", cfg.ChangelogDir, DefaultChangelogDir)
	}

	if cfg.VersionPrefix != DefaultVersionPrefix {
		t.Errorf("VersionPrefix = %q, want %q", cfg.VersionPrefix, DefaultVersionPrefix)
	}

	if !cfg.PreRelease.FoldPreReleases {
		t.Error("fold_pre_releases should default to true")
	}

	if !cfg.PreRelease.Detect {
		t.Error("detect should default to true")
	}

	if cfg.PreRelease.DeletePreRelease {
		t.Error("delete_pre_release should default to false")
	}
}

// An explicit value has to beat the default, including when it is the zero
// value -- that is how a project turns the version prefix off.
func TestExplicitZeroValuesBeatDefaults(t *testing.T) {
	path := writeConfig(t, "version = '2'\nversion_prefix = ''\nfold_pre_releases = false\n\n[pre_release]\nfold_pre_releases = false\n")

	cfg, err := Load(path)
	if err == nil {
		t.Fatalf("a top-level fold_pre_releases should be rejected, got %+v", cfg)
	}

	path = writeConfig(t, "version = '2'\nversion_prefix = ''\n\n[pre_release]\nfold_pre_releases = false\n")

	cfg, err = Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.VersionPrefix != "" {
		t.Errorf("VersionPrefix = %q, want the empty string", cfg.VersionPrefix)
	}

	if cfg.PreRelease.FoldPreReleases {
		t.Error("fold_pre_releases = false was ignored")
	}
}

func TestDefaultProjectConfigIsValid(t *testing.T) {
	cfg, err := DefaultProjectConfig()
	if err != nil {
		t.Fatalf("the built-in default config does not load: %v", err)
	}

	if !cfg.Version.IsValid() {
		t.Errorf("default config version = %q, want a valid one", cfg.Version)
	}

	if len(cfg.Entry.Types) != 6 {
		t.Errorf("expected 6 default change types, got %d", len(cfg.Entry.Types))
	}

	for _, id := range []string{
		DefaultEntryNewFeatureID, DefaultEntryBugFixID, DefaultEntryFeatureChangeID,
		DefaultEntryDeprecateID, DefaultEntryRemovalID, DefaultEntrySecurityID,
	} {
		if _, ok := cfg.ResolveType(id); !ok {
			t.Errorf("default config is missing the %q type", id)
		}
	}
}

// Renaming a key must not fail silently: a config that still uses the
// v2.0.0-rc.1 spelling has to say what to replace it with.
func TestLegacyCamelCaseKeysAreRejected(t *testing.T) {
	path := writeConfig(t, `version = '2'

[preRelease]
deletePreRelease = false
detect = true
foldPreReleases = false
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected the legacy camelCase keys to be rejected")
	}

	for _, want := range []string{"preRelease", "pre_release"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not mention %q:\n%v", want, err)
		}
	}
}

func TestUnknownKeyIsRejectedWithItsLine(t *testing.T) {
	path := writeConfig(t, "version = '2'\ntypo_key = 'oops'\n")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected an unknown key to be rejected")
	}

	if !strings.Contains(err.Error(), "typo_key") {
		t.Errorf("error does not name the offending key:\n%v", err)
	}
}

func TestDuplicateTypeIDsAreRejected(t *testing.T) {
	path := writeConfig(t, `version = '2'

[entry]
  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'
  title = 'Bug Fixed'
  effect = 'patch'

  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Also fixed'
  title = 'Bug Fixed'
  effect = 'patch'
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected duplicate type IDs to be rejected")
	}

	if !strings.Contains(err.Error(), "bug_fix") {
		t.Errorf("error does not name the duplicated ID:\n%v", err)
	}
}

func TestInvalidEffectIsRejected(t *testing.T) {
	path := writeConfig(t, `version = '2'

[entry]
  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'
  title = 'Bug Fixed'
  effect = 'enormous'
`)

	if _, err := Load(path); err == nil {
		t.Fatal("expected an unknown version effect to be rejected")
	}
}

func TestMissingEffectIsRejected(t *testing.T) {
	path := writeConfig(t, `version = '2'

[entry]
  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'
  title = 'Bug Fixed'
`)

	if _, err := Load(path); err == nil {
		t.Fatal("expected a type without an effect to be rejected")
	}
}

// A config that is written back out has to stay loadable, which means the
// effect renders as text and an unset option does not turn into an empty
// string.
func TestConfigSurvivesAMarshalRoundTrip(t *testing.T) {
	cfg, err := DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	raw, err := Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if strings.Contains(string(raw), "effect = 2") {
		t.Errorf("effect was written numerically:\n%s", raw)
	}

	if !strings.Contains(string(raw), "effect = 'minor'") {
		t.Errorf("effect was not written as text:\n%s", raw)
	}

	if strings.Contains(string(raw), "description = ''") {
		t.Errorf("an unset description was written out:\n%s", raw)
	}

	back := Default()
	if err := tomlx.Unmarshal(raw, &back); err != nil {
		t.Fatalf("the marshaled config no longer decodes: %v\n%s", tomlx.Explain(err), raw)
	}

	if err := back.Validate(); err != nil {
		t.Fatalf("the marshaled config no longer validates: %v\n%s", err, raw)
	}

	if len(back.Entry.Types) != len(cfg.Entry.Types) {
		t.Errorf("types changed across the round trip: %d -> %d",
			len(cfg.Entry.Types), len(back.Entry.Types))
	}

	for i := range cfg.Entry.Types {
		if back.Entry.Types[i] != cfg.Entry.Types[i] {
			t.Errorf("type %d changed: %+v -> %+v", i, cfg.Entry.Types[i], back.Entry.Types[i])
		}
	}
}

func TestUnsetOutputPathIsNotWrittenOut(t *testing.T) {
	cfg := Default()
	cfg.Version = "2"
	cfg.OutputPath = typact.None[string]()

	raw, err := Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if strings.Contains(string(raw), "output_path = ''") {
		t.Errorf("an unset output_path was written out:\n%s", raw)
	}
}

func TestVersionEffectMarshalText(t *testing.T) {
	for effect, want := range map[VersionEffect]string{
		VersionEffectMajor: "major",
		VersionEffectMinor: "minor",
		VersionEffectPatch: "patch",
	} {
		got, err := effect.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v): %v", effect, err)
		}

		if string(got) != want {
			t.Errorf("MarshalText(%v) = %q, want %q", effect, got, want)
		}

		var back VersionEffect
		if err := back.UnmarshalText(got); err != nil {
			t.Fatalf("UnmarshalText(%q): %v", got, err)
		}

		if back != effect {
			t.Errorf("round trip changed %v into %v", effect, back)
		}
	}

	// must report instead of panicking the way String() does
	if _, err := VersionEffect(0).MarshalText(); err == nil {
		t.Error("expected an error for the zero effect")
	}
}

func TestLoadReportsAMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "nope.toml")); err == nil {
		t.Fatal("expected an error for a missing config file")
	}
}
