package changelog

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

// VersionFindMode is the mode which is used to find the next version.
type VersionFindMode uint8

const (
	VersionModeAuto VersionFindMode = iota + 1
)

type ProposedVersion struct {
	Version     string
	VersionType config.VersionEffect
}

// NextVersion returns the next version based on the mode.
func NextVersion(mode VersionFindMode) (ProposedVersion, error) {
	if mode != VersionModeAuto {
		return ProposedVersion{}, errors.Errorf("unsupported mode %q", mode)
	}

	released, err := ParseReleased()
	if err != nil {
		return ProposedVersion{}, err
	}

	if len(released.Releases) == 0 {
		return ProposedVersion{
			Version:     "0.1.0",
			VersionType: config.VersionEffectMinor,
		}, nil
	}

	released.SortByRelease()

	latest := released.Releases[0]
	lastVersion := semver.MustParse(latest.Info.Version)

	eff := latest.VersionEffect()
	switch eff {
	case config.VersionEffectMajor:
		lastVersion.IncrementMajor()

	case config.VersionEffectMinor:
		lastVersion.IncrementMinor()

	case config.VersionEffectPatch:
		lastVersion.IncrementPatch()

	default:
		panic(fmt.Sprintf("unkown version effect %q", latest.VersionEffect()))
	}

	return ProposedVersion{
		Version:     lastVersion.String(),
		VersionType: eff,
	}, nil
}
