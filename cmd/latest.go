package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
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
	project, err := loadProject()
	if err != nil {
		return err
	}

	release, err := project.LatestRelease(c.Bool("skip-prereleases"))
	if err != nil {
		return err
	}

	fmt.Println(printableVersion(project, release.Info.Version, c.Bool("remove-prefix")))

	return nil
}
