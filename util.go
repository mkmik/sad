package main

import (
	"strings"
	"unicode/utf8"
)

// runeAt returns the rune at byte index i.
func runeAt(str string, i int) rune {
	r, _ := utf8.DecodeRuneInString(str[i:])
	return r
}

// until scans src until it founds terminator ch
// it returns the length of the match in the input sequence, including the terminator.
// It also returns a possibly unescaped body.
func until(src string, ch rune) (int, string) {
	esc := false
	var res strings.Builder
	for i, r := range src {
		if r == '\\' {
			esc = true
		} else {
			if r == ch {
				if !esc {
					return i + utf8.RuneLen(r), res.String()
				}
			}
			if esc {
				esc = false
				res.WriteRune('\\')
			}
			res.WriteRune(r)
		}
	}
	return len(src), res.String()
}