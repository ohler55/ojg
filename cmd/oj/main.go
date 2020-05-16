// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohler55/ojg/oj"
)

var indent = 2
var color = false
var bright = false

// If true wrap extracts with an array.
var wrapExtract = false

// TBD extract []oj.Expr
// TBD match []oj.Expr

func init() {
	flag.IntVar(&indent, "i", indent, "indent")
	flag.BoolVar(&color, "c", color, "color")
	flag.BoolVar(&bright, "b", bright, "bright color")
	// TBD -x extract into an array
	// TBD -m match into an array
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `

usage: %s [<options>] [@<extraction>]... [<json-file>]...

`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}
	flag.Parse()

	var files []string
	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "@") {
			// TBD an extract
		} else {
			files = append(files, arg)
		}
	}
	var p oj.Parser
	if 0 < len(files) {
		// TBD
		for _, f := range files {
			fmt.Printf("*** %s\n", f)
		}
	} else {
		p.ParseReader(os.Stdin, write)
	}
}

func write(v interface{}) bool {
	if bright {
		o := oj.BrightOptions
		o.Indent = indent
		o.Color = true
		oj.Write(os.Stdout, v, &o)
	} else if color {
		o := oj.DefaultOptions
		o.Indent = indent
		o.Color = true
		oj.Write(os.Stdout, v, &o)
	} else {
		oj.Write(os.Stdout, v, indent)
	}
	os.Stdout.Write([]byte{'\n'})
	return false
}
