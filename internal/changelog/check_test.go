package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/l0nax/changelog-go/v2/internal/config"
	"gitlab.com/l0nax/changelog-go/v2/internal/vcs"
)

// fakeGit answers the git commands [Project.Check] issues from a canned script,
// so the decision logic can be exercised without a repository.
type fakeGit struct {
	// root is what "rev-parse --show-toplevel" reports.
	root string
	// changed is what the diff reports, as "<status>\t<path>" lines, exactly
	// the shape "git diff --name-status" produces.
	changed []string
	// refs are the refs "rev-parse --verify" accepts.
	refs map[string]bool
	// originHEAD is what "symbolic-ref" reports, if anything.
	originHEAD string

	// calls records every invocation, so a test can assert on the git
	// command that was actually built.
	calls [][]string
}

func (f *fakeGit) run(_ string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, args)

	switch {
	case args[0] == "rev-parse" && args[1] == "--show-toplevel":
		return []byte(f.root + "\n"), nil

	case args[0] == "symbolic-ref":
		if f.originHEAD == "" {
			return nil, fmt.Errorf("git symbolic-ref: no such ref")
		}

		return []byte("refs/remotes/" + f.originHEAD + "\n"), nil

	case args[0] == "rev-parse" && args[1] == "--verify":
		if f.refs[args[len(args)-1]] {
			return []byte("deadbeef\n"), nil
		}

		return nil, fmt.Errorf("git rev-parse: unknown revision")

	case args[0] == "diff":
		return []byte(strings.Join(f.changed, "\n")), nil
	}

	return nil, fmt.Errorf("unexpected git call: %v", args)
}

// checkProject returns a project rooted at the fake repository root.
func checkProject(t *testing.T, root string, mutate func(*config.Config)) *Project {
	t.Helper()

	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	if mutate != nil {
		mutate(&cfg)
	}

	return NewProject(cfg, root)
}

func TestCheck(t *testing.T) {
	const root = "/repo"

	tests := map[string]struct {
		changed      []string
		staged       bool
		compareWith  string
		wantOK       bool
		wantSkipped  bool
		wantFragment string
	}{
		"a branch with an entry passes": {
			changed:      []string{"M\tinternal/thing.go", "A\t.changelogs/unreleased/1733309257759"},
			wantOK:       true,
			wantFragment: ".changelogs/unreleased/1733309257759",
		},
		"a branch without an entry fails": {
			changed: []string{"M\tinternal/thing.go", "A\tREADME.md"},
			wantOK:  false,
		},
		"a branch that changed nothing fails": {
			changed: nil,
			wantOK:  false,
		},
		"editing an existing entry is not adding one": {
			changed: []string{"M\t.changelogs/unreleased/1733309257759"},
			wantOK:  false,
		},
		"a renamed entry counts by its new path": {
			changed:      []string{"R100\t.changelogs/unreleased/old\t.changelogs/unreleased/new"},
			wantOK:       false,
			wantFragment: "",
		},
		// A release branch regenerates the changelog rather than adding an
		// entry. Note the changelog file is modified, not added, which is
		// exactly what a status-blind filter missed.
		"a release branch is skipped": {
			changed:     []string{"M\tCHANGELOG.md", "A\t.changelogs/released/2.0.0/ReleaseInfo"},
			wantOK:      true,
			wantSkipped: true,
		},
		"a released entry is not an added entry": {
			changed: []string{"A\t.changelogs/released/2.0.0/1733309257759"},
			wantOK:  false,
		},
		"hidden files are not entries": {
			changed: []string{"A\t.changelogs/unreleased/.gitkeep"},
			wantOK:  false,
		},
		"the index can be checked instead": {
			changed:      []string{"A\t.changelogs/unreleased/1733309257759"},
			staged:       true,
			wantOK:       true,
			wantFragment: ".changelogs/unreleased/1733309257759",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			git := &fakeGit{root: root, changed: tc.changed, refs: map[string]bool{"origin/develop": true}}
			project := checkProject(t, root, nil)

			got, err := project.Check(CheckOptions{
				Run:         git.run,
				Staged:      tc.staged,
				CompareWith: tc.compareWith,
			})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}

			if got.OK() != tc.wantOK {
				t.Errorf("OK() = %t, want %t (fragments %v, skip %q)",
					got.OK(), tc.wantOK, got.Fragments, got.SkipReason)
			}

			if (got.SkipReason != "") != tc.wantSkipped {
				t.Errorf("SkipReason = %q, want skipped = %t", got.SkipReason, tc.wantSkipped)
			}

			if tc.wantFragment != "" {
				if len(got.Fragments) != 1 || got.Fragments[0] != tc.wantFragment {
					t.Errorf("Fragments = %v, want [%s]", got.Fragments, tc.wantFragment)
				}
			}
		})
	}
}

// The diff has to be three-dot, or a branch fails whenever the base moved on
// without it.
func TestCheckComparesSinceTheMergeBase(t *testing.T) {
	git := &fakeGit{root: "/repo", refs: map[string]bool{"origin/develop": true}}
	project := checkProject(t, "/repo", nil)

	if _, err := project.Check(CheckOptions{Run: git.run, CompareWith: "origin/main"}); err != nil {
		t.Fatalf("Check: %v", err)
	}

	var diff []string

	for _, call := range git.calls {
		if call[0] == "diff" {
			diff = call
		}
	}

	if diff == nil {
		t.Fatal("no diff was run")
	}

	joined := strings.Join(diff, " ")

	for _, want := range []string{"--name-status", "origin/main...HEAD"} {
		if !strings.Contains(joined, want) {
			t.Errorf("diff %q is missing %q", joined, want)
		}
	}
}

func TestCheckStagedDoesNotUseABaseRef(t *testing.T) {
	git := &fakeGit{root: "/repo"}
	project := checkProject(t, "/repo", nil)

	got, err := project.Check(CheckOptions{Run: git.run, Staged: true})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if got.BaseRef != "" {
		t.Errorf("BaseRef = %q, want empty when checking the index", got.BaseRef)
	}

	var diff []string

	for _, call := range git.calls {
		if call[0] == "diff" {
			diff = call
		}
	}

	if !strings.Contains(strings.Join(diff, " "), "--cached") {
		t.Errorf("staged check did not use --cached: %v", diff)
	}
}

func TestCheckBaseRefResolution(t *testing.T) {
	tests := map[string]struct {
		configured string
		override   string
		originHEAD string
		refs       map[string]bool
		want       string
	}{
		"--compare-with wins": {
			configured: "origin/develop",
			override:   "origin/release",
			want:       "origin/release",
		},
		"config beats detection": {
			configured: "origin/trunk",
			originHEAD: "origin/main",
			want:       "origin/trunk",
		},
		"origin/HEAD is preferred": {
			originHEAD: "origin/main",
			refs:       map[string]bool{"origin/develop": true},
			want:       "origin/main",
		},
		"falls back to a known branch": {
			refs: map[string]bool{"origin/master": true},
			want: "origin/master",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			git := &fakeGit{root: "/repo", originHEAD: tc.originHEAD, refs: tc.refs}

			project := checkProject(t, "/repo", func(c *config.Config) {
				c.Check.BaseBranch = tc.configured
			})

			got, err := project.Check(CheckOptions{Run: git.run, CompareWith: tc.override})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}

			if got.BaseRef != tc.want {
				t.Errorf("BaseRef = %q, want %q", got.BaseRef, tc.want)
			}
		})
	}
}

func TestCheckWithoutAnyBaseRef(t *testing.T) {
	git := &fakeGit{root: "/repo"}
	project := checkProject(t, "/repo", nil)

	_, err := project.Check(CheckOptions{Run: git.run})
	if err == nil {
		t.Fatal("expected an error when no base branch can be determined")
	}

	if !strings.Contains(err.Error(), "--compare-with") {
		t.Errorf("error should say how to fix it: %v", err)
	}
}

// A changelog directory outside the repository root still has to be matched
// against the paths git reports.
func TestCheckWithANestedProject(t *testing.T) {
	git := &fakeGit{
		root:    "/repo",
		changed: []string{"A\tsub/project/.changelogs/unreleased/1733309257759"},
		refs:    map[string]bool{"origin/develop": true},
	}

	project := checkProject(t, "/repo/sub/project", nil)

	got, err := project.Check(CheckOptions{Run: git.run})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if !got.OK() {
		t.Errorf("entry under a nested project was not found: %+v", got)
	}
}

