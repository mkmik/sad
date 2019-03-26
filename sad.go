package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type cmd interface {
	process(src []byte) ([]byte, error)
}

func runeAt(str string, i int) rune {
	return []rune(str[i:])[0]
}

func parse(src string) (cmd, error) {
	if src[0] == 'd' {
		return &del{}, nil
	}

	ch, term := src[0], runeAt(src, 1)
	rest := src[1+utf8.RuneLen(term):]

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

func run(w io.Writer, r io.Reader, src string) error {
	cmd, err := parse(src)
	if err != nil {
		return err
	}

	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	b, err := cmd.process(all)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func main() {
	flag.Parse()

	cmd := strings.Join(flag.Args(), " ")

	if err := run(os.Stdout, os.Stdin, cmd); err != nil {
		log.Fatalf("%+v", err)
	}
}