package cmd

import (
	"log/slog"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
)

func newCheckCmd() *cli.Command {
	return &cli.Command{
		Name: "check",
		UsageText: `Check verifies that the branch adds at least one changelog entry.

Run it on every merge request to keep the changelog honest:

    changelog check --compare-with origin/develop

It exits 0 when an entry was added and 1 when none was. A branch that changes the
changelog file itself is a release branch rather than a change needing an entry of
its own, so it is skipped instead of failed.`,
		Action: checkAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "compare-with",
				Usage: "Ref to compare the branch against (default: check.base_branch, else detected)",
			},
			&cli.BoolFlag{
				Name:  "staged",
				Usage: "Checks the index instead of the branch, for use in a pre-commit hook",
			},
		},
	}
}

func checkAction(c *cli.Context) error {
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	result, err := project.Check(changelog.CheckOptions{
		CompareWith: c.String("compare-with"),
		Staged:      c.Bool("staged"),
	})
	if err != nil {
		return err
	}

	if result.SkipReason != "" {
		slog.Info("Skipping check", slog.String("reason", result.SkipReason))

		return nil
	}

	if !result.OK() {
		return cli.Exit(missingEntryMessage(result), ExitError)
	}

	slog.Info("Found new changelog entries",
		slog.Int("count", len(result.Fragments)),
		slog.Any("entries", result.Fragments))

	return nil
}

// missingEntryMessage tells the author what to do about a branch with no entry.
func missingEntryMessage(result changelog.CheckResult) string {
	where := "the index"
	if result.BaseRef != "" {
		where = result.BaseRef
	}

	return "no changelog entry added relative to " + where +
		`. Run "changelog new" to add one, e.g. changelog new -t bug_fix --title "..."`
}
