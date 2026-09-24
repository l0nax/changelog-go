package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.l0nax.org/typact"
	"gopkg.in/yaml.v3"

	"gitlab.com/l0nax/changelog-go/v2/internal/changelog"
	"gitlab.com/l0nax/changelog-go/v2/internal/config"
)

// LegacyConfigFileName is the name of the changelog-go v1 config file.
const LegacyConfigFileName = ".changelog-go.yaml"

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

// legacyConfig mirrors the changelog-go v1 configuration file.
//
// v1 had no user-defined change types; a type was an integer from 0 to 6 in the
// entry files. [legacyChangeTypeIDs] maps those onto the v2 type IDs.
type legacyConfig struct {
	Version    string `yaml:"version"`
	PreRelease struct {
		Detect           bool `yaml:"detect"`
		DeletePreRelease bool `yaml:"deletePreRelease"`
		FoldPreReleases  bool `yaml:"foldPreReleases"`
	} `yaml:"preRelease"`
	Entry struct {
		Author bool `yaml:"author"`
	} `yaml:"entry"`
	Changelog struct {
		EntryPath    string `yaml:"entryPath"`
		Changelog    string `yaml:"changelog"`
		CustomScheme bool   `yaml:"customScheme"`
	} `yaml:"changelog"`
}

func migrateCmd(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return migrateProject(wd, c.Bool("force"))
}

// migrateProject migrates the project rooted at wd from v1 to v2.
func migrateProject(wd string, force bool) error {
	if !force {
		if cfg, err := config.Load(filepath.Join(wd, ConfigFileName)); err == nil && cfg.Version.IsValid() {
			return fmt.Errorf("project already migrated to v2")
		}
	}

	legacyPath := filepath.Join(wd, LegacyConfigFileName)

	cfg, err := migrateConfig(legacyPath)
	if err != nil {
		return err
	}

	configPath := filepath.Join(wd, ConfigFileName)

	if err := writeConfig(configPath, cfg); err != nil {
		return err
	}

	_ = os.Remove(legacyPath)

	project := changelog.NewProject(cfg, wd)

	if err := migrateOldReleased(project.ReleasedDir()); err != nil {
		return err
	}

	if err := migrateOldUnreleasedEntries(project.UnreleasedDir()); err != nil {
		return err
	}

	slog.Info("Migrated project to changelog-go v2", slog.String("config", configPath))

	return nil
}

// migrateConfig returns the v2 config for the v1 config file at legacyPath.
//
// The settings v1 shares with v2 are carried over, and the "other" change type
// that v1 knew only as the numeric type 6 is added.
func migrateConfig(legacyPath string) (config.Config, error) {
	cfg, err := config.DefaultProjectConfig()
	if err != nil {
		return config.Config{}, err
	}

	cfg.Entry.Types = append(cfg.Entry.Types, config.ChangeType{
		ID:         config.DefaultEntryOtherID,
		Title:      "Other change",
		GroupTitle: "Other",
		// a catch-all should not force a minor bump
		Effect: config.VersionEffectPatch,
	})

	raw, err := os.ReadFile(legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("No v1 config file found, migrating with the default settings",
				slog.String("path", legacyPath))

			return cfg, cfg.Validate()
		}

		return config.Config{}, errors.Wrapf(err, "unable to read v1 config file %q", legacyPath)
	}

	var legacy legacyConfig
	if err := yaml.Unmarshal(raw, &legacy); err != nil {
		return config.Config{}, errors.Wrapf(err, "unable to parse v1 config file %q", legacyPath)
	}

	cfg.PreRelease.Detect = legacy.PreRelease.Detect
	cfg.PreRelease.DeletePreRelease = legacy.PreRelease.DeletePreRelease
	cfg.PreRelease.FoldPreReleases = legacy.PreRelease.FoldPreReleases

	if legacy.Changelog.EntryPath != "" {
		cfg.ChangelogDir = legacy.Changelog.EntryPath
	}

	if legacy.Changelog.Changelog != "" {
		cfg.OutputPath = typact.Some(legacy.Changelog.Changelog)
	}

	if legacy.Changelog.CustomScheme {
		slog.Warn("The v1 'changelog.customScheme' setting has no v2 equivalent yet and was dropped")
	}

	if legacy.Entry.Author {
		slog.Warn("The v1 'entry.author' setting has no v2 equivalent yet and was dropped")
	}

	return cfg, cfg.Validate()
}

