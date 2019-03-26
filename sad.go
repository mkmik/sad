package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func run(w io.Writer, r io.Reader, src string) error {
	cmd, err := parseCmd(src)
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