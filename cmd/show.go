package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
)

// UnreleasedArg selects the pending entries instead of a released version.
const UnreleasedArg = "unreleased"

func newShowCmd() *cli.Command {
	return &cli.Command{
		Name:    "show",
		Aliases: []string{"ci-get"},
		UsageText: `Show renders the notes of a single release to stdout.

This is the body a forge publishes as the release notes, e.g.

    changelog show $(changelog latest)

Pass "unreleased" instead of a version to render the entries that are still pending.`,
		ArgsUsage: "<version|unreleased>",
		Args:      true,
		Action:    showAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-heading",
				Usage: "Omits the '## version (date)' heading, leaving only the notes",
			},
		},
	}
}

func showAction(c *cli.Context) error {
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	version := c.Args().First()
	if version == "" {
		return usageError("no version specified: pass a version or %q", UnreleasedArg)
	}

	var cl *changelog.Changelog

	if version == UnreleasedArg {
		cl, err = project.UnreleasedChangelog()
	} else {
		cl, err = project.ReleaseChangelog(version)
	}

	if err != nil {
		var notFound *changelog.VersionNotFoundError
		if errors.As(err, &notFound) {
			return notFoundError("%s", err)
		}

		return err
	}

	out, err := cl.RenderReleases(!c.Bool("no-heading"))
	if err != nil {
		return err
	}

	// The surrounding blank lines only matter inside the changelog file; on
	// its own the body should start and end cleanly.
	fmt.Println(strings.TrimSpace(string(out)))

	return nil
}
