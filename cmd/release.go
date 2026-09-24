package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"gitlab.com/l0nax/changelog-go/internal/changelog"
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
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Prints the release notes that would be written, changing nothing on disk",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:  "allow-empty",
				Usage: "Releases even though there are no unreleased entries",
			},
		},
	}
}

var semverRegex = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func releaseAction(c *cli.Context) error {
	project, err := loadMigratedProject(c)
	if err != nil {
		return err
	}

	cfg := project.Config()

	rawVersion := c.Args().First()
	if rawVersion == "" {
		return errors.New("no version specified")
	}

	versionMatch := semverRegex.FindStringSubmatch(rawVersion)

	isPreRelease := c.Bool("pre-release")
	if !isPreRelease && cfg.PreRelease.Detect { // only try detection if not manually defined
		if len(versionMatch) == 0 {
			return errors.New("pre-release detection is enabled but the version is not valid SemVer")
		}

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

	entries, err := project.LoadUnreleasedEntries()
	if err != nil {
		return err
	}

	// Releasing nothing is nearly always a mistake -- the wrong branch, or
	// entries a previous release already consumed.
	if len(entries) == 0 && !c.Bool("allow-empty") {
		return nothingToDoError(
			"no unreleased entries found, so %q would be an empty release. Pass --allow-empty to release anyway",
			rawVersion)
	}

	slog.Info("Releasing new version",
		slog.String("version", rawVersion), slog.Bool("is_pre_release", isPreRelease),
		slog.Int("num_entries", len(entries)))

	release := changelog.Release{
		Info: changelog.ReleaseInfo{
			Version: rawVersion,
			// truncated: the timestamp lands in a committed file
			ReleaseDate:  time.Now().Truncate(time.Second),
			IsPreRelease: isPreRelease,
		},
		Entries: entries,
	}

	if c.Bool("dry-run") {
		return printDryRunRelease(project, release)
	}

	if err := project.CreateRelease(release); err != nil {
		return err
	}

	// before rendering, which is what gives deletion precedence over folding
	if !isPreRelease && cfg.PreRelease.DeletePreRelease {
		if err := project.RemoveSupersededPreReleases(rawVersion); err != nil {
			return err
		}
	}

	released, err := project.ParseReleased()
	if err != nil {
		return err
	}

	return released.SaveToFile(project.OutputPath())
}

// printDryRunRelease prints the notes the release would get, without touching
// the project.
func printDryRunRelease(project *changelog.Project, release changelog.Release) error {
	cl := &changelog.Changelog{
		VersionPrefix: project.Config().VersionPrefix,
		Releases:      []changelog.Release{release},
	}

	out, err := cl.RenderReleases(true)
	if err != nil {
		return err
	}

	slog.Info("Dry run, nothing was written",
		slog.String("release_dir", filepath.Join(project.ReleasedDir(), release.Info.Version)),
		slog.String("output_file", project.OutputPath()))

	fmt.Println(strings.TrimSpace(string(out)))

	return nil
}
