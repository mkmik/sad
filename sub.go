package main

import (
	"io"
	"io/ioutil"
	"regexp"
)

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