package changelog

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/v2/internal/vcs"
)

// baseRefCandidates are tried in order when neither the config nor the caller
// names a base branch.
var baseRefCandidates = []string{
	"origin/develop", "origin/main", "origin/master",
	"develop", "main", "master",
}

// CheckOptions configures [Project.Check].
type CheckOptions struct {
	// Run executes git. The real binary is used when it is nil.
	Run vcs.Runner

	// CompareWith is the ref the branch is compared against. Empty means the
	// configured base branch, and failing that one detected from the
	// repository.
	CompareWith string

	// Staged compares the index instead of the branch, which is what a
	// pre-commit hook wants.
	Staged bool
}

// CheckResult is the outcome of [Project.Check].
type CheckResult struct {
	// BaseRef is what the branch was compared against. It is empty when the
	// index was checked instead.
	BaseRef string

	// Fragments are the changelog entries the branch adds, as repository
	// relative paths.
	Fragments []string

	// SkipReason is set when the check does not apply to this branch, which
	// is not a failure.
	SkipReason string
}

// OK reports whether the check passed.
func (r CheckResult) OK() bool {
	return r.SkipReason != "" || len(r.Fragments) > 0
}

// Check reports whether the branch adds at least one changelog entry.
//
// A branch that regenerates the changelog file is a release branch rather than
// a change that needs an entry of its own, so it is skipped instead of failed.
func (p *Project) Check(opts CheckOptions) (CheckResult, error) {
	run := opts.Run
	if run == nil {
		run = vcs.Git
	}

	root, err := repoRoot(run, p.rootDir)
	if err != nil {
		return CheckResult{}, err
	}

	// git reports paths relative to the repository root, so the directories
	// being looked for have to be expressed the same way.
	unreleasedDir, err := repoRelative(root, p.UnreleasedDir())
	if err != nil {
		return CheckResult{}, err
	}

	outputPath, err := repoRelative(root, p.OutputPath())
	if err != nil {
		return CheckResult{}, err
	}

	var result CheckResult

	// --name-status rather than --name-only: an entry has to be *added* to
	// count, while the changelog file counts however it changed, and one diff
	// carrying both statuses is cheaper than two.
	args := []string{"diff", "--name-status"}

	if opts.Staged {
		args = append(args, "--cached")
	} else {
		base, err := p.baseRef(run, opts.CompareWith)
		if err != nil {
			return CheckResult{}, err
		}

		result.BaseRef = base

		// three dots: what this branch added since it forked, ignoring what
		// happened on the base branch in the meantime
		args = append(args, base+"...HEAD")
	}

	out, err := run(p.rootDir, args...)
	if err != nil {
		return CheckResult{}, err
	}

	for _, line := range vcs.Lines(out) {
		status, path, ok := parseNameStatus(line)
		if !ok {
			continue
		}

		if path == outputPath {
			result.SkipReason = "the changelog file itself changed, so this looks like a release branch"

			continue
		}

		if !strings.HasPrefix(status, "A") || !isUnder(path, unreleasedDir) {
			continue
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			// hidden files are not entries, e.g. a .gitkeep
			continue
		}

		result.Fragments = append(result.Fragments, path)
	}

	return result, nil
}

// parseNameStatus splits one line of "git diff --name-status" into its status
// and the path it applies to.
//
// A rename carries both the old and the new path; the new one is what matters,
// and it is always last.
func parseNameStatus(line string) (status, path string, ok bool) {
	fields := strings.Split(line, "\t")
	if len(fields) < 2 {
		return "", "", false
	}

	return fields[0], fields[len(fields)-1], true
}

// baseRef returns the ref the branch should be compared against.
func (p *Project) baseRef(run vcs.Runner, override string) (string, error) {
	if override != "" {
		return override, nil
	}

	if configured := p.cfg.Check.BaseBranch; configured != "" {
		return configured, nil
	}

	// origin/HEAD is what the remote itself calls its default branch
	if out, err := run(p.rootDir, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD"); err == nil {
		if ref := strings.TrimSpace(string(out)); ref != "" {
			return strings.TrimPrefix(ref, "refs/remotes/"), nil
		}
	}

	for _, candidate := range baseRefCandidates {
		if _, err := run(p.rootDir, "rev-parse", "--verify", "--quiet", candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New(
		"unable to determine a base branch: set check.base_branch in the config or pass --compare-with")
}

// repoRoot returns the root of the repository holding dir.
func repoRoot(run vcs.Runner, dir string) (string, error) {
	out, err := run(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", errors.Wrap(err, "unable to determine the repository root")
	}

	return strings.TrimSpace(string(out)), nil
}

// repoRelative expresses path the way git reports it: relative to the
// repository root and separated by forward slashes.
func repoRelative(root, path string) (string, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to locate %q inside the repository at %q", path, root)
	}

	return filepath.ToSlash(rel), nil
}

// isUnder reports whether the repository relative path lies inside dir.
func isUnder(path, dir string) bool {
	if dir == "." {
		return true
	}

	return strings.HasPrefix(path, dir+"/")
}
