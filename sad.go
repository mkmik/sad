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
		return &appe{until(src[2:], rune(del))}, nil
	case 'i':
		return &inse{until(src[2:], rune(del))}, nil
	case 's':
		reg := until(src[2:], rune(del))
		rep := until(src[2+1+len(reg):], rune(del))
		g := false
		if i := 2 + 1 + len(reg) + 1 + len(rep); len(src) > i {
			g = src[i] == 'g'
		}
		return &sub{reg, rep, g}, nil
	default:
		return nil, fmt.Errorf("unknown cmd %q", ch)
	}
}

func until(src string, ch rune) string {
	i := strings.IndexRune(src, ch)
	if i == -1 {
		return src
	}
	return src[:i]
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