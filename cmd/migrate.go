package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
	"gitlab.com/l0nax/changelog-go/internal/config"
)

func newMigrateCmd() *cli.Command {
	return &cli.Command{
		Name:      "migrate",
		UsageText: "Migrates the project from changelog-go v1 to v2",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force ignores the security checks and continues with the migration.",
			},
		},
		Action: migrateCmd,
	}
}

func migrateCmd(c *cli.Context) error {
	if err := loadConfig(); err == nil {
		if config.C.Version.IsValid() && !c.Bool("force") {
			return fmt.Errorf("project already migrated to v2")
		}

		// this is expected!
	}

	// security check: do not override the current config
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_ = os.Remove(filepath.Join(wd, ".changelog-go.yaml"))

	filePath := filepath.Join(wd, ".changelog-go.toml")
	if err = config.CreateDefault(filePath, true); err != nil {
		return err
	}

	if err := loadConfig(); err != nil {
		return err
	}

	// 1. Migrate released entry
	if err := migrateOldReleased(config.C.ChangelogDir); err != nil {
		return err
	}

	// 2. Migrate unreleased entries
	if err := migrateOldUnreleasedEntries(config.C.ChangelogDir); err != nil {
		return err
	}

	return nil
}

func migrateOldReleased(path string) error {
	base := filepath.Join(path, changelog.ReleasedDir)

	entries, err := os.ReadDir(base)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if err = migrateOldReleasedEntry(filepath.Join(base, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

func migrateOldUnreleasedEntries(path string) error {
	base := filepath.Join(path, changelog.UnreleasedDir)

	slog.Debug("Migrating unreleased entries", slog.String("path", base))

	entries, err := os.ReadDir(base)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		} else if entry.Name() == changelog.ReleaseInfoFileName {
			continue
		}

		if err := migrateOldChangeEntry(filepath.Join(base, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

func migrateOldReleasedEntry(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		} else if entry.Name() == changelog.ReleaseInfoFileName {
			continue
		}

		if err := migrateOldChangeEntry(filepath.Join(path, entry.Name())); err != nil {
			return err
		}
	}

	slog.Debug("Successfully migrated release",
		slog.String("release_dir_name", filepath.Base(path)))

	return nil
}

func migrateOldChangeEntry(path string) error {
	oldRaw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var entry changelog.Entry

	scanner := bufio.NewScanner(bytes.NewReader(oldRaw))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "author: "):
			entry.Author = typact.Some(strings.TrimPrefix(line, "author: "))

		case strings.HasPrefix(line, "title: "):
			entry.Title = strings.TrimPrefix(line, "title: ")

		case strings.HasPrefix(line, "type: "):
			switch strings.TrimPrefix(line, "type: ") {
			case "0":
				entry.ChangeTypeID = config.DefaultEntryNewFeatureID
			case "1":
				entry.ChangeTypeID = config.DefaultEntryBugFixID
			case "2":
				entry.ChangeTypeID = config.DefaultEntryFeatureChangeID
			case "3":
				entry.ChangeTypeID = config.DefaultEntryDeprecateID
			case "4":
				entry.ChangeTypeID = config.DefaultEntryRemovalID
			case "5":
				entry.ChangeTypeID = config.DefaultEntrySecurityID
			case "6":
				entry.ChangeTypeID = config.DefaultEntryOtherID
			default:
				slog.Error("Unknown change type: please migrate manually",
					slog.String("entry_path", path), slog.String("line", line))
			}

		default:
			slog.Warn("Changelog entry has unknown line",
				slog.String("entry_path", path), slog.String("line", line))
		}
	}

	return entry.SaveToFile(path)
}
