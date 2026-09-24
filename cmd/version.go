package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/version"
)

func newVersionCmd() *cli.Command {
	return &cli.Command{
		Name:      "version",
		UsageText: "Version prints the version of changelog-go itself.",
		Action:    versionAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "short",
				Usage: "Prints only the version, without the commit and build time",
			},
		},
	}
}

func versionAction(c *cli.Context) error {
	info := version.Get()

	if c.Bool("short") {
		fmt.Println(info.Version)

		return nil
	}

	fmt.Println(info.String())

	return nil
}
