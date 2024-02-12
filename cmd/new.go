package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/tui/create"
)

func newNewCmd() *cli.Command {
	return &cli.Command{
		Name:      "new",
		UsageText: "creates a new Changelog-Entry so you can easily commit your entry",
		ArgsUsage: "<title>",
		Args:      true,
		Action:    newAction,
	}
}

func newAction(c *cli.Context) error {
	if err := loadConfig(); err != nil {
		return err
	}

	input, err := create.Run()
	if err != nil {
		return err
	}

	fmt.Printf("* Got: %+v\n", input)

	return nil
}
