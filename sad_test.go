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
		/*	{
				input:  "f/o",
				cmd:    `s/\//o/`,
				output: "foo",
			},
		*/
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