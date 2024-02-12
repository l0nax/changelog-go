package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

func newInitCmd() *cli.Command {
	return &cli.Command{
		Name:      "init",
		UsageText: "Initializes changelog-go at the current CWD",
		Action:    initAction,
	}
}

func initAction(c *cli.Context) error {
	// security check: do not override the current config
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(wd, ".changelog-go.yaml")
	if _, err := os.Stat(filePath); err == nil {
		slog.Error("A config file already exists!")

		return nil
	}

	return config.CreateDefault(filePath)
}
