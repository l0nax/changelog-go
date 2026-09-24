package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
)

func newLatestCmd() *cli.Command {
	return &cli.Command{
		Name:      "latest",
		UsageText: "Latest returns the latest released version, i.e. the release with the highest number.",
		Action:    latestAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "remove-prefix",
				Usage: "Prints the bare version, without the configured version prefix",
			},
			&cli.BoolFlag{
				Name:  "skip-prereleases",
				Usage: "Ignores pre-releases, returning the latest finalized version",
			},
		},
	}
}

func latestAction(c *cli.Context) error {
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	release, err := project.LatestRelease(c.Bool("skip-prereleases"))
	if err != nil {
		var notFound *changelog.NoReleasesError
		if errors.As(err, &notFound) {
			return notFoundError("%s", err)
		}

		return err
	}

	fmt.Println(printableVersion(project, release.Info.Version, c.Bool("remove-prefix")))

	return nil
}
