package changelog

import (
	"strings"
	"time"
)

func formatTime(tt time.Time) string {
	return tt.Format("2006-01-02")
}

// indentBody indents every line of an entry body by two spaces.
//
// That indent is what keeps a multi-line body inside the list item it belongs
// to; without it Markdown ends the list at the first unindented line.
func indentBody(body string) string {
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			// no trailing whitespace on an otherwise blank line
			lines[i] = ""

			continue
		}

		lines[i] = "  " + line
	}

	return strings.Join(lines, "\n")
}
