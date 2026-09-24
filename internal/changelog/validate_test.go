package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// validateProject lays out a project from a map of repository relative paths.
func validateProject(t *testing.T, files map[string]string) *Project {
	t.Helper()

	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		t.Fatalf("default config: %v", err)
	}

	root := t.TempDir()

	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write %q: %v", name, err)
		}
	}

	return NewProject(cfg, root)
}

// releaseInfo returns a well-formed ReleaseInfo file.
func releaseInfo(version string) string {
	return "version = '" + version + "'\ndate = 2024-01-01T00:00:00Z\npre_release = false\n"
}

func TestValidateAcceptsAHealthyProject(t *testing.T) {
	project := validateProject(t, map[string]string{
		".changelogs/released/1.0.0/ReleaseInfo": releaseInfo("1.0.0"),
		".changelogs/released/1.0.0/a":           "change_type_id = 'new_feat'\ntitle = 'Add the widget'\n",
		// a directory stored with the prefix still matches its version
		".changelogs/released/v1.1.0-rc.1/ReleaseInfo": releaseInfo("v1.1.0-rc.1"),
		".changelogs/unreleased/b":                     "change_type_id = 'bug_fix'\ntitle = 'Fix it'\n",
		".changelogs/unreleased/.gitkeep":              "",
	})

	problems, err := project.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}

	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestValidateOnAnEmptyProject(t *testing.T) {
	project := validateProject(t, nil)

	problems, err := project.Validate()
	if err != nil {
		t.Fatalf("Validate on a fresh project: %v", err)
	}

	if len(problems) != 0 {
		t.Errorf("a project with nothing in it is valid, got %v", problems)
	}
}

func TestValidateReportsEveryProblem(t *testing.T) {
	project := validateProject(t, map[string]string{
		// missing ReleaseInfo
		".changelogs/released/1.0.0/a": "change_type_id = 'new_feat'\ntitle = 'Add the widget'\n",

		// version disagrees with the directory name
		".changelogs/released/2.0.0/ReleaseInfo": releaseInfo("3.0.0"),

		// malformed TOML
		".changelogs/released/4.0.0/ReleaseInfo": releaseInfo("4.0.0"),
		".changelogs/released/4.0.0/broken":      "change_type_id = 'new_feat'\ntitle = \n",

		// a typo'd key
		".changelogs/unreleased/typo": "change_type = 'bug_fix'\ntitle = 'Fix it'\n",

		// a type the config does not define
		".changelogs/unreleased/retired": "change_type_id = 'retired'\ntitle = 'Fix it'\n",

		// no title
		".changelogs/unreleased/untitled": "change_type_id = 'bug_fix'\n",
	})

	problems, err := project.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}

	// the point of validate is one pass, not one error at a time
	if len(problems) < 6 {
		t.Errorf("expected every problem to be reported, got %d:\n%v", len(problems), problems)
	}

	report := make([]string, len(problems))
	for i, problem := range problems {
		report[i] = problem.String()
	}

	joined := strings.Join(report, "\n")

	for _, want := range []string{
		"released/1.0.0: release directory has no ReleaseInfo file",
		"does not match the directory name",
		"unreleased/typo",
		"unknown change type retired",
		"no title set",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("report is missing %q:\n%s", want, joined)
		}
	}
}

// The report has to name files the way the person sees them, whatever
// directory the command ran from.
func TestValidateReportsProjectRelativePaths(t *testing.T) {
	project := validateProject(t, map[string]string{
		".changelogs/unreleased/bad": "change_type_id = 'retired'\ntitle = 'Fix it'\n",
	})

	problems, err := project.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}

	if len(problems) != 1 {
		t.Fatalf("expected one problem, got %v", problems)
	}

	if got := problems[0].Path; got != ".changelogs/unreleased/bad" {
		t.Errorf("Path = %q, want a project relative path", got)
	}
}

func TestValidateReportsAMissingVersion(t *testing.T) {
	project := validateProject(t, map[string]string{
		".changelogs/released/1.0.0/ReleaseInfo": "date = 2024-01-01T00:00:00Z\npre_release = false\n",
	})

	problems, err := project.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}

	if len(problems) != 1 || !strings.Contains(problems[0].Message, "no version recorded") {
		t.Errorf("expected a missing-version problem, got %v", problems)
	}
}
