// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

const version = "1.5.0"

var (
	indent = 2
	color  = false
	bright = false
	sort   = false
	lazy   = false
	senOut = false
	tab    = false

	// If true wrap extracts with an array.
	wrapExtract = false
	extracts    = []jp.Expr{}
	matches     = []*jp.Script{}
	showVersion bool
)

func init() {
	flag.IntVar(&indent, "i", indent, "indent")
	flag.BoolVar(&color, "c", color, "color")
	flag.BoolVar(&sort, "s", sort, "sort")
	flag.BoolVar(&bright, "b", bright, "bright color")
	flag.BoolVar(&wrapExtract, "w", wrapExtract, "wrap extracts in an array")
	flag.BoolVar(&lazy, "z", lazy, "lazy mode accepts Simple Encoding Notation (quotes and commas mostly optional)")
	flag.BoolVar(&senOut, "sen", senOut, "output in Simple Encoding Notation")
	flag.BoolVar(&tab, "t", tab, "indent with tabs")
	flag.Var(&exValue{}, "x", "extract path")
	flag.Var(&matchValue{}, "m", "match equation/script")
	flag.BoolVar(&showVersion, "version", showVersion, "display version and exit")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
usage: %s [<options>] [@<extraction>]... [(<match>)]... [<json-file>]...

The default behavior it to write the JSON formatted according to the color
options and the indentation option. If no files are specified JSON input is
expected from stdin.

Filtering and extraction of elements is supported using JSONPath and the
scripting that is part of JSONPath filters.

Extraction paths can be provided either with the -x option or an argument
starting with a $ or @. A Expr.Get() is executed and all the results are
either written or wrapped with an array and written depending on the value of
the wrap option (-w).

  oj -x abc.def myfile.json "@.x[?(@.y > 1)]"

To filter JSON documents the match option (-m) is used. If a JSON document
matches at least one match option the JSON will be written. In addition to the
-m option an argument starting with a '(' is assumed to be a match script that
follows the oj.Script format.

  oj -m "(@.name == 'Pete')" myfile.json "(@.name == "Makie")"

An argument that starts with a { or [ marks the start of a JSON document that
is composed of the remaining argument concatenated together. That document is
then used as the input.

  oj -i 0 -z {a:1, b:two}
  => {"a":1,"b":"two"}

`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("oj %s\n", version)
		os.Exit(0)
	}
	var input []byte
	var files []string
	for _, arg := range flag.Args() {
		if len(arg) == 0 {
			continue
		}
		if 0 < len(input) {
			input = append(input, arg...)
			continue
		}
		switch arg[0] {
		case '@', '$':
			x, err := jp.ParseString(arg)
			if err == nil {
				extracts = append(extracts, x)
			}
		case '(':
			script, err := jp.NewScript(arg)
			if err == nil {
				matches = append(matches, script)
			}
		case '{', '[':
			input = append(input, arg...)
		default:
			files = append(files, arg)
		}
	}
	var p oj.SimpleParser
	var err error
	if lazy {
		p = &sen.Parser{}
	} else {
		p = &oj.Parser{Reuse: true}
	}
	if 0 < len(files) {
		var f *os.File
		for _, file := range files {
			if f, err = os.Open(file); err == nil {
				_, err = p.ParseReader(f, write)
				f.Close()
			}
			if err != nil {
				break
			}
		}
	}
	if 0 < len(input) {
		_, err = p.Parse(input, write)
	}
	if len(files) == 0 && len(input) == 0 {
		_, err = p.ParseReader(os.Stdin, write)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "*-*-* %s\n", err)
	}
}

func write(v interface{}) bool {
	if 0 < len(matches) {
		match := false
		for _, m := range matches {
			if m.Match(v) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	if 0 < len(extracts) {
		if wrapExtract {
			var w []interface{}
			for _, x := range extracts {
				w = append(w, x.Get(v)...)
			}
			if senOut {
				writeSEN(w)
			} else {
				writeJSON(w)
			}
		} else {
			for _, x := range extracts {
				for _, v2 := range x.Get(v) {
					if senOut {
						writeSEN(v2)
					} else {
						writeJSON(v2)
					}
				}
			}
		}
	} else if senOut {
		writeSEN(v)
	} else {
		writeJSON(v)
	}
	return false
}

func writeJSON(v interface{}) {
	if bright {
		o := oj.BrightOptions
		o.Indent = indent
		o.Color = true
		o.Sort = sort
		o.Tab = tab
		_ = oj.Write(os.Stdout, v, &o)
	} else if color || sort || tab {
		o := oj.DefaultOptions
		o.Indent = indent
		o.Color = color
		o.Sort = sort
		o.Tab = tab
		_ = oj.Write(os.Stdout, v, &o)
	} else {
		_ = oj.Write(os.Stdout, v, indent)
	}
	os.Stdout.Write([]byte{'\n'})
}

func writeSEN(v interface{}) {
	if bright {
		o := sen.BrightOptions
		o.Indent = indent
		o.Color = true
		o.Sort = sort
		o.Tab = tab
		_ = sen.Write(os.Stdout, v, &o)
	} else if color || sort || tab {
		o := sen.DefaultOptions
		o.Indent = indent
		o.Color = color
		o.Sort = sort
		o.Tab = tab
		_ = sen.Write(os.Stdout, v, &o)
	} else {
		_ = sen.Write(os.Stdout, v, indent)
	}
	os.Stdout.Write([]byte{'\n'})
}

type exValue struct {
}

func (xv exValue) String() string {
	return ""
}

func (xv exValue) Set(s string) error {
	x, err := jp.ParseString(s)
	if err == nil {
		extracts = append(extracts, x)
	}
	return err
}

type matchValue struct {
}

func (mv matchValue) String() string {
	return ""
}

func (mv matchValue) Set(s string) error {
	script, err := jp.NewScript(s)
	if err == nil {
		matches = append(matches, script)
	}
	return err
}
