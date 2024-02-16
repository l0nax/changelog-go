package cmd

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
)

func newReleaseCmd() *cli.Command {
	return &cli.Command{
		Name:      "release",
		UsageText: "Releases a new version and regenerates the changelog file",
		ArgsUsage: "<version>",
		Args:      true,
		Action:    releaseAction,
	}
}

func releaseAction(c *cli.Context) error {
	if err := loadConfig(); err != nil {
		return err
	} else if err = checkVersion(); err != nil {
		return err
	}

	released, err := changelog.ParseReleased()
	if err != nil {
		return err
	}

	_ = released
	_ = c

	return nil
}