// writeConfig renders cfg and writes it to path.
func writeConfig(path string, cfg config.Config) error {
	raw, err := config.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to render the migrated config")
	}

	return os.WriteFile(path, raw, 0644)
}

// migrateOldReleased migrates every release directory below base.
func migrateOldReleased(base string) error {
	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if err := migrateOldReleasedEntry(filepath.Join(base, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

// migrateOldUnreleasedEntries migrates every entry file in base.
func migrateOldUnreleasedEntries(base string) error {
	slog.Debug("Migrating unreleased entries", slog.String("path", base))

	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

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

// migrateOldReleasedEntry migrates the release directory at path.
func migrateOldReleasedEntry(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		} else if entry.Name() == changelog.ReleaseInfoFileName {
			infoPath := filepath.Join(path, entry.Name())

			err = migrateReleaseInfoFile(infoPath)
			if err != nil {
				return errors.Wrapf(err, "unable to migrate ReleaseInfo file at %q", infoPath)
			}

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

// migrateReleaseInfoFile rewrites the v1 ReleaseInfo file at path as TOML.
func migrateReleaseInfoFile(path string) error {
	slog.Debug("Migrating release info file", slog.String("path", path))

	oldRaw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var entry changelog.ReleaseInfo

	scanner := bufio.NewScanner(bytes.NewReader(oldRaw))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "releasedate: "):
			rawStr := strings.TrimPrefix(line, "releasedate: ")
			rawStr = strings.TrimSpace(rawStr)
			rawStr = strings.Trim(rawStr, `"`)

			oldDate, err := time.Parse("2006-01-02", rawStr)
			if err != nil {
				return errors.Wrapf(err, "unable to parse 'releasedate': %q", rawStr)
			}

			entry.ReleaseDate = oldDate

		case strings.HasPrefix(line, "isprerelease: "):
			rawStr := strings.TrimPrefix(line, "isprerelease: ")
			rawStr = strings.TrimSpace(rawStr)

			bl, err := strconv.ParseBool(rawStr)
			if err != nil {
				return errors.Wrapf(err, "unable to parse 'isprerelease': %q", rawStr)
			}

			entry.IsPreRelease = bl

		default:
			slog.Warn("Release Info entry has unknown line",
				slog.String("entry_path", path), slog.String("line", line))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	entry.Version = filepath.Base(filepath.Dir(path))

	return entry.SaveToFile(path)
}

// migrateOldChangeEntry rewrites the v1 changelog entry at path as TOML.
func migrateOldChangeEntry(path string) error {
	oldRaw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

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
			rawType := strings.TrimPrefix(line, "type: ")

			id, ok := legacyChangeTypeIDs[rawType]
			if !ok {
				slog.Error("Unknown change type: please migrate manually",
					slog.String("entry_path", path), slog.String("line", line))

				continue
			}

			entry.ChangeTypeID = id

		default:
			slog.Warn("Changelog entry has unknown line",
				slog.String("entry_path", path), slog.String("line", line))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return entry.SaveToFile(path)
}

// legacyChangeTypeIDs maps the numeric change types of v1 onto the v2 type IDs.
var legacyChangeTypeIDs = map[string]string{
	"0": config.DefaultEntryNewFeatureID,
	"1": config.DefaultEntryBugFixID,
	"2": config.DefaultEntryFeatureChangeID,
	"3": config.DefaultEntryDeprecateID,
	"4": config.DefaultEntryRemovalID,
	"5": config.DefaultEntrySecurityID,
	"6": config.DefaultEntryOtherID,
}
