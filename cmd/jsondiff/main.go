package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/aereal/jsondiff"
	"github.com/itchyny/gojq"
)

func main() {
	os.Exit((&app{outstream: os.Stdout, errStream: os.Stderr}).run(os.Args))
}

type app struct {
	outstream io.Writer
	errStream io.Writer

	only   string
	ignore string
}

func (a *app) run(argv []string) int {
	fls := flag.NewFlagSet(argv[0], flag.ContinueOnError)
	fls.StringVar(&a.only, "only", "", "gojq query to point the structure to calculate differences")
	fls.StringVar(&a.ignore, "ignore", "", "gojq query to ignore the structure to calculate differences")
	switch err := fls.Parse(argv[1:]); err {
	case flag.ErrHelp:
		return 0
	case nil:
		// no-op
	default:
		return 1
	}
	fromPath := fls.Arg(0)
	toPath := fls.Arg(1)
	if fromPath == "" && toPath == "" {
		return a.abort("2 file paths must be passed")
	}
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return a.abort("cannot read file %s: %s", fromPath, err)
	}
	defer fromFile.Close()
	toFile, err := os.Open(toPath)
	if err != nil {
		return a.abort("cannot read file %s: %s", toPath, err)
	}
	defer toFile.Close()
	var opts []jsondiff.Option
	if a.only != "" {
		query, err := gojq.Parse(a.only)
		if err != nil {
			return a.abort("failed to parse query (%s): %s", a.only, err)
		}
		opts = append(opts, jsondiff.Only(query))
	}
	if a.ignore != "" {
		query, err := gojq.Parse(a.ignore)
		if err != nil {
			return a.abort("failed to parse query (%s): %s", a.ignore, err)
		}
		opts = append(opts, jsondiff.Ignore(query))
	}
	diff, err := jsondiff.DiffFromFiles(fromFile, toFile, opts...)
	if err != nil {
		return a.abort("cannot calculate diff: %s", err)
	}
	fmt.Fprint(a.outstream, diff)
	return 0
}

func (a *app) abort(format string, x ...interface{}) int {
	fmt.Fprintf(a.errStream, format, x...)
	fmt.Fprintln(a.errStream)
	return 1
}
