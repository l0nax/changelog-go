package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
)

func newNextCmd() *cli.Command {
	return &cli.Command{
		Name: "next",
		UsageText: `Next returns the next version.

The "auto" modes parses all unreleased entries and chooses the most appropriate next version.`,
		Usage:     "[flags...] <args...>",
		ArgsUsage: "<auto>",
		Args:      true,
		Action:    nextAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "type-only",
				Usage: "If set to true, only the version type is printed, instead of the version. This can be useful with semantic-release",
			},
		},
	}
}

func nextAction(c *cli.Context) error {
	if err := loadConfig(); err != nil {
		return err
	}

	rawArg := c.Args().First()
	if rawArg != "auto" {
		return errors.New("Please specify a valid mode")
	}

	res, err := changelog.NextVersion(changelog.VersionModeAuto)
	if err != nil {
		return err
	}

	if c.Bool("type-only") {
		fmt.Println(res.VersionType.String())
	} else {
		fmt.Println(res.Version)
	}

	return nil
}
