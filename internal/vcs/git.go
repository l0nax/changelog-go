// Package vcs is the thin layer over the git commands the tool needs.
package vcs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Runner runs a git command inside dir and returns its standard output.
//
// It is a function rather than a hard call to exec so that the logic built on
// top of it can be tested without a repository to set up.
type Runner func(dir string, args ...string) ([]byte, error)

// Git runs the real git binary.
func Git(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		// git says what went wrong on stderr; the exec error alone is just
		// "exit status 128".
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}

		return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), message)
	}

	return out, nil
}

// Available reports whether the git binary is on PATH.
func Available() bool {
	_, err := exec.LookPath("git")

	return err == nil
}

// Lines splits git output into its non-empty lines.
func Lines(out []byte) []string {
	raw := strings.Split(strings.TrimSpace(string(out)), "\n")

	lines := make([]string, 0, len(raw))

	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}
