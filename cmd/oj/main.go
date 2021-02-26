// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
)

const version = "1.7.0"

var (
	indent     = 2
	color      = false
	bright     = false
	sortKeys   = false
	lazy       = false
	senOut     = false
	tab        = false
	showFnDocs = false
	showConf   = false
	safe       = false

	// If true wrap extracts with an array.
	wrapExtract = false
	extracts    = []jp.Expr{}
	matches     = []*jp.Script{}
	planDef     = ""
	showVersion bool
	plan        *asm.Plan
	root        = map[string]interface{}{}
	showRoot    bool
	edgeDepth   = 0.0
	html        = false
	confFile    = ""
)

func init() {
	flag.IntVar(&indent, "i", indent, "indent")
	flag.BoolVar(&color, "c", color, "color")
	flag.BoolVar(&sortKeys, "s", sortKeys, "sort")
	flag.BoolVar(&bright, "b", bright, "bright color")
	flag.BoolVar(&wrapExtract, "w", wrapExtract, "wrap extracts in an array")
	flag.BoolVar(&lazy, "z", lazy, "lazy mode accepts Simple Encoding Notation (quotes and commas mostly optional)")
	flag.BoolVar(&senOut, "sen", senOut, "output in Simple Encoding Notation")
	flag.BoolVar(&tab, "t", tab, "indent with tabs")
	flag.Var(&exValue{}, "x", "extract path")
	flag.Var(&matchValue{}, "m", "match equation/script")
	flag.BoolVar(&showVersion, "version", showVersion, "display version and exit")
	flag.StringVar(&planDef, "a", planDef, "assembly plan or plan file using @<plan>")
	flag.BoolVar(&showRoot, "r", showRoot, "print root if an assemble plan provided")
	flag.Float64Var(&edgeDepth, "p", edgeDepth, "pretty print with the edge and depth as a float <edge>.<max-depth>")
	flag.BoolVar(&html, "html", html, "output colored output as HTML")
	flag.BoolVar(&safe, "safe", safe, "escape &, <, and > for HTML inclusion")
	flag.StringVar(&confFile, "f", confFile, "configuration file (see -help-config)")
	flag.BoolVar(&showFnDocs, "fn", showFnDocs, "describe assembly plan functions")
	flag.BoolVar(&showFnDocs, "help-fn", showFnDocs, "describe assembly plan functions")
	flag.BoolVar(&showConf, "help-config", showConf, "describe .oj-config.sen format")
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

Oj can also be used to assemble new JSON output from input data. An assembly
plan that describes how to assemble the new JSON if specified by the -a
option. The -fn option will display the documentation for assembly.

Pretty mode output can be used with JSON or the -sen option. It indents
according to a defined edge and maximum depth in a best effort approach. The
-p takes and edge and a maximum depth as the numerator and fractional part of
a decimal.

`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse() // get config file if specified
	if showVersion {
		fmt.Printf("oj %s\n", version)
		os.Exit(0)
	}
	if showConf {
		displayConf()
		os.Exit(0)
	}
	if showFnDocs {
		displayFnDocs()
		os.Exit(0)
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "*-*-* %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()
	loadConfig()

	flag.Parse() // load again to over-ride loaded config

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
	if lazy {
		p = &sen.Parser{}
	} else {
		p = &oj.Parser{Reuse: true}
	}
	planDef = strings.TrimSpace(planDef)
	if 0 < len(planDef) {
		if planDef[0] != '[' {
			var b []byte
			if b, err = ioutil.ReadFile(planDef); err != nil {
				fmt.Fprintf(os.Stderr, "*-*-* %s\n", err)
				os.Exit(1)
			}
			planDef = string(b)
		}
		var pd interface{}
		if pd, err = (&sen.Parser{}).Parse([]byte(planDef)); err != nil {
			panic(err)
		}
		plist, _ := pd.([]interface{})
		if len(plist) == 0 {
			panic(fmt.Errorf("assemble plan not an array"))
		}
		plan = asm.NewPlan(plist)
	}
	// TBD define writer of correct type and use for write
	//   take advantage of buffer reuse
	if 0 < len(files) {
		var f *os.File
		for _, file := range files {
			if f, err = os.Open(file); err == nil {
				_, err = p.ParseReader(f, write)
				f.Close()
			}
			if err != nil {
				panic(err)
			}
		}
	}
	if 0 < len(input) {
		if _, err = p.Parse(input, write); err != nil {
			panic(err)
		}
	}
	if len(files) == 0 && len(input) == 0 {
		if _, err = p.ParseReader(os.Stdin, write); err != nil {
			panic(err)
		}
	}
	if showRoot && plan != nil {
		plan = nil
		delete(root, "src")
		delete(root, "asm")
		write(root)
	}
	return
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
		if plan != nil {
			root["src"] = v
			if err := plan.Execute(root); err != nil {
				fmt.Fprintf(os.Stderr, "*-*-* %s\n", err)
				os.Exit(1)
			} else {
				v = root["asm"]
			}
		}
		writeJSON(v)
	}
	return false
}

func writeJSON(v interface{}) {
	var o oj.Options
	if bright {
		o = oj.BrightOptions
		o.Color = true
		o.Sort = sortKeys
		o.Tab = tab
	} else if color || sortKeys || tab {
		o = oj.DefaultOptions
		o.Color = color
		o.Sort = sortKeys
		o.Tab = tab
	}
	o.Indent = indent
	o.HTMLUnsafe = !safe
	if html {
		o.HTMLUnsafe = false
		if color {
			o.SyntaxColor = sen.HTMLOptions.SyntaxColor
			o.KeyColor = sen.HTMLOptions.KeyColor
			o.NullColor = sen.HTMLOptions.NullColor
			o.BoolColor = sen.HTMLOptions.BoolColor
			o.NumberColor = sen.HTMLOptions.NumberColor
			o.StringColor = sen.HTMLOptions.StringColor
			o.NoColor = sen.HTMLOptions.NoColor
		}
	}
	if 0.0 < edgeDepth {
		_ = pretty.WriteJSON(os.Stdout, v, &o, edgeDepth)
	} else {
		_ = oj.Write(os.Stdout, v, &o)
	}
	os.Stdout.Write([]byte{'\n'})
}

func writeSEN(v interface{}) {
	var o sen.Options
	switch {
	case html:
		o = sen.HTMLOptions
		o.Color = true
		o.Sort = sortKeys
		o.Tab = tab
		o.HTMLSafe = true
	case bright:
		o = sen.BrightOptions
		o.Color = true
		o.Sort = sortKeys
		o.Tab = tab
	case color || sortKeys || tab:
		o = sen.DefaultOptions
		o.Color = color
		o.Sort = sortKeys
		o.Tab = tab
	}
	o.Indent = indent
	o.HTMLSafe = safe
	if 0.0 < edgeDepth {
		_ = pretty.WriteSEN(os.Stdout, v, &o, edgeDepth)
	} else {
		_ = sen.Write(os.Stdout, v, &o)
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

func loadConfig() {
	var conf interface{}
	if 0 < len(confFile) {
		f, err := os.Open(confFile)
		if err != nil {
			panic(err)
		}
		if conf, err = sen.ParseReader(f); err != nil {
			panic(err)
		}
		applyConf(conf)
	}
	home := os.Getenv("HOME")
	for _, path := range []string{
		"./.oj-config.sen",
		"./.oj-config.json",
		home + "/.oj-config.sen",
		home + "/.oj-config.json",
	} {
		f, err := os.Open(path)
		if err == nil {
			if conf, err = sen.ParseReader(f); err == nil {
				applyConf(conf)
				return
			}
		}
	}
}

func applyConf(conf interface{}) {
	bright, _ = jp.C("bright").First(conf).(bool)
	color, _ = jp.C("color").First(conf).(bool)
	for _, v := range jp.C("format").C("indent").Get(conf) {
		indent = int(alt.Int(v))
	}
	for _, v := range jp.C("format").C("tab").Get(conf) {
		tab = alt.Bool(v)
	}
	for _, v := range jp.C("format").C("pretty").Get(conf) {
		edgeDepth = alt.Float(v)
	}
	safe, _ = jp.C("html-safe").First(conf).(bool)
	lazy, _ = jp.C("lazy").First(conf).(bool)
	senOut, _ = jp.C("sen").First(conf).(bool)

	setOptionsColor(conf, "bool", setBoolColor)
	setOptionsColor(conf, "key", setKeyColor)
	setOptionsColor(conf, "no-color", setNoColor)
	setOptionsColor(conf, "null", setNullColor)
	setOptionsColor(conf, "number", setNumberColor)
	setOptionsColor(conf, "string", setStringColor)
	setOptionsColor(conf, "syntax", setSyntaxColor)

	setHTMLColor(conf, "bool", &sen.HTMLOptions.BoolColor)
	setHTMLColor(conf, "key", &sen.HTMLOptions.KeyColor)
	setHTMLColor(conf, "no-color", &sen.HTMLOptions.NoColor)
	setHTMLColor(conf, "null", &sen.HTMLOptions.NullColor)
	setHTMLColor(conf, "number", &sen.HTMLOptions.NumberColor)
	setHTMLColor(conf, "string", &sen.HTMLOptions.StringColor)
	setHTMLColor(conf, "syntax", &sen.HTMLOptions.SyntaxColor)
}

func setOptionsColor(conf interface{}, key string, fun func(color string)) {
	for _, v := range jp.C("colors").C(key).Get(conf) {
		fun(pickColor(alt.String(v)))
	}
}

func setBoolColor(color string) {
	oj.DefaultOptions.BoolColor = color
	oj.BrightOptions.BoolColor = color
	sen.DefaultOptions.BoolColor = color
	sen.BrightOptions.BoolColor = color
}

func setKeyColor(color string) {
	oj.DefaultOptions.KeyColor = color
	oj.BrightOptions.KeyColor = color
	sen.DefaultOptions.KeyColor = color
	sen.BrightOptions.KeyColor = color
}

func setNoColor(color string) {
	oj.DefaultOptions.NoColor = color
	oj.BrightOptions.NoColor = color
	sen.DefaultOptions.NoColor = color
	sen.BrightOptions.NoColor = color
}

func setNullColor(color string) {
	oj.DefaultOptions.NullColor = color
	oj.BrightOptions.NullColor = color
	sen.DefaultOptions.NullColor = color
	sen.BrightOptions.NullColor = color
}

func setNumberColor(color string) {
	oj.DefaultOptions.NumberColor = color
	oj.BrightOptions.NumberColor = color
	sen.DefaultOptions.NumberColor = color
	sen.BrightOptions.NumberColor = color
}

func setStringColor(color string) {
	oj.DefaultOptions.StringColor = color
	oj.BrightOptions.StringColor = color
	sen.DefaultOptions.StringColor = color
	sen.BrightOptions.StringColor = color
}

func setSyntaxColor(color string) {
	oj.DefaultOptions.SyntaxColor = color
	oj.BrightOptions.SyntaxColor = color
	sen.DefaultOptions.SyntaxColor = color
	sen.BrightOptions.SyntaxColor = color
}

func setHTMLColor(conf interface{}, key string, sp *string) {
	for _, v := range jp.C("colors").C(key).Get(conf) {
		*sp = pickColor(alt.String(v))
	}
}

func pickColor(s string) (color string) {
	switch strings.ToLower(s) {
	case "normal":
		color = "\x1b[m"
	case "black":
		color = "\x1b[30m"
	case "red":
		color = "\x1b[31m"
	case "green":
		color = "\x1b[32m"
	case "yellow":
		color = "\x1b[33m"
	case "blue":
		color = "\x1b[34m"
	case "magenta":
		color = "\x1b[35m"
	case "cyan":
		color = "\x1b[36m"
	case "white":
		color = "\x1b[37m"
	case "gray":
		color = "\x1b[90m"
	case "bright-red":
		color = "\x1b[91m"
	case "bright-green":
		color = "\x1b[92m"
	case "bright-yellow":
		color = "\x1b[93m"
	case "bright-blue":
		color = "\x1b[94m"
	case "bright-magenta":
		color = "\x1b[95m"
	case "bright-cyan":
		color = "\x1b[96m"
	case "bright-white":
		color = "\x1b[97m"
	default:
		panic(fmt.Errorf("%s is not a valid color choice", s))
	}
	return
}

func displayFnDocs() {
	fmt.Printf(`
An assembly plan is described by a JSON document or a SEN document. The format
is much like LISP but with brackets instead of parenthesis. A plan is
evaluated by evaluating the plan function which is usually an 'asm'
function. The plan operates on a data map which is the root during
evaluation. The source data is in the $.src and the expected assembled output
should be in $.asm.

An example of a plan in SEN format is (the first asm is optional):

  [ asm
    [set $.asm {good: bye}]
    [set $.asm.hello world]
  ]

The functions available are:

`)
	var b []byte
	var keys []string
	docs := asm.FnDocs()
	for k := range docs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b = append(b, fmt.Sprintf("  %10s: %s\n\n", k, strings.ReplaceAll(docs[k], "\n", "\n              "))...)
	}
	fmt.Println(string(b))
}

func displayConf() {
	fmt.Printf(`
If an oj configuration file is present in the local directory or the home
directory that file is used to set the defaults for oj. The file can be in
either SEN or JSON format. The paths check, in order are:

  ./.oj-config.sen
  ./.oj-config.json
  ~/.oj-config.sen
  ~/.oj-config.json

The file format (SEN with comments) is:

{
  bright: true // Color if true will colorize the output with bright colors.
  color: false // Color if true will colorize the output. The bright option takes precedence.
  colors: {
    // Color values can be one of the following:
    //   normal
    //   black
    //   red
    //   green
    //   yellow
    //   blue
    //   magenta
    //   cyan
    //   white
    //   gray
    //   bright-red
    //   bright-green
    //   bright-yellow
    //   bright-blue
    //   bright-magenta
    //   bright-cyan
    //   bright-white
    syntax: normal
    key: bright-blue
    null: bright-red
    bool: bright-yellow
    number: bright-cyan
    string: bright-green
    no-color: normal // NoColor turns the color off.
  }
  format: {indent: 2 tab: false pretty: 80.3}
  html: {
    syntax: "<span>"
    key: "<span style=\"color:#44f\">"
    null: "<span style=\"color:red\">"
    bool: "<span style=\"color:#a40\">"
    number: "<span style=\"color:#04a\">"
    string: "<span style=\"color:green\">"
    no-color: "</span>"
  }
  html-safe: false
  lazy: true // -z option, lazy read for SEN format
  sen: true
}
`)
}
