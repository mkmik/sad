package main

import (
	"regexp"
)

type sub struct {
	reg    string
	rep    string
	global bool
}

func (s *sub) process(src []byte) ([]byte, error) {
	reg, err := regexp.Compile(s.reg)
	if err != nil {
		return nil, err
	}

	rep := []byte(s.rep)

	if s.global {
		return reg.ReplaceAll(src, rep), nil
	}

	match := reg.FindSubmatchIndex(src)
	dst := reg.Expand(nil, rep, src, match)

	var res []byte

	res = append(res, src[:match[0]]...)
	res = append(res, dst...)
	res = append(res, src[match[1]:]...)
	return res, nil
}
