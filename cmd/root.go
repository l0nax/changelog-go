package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		Name: "changelog-go",
		Authors: []*cli.Author{
			{
				Name:  "Emanuel Bennici",
				Email: "emanuel@l0nax.org",
			},
		},
		Usage: `changelog-go helps you to keep track of your
Changelog (and changes) and its fully compatible with (eg) the Git Flow.

It extends your DevOps Workflow and gives other people
the possibility to read a beautiful formatted CHANGELOG.md
file.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "level",
				Usage:   "logging level (allowed values: debug, info, warn, error)",
				Aliases: []string{"log-level", "l"},
				Value:   "info",
			},
		},
		Before: func(c *cli.Context) error {
			var lvl slog.Level
			switch strings.ToLower(c.String("level")) {
			case "error":
				lvl = slog.LevelError
			case "warn":
				lvl = slog.LevelWarn
			case "info":
				lvl = slog.LevelInfo
			case "debug":
				lvl = slog.LevelDebug
			default:
				return fmt.Errorf("unknown log level %q", c.String("level"))
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: lvl,
			}))
			slog.SetDefault(logger)

			return nil
		},
		Commands: []*cli.Command{
			newNewCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("An error occurred", err)
		os.Exit(1)
	}
}
