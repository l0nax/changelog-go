// Package common holds helpers shared by the terminal interfaces.
package common

import "github.com/muesli/termenv"

var term = termenv.ColorProfile()

// FontColor returns str in the given color, bolded.
func FontColor(str, color string) string {
	return termenv.String(str).
		Foreground(term.Color(color)).
		Bold().
		String()
}
