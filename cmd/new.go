package cmd

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
	"gitlab.com/l0nax/changelog-go/internal/config"
	"gitlab.com/l0nax/changelog-go/internal/tui/create"
)

func newNewCmd() *cli.Command {
	return &cli.Command{
		Name:      "new",
		UsageText: "Creates a new Changelog-Entry so you can easily commit your entry",
		ArgsUsage: "<title>",
		Args:      true,
		Action:    newAction,
	}
}

func newAction(c *cli.Context) error {
	if err := loadConfig(); err != nil {
		return err
	}

	askTitle := c.Args().First() == ""

	input, err := create.Run(askTitle)
	if err != nil {
		if errors.Is(err, create.ErrCanceled) {
			slog.Info("Operation canceled")

			return nil
		}

		return err
	}

	title := c.Args().First()
	if title == "" {
		title = input.Title
	}

	entry := changelog.Entry{
		ChangeTypeID: input.Type.ID,
		Title:        title,
		Author:       typact.None[string](), // TODO: Support author
	}

	dir := filepath.Join(config.C.ChangelogDir, "unreleased")
	if err = os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	path := filepath.Join(dir, strconv.FormatInt(time.Now().UnixMilli(), 10))
	slog.Debug("Saving new entry in file", slog.String("file_path", path))

	if err = entry.SaveToFile(path); err != nil {
		return err
	}

	return nil
}
