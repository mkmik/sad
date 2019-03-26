package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		input  string
		cmd    string
		output string
	}{
		{
			input:  "nope",
			cmd:    `s/foo/bar/`,
			output: "nope",
		},
		{
			input:  "foo",
			cmd:    `s/foo/bar/`,
			output: "bar",
		},
		{
			input:  "foofoo",
			cmd:    `s/foo/bar/`,
			output: "barfoo",
		},
		{
			input:  "foofoo",
			cmd:    `s/foo/bar/g`,
			output: "barbar",
		},
		{
			input:  "foofoo",
			cmd:    `s/(oo)/[$1]/g`,
			output: "f[oo]f[oo]",
		},
		{
			input:  "foofoo",
			cmd:    `s/(oo)/[$1]/`,
			output: "f[oo]foo",
		},
		{
			input:  "foo\nbar\nfoo\nbar",
			cmd:    `s/foo/FOO/`,
			output: "FOO\nbar\nfoo\nbar",
		},
		{
			input:  "foo\nbar\nfoo\nbar\n",
			cmd:    `s/foo\nbar/FOO/g`,
			output: "FOO\nFOO\n",
		},
		{
			input:  "foo\nbar\nfoo\nbar\n",
			cmd:    `s/foo\nbar/FOO/`,
			output: "FOO\nfoo\nbar\n",
		},
		{
			input:  "foo",
			cmd:    `s/foo/bar`,
			output: "bar",
		},
		{
			input:  "foo",
			cmd:    `a/bar`,
			output: "foobar",
		},
		{
			input:  "foo",
			cmd:    `i/bar`,
			output: "barfoo",
		},
		{
			input:  "foo",
			cmd:    `d`,
			output: "",
		},
		{
			input:  "f/o",
			cmd:    `s/\//o/`,
			output: "foo",
		},
		{
			// deviate from sam and align to modernity: only escape the delimiter
			// and preserve all escape sequences supported by golang's regexp
			input: "f/	o",
			cmd:    `s;/\t;o;`,
			output: "foo",
		},
		{ // utf-8 separators
			input:  "foo",
			cmd:    `s⌘foo⌘bar⌘`,
			output: "bar",
		},
		{
			input:  "foo",
			cmd:    `i/⌘`,
			output: "⌘foo",
		},
		{
			input:  "f⌘a",
			cmd:    `s/⌘a/oo`,
			output: "foo",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var buf strings.Builder
			if err := run(&buf, strings.NewReader(tc.input), tc.cmd); err != nil {
				t.Fatalf("%+v", err)
			}
			if got, want := buf.String(), tc.output; got != want {
				t.Errorf("got: %q, want: %q", got, want)
			}
		})
	}
}

func TestAddr(t *testing.T) {
	testCases := []struct {
		input string
		addr  string
		sel   string
	}{
		{
			input: "foo",
			addr:  ``,
			sel:   "foo",
		},
		{
			input: "foo bar",
			addr:  `/foo/`,
			sel:   "foo",
		},
		{
			input: "foo bar",
			addr:  `/f/,/b`,
			sel:   "foo b",
		},
		{
			input: "foo barf",
			addr:  `/f/;/f`,
			sel:   "foo barf",
		},
		{
			input: "⌘oo bar⌘",
			addr:  `/⌘/;/⌘`,
			sel:   "⌘oo bar⌘",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			a, err := parseAddr(tc.addr)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			sel, err := a.address([]byte(tc.input), dot{0, len(tc.input)})
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if got, want := tc.input[sel.low:sel.hi], tc.sel; got != want {
				t.Errorf("got: %q, want: %q", got, want)
			}
		})
	}
}
