package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

func newInitCmd() *cli.Command {
	return &cli.Command{
		Name:      "init",
		UsageText: "Initializes changelog-go at the current CWD",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Forces the creation of a new config file",
			},
		},
		Action: initAction,
	}
}

func initAction(c *cli.Context) error {
	// security check: do not override the current config
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(wd, ConfigFileName)
	if _, err := os.Stat(filePath); err == nil && !c.Bool("force") {
		slog.Error("A config file already exists!")

		return nil
	}

	if err := config.CreateDefault(filePath, c.Bool("force")); err != nil {
		return err
	}

	cfg, err := config.Load(filePath)
	if err != nil {
		return err
	}

	project := changelog.NewProject(cfg, wd)

	// created up front so that "release" and "next" work on a fresh project
	for _, dir := range []string{project.ReleasedDir(), project.UnreleasedDir()} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	slog.Info("Initialized changelog-go",
		slog.String("config", filePath), slog.String("changelog_dir", project.ChangelogDir()))

	return nil
}
