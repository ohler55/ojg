// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

const (
	blocks    = " ▏▎▍▌▋▊▉█"
	darkBlock = "▓"
	//mediumBlock = "▒"
	//lightBlock  = "░"
)

var filename = "test/patient.json"

type bench struct {
	pkg  string
	name string
	fun  func(b *testing.B)

	res    testing.BenchmarkResult
	ns     int64 // base adjusted
	bytes  int64 // base adjusted
	allocs int64 // base adjusted
}

type noWriter int

func (w noWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func main() {
	testing.Init()
	flag.Parse()
	if 0 < len(flag.Args()) {
		filename = flag.Args()[0]
	}
	gen.TimeFormat = "nano"
	/*
		benchSuite("Parse string/[]byte", []*bench{
			{pkg: "json", name: "Unmarshal", fun: goParse},
			{pkg: "oj", name: "Parse", fun: ojParse},
			{pkg: "gen", name: "Parse", fun: genParse},
			{pkg: "sen", name: "Parse", fun: senParse},
		})
		benchSuite("Parse io.Reader", []*bench{
			{pkg: "json", name: "Decode", fun: goDecodeReader},
			{pkg: "oj", name: "ParseReader", fun: ojParseReader},
			{pkg: "gen", name: "ParseReder", fun: genParseReader},
			{pkg: "sen", name: "ParseReader", fun: senParseReader},
		})
	*/
	benchSuite("Validate string/[]byte", []*bench{
		{pkg: "json", name: "Valid", fun: goValidate},
		{pkg: "oj", name: "Valdate", fun: ojValidate},
		{pkg: "oj2", name: "Valdate2", fun: ojValidate2},
	})

	benchSuite("Validate io.Reader", []*bench{
		{pkg: "json", name: "Decode", fun: goDecodeReader},
		{pkg: "oj", name: "Valdate", fun: ojValidateReader},
	})

	benchSuite("to JSON", []*bench{
		{pkg: "json", name: "Marshal", fun: marshalJSON},
		{pkg: "oj", name: "JSON", fun: ojJSON},
		{pkg: "sen", name: "String", fun: senString},
	})
	benchSuite("to JSON with indentation", []*bench{
		{pkg: "json", name: "Marshal", fun: marshalJSONIndent},
		{pkg: "oj", name: "JSON", fun: ojJSONIndent},
		{pkg: "sen", name: "String", fun: senStringIndent},
	})
	benchSuite("to JSON with indentation and sorted keys", []*bench{
		{pkg: "oj", name: "JSON", fun: ojJSONSort},
		{pkg: "sen", name: "String", fun: senStringSort},
	})

	benchSuite("Write indented JSON", []*bench{
		{pkg: "json", name: "Encode", fun: jsonEncodeIndent},
		{pkg: "oj", name: "Write", fun: ojWriteIndent},
	})

	benchSuite("Convert or Alter", []*bench{
		{pkg: "alt", name: "Generify", fun: altGenerify},
		{pkg: "alt", name: "Alter", fun: altGenAlter},
	})

	benchSuite("JSONPath Get $..a[2].c", []*bench{
		{pkg: "jp", name: "Get", fun: jpGet},
	})
	benchSuite("JSONPath First  $..a[2].c", []*bench{
		{pkg: "jp", name: "First", fun: jpFirst},
	})

	fmt.Println()
	fmt.Println(" Higher values (longer bars) are better in all cases. The bar graph compares the")
	fmt.Println(" parsing performance. The lighter colored bar is the reference, usually the go")
	fmt.Println(" json package.")
	fmt.Println()
	fmt.Println(" The Benchmarks reflect a use case where JSON is either provided as a string or")
	fmt.Println(" read from a file (io.Reader) then parsed into simple go types of nil, bool, int64")
	fmt.Println(" float64, string, []interface{}, or map[string]interface{}. When supported, an")
	fmt.Println(" io.Writer benchmark is also included along with some miscellaneous operations.")
	fmt.Println()
}

func benchSuite(title string, suite []*bench) {
	fmt.Println()
	fmt.Println(title)

	for _, b := range suite {
		b.res = testing.Benchmark(b.fun)
		b.ns = b.res.NsPerOp()
		b.bytes = b.res.AllocedBytesPerOp()
		b.allocs = b.res.AllocsPerOp()
		fmt.Printf(" %8s.%-12s %6d ns/op  %6d B/op  %6d allocs/op\n",
			b.pkg, b.name, b.ns, b.bytes, b.allocs)
	}
	fmt.Println()

	scale := 10 // TBD adjust to fit screen better?
	ss := make([]*bench, len(suite))
	copy(ss, suite)
	sort.Slice(ss, func(i, j int) bool { return ss[i].ns < ss[j].ns })
	ref := suite[0]
	for _, b := range ss {
		x := 1.0
		var bar string
		if ref == b {
			bar = strings.Repeat(darkBlock, scale)
		} else {
			x = float64(ref.ns) / float64(b.ns)
			size := x * float64(scale)
			bar = strings.Repeat(string([]rune(blocks)[8:]), int(size))
			frac := int(size*8.0) - (int(size) * 8)
			bar += string([]rune(blocks)[frac : frac+1])
		}
		fmt.Printf(" %8s %s %3.2f\n", b.pkg, bar, x)
	}
}

// data
func benchmarkData(tm time.Time) interface{} {
	return map[string]interface{}{
		"a": []interface{}{1, 2, true, tm},
		"b": 2.3,
		"c": map[string]interface{}{
			"x": "xxx",
		},
		"d": nil,
	}
}

func buildTree(size, depth, iv int) interface{} {
	if depth%2 == 0 {
		list := []interface{}{}
		for i := 0; i < size; i++ {
			nv := iv*10 + i + 1
			if 1 < depth {
				list = append(list, buildTree(size, depth-1, nv))
			} else {
				list = append(list, nv)
			}
		}
		return list
	}
	obj := map[string]interface{}{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 1 < depth {
			obj[k] = buildTree(size, depth-1, nv)
		} else {
			obj[k] = nv
		}
	}
	return obj
}

func loadSample() (data interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to load %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()

	var p oj.Parser
	if data, err = p.ParseReader(f); err != nil {
		log.Fatalf("Failed to parse %s. %s\n", filename, err)
	}
	return
}
