package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type cmd interface {
	process(src []byte) ([]byte, error)
}

func parse(src string) (cmd, error) {
	if src[0] == 'd' {
		return &del{}, nil
	}
	switch ch, del := src[0], src[1]; ch {
	case 'a':
		_, body := until(src[2:], rune(del))
		return &appe{body}, nil
	case 'i':
		_, body := until(src[2:], rune(del))
		return &inse{body}, nil
	case 's':
		a, reg := until(src[2:], rune(del))
		b, rep := until(src[2+1+a:], rune(del))
		g := false
		if i := 2 + 1 + len(reg) + 1 + b; len(src) > i {
			g = src[i] == 'g'
		}
		return &sub{reg, rep, g}, nil
	default:
		return nil, fmt.Errorf("unknown cmd %q", ch)
	}
}

// until scans src until it founds terminator ch
// it returns the index of the terminator in the input sequence
// and a possibly unescaped body.
func until(src string, ch rune) (int, string) {
	esc := false
	var res strings.Builder
	for i, r := range src {
		if r == '\\' {
			esc = true
		} else {
			if r == ch {
				if !esc {
					return i, res.String()
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