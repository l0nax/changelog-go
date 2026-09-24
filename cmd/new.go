package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
	"gitlab.com/l0nax/changelog-go/internal/config"
	"gitlab.com/l0nax/changelog-go/internal/tui/create"
)

func newNewCmd() *cli.Command {
	return &cli.Command{
		Name: "new",
		UsageText: `New creates a changelog entry so you can commit it alongside your change.

Given both --type and --title it runs without any terminal interaction, which is
what makes it usable from CI, a git hook or an editor plugin:

    changelog new -t bug_fix --title "Fix the sprocket"`,
		ArgsUsage: "<title>",
		Args:      true,
		Action:    newAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Usage:   "ID of the change type, as configured under [[entry.types]]",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:  "title",
				Usage: "Title of the entry, taking precedence over the positional argument",
			},
			&cli.StringFlag{
				Name:    "body",
				Usage:   "Optional long-form description, rendered underneath the title",
				Aliases: []string{"b"},
			},
			&cli.StringFlag{
				Name:  "author",
				Usage: "Author to record on the entry",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Prints the entry that would be written, without writing it",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Usage:   "Prompts for anything not given as a flag (defaults to false when stdin is not a TTY)",
				Aliases: []string{"i"},
				Value:   true,
			},
		},
	}
}

func newAction(c *cli.Context) error {
	project, err := loadProject(c)
	if err != nil {
		return err
	}

	cfg := project.Config()

	title := c.String("title")
	if title == "" {
		title = c.Args().First()
	}

	changeType, err := selectedChangeType(c, cfg)
	if err != nil {
		return err
	}

	// Prompt only for what is still missing, and only if we may.
	if changeType.IsNone() || title == "" {
		if !isInteractive(c) {
			return usageError("%s. Pass them as flags or run interactively",
				strings.Join(missingInputs(changeType, title), " and "))
		}

		input, err := create.Run(cfg, create.Options{
			PreselectedType: changeType,
			AskTitle:        title == "",
		})
		if err != nil {
			if errors.Is(err, create.ErrCanceled) {
				slog.Info("Operation canceled")

				return nil
			}

			return err
		}

		changeType = typact.Some(input.Type)

		if title == "" {
			title = input.Title
		}
	}

	entry := changelog.Entry{
		ChangeTypeID: changeType.UnwrapOrZero().ID,
		Title:        title,
		Body:         optionalFlag(c, "body"),
		Author:       optionalFlag(c, "author"),
	}

	dir := project.UnreleasedDir()
	path := filepath.Join(dir, strconv.FormatInt(time.Now().UnixMilli(), 10))

	if c.Bool("dry-run") {
		raw, err := entry.Marshal()
		if err != nil {
			return err
		}

		slog.Info("Would write changelog entry", slog.String("file_path", path))
		fmt.Print(string(raw))

		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	slog.Debug("Saving new entry in file", slog.String("file_path", path))

	return entry.SaveToFile(path)
}

// selectedChangeType resolves --type against the configured types.
func selectedChangeType(c *cli.Context, cfg config.Config) (typact.Option[config.ChangeType], error) {
	id := c.String("type")
	if id == "" {
		return typact.None[config.ChangeType](), nil
	}

	changeType, ok := cfg.ResolveType(id)
	if !ok {
		return typact.None[config.ChangeType](),
			usageError("unknown change type %q. Configured types: %s", id, configuredTypeIDs(cfg))
	}

	return typact.Some(changeType), nil
}

// configuredTypeIDs lists the type IDs an entry may use.
func configuredTypeIDs(cfg config.Config) string {
	ids := make([]string, 0, len(cfg.Entry.Types))
	for _, typ := range cfg.Entry.Types {
		ids = append(ids, typ.ID)
	}

	return strings.Join(ids, ", ")
}

// missingInputs names what the caller still has to supply.
func missingInputs(changeType typact.Option[config.ChangeType], title string) []string {
	var missing []string

	if changeType.IsNone() {
		missing = append(missing, "no change type given (--type)")
	}

	if title == "" {
		missing = append(missing, "no title given (--title)")
	}

	return missing
}

// isInteractive reports whether the command may prompt.
//
// Prompting is off whenever stdin is not a terminal, so that a CI job, a git
// hook or a test never hangs waiting for input it cannot give. An explicit
// --interactive overrides the detection.
func isInteractive(c *cli.Context) bool {
	if c.IsSet("interactive") {
		return c.Bool("interactive")
	}

	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

// optionalFlag returns the flag value as an option, absent when unset.
func optionalFlag(c *cli.Context, name string) typact.Option[string] {
	if value := c.String(name); value != "" {
		return typact.Some(value)
	}

	return typact.None[string]()
}
