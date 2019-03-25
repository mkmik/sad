package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/juju/errors"
)

type cmd interface {
	process(w io.Writer, r io.Reader) error
}

func parse(src string) (cmd, error) {
	switch ch, del := src[0], src[1]; ch {
	case 's':
		reg, err := until(src[2:], rune(del))
		if err != nil {
			return nil, errors.Trace(err)
		}
		rep, err := until(src[2+1+len(reg):], rune(del))
		if err != nil {
			return nil, errors.Trace(err)
		}
		g := false
		if i := 2 + 1 + len(reg) + 1 + len(rep); len(src) > i {
			g = src[i] == 'g'
		}
		return &sub{reg, rep, g}, nil
	default:
		return nil, fmt.Errorf("unknown cmd %q", ch)
	}
}

func until(src string, ch rune) (string, error) {
	i := strings.IndexRune(src, ch)
	if i == -1 {
		return "", fmt.Errorf("undelimited cmd")
	}
	return src[:i], nil
}

func run(w io.Writer, r io.Reader, src string) error {
	cmd, err := parse(src)
	if err != nil {
		return err
	}

	return cmd.process(w, r)
}

func main() {
	flag.Parse()

	cmd := strings.Join(flag.Args(), " ")

	if err := run(os.Stdout, os.Stdin, cmd); err != nil {
		log.Fatalf("%+v", err)
	}
}