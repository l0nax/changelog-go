package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/config"
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
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disables color output",
				Value: false,
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

			noColor := !isatty.IsTerminal(os.Stdout.Fd()) || c.Bool("no-color")

			logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
				Level:   lvl,
				NoColor: noColor,
			}))
			slog.SetDefault(logger)

			return nil
		},
		Commands: []*cli.Command{
			newNewCmd(),
			newInitCmd(),
			newReleaseCmd(),
			newMigrateCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("An error occurred", slog.Any("error", err))
		os.Exit(1)
	}
}

func loadConfig() error {
	path, err := findConfig()
	if err != nil {
		return err
	}

	return config.Load(path)
}

func findConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	check := func(p string) (string, bool, error) {
		full := filepath.Join(p, ".changelog-go.toml")
		slog.Debug("Checking existence of config", slog.String("path", full))

		info, err := os.Stat(full)
		if err != nil {
			if os.IsNotExist(err) {
				return "", false, nil
			}

			return "", false, err
		}

		return full, !info.IsDir(), nil
	}

	path, c, err := check(cwd)
	if err != nil {
		return "", err
	} else if c {
		return path, nil
	}

	var prevPath string

	for prevPath != cwd {
		prevPath = cwd
		cwd = filepath.Join(cwd, "../")

		path, c, err := check(cwd)
		if err != nil {
			return "", err
		} else if c {
			return path, nil
		}
	}

	return "", fmt.Errorf("unable to find config: did you forget to run `changelog-go init`?")
}

func checkVersion() error {
	if !config.C.Version.IsValid() {
		return nil
	}

	return fmt.Errorf("project is not migrated to v2. Please execute 'changelog migrate'")
}
