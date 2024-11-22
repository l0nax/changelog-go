package cmd

import (
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
	"gitlab.com/l0nax/changelog-go/internal/config"
)

func newReleaseCmd() *cli.Command {
	return &cli.Command{
		Name:      "release",
		UsageText: "Releases a new version and regenerates the changelog file",
		ArgsUsage: "<version>",
		Args:      true,
		Action:    releaseAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "pre-release",
				Usage: "If set to true, it will treat the version as a pre-release",
			},
		},
	}
}

var semverRegex = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func releaseAction(c *cli.Context) error {
	if err := loadConfig(); err != nil {
		return err
	} else if err = checkVersion(); err != nil {
		return err
	}

	released, err := changelog.ParseReleased()
	if err != nil {
		return err
	}

	rawVersion := c.Args().First()
	if rawVersion == "" {
		return errors.New("No version specified")
	}

	versionMatch := semverRegex.FindStringSubmatch(rawVersion)

	isPreRelease := c.Bool("pre-release")
	if !isPreRelease && config.C.PreRelease.Detect { // only try detection if not manually defined
		if len(versionMatch) == 0 {
			return errors.New("Pre-Release detection is enabled but version is not a valid SemVer")
		}

		// check if Version is a pre-release
		for i, group := range semverRegex.SubexpNames() {
			if group == "prerelease" {
				if versionMatch[i] != "" {
					isPreRelease = true
				}

				break
			}
		}
	}

	slog.Debug("Loading unreleased changelog entries")

	entries, err := changelog.LoadUnreleasedEntries()
	if err != nil {
		return err
	}

	slog.Info("Releasing new version",
		slog.String("version", rawVersion), slog.Bool("is_pre_release", isPreRelease),
		slog.Int("num_entries", len(entries)))

	release := changelog.Release{
		Info: changelog.ReleaseInfo{
			Version:      rawVersion,
			ReleaseDate:  time.Now(),
			IsPreRelease: isPreRelease,
		},
		Entries:  entries,
		Collapse: config.C.PreRelease.FoldPreReleases,
	}

	if err := release.Create(); err != nil {
		return err
	}

	return nil
}
