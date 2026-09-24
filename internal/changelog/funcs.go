package changelog

import "time"

// formatTime renders tt the way release headings show a date.
func formatTime(tt time.Time) string {
	return tt.Format("2006-01-02")
}
