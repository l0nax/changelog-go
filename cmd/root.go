// Package cmd implements the changelog-go command line interface.
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
	"gitlab.com/l0nax/changelog-go/internal/config"
)

// ConfigFileName is the name of the project configuration file.
const ConfigFileName = ".changelog-go.toml"

// Run executes the command line interface.
func Run() {
	app := &cli.App{
		Name: "changelog",
		Authors: []*cli.Author{
			{
				Name:  "Emanuel Bennici",
				Email: "emanuel@l0nax.org",
			},
		},
		Usage: `Keep a changelog without merge conflicts.

Every change gets its own file under .changelogs/unreleased/, committed on the
branch that made it, so two branches never conflict over the changelog. On
release these are collected into a CHANGELOG.md and filed under the version
they shipped in.`,
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
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to the config file, skipping the upward search",
				Aliases: []string{"c"},
				EnvVars: []string{"CHANGELOG_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "changelog-dir",
				Usage:   "Overrides the configured changelog directory",
				EnvVars: []string{"CHANGELOG_DIR"},
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

			// stderr keeps stdout free for the value a command prints,
			// e.g. $(changelog latest)
			noColor := !isatty.IsTerminal(os.Stderr.Fd()) || c.Bool("no-color")

			logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
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
			newNextCmd(),
			newLatestCmd(),
			newShowCmd(),
			newVersionCmd(),
			newCheckCmd(),
			newValidateCmd(),
			newCompletionCmd(),
			newManCmd(),
		},
		EnableBashCompletion: true,
		// Errors are reported here so that every command exits through the
		// documented codes in exit.go.
		ExitErrHandler: func(_ *cli.Context, err error) {},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(exitCodeOf(err))
	}
}

// exitCodeOf returns the documented exit code err asks for, defaulting to
// [ExitError].
func exitCodeOf(err error) int {
	var coder cli.ExitCoder
	if errors.As(err, &coder) {
		return coder.ExitCode()
	}

	return ExitError
}

// loadProject returns the project described by the nearest config file.
func loadProject(c *cli.Context) (*changelog.Project, error) {
	path := c.String("config")

	if path == "" {
		found, err := findConfig()
		if err != nil {
			return nil, err
		}

		path = found
	}

	cfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}

	rootDir := filepath.Dir(path)

	// An explicit --changelog-dir is what the user typed, so it is relative to
	// the working directory rather than to the config file.
	if override := c.String("changelog-dir"); override != "" {
		abs, err := filepath.Abs(override)
		if err != nil {
			return nil, err
		}

		cfg.ChangelogDir = abs
	}

	return changelog.NewProject(cfg, rootDir), nil
}

// loadMigratedProject returns the project, or an error if it still uses the v1
// layout.
func loadMigratedProject(c *cli.Context) (*changelog.Project, error) {
	project, err := loadProject(c)
	if err != nil {
		return nil, err
	}

	if !project.Config().Version.IsValid() {
		return nil, fmt.Errorf("project is not migrated to v2 (got %q). Please execute 'changelog migrate'",
			project.Config().Version)
	}

	return project, nil
}

// findConfig returns the path of the nearest config file, walking up from the
// working directory.
func findConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	check := func(p string) (string, bool, error) {
		full := filepath.Join(p, ConfigFileName)
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
		cwd = filepath.Dir(cwd)

		path, c, err := check(cwd)
		if err != nil {
			return "", err
		} else if c {
			return path, nil
		}
	}

	return "", fmt.Errorf("unable to find config: did you forget to run `changelog init`?")
}
