package cmd

import (
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func newValidateCmd() *cli.Command {
	return &cli.Command{
		Name: "validate",
		UsageText: `Validate parses the config and every changelog file, reporting what is wrong.

It reports every problem it finds rather than stopping at the first, so that a
broken project can be fixed in one pass. It exits 1 if anything is wrong.`,
		Action: validateAction,
	}
}

func validateAction(c *cli.Context) error {
	// A config that does not load at all is reported by loadProject, with the
	// offending line.
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	problems, err := project.Validate()
	if err != nil {
		return err
	}

	for _, problem := range problems {
		fmt.Println(problem)
	}

	if len(problems) > 0 {
		return cli.Exit(fmt.Sprintf("found %d problem(s)", len(problems)), ExitError)
	}

	slog.Info("Everything checks out", slog.String("changelog_dir", project.ChangelogDir()))

	return nil
}
