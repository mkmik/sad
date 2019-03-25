package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/juju/errors"
)

type cmd interface {
	process(w io.Writer, r io.Reader) error
}

type sub struct {
	reg    string
	rep    string
	global bool
}

func (s *sub) process(w io.Writer, r io.Reader) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	reg, err := regexp.Compile(s.reg)
	if err != nil {
		return err
	}

	rep := []byte(s.rep)

	var res []byte
	if s.global {
		res = reg.ReplaceAll(src, rep)
	} else {
		match := reg.FindSubmatchIndex(src)
		dst := reg.Expand(nil, rep, src, match)

		res = append(res, src[:match[0]]...)
		res = append(res, dst...)
		res = append(res, src[match[1]:]...)
	}
	_, err = w.Write(res)
	return err
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