package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
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
			&cli.BoolFlag{
				Name:  "remove-prefix",
				Usage: "Prints the bare version, without the configured version prefix",
			},
		},
	}
}

func nextAction(c *cli.Context) error {
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	rawArg := c.Args().First()
	if rawArg != "auto" {
		return errors.New("please specify a valid mode")
	}

	res, err := project.NextVersion(changelog.VersionModeAuto)
	if err != nil {
		return err
	}

	if c.Bool("type-only") {
		fmt.Println(res.VersionType.String())

		return nil
	}

	fmt.Println(printableVersion(project, res.Version, c.Bool("remove-prefix")))

	return nil
}

// printableVersion returns version with the configured prefix applied, or bare
// when removePrefix is set.
func printableVersion(project *changelog.Project, version string, removePrefix bool) string {
	if removePrefix {
		return changelog.ApplyVersionPrefix("", version)
	}

	return project.DisplayVersion(version)
}
