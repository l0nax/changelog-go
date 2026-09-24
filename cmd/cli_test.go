package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

// result is the observable outcome of one command invocation.
type result struct {
	stdout string
	code   int
}

// runCLI invokes the app the way a shell would, with stdout captured and the
// config pinned via the global flag so that no test depends on the working
// directory.
func runCLI(t *testing.T, root string, args ...string) result {
	t.Helper()

	stdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = w

	argv := append([]string{
		"changelog", "--config", filepath.Join(root, ConfigFileName),
	}, args...)

	runErr := newTestApp().Run(argv)

	_ = w.Close()

	os.Stdout = stdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}

	code := ExitOK
	if runErr != nil {
		code = exitCodeOf(runErr)
	}

	return result{stdout: buf.String(), code: code}
}

// newTestApp mirrors the real app without the logging setup, which would
// otherwise fight the captured stdout.
func newTestApp() *cli.App {
	return &cli.App{
		Name: "changelog",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "changelog-dir"},
		},
		Commands: []*cli.Command{
			newNewCmd(), newInitCmd(), newReleaseCmd(),
			newNextCmd(), newLatestCmd(), newShowCmd(), newVersionCmd(),
		},
		ExitErrHandler: func(_ *cli.Context, _ error) {},
	}
}

// project lays out a v2 project with the given pending entries.
func project(t *testing.T, pending map[string]string) string {
	t.Helper()

	root := t.TempDir()

	writeFile(t, filepath.Join(root, ConfigFileName), config.DefaultConfigFile())

	for name, body := range pending {
		writeFile(t, filepath.Join(root, ".changelogs", "unreleased", name), body)
	}

	if len(pending) == 0 {
		if err := os.MkdirAll(filepath.Join(root, ".changelogs", "unreleased"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	return root
}

func entryFile(typeID, title string) string {
	return "change_type_id = '" + typeID + "'\ntitle = '" + title + "'\n"
}

// pendingFiles lists the unreleased entry files.
func pendingFiles(t *testing.T, root string) []string {
	t.Helper()

	dirEntries, err := os.ReadDir(filepath.Join(root, ".changelogs", "unreleased"))
	if err != nil {
		t.Fatalf("read unreleased: %v", err)
	}

	names := make([]string, 0, len(dirEntries))
	for _, e := range dirEntries {
		names = append(names, e.Name())
	}

	return names
}

// The headline of Stage 2: no terminal, no prompt, no hang.
func TestNewIsFullyScriptable(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "new", "-t", "bug_fix", "--title", "Fix the sprocket",
		"--body", "Why it was broken.", "--author", "emanuel")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d\n%s", got.code, ExitOK, got.stdout)
	}

	names := pendingFiles(t, root)
	if len(names) != 1 {
		t.Fatalf("expected exactly one entry, got %v", names)
	}

	raw, err := os.ReadFile(filepath.Join(root, ".changelogs", "unreleased", names[0]))
	if err != nil {
		t.Fatalf("read entry: %v", err)
	}

	for _, want := range []string{
		"change_type_id = 'bug_fix'",
		"title = 'Fix the sprocket'",
		"Why it was broken.",
		"author = 'emanuel'",
	} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("entry is missing %q:\n%s", want, raw)
		}
	}
}

func TestNewTitleFromPositionalArg(t *testing.T) {
	root := project(t, nil)

	if got := runCLI(t, root, "new", "-t", "bug_fix", "Fix it"); got.code != ExitOK {
		t.Fatalf("exit = %d, want %d", got.code, ExitOK)
	}

	names := pendingFiles(t, root)
	if len(names) != 1 {
		t.Fatalf("expected one entry, got %v", names)
	}
}

func TestNewUnknownTypeListsTheValidIDs(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "new", "--type", "no_such_type", "--title", "x")

	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}

	if len(pendingFiles(t, root)) != 0 {
		t.Error("an entry was written despite the unknown type")
	}
}

// Without a TTY the command has to fail rather than block on a prompt.
func TestNewWithoutTTYAndWithoutFlagsFails(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "new", "--interactive=false")

	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}
}

func TestNewDryRunWritesNothing(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "new", "-t", "bug_fix", "--title", "Fix it", "--dry-run")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d", got.code, ExitOK)
	}

	if !strings.Contains(got.stdout, "change_type_id = 'bug_fix'") {
		t.Errorf("dry run did not print the entry:\n%s", got.stdout)
	}

	if len(pendingFiles(t, root)) != 0 {
		t.Error("--dry-run wrote an entry")
	}
}

