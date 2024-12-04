package common

import "github.com/muesli/termenv"

var term = termenv.ColorProfile()

// FontColor sets the color of the given string and bolds the font
func FontColor(str, color string) string {
	return termenv.String(str).
		Foreground(term.Color(color)).
		Bold().
		String()
}