// TestCheckAgainstRealGit exercises the git plumbing the fake stands in for.
//
// It builds throwaway repositories, so it is skipped wherever git is not
// available.
func TestCheckAgainstRealGit(t *testing.T) {
	if !vcs.Available() {
		t.Skip("git is not available")
	}

	// newRepo returns a repository with one commit on "main" and a branch
	// checked out on top of it.
	newRepo := func(t *testing.T) string {
		t.Helper()

		root := t.TempDir()

		git := func(args ...string) {
			t.Helper()

			cmd := exec.Command("git", args...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(),
				"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
				"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null",
			)

			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
			}
		}

		write := func(path, body string) {
			t.Helper()

			full := filepath.Join(root, path)
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}

			if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
				t.Fatalf("write %q: %v", path, err)
			}
		}

		git("init", "--initial-branch=main", ".")

		write(".changelog-go.toml", config.DefaultConfigFile())
		write("CHANGELOG.md", "# Changelog\n")
		write(".changelogs/unreleased/.gitkeep", "")

		git("add", "-A")
		git("commit", "-m", "initial")
		git("checkout", "-b", "feature")

		return root
	}

	load := func(t *testing.T, root string) *Project {
		t.Helper()

		cfg, err := config.Load(filepath.Join(root, ".changelog-go.toml"))
		if err != nil {
			t.Fatalf("load config: %v", err)
		}

		return NewProject(cfg, root)
	}

	commit := func(t *testing.T, root string, paths ...string) {
		t.Helper()

		for _, args := range [][]string{append([]string{"add"}, paths...), {"commit", "-m", "change"}} {
			cmd := exec.Command("git", args...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(),
				"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
				"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null",
			)

			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
			}
		}
	}

	writeEntry := func(t *testing.T, root, name string) string {
		t.Helper()

		rel := filepath.Join(".changelogs", "unreleased", name)
		if err := os.WriteFile(filepath.Join(root, rel),
			[]byte("change_type_id = 'bug_fix'\ntitle = 'Fix it'\n"), 0o644); err != nil {
			t.Fatalf("write entry: %v", err)
		}

		return rel
	}

	t.Run("no entry fails", func(t *testing.T) {
		root := newRepo(t)

		if err := os.WriteFile(filepath.Join(root, "thing.go"), []byte("package main\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}

		commit(t, root, "thing.go")

		got, err := load(t, root).Check(CheckOptions{CompareWith: "main"})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}

		if got.OK() {
			t.Errorf("a branch with no entry should fail: %+v", got)
		}
	})

	t.Run("one entry passes", func(t *testing.T) {
		root := newRepo(t)
		rel := writeEntry(t, root, "1733309257759")
		commit(t, root, rel)

		got, err := load(t, root).Check(CheckOptions{CompareWith: "main"})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}

		if !got.OK() || len(got.Fragments) != 1 {
			t.Errorf("a branch with an entry should pass: %+v", got)
		}
	})

	t.Run("staged only", func(t *testing.T) {
		root := newRepo(t)
		rel := writeEntry(t, root, "1733309257759")

		cmd := exec.Command("git", "add", rel)
		cmd.Dir = root

		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git add: %v\n%s", err, out)
		}

		project := load(t, root)

		staged, err := project.Check(CheckOptions{Staged: true})
		if err != nil {
			t.Fatalf("Check --staged: %v", err)
		}

		if !staged.OK() {
			t.Errorf("a staged entry should pass --staged: %+v", staged)
		}

		// not committed, so the branch itself still has nothing
		branch, err := project.Check(CheckOptions{CompareWith: "main"})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}

		if branch.OK() {
			t.Errorf("an uncommitted entry should not pass the branch check: %+v", branch)
		}
	})

	t.Run("release branch is skipped", func(t *testing.T) {
		root := newRepo(t)

		if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"),
			[]byte("# Changelog\n\n## v1.0.0 (2024-01-01)\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}

		commit(t, root, "CHANGELOG.md")

		got, err := load(t, root).Check(CheckOptions{CompareWith: "main"})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}

		if got.SkipReason == "" {
			t.Errorf("a release branch should be skipped: %+v", got)
		}

		if !got.OK() {
			t.Errorf("a skipped check must not fail: %+v", got)
		}
	})
}