func TestShowRendersTheReleaseNotes(t *testing.T) {
	root := project(t, map[string]string{
		"a": entryFile("new_feat", "Add the widget"),
		"b": entryFile("bug_fix", "Fix the sprocket"),
	})

	if got := runCLI(t, root, "release", "1.0.0"); got.code != ExitOK {
		t.Fatalf("release exit = %d: %s", got.code, got.stdout)
	}

	latest := runCLI(t, root, "latest")
	if latest.code != ExitOK {
		t.Fatalf("latest exit = %d", latest.code)
	}

	version := strings.TrimSpace(latest.stdout)
	if version != "v1.0.0" {
		t.Fatalf("latest = %q, want %q", version, "v1.0.0")
	}

	// the documented pipeline: changelog show $(changelog latest)
	got := runCLI(t, root, "show", version)
	if got.code != ExitOK {
		t.Fatalf("show exit = %d: %s", got.code, got.stdout)
	}

	for _, want := range []string{"## v1.0.0 (", "### Added (1 change)", "- Add the widget"} {
		if !strings.Contains(got.stdout, want) {
			t.Errorf("show output is missing %q:\n%s", want, got.stdout)
		}
	}

	// no leading or trailing blank lines on a body meant for a release page
	if strings.TrimSpace(got.stdout) != strings.TrimSuffix(got.stdout, "\n") {
		t.Errorf("show output is not trimmed:\n%q", got.stdout)
	}
}

func TestShowNoHeading(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	if got := runCLI(t, root, "release", "1.0.0"); got.code != ExitOK {
		t.Fatalf("release exit = %d", got.code)
	}

	got := runCLI(t, root, "show", "--no-heading", "1.0.0")
	if got.code != ExitOK {
		t.Fatalf("show exit = %d", got.code)
	}

	// the version heading goes, the group headings stay
	if strings.Contains(got.stdout, "## v1.0.0") {
		t.Errorf("--no-heading still rendered the version heading:\n%s", got.stdout)
	}

	if !strings.HasPrefix(got.stdout, "### Added") {
		t.Errorf("body should start at the first group:\n%q", got.stdout)
	}
}

func TestShowUnreleased(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	got := runCLI(t, root, "show", "unreleased")
	if got.code != ExitOK {
		t.Fatalf("show exit = %d: %s", got.code, got.stdout)
	}

	if !strings.Contains(got.stdout, "## Unreleased") {
		t.Errorf("expected an Unreleased heading:\n%s", got.stdout)
	}

	if !strings.Contains(got.stdout, "- Add the widget") {
		t.Errorf("pending entry missing:\n%s", got.stdout)
	}
}

func TestShowMissingVersion(t *testing.T) {
	root := project(t, nil)

	if got := runCLI(t, root, "show", "9.9.9"); got.code != ExitNotFound {
		t.Fatalf("exit = %d, want %d", got.code, ExitNotFound)
	}

	if got := runCLI(t, root, "show"); got.code != ExitUsage {
		t.Fatalf("exit without an argument = %d, want %d", got.code, ExitUsage)
	}
}

func TestReleaseRefusesAnEmptyRelease(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "release", "1.0.0")
	if got.code != ExitNothingToDo {
		t.Fatalf("exit = %d, want %d", got.code, ExitNothingToDo)
	}

	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released", "1.0.0")); !os.IsNotExist(err) {
		t.Error("an empty release was created anyway")
	}

	if got := runCLI(t, root, "release", "--allow-empty", "1.0.0"); got.code != ExitOK {
		t.Fatalf("--allow-empty exit = %d, want %d", got.code, ExitOK)
	}

	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released", "1.0.0")); err != nil {
		t.Errorf("--allow-empty did not create the release: %v", err)
	}
}

func TestReleaseDryRunChangesNothing(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	got := runCLI(t, root, "release", "--dry-run", "1.0.0")
	if got.code != ExitOK {
		t.Fatalf("exit = %d: %s", got.code, got.stdout)
	}

	if !strings.Contains(got.stdout, "- Add the widget") {
		t.Errorf("dry run did not print the notes:\n%s", got.stdout)
	}

	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released", "1.0.0")); !os.IsNotExist(err) {
		t.Error("--dry-run created the release directory")
	}

	if len(pendingFiles(t, root)) != 1 {
		t.Error("--dry-run consumed the pending entry")
	}

	if _, err := os.Stat(filepath.Join(root, "CHANGELOG.md")); !os.IsNotExist(err) {
		t.Error("--dry-run wrote the changelog file")
	}
}

