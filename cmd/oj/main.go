// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
)

var version = "unknown"

var (
	indent         = 2
	color          = false
	bright         = false
	sortKeys       = false
	lazy           = false
	senOut         = false
	tab            = false
	showFnDocs     = false
	showFilterDocs = false
	showConf       = false
	safe           = false
	mongo          = false
	omit           = false
	dig            = false
	annotate       = false

	// If true wrap extracts with an array.
	wrapExtract = false
	extracts    = []jp.Expr{}
	matches     = []*jp.Script{}
	dels        = []jp.Expr{}
	planDef     = ""
	showVersion bool
	plan        *asm.Plan
	root        = map[string]any{}
	showRoot    bool
	prettyOpt   = ""
	width       = 80
	maxDepth    = 3
	prettyOn    = false
	align       = false
	html        = false
	convName    = ""
	confFile    = ""

	conv    *alt.Converter
	options *ojg.Options
)

func init() {
	flag.IntVar(&indent, "i", indent, "indent")
	flag.BoolVar(&color, "c", color, "color")
	flag.BoolVar(&sortKeys, "s", sortKeys, "sort")
	flag.BoolVar(&bright, "b", bright, "bright color")
	flag.BoolVar(&omit, "o", omit, "omit nil and empty")
	flag.BoolVar(&wrapExtract, "w", wrapExtract, "wrap extracts in an array")
	flag.BoolVar(&lazy, "z", lazy, "lazy mode accepts Simple Encoding Notation (quotes and commas mostly optional)")
	flag.BoolVar(&senOut, "sen", senOut, "output in Simple Encoding Notation")
	flag.BoolVar(&tab, "t", tab, "indent with tabs")
	flag.BoolVar(&annotate, "annotate", annotate, "annotate dig extracts with a path comment")
	flag.Var(&exValue{}, "x", "extract path")
	flag.Var(&matchValue{}, "m", "match equation/script")
	flag.Var(&delValue{}, "d", "delete path")
	flag.BoolVar(&dig, "dig", dig, "dig into a large document using the tokenizer")
	flag.BoolVar(&showVersion, "version", showVersion, "display version and exit")
	flag.StringVar(&planDef, "a", planDef, "assembly plan or plan file using @<plan>")
	flag.BoolVar(&showRoot, "r", showRoot, "print root if an assemble plan provided")
	flag.StringVar(&prettyOpt, "p", prettyOpt,
		`pretty print with the width, depth, and align as <width>.<max-depth>.<align>`)
	flag.BoolVar(&html, "html", html, "output colored output as HTML")
	flag.BoolVar(&safe, "safe", safe, "escape &, <, and > for HTML inclusion")
	flag.StringVar(&confFile, "f", confFile, "configuration file (see -help-config), - indicates no file")
	flag.BoolVar(&showFnDocs, "fn", showFnDocs, "describe assembly plan functions")
	flag.BoolVar(&showFnDocs, "help-fn", showFnDocs, "describe assembly plan functions")
	flag.BoolVar(&showFilterDocs, "help-filter", showFilterDocs, "describe filter operators like [?(@.x == 3)]")
	flag.BoolVar(&showConf, "help-config", showConf, "describe .oj-config.sen format")
	flag.BoolVar(&mongo, "mongo", mongo, "parse mongo Javascript output")
	flag.StringVar(&convName, "conv", convName, `apply converter before writing. Supported values are:
  nano - converts integers over 946684800000000000 (2000-01-01) to time
  rcf3339 - converts string in RFC3339 or RFC3339Nano to time
  mongo - converts mongo wrapped values e.g.,  {$numberLong: "123"} => 123
  <with-numbers> - if digits are included then time layout is assumed
  <other> - any other is taken to be a key in a map with a string or nano time
`)
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

Elements can be deleted from the JSON using the -d option. Multiple
occurrences of -d are supported.

Oj can also be used to assemble new JSON output from input data. An assembly
plan that describes how to assemble the new JSON if specified by the -a
option. The -fn option will display the documentation for assembly.

Pretty mode output can be used with JSON or the -sen option. It indents
according to a defined width and maximum depth in a best effort approach. The
-p takes a pattern of <width>.<max-depth>.<align> where width and max-depth
are integers and align is a boolean.

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
	if showFilterDocs {
		displayFilterDocs()
		os.Exit(0)
	}
	extracts = extracts[:0]
	matches = matches[:0]
	dels = dels[:0]
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
	if 0 < len(convName) {
		switch strings.ToLower(convName) {
		case "nano":
			conv = &alt.TimeNanoConverter
		case "rfc3339":
			conv = &alt.TimeRFC3339Converter
		case "mongo":
			conv = &alt.MongoConverter
		default:
			if strings.ContainsAny(convName, "0123456789") {
				conv = &alt.Converter{
					String: []func(val string) (any, bool){
						func(val string) (any, bool) {
							if len(val) == len(convName) {
								if t, err := time.ParseInLocation(convName, val, time.UTC); err == nil {
									return t, true
								}
							}
							return val, false
						},
					},
				}
			} else {
				conv = &alt.Converter{
					Map: []func(val map[string]any) (any, bool){
						func(val map[string]any) (any, bool) {
							if len(val) == 1 {
								switch tv := val[convName].(type) {
								case string:
									for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"} {
										if t, err := time.ParseInLocation(layout, tv, time.UTC); err == nil {
											return t, true
										}
									}
								case int64:
									return time.Unix(0, tv), true
								}
							}
							return val, false
						},
					},
				}
			}
		}
	}
	var p oj.SimpleParser
	switch {
	case mongo:
		sp := &sen.Parser{}
		sp.AddMongoFuncs()
		p = sp
		if conv == nil {
			conv = &alt.MongoConverter
		}
	case lazy:
		p = &sen.Parser{}
	default:
		p = &oj.Parser{Reuse: true}
	}
	planDef = strings.TrimSpace(planDef)
	if 0 < len(planDef) {
		if planDef[0] != '[' {
			var b []byte
			if b, err = ioutil.ReadFile(planDef); err != nil {
				return err
			}
			planDef = string(b)
		}
		var pd any
		if pd, err = (&sen.Parser{}).Parse([]byte(planDef)); err != nil {
			panic(err)
		}
		plist, _ := pd.([]any)
		if len(plist) == 0 {
			panic(fmt.Errorf("assembly plan not an array"))
		}
		plan = asm.NewPlan(plist)
	}
	if 0 < len(files) {
		var f *os.File
		for _, file := range files {
			if f, err = os.Open(file); err == nil {
				if dig {
					err = digParse(f)
				} else {
					_, err = p.ParseReader(f, write)
				}
				_ = f.Close()
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
		if dig {
			err = digParse(os.Stdin)
		} else {
			_, err = p.ParseReader(os.Stdin, write)
		}
		if err != nil {
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

func digParse(r io.Reader) error {
	var fn func(path jp.Expr, data any)
	annotateColor := ""

	if color {
		annotateColor = ojg.Gray
	}
	// Pick a function that satisfies omit, annotate, and senOut
	// values. Determining the function before the actual calling means few
	// conditional paths during the repeated calls later.
	if omit {
		if annotate {
			if senOut {
				fn = func(path jp.Expr, data any) {
					if data != nil && data != "" {
						fmt.Printf("%s// %s\n", annotateColor, path)
						writeSEN(data)
					}
				}
			} else {
				fn = func(path jp.Expr, data any) {
					if data != nil && data != "" {
						fmt.Printf("%s// %s\n", annotateColor, path)
						writeJSON(data)
					}
				}
			}
		} else {
			if senOut {
				fn = func(path jp.Expr, data any) {
					if data != nil && data != "" {
						writeSEN(data)
					}
				}
			} else {
				fn = func(path jp.Expr, data any) {
					if data != nil && data != "" {
						writeJSON(data)
					}
				}
			}
		}
	} else {
		if annotate {
			if senOut {
				fn = func(path jp.Expr, data any) {
					fmt.Printf("%s// %s\n", annotateColor, path)
					writeSEN(data)
				}
			} else {
				fn = func(path jp.Expr, data any) {
					fmt.Printf("%s// %s\n", annotateColor, path)
					writeJSON(data)
				}
			}
		} else {
			if senOut {
				fn = func(path jp.Expr, data any) {
					writeSEN(data)
				}
			} else {
				fn = func(path jp.Expr, data any) {
					writeJSON(data)
				}
			}
		}
	}
	if lazy {
		return sen.MatchLoad(r, fn, extracts...)
	}
	return oj.MatchLoad(r, fn, extracts...)
}

func write(v any) bool {
	if conv != nil {
		v = conv.Convert(v)
	}
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
	for _, x := range dels {
		_ = x.Del(v)
	}
	switch {
	case 0 < len(extracts):
		if wrapExtract {
			var w []any
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
	case senOut:
		writeSEN(v)
	default:
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

func writeJSON(v any) {
	if options == nil {
		o := ojg.Options{}
		if bright {
			o = oj.BrightOptions
			o.Color = true
			o.Sort = sortKeys
		} else if color || sortKeys || tab {
			o = ojg.DefaultOptions
			o.Color = color
		}
		o.Indent = indent
		o.Tab = tab
		o.HTMLUnsafe = !safe
		o.TimeFormat = time.RFC3339Nano
		o.Sort = sortKeys
		if html {
			o.HTMLUnsafe = false
			if color {
				o.SyntaxColor = ojg.HTMLOptions.SyntaxColor
				o.KeyColor = ojg.HTMLOptions.KeyColor
				o.NullColor = ojg.HTMLOptions.NullColor
				o.BoolColor = ojg.HTMLOptions.BoolColor
				o.NumberColor = ojg.HTMLOptions.NumberColor
				o.StringColor = ojg.HTMLOptions.StringColor
				o.TimeColor = ojg.HTMLOptions.TimeColor
				o.NoColor = ojg.HTMLOptions.NoColor
			}
		}
		options = &o
	}
	if omit {
		// Use alt.Alter to remove empty since it handles recursive removal.
		v = alt.Alter(v, &ojg.Options{OmitNil: true, OmitEmpty: true})
	}
	if 0 < len(prettyOpt) {
		parsePrettyOpt()
	}
	if prettyOn {
		_ = pretty.WriteJSON(os.Stdout, v, options, float64(width)+float64(maxDepth)/10.0, align)
	} else {
		_ = oj.Write(os.Stdout, v, options)
	}
	_, _ = os.Stdout.Write([]byte{'\n'})
}

func writeSEN(v any) {
	if options == nil {
		o := ojg.Options{}
		switch {
		case html:
			o = ojg.HTMLOptions
			o.Color = true
			o.HTMLUnsafe = false
		case bright:
			o = ojg.BrightOptions
			o.Color = true
		case color || sortKeys || tab:
			o = ojg.DefaultOptions
			o.Color = color
		}
		o.Indent = indent
		o.Tab = tab
		o.HTMLUnsafe = !safe
		o.TimeFormat = time.RFC3339Nano
		o.Sort = sortKeys
		options = &o
	}
	if omit {
		// Use alt.Alter to remove empty since it handles recursive removal.
		v = alt.Alter(v, &ojg.Options{OmitNil: true, OmitEmpty: true})
	}
	if 0 < len(prettyOpt) {
		parsePrettyOpt()
	}
	if prettyOn {
		_ = pretty.WriteSEN(os.Stdout, v, options, float64(width)+float64(maxDepth)/10.0, align)
	} else {
		_ = sen.Write(os.Stdout, v, options)
	}
	_, _ = os.Stdout.Write([]byte{'\n'})
}

func parsePrettyOpt() {
	if 0 < len(prettyOpt) {
		parts := strings.Split(prettyOpt, ".")
		if 0 < len(parts[0]) {
			if i, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
				width = int(i)
				prettyOn = true
			} else {
				panic(err)
			}
		}
		if 1 < len(parts) && 0 < len(parts[1]) {
			if i, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				maxDepth = int(i)
				prettyOn = true
			} else {
				panic(err)
			}
		}
		if 2 < len(parts) && 0 < len(parts[2]) {
			var err error
			if align, err = strconv.ParseBool(parts[2]); err != nil {
				panic(err)
			}
			prettyOn = true
		}
	}
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

type delValue struct {
}

func (dv delValue) String() string {
	return ""
}

func (dv delValue) Set(s string) error {
	x, err := jp.ParseString(s)
	if err == nil {
		dels = append(dels, x)
	}
	return err
}

func loadConfig() {
	var conf any
	if 0 < len(confFile) {
		if confFile == "-" { // special case
			return
		}
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

func applyConf(conf any) {
	bright, _ = jp.C("bright").First(conf).(bool)
	color, _ = jp.C("color").First(conf).(bool)
	for _, v := range jp.C("format").C("indent").Get(conf) {
		indent = int(alt.Int(v))
	}
	for _, v := range jp.C("format").C("tab").Get(conf) {
		tab = alt.Bool(v)
	}
	for _, v := range jp.C("format").C("pretty").Get(conf) {
		prettyOpt, _ = v.(string)
		parsePrettyOpt()
	}
	for _, v := range jp.C("format").C("width").Get(conf) {
		width = int(alt.Int(v))
		prettyOn = true
	}
	for _, v := range jp.C("format").C("depth").Get(conf) {
		maxDepth = int(alt.Int(v))
		prettyOn = true
	}
	for _, v := range jp.C("format").C("align").Get(conf) {
		align = alt.Bool(v)
		prettyOn = true
	}
	safe, _ = jp.C("html-safe").First(conf).(bool)
	lazy, _ = jp.C("lazy").First(conf).(bool)
	senOut, _ = jp.C("sen").First(conf).(bool)
	convName, _ = jp.C("conv").First(conf).(string)
	mongo, _ = jp.C("mongo").First(conf).(bool)

	setOptionsColor(conf, "bool", setBoolColor)
	setOptionsColor(conf, "key", setKeyColor)
	setOptionsColor(conf, "no-color", setNoColor)
	setOptionsColor(conf, "null", setNullColor)
	setOptionsColor(conf, "number", setNumberColor)
	setOptionsColor(conf, "string", setStringColor)
	setOptionsColor(conf, "time", setTimeColor)
	setOptionsColor(conf, "syntax", setSyntaxColor)

	setHTMLColor(conf, "bool", &sen.HTMLOptions.BoolColor)
	setHTMLColor(conf, "key", &sen.HTMLOptions.KeyColor)
	setHTMLColor(conf, "no-color", &sen.HTMLOptions.NoColor)
	setHTMLColor(conf, "null", &sen.HTMLOptions.NullColor)
	setHTMLColor(conf, "number", &sen.HTMLOptions.NumberColor)
	setHTMLColor(conf, "string", &sen.HTMLOptions.StringColor)
	setHTMLColor(conf, "syntax", &sen.HTMLOptions.SyntaxColor)
}

func setOptionsColor(conf any, key string, fun func(color string)) {
	for _, v := range jp.C("colors").C(key).Get(conf) {
		fun(pickColor(alt.String(v)))
	}
}

func setBoolColor(color string) {
	ojg.DefaultOptions.BoolColor = color
	ojg.BrightOptions.BoolColor = color
}

func setKeyColor(color string) {
	ojg.DefaultOptions.KeyColor = color
	ojg.BrightOptions.KeyColor = color
}

func setNoColor(color string) {
	ojg.DefaultOptions.NoColor = color
	ojg.BrightOptions.NoColor = color
}

func setNullColor(color string) {
	ojg.DefaultOptions.NullColor = color
	ojg.BrightOptions.NullColor = color
}

func setNumberColor(color string) {
	ojg.DefaultOptions.NumberColor = color
	ojg.BrightOptions.NumberColor = color
}

func setStringColor(color string) {
	ojg.DefaultOptions.StringColor = color
	ojg.BrightOptions.StringColor = color
}

func setTimeColor(color string) {
	ojg.DefaultOptions.TimeColor = color
	ojg.BrightOptions.TimeColor = color
}

func setSyntaxColor(color string) {
	ojg.DefaultOptions.SyntaxColor = color
	ojg.BrightOptions.SyntaxColor = color
}

func setHTMLColor(conf any, key string, sp *string) {
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
    [set $.asm {good: bye}]  // set output to {good: bye}
    [set $.asm.hello world]  // output is now {good: bye, hello: world}
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

func displayFilterDocs() {
	fmt.Printf(`

JSONPaths can include filters such as $.x[?(@.y == 'z')].value. As with other
square bracket operators it applies to arrays. The general form of a filter is
[?(left operator right)]. Both left and right can be constants or JSONPaths
where @ is each array element. Nested filter are supported. Operators
supported are:

 ==    returns true if left is equal to right.

 !=    returns true if left is not equal to right.

 <     returns true if left is less than right.

 <=    returns true if left is less than or equal to right.

 >     returns true if left is greater than right.

 >=    returns true if left is greater than or equal to right.

 ||    returns true if either left or right is true

 &&    returns true if both left and right are true.

 !     inverts the boolean value of the right. No left should be
       present. Examples are !@.x or !(@.x == 2).

 empty returns true if the left empty condition (length is zero) matches the
       right which must be a boolean.

 has   returns true if the left has condition is null or missing matches the
       right which must be a boolean.

 +     returns the sum of left and right.

 -     returns the difference of left and right. (left - right)

 *     returns the product of left and right.

 /     returns left divided by right.

 in    returns true if left is in right. Right must be an array either as a
       constant of the form [1,'a'] or as a path that evaluates to an array.

 =~    returns true if left is a string and matches the right regex which can be
       either a regex delimited by / or a string.

Functions are also support and take the for of [?length(@.x) == 3]. The
supported functions are:

 length(path)        returns the length of the list, object, or string at the
                     path. If the element does not exist or is not a list,
                     object, or string then Nothing is returned.

 count(path)         returns the number of elements that match the path which
                     should return a node list.

 match(path, regex)  the path should return a string which is then compared to
                     the regex string. If there is a match to on the entirety of
                     the string at path then true is returned otherwise if the
                     string does not match false is returned. I the value at
                     path is not a string or does not exist then Nothing is
                     returned.

 search(path, regex) the path should return a string which is then compared to
                     the regex string. If there is a match to on a substring of
                     the string at path then true is returned otherwise if the
                     string does not match false is returned. I the value at
                     path is not a string or does not exist then Nothing is
                     returned.

`)
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
    time: bright-magenta
    no-color: normal // NoColor turns the color off.
  }
  // Either the pretty element can be used or the individual width, depth, and
  // align options can be specified separately.
  format: {indent: 2 tab: false pretty: 80.3.false}
  // format: {indent: 2 tab: false width: 80 depth: 3 align: false}
  html: {
    syntax: "<span>"
    key: '<span style="color:#44f">'
    null: '<span style="color:red">'
    bool: '<span style="color:#a40">"
    number: '<span style="color:#04a">'
    string: '<span style="color:green">'
    time: '<span style="color:#f0f">'
    no-color: "</span>"
  }
  html-safe: false
  lazy: true // -z option, lazy read for SEN format
  sen: true
  conv: rfc3339
  mongo: false
}
`)
}
