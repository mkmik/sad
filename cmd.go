package main

import (
	"fmt"
	"unicode/utf8"
)

type cmd interface {
	process(src []byte) ([]byte, error)
}

func parseCmd(src string) (cmd, error) {
	if src[0] == 'd' {
		return &del{}, nil
	}

	ch := src[0]
	term, lterm := utf8.DecodeRuneInString(src[1:])
	rest := src[1+lterm:]

	switch ch {
	case 'a':
		_, body := until(rest, term)
		return &appe{body}, nil
	case 'i':
		_, body := until(rest, term)
		return &inse{body}, nil
	case 's':
		a, reg := until(rest, term)
		b, rep := until(rest[a:], term)
		g := false
		if i := a + b; len(rest) > i {
			g = rest[i] == 'g'
		}
		return &sub{reg, rep, g}, nil
	default:
		return nil, fmt.Errorf("unknown cmd %q", ch)
	}
}