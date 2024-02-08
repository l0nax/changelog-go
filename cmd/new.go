package cmd

import "github.com/urfave/cli/v2"

func newNewCmd() *cli.Command {
	return &cli.Command{
		Name: "new",
		UsageText: "creates a new Changelog-Entry so you can easily commit your entry",
		ArgsUsage: "<title>",
		Args: true,
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
