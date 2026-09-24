// Package version reports the build information of the running binary.
//
// Everything here comes from [debug.ReadBuildInfo], which the Go toolchain
// fills in on its own: since Go 1.24 "go build" stamps the main module's
// version from the VCS tag, and "go install module@version" stamps the module
// version. Nothing is passed through -ldflags, so there is no set of stamps
// that can silently stop matching the code -- an -X pointing at a package that
// no longer exists is dropped by the linker without a word, which is exactly
// how the v1 stamps came to report nothing at all.
package version

import (
	"runtime/debug"
	"strings"
)

// devel is the version of a binary built outside any tagged commit.
const devel = "devel"

// Info is the build information of the running binary.
type Info struct {
	// Version is the module version, e.g. "v2.0.0". It is [devel] when the
	// toolchain could not derive one.
	Version string
	// Revision is the commit the binary was built from, if known.
	Revision string
	// Modified reports whether the working tree had uncommitted changes.
	Modified bool
	// CommitTime is the time of Revision, not the time of the build: two
	// builds of the same commit report the same value.
	CommitTime string
	// GoVersion is the toolchain that built the binary.
	GoVersion string
}

// Get returns the build information.
func Get() Info {
	build, ok := debug.ReadBuildInfo()
	if !ok {
		return Info{Version: devel}
	}

	info := Info{
		Version:   build.Main.Version,
		GoVersion: build.GoVersion,
	}

	for _, setting := range build.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.Revision = setting.Value
		case "vcs.time":
			info.CommitTime = setting.Value
		case "vcs.modified":
			info.Modified = setting.Value == "true"
		}
	}

	// "(devel)" is what the toolchain reports when there is no version to
	// derive, e.g. a build from an untagged repository.
	if info.Version == "" || info.Version == "(devel)" {
		info.Version = devel
	}

	return info
}

// String renders the build information as a single line.
func (i Info) String() string {
	var buf strings.Builder

	buf.WriteString(i.Version)

	if i.Revision != "" {
		buf.WriteString(" (")
		buf.WriteString(i.Revision)

		if i.Modified {
			buf.WriteString(", dirty")
		}

		buf.WriteString(")")
	}

	if i.CommitTime != "" {
		buf.WriteString(" committed ")
		buf.WriteString(i.CommitTime)
	}

	if i.GoVersion != "" {
		buf.WriteString(" built with ")
		buf.WriteString(i.GoVersion)
	}

	return buf.String()
}
