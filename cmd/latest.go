package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
)

func newLatestCmd() *cli.Command {
	return &cli.Command{
		Name:      "latest",
		UsageText: "Latest returns the latest released version, i.e. the release with the highest number.",
		Action:    latestAction,
	}
}

func latestAction(c *cli.Context) error {
	releases, err := changelog.ParseReleased()
	if err != nil {
		return err
	}

	if len(releases.Releases) == 0 {
		return errors.New("No released versions found")
	}

	releases.SortByRelease()

	fmt.Println(releases.Releases[0].Info.Version)

	return nil
}