func TestNextRemovePrefix(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	prefixed := runCLI(t, root, "next", "auto")
	if strings.TrimSpace(prefixed.stdout) != "v0.1.0" {
		t.Errorf("next auto = %q, want %q", strings.TrimSpace(prefixed.stdout), "v0.1.0")
	}

	bare := runCLI(t, root, "next", "--remove-prefix", "auto")
	if strings.TrimSpace(bare.stdout) != "0.1.0" {
		t.Errorf("next --remove-prefix = %q, want %q", strings.TrimSpace(bare.stdout), "0.1.0")
	}
}

func TestChangelogDirOverride(t *testing.T) {
	root := project(t, nil)

	elsewhere := t.TempDir()
	writeFile(t, filepath.Join(elsewhere, "unreleased", "a"), entryFile("new_feat", "Add the widget"))

	got := runCLI(t, root, "--changelog-dir", elsewhere, "show", "unreleased")
	if got.code != ExitOK {
		t.Fatalf("exit = %d: %s", got.code, got.stdout)
	}

	if !strings.Contains(got.stdout, "- Add the widget") {
		t.Errorf("--changelog-dir was ignored:\n%s", got.stdout)
	}
}

func TestVersionCommand(t *testing.T) {
	root := project(t, nil)

	got := runCLI(t, root, "version", "--short")
	if got.code != ExitOK {
		t.Fatalf("exit = %d", got.code)
	}

	// never blank, even without the -ldflags stamps
	if strings.TrimSpace(got.stdout) == "" {
		t.Error("version printed nothing")
	}
}

func TestExitCodeOfPlainError(t *testing.T) {
	if got := exitCodeOf(errors.New("boom")); got != ExitError {
		t.Errorf("exitCodeOf(plain) = %d, want %d", got, ExitError)
	}
}

func TestReleaseAuto(t *testing.T) {
	root := project(t, map[string]string{
		"a": entryFile("new_feat", "Add the widget"),
		"b": entryFile("bug_fix", "Fix the sprocket"),
	})

	// no releases yet, so the pending entries make it 0.1.0
	if got := runCLI(t, root, "release", "--auto"); got.code != ExitOK {
		t.Fatalf("release --auto exit = %d: %s", got.code, got.stdout)
	}

	latest := runCLI(t, root, "latest")
	if strings.TrimSpace(latest.stdout) != "v0.1.0" {
		t.Fatalf("latest = %q, want %q", strings.TrimSpace(latest.stdout), "v0.1.0")
	}

	// the directory is named by the bare version; the prefix is rendering
	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released", "0.1.0")); err != nil {
		t.Errorf("expected a bare 0.1.0 release directory: %v", err)
	}

	// a removal is a major bump
	writeFile(t, filepath.Join(root, ".changelogs", "unreleased", "c"),
		entryFile("rem_feat", "Remove the doohickey"))

	if got := runCLI(t, root, "release", "--auto"); got.code != ExitOK {
		t.Fatalf("second release --auto exit = %d: %s", got.code, got.stdout)
	}

	latest = runCLI(t, root, "latest")
	if strings.TrimSpace(latest.stdout) != "v1.0.0" {
		t.Errorf("latest = %q, want %q", strings.TrimSpace(latest.stdout), "v1.0.0")
	}
}

func TestReleaseAutoRejectsAVersionToo(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	got := runCLI(t, root, "release", "--auto", "1.0.0")
	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}

	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released")); !os.IsNotExist(err) {
		t.Error("a release was created despite the bad invocation")
	}
}

func TestReleaseWithoutAVersionOrAuto(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	if got := runCLI(t, root, "release"); got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}
}

func TestReleaseAutoWithNothingPending(t *testing.T) {
	root := project(t, nil)

	// nothing to derive a version from
	if got := runCLI(t, root, "release", "--auto"); got.code == ExitOK {
		t.Fatal("expected release --auto to fail with no pending entries")
	}
}

func TestReleaseAutoDryRunChangesNothing(t *testing.T) {
	root := project(t, map[string]string{"a": entryFile("new_feat", "Add the widget")})

	got := runCLI(t, root, "release", "--auto", "--dry-run")
	if got.code != ExitOK {
		t.Fatalf("exit = %d: %s", got.code, got.stdout)
	}

	if !strings.Contains(got.stdout, "## v0.1.0") {
		t.Errorf("dry run did not render the derived version:\n%s", got.stdout)
	}

	if _, err := os.Stat(filepath.Join(root, ".changelogs", "released")); !os.IsNotExist(err) {
		t.Error("--dry-run created a release")
	}
}
