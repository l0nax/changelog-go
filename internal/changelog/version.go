package changelog

import (
	"strings"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// VersionFindMode selects how [Project.NextVersion] derives the next version.
type VersionFindMode uint8

const (
	VersionModeAuto VersionFindMode = iota + 1
)

// ProposedVersion is the version [Project.NextVersion] suggests.
type ProposedVersion struct {
	Version     string
	VersionType config.VersionEffect
}

// ParseVersion parses a stored version string. Both "1.4.0" and "v2.0.0-rc.1"
// are accepted.
func ParseVersion(version string) (semver.Version, error) {
	ver, err := semver.ParseTolerant(version)
	if err != nil {
		return semver.Version{}, errors.Wrapf(err, "unable to parse version %q", version)
	}

	return ver, nil
}

// ApplyVersionPrefix returns version carrying prefix in place of any prefix it
// already has.
func ApplyVersionPrefix(prefix, version string) string {
	return prefix + strings.TrimPrefix(version, "v")
}

// BaseVersion returns version without its pre-release and build metadata.
func BaseVersion(version semver.Version) semver.Version {
	version.Pre = nil
	version.Build = nil

	return version
}

// versionEffect returns the strongest version effect across entries.
func versionEffect(entries []Entry) config.VersionEffect {
	// major > minor > patch, and config.VersionEffect ascends in that order,
	// so the strongest effect is the minimum.
	eff := config.VersionEffectPatch

	for i := range entries {
		effect := entries[i].ChangeType.Effect
		if !effect.IsValid() {
			continue
		}

		eff = min(eff, effect)
	}

	return eff
}

// NextVersion returns the version proposed for the entries that are still
// unreleased.
func (p *Project) NextVersion(mode VersionFindMode) (ProposedVersion, error) {
	if mode != VersionModeAuto {
		return ProposedVersion{}, errors.Errorf("unsupported mode %d", uint8(mode))
	}

	unreleased, err := p.LoadUnreleasedEntries()
	if err != nil {
		return ProposedVersion{}, err
	}

	if len(unreleased) == 0 {
		return ProposedVersion{}, errors.New("no unreleased changelog entries found")
	}

	eff := versionEffect(unreleased)

	released, err := p.ParseReleased()
	if err != nil {
		return ProposedVersion{}, err
	}

	if len(released.Releases) == 0 {
		return ProposedVersion{
			Version:     "0.1.0",
			VersionType: config.VersionEffectMinor,
		}, nil
	}

	if err := released.SortByRelease(); err != nil {
		return ProposedVersion{}, err
	}

	latest := released.Releases[0]

	version, err := ParseVersion(latest.Info.Version)
	if err != nil {
		return ProposedVersion{}, err
	}

	// A pre-release keeps its fragments, so it holds the same entries as the
	// unreleased set; bumping from it would skip a version that was never
	// released.
	if latest.Info.IsPreRelease {
		return ProposedVersion{
			Version:     BaseVersion(version).String(),
			VersionType: eff,
		}, nil
	}

	next := BaseVersion(version)

	// the Increment* methods only ever report nil
	switch eff {
	case config.VersionEffectMajor:
		_ = next.IncrementMajor()

	case config.VersionEffectMinor:
		_ = next.IncrementMinor()

	case config.VersionEffectPatch:
		_ = next.IncrementPatch()

	default:
		return ProposedVersion{}, errors.Errorf("unknown version effect %d", uint8(eff))
	}

	return ProposedVersion{
		Version:     next.String(),
		VersionType: eff,
	}, nil
}

// LatestRelease returns the released version with the highest number.
// If skipPreReleases is set, pre-releases are ignored.
func (p *Project) LatestRelease(skipPreReleases bool) (Release, error) {
	released, err := p.ParseReleased()
	if err != nil {
		return Release{}, err
	}

	if err := released.SortByRelease(); err != nil {
		return Release{}, err
	}

	for _, release := range released.Releases {
		if skipPreReleases && release.Info.IsPreRelease {
			continue
		}

		return release, nil
	}

	if skipPreReleases {
		return Release{}, errors.New("no released versions found (excluding pre-releases)")
	}

	return Release{}, errors.New("no released versions found")
}
