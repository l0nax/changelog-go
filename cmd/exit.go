package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Exit codes. These are part of the tool's interface: a pipeline may branch on
// them, so they only ever gain new values, never new meanings.
const (
	// ExitOK is returned when the command succeeded.
	ExitOK = 0
	// ExitError is returned for any failure without a more specific code,
	// e.g. unreadable files or a malformed config.
	ExitError = 1
	// ExitUsage is returned when the invocation itself is wrong: a missing
	// argument, or a value the config does not know such as a change type ID.
	ExitUsage = 2
	// ExitNotFound is returned when the requested version does not exist, or
	// nothing has been released at all.
	ExitNotFound = 3
	// ExitNothingToDo is returned when the command found no work: releasing
	// with no unreleased entries and without --allow-empty.
	ExitNothingToDo = 4
)

// usageError reports a wrong invocation, i.e. [ExitUsage].
func usageError(format string, a ...any) error {
	return cli.Exit(fmt.Sprintf(format, a...), ExitUsage)
}

// notFoundError reports a missing version, i.e. [ExitNotFound].
func notFoundError(format string, a ...any) error {
	return cli.Exit(fmt.Sprintf(format, a...), ExitNotFound)
}

// nothingToDoError reports that there was no work, i.e. [ExitNothingToDo].
func nothingToDoError(format string, a ...any) error {
	return cli.Exit(fmt.Sprintf(format, a...), ExitNothingToDo)
}
