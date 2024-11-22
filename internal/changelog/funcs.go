package changelog

import "time"

func formatTime(tt time.Time) string {
	return tt.Format("2006-01-02")
}
