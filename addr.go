package main

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

type dot struct {
	low int
	hi  int
}

type addresser interface {
	address(body []byte, in dot) (dot, error)
}

type regexpAddresser struct {
	body string
}

func (r *regexpAddresser) address(body []byte, in dot) (dot, error) {
	reg, err := regexp.Compile(r.body)
	if err != nil {
		return in, err
	}
	match := reg.FindIndex(body[in.low:])
	if match == nil {
		return in, fmt.Errorf("no match for regexp")
	}

	return dot{match[0], match[1]}, nil
}

type comma struct {
	prev addresser
	next addresser
}

func (c *comma) address(body []byte, in dot) (dot, error) {
	pr, err := c.prev.address(body, in)
	if err != nil {
		return in, err
	}

	ne, err := c.next.address(body, pr)
	if err != nil {
		return in, err
	}

	return dot{pr.low, ne.hi}, nil
}

type semicolon struct {
	prev addresser
	next addresser
}

func (c *semicolon) address(body []byte, in dot) (dot, error) {
	pr, err := c.prev.address(body, in)
	if err != nil {
		return in, err
	}

	ne, err := c.next.address(body, dot{pr.hi, pr.hi})
	if err != nil {
		return in, err
	}

	return dot{pr.low, ne.hi + utf8.RuneLen(runeAt(string(body), ne.hi))}, nil
}

func parseAddr(addr string) (addresser, error) {
	ch := addr[0]
	rest := addr[1:]
	switch {
	case ch == '/':
		n, body := until(rest, '/')
		return continueAddr(addr[n+1:], &regexpAddresser{body})
	default:
		return nil, fmt.Errorf("unknown address %q", addr)
	}
}

func continueAddr(addr string, prev addresser) (addresser, error) {
	if addr == "" {
		return prev, nil
	}

	switch ch := addr[0]; ch {
	case ',':
		next, err := parseAddr(addr[1:])
		if err != nil {
			return nil, err
		}
		return &comma{prev, next}, nil
	case ';':
		next, err := parseAddr(addr[1:])
		if err != nil {
			return nil, err
		}
		return &semicolon{prev, next}, nil

	default:
		return prev, nil
	}
}
