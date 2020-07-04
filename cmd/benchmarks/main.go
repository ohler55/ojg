// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

type bench struct {
	pkg  string
	name string
	fun  func(b *testing.B)

	res    testing.BenchmarkResult
	ns     int64 // base adjusted
	bytes  int64 // base adjusted
	allocs int64 // base adjusted
}

func main() {
	testing.Init()
	flag.Parse()
	gen.TimeFormat = "nano"

	fmt.Println()
	fmt.Println(" The number in parenthesis are the ratio of results between the reference and")
	fmt.Println(" the listed. Higher values are better.")

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

	benchSuite("Validate string/[]byte", []*bench{
		{pkg: "json", name: "Valid", fun: goValidate},
		{pkg: "oj", name: "Valdate", fun: ojValidate},
	})

	benchSuite("Validate io.Reader", []*bench{
		{pkg: "json", name: "Decode", fun: goDecodeReader},
		{pkg: "oj", name: "Valdate", fun: ojValidateReader},
	})

	benchSuite("to JSON", []*bench{
		{pkg: "json", name: "Marshal", fun: marshalJSON},
		{pkg: "oj", name: "JSON", fun: ojJSON},
		{pkg: "oj", name: "Write", fun: ojWrite},
		{pkg: "sen", name: "String", fun: senString},
	})
	benchSuite("to JSON with indentation", []*bench{
		{pkg: "json", name: "Marshal", fun: marshalJSONIndent},
		{pkg: "oj", name: "JSON", fun: ojJSONIndent},
		{pkg: "oj", name: "Write", fun: ojWriteIndent},
		{pkg: "sen", name: "String", fun: senStringIndent},
	})
	benchSuite("to JSON with indentation and sorted keys", []*bench{
		{pkg: "oj", name: "JSON", fun: ojJSONSort},
		{pkg: "oj", name: "Write", fun: ojWriteSort},
		{pkg: "sen", name: "String", fun: senStringSort},
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
}

func benchSuite(title string, suite []*bench) {
	fmt.Println()
	fmt.Println(title)

	for _, b := range suite {
		b.res = testing.Benchmark(b.fun)
		b.ns = b.res.NsPerOp()
		b.bytes = b.res.AllocedBytesPerOp()
		b.allocs = b.res.AllocsPerOp()
		fmt.Printf(" %4s.%-12s %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
			b.pkg, b.name,
			b.ns, float64(suite[0].ns)/float64(b.ns),
			b.bytes, float64(suite[0].bytes)/float64(b.bytes),
			b.allocs, float64(suite[0].allocs)/float64(b.allocs))
	}
}

// Parse functions
func goParse(b *testing.B) {
	var result interface{}
	for n := 0; n < b.N; n++ {
		_ = json.Unmarshal([]byte(sampleJSON), &result)
	}
}

func ojParse(b *testing.B) {
	p := &oj.Parser{}
	for n := 0; n < b.N; n++ {
		_, _ = p.Parse([]byte(sampleJSON))
		//_, err := p.Parse([]byte(sampleJSON))
		//fmt.Println(err)
	}
}

func genParse(b *testing.B) {
	p := &gen.Parser{}
	for n := 0; n < b.N; n++ {
		_, _ = p.Parse([]byte(sampleJSON))
		//_, err := p.Parse([]byte(sampleJSON))
		//fmt.Println(err)
	}
}

func senParse(b *testing.B) {
	p := &sen.Parser{}
	for n := 0; n < b.N; n++ {
		_, _ = p.Parse([]byte(sampleSen))
		//_, err := p.Parse([]byte(sampleJSON))
		//fmt.Println(err)
	}
}

// Parse io.Reader
func ojParseReader(b *testing.B) {
	var p oj.Parser
	f, err := os.Open("test/sample.json")
	if err != nil {
		fmt.Printf("Failed to read test/sample.json. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		_, _ = p.ParseReader(f)
		//_, err = p.ParseReader(f)
		//fmt.Println(err)
	}
}

func genParseReader(b *testing.B) {
	var p gen.Parser
	f, err := os.Open("test/sample.json")
	if err != nil {
		fmt.Printf("Failed to read test/sample.json. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		_, _ = p.ParseReader(f)
		//_, err = p.ParseReader(f)
		//fmt.Println(err)
	}
}

func senParseReader(b *testing.B) {
	var p sen.Parser
	f, err := os.Open("test/sample.sen")
	if err != nil {
		fmt.Printf("Failed to read test/sample.sen. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		_, _ = p.ParseReader(f)
		//_, err = p.ParseReader(f)
		//fmt.Println(err)
	}
}

// Validate string/[]byte
func goValidate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = json.Valid([]byte(sampleJSON))
	}
}

func ojValidate(b *testing.B) {
	var v oj.Validator
	for n := 0; n < b.N; n++ {
		_ = v.Validate([]byte(sampleJSON))
		//err := v.Validate([]byte(sampleJSON))
		//fmt.Println(err)
	}
}

func goDecodeReader(b *testing.B) {
	f, err := os.Open("test/sample.json")
	if err != nil {
		fmt.Printf("Failed to read test/sample.json. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := json.NewDecoder(f)
		for {
			var data interface{}
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ojValidateReader(b *testing.B) {
	var v oj.Validator
	f, err := os.Open("test/sample.json")
	if err != nil {
		fmt.Printf("Failed to read test/sample.json. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		_ = v.ValidateReader(f)
		//err = v.ValidateReader(f)
		//fmt.Println(err)
	}
}

// JSON functions
func marshalJSON(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(data)
	}
}

func ojJSON(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojWrite(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = oj.Write(&buf, data, &opt)
		_ = buf.String()
	}
}

func senString(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := sen.Options{OmitNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sen.String(data, &opt)
	}
}

// JSON with indent functions
func marshalJSONIndent(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = json.MarshalIndent(data, "", "  ")
	}
}

func ojJSONIndent(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojWriteIndent(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = oj.Write(&buf, data, &opt)
		_ = buf.String()
	}
}

func senStringIndent(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := sen.Options{OmitNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sen.String(data, &opt)
	}
}

// JSON indented and sorted
func ojJSONSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2, Sort: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojWriteSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2, Sort: true}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = oj.Write(&buf, data, &opt)
		_ = buf.String()
	}
}

func senStringSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := sen.Options{OmitNil: true, Indent: 2, Sort: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sen.String(data, &opt)
	}
}

// Alter functions
func altGenerify(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_ = alt.Generify(native)
	}
}

func altGenAlter(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_ = alt.GenAlter(native)
	}
}

// jp.Get
func jpGet(b *testing.B) {
	p := jp.R().D().C("a").N(2).C("c")
	data := buildTree(10, 4, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.Get(data)
		//x := p.Get(data)
		//fmt.Printf("*** %s\n", oj.JSON(x))
	}
}

// jp.First
func jpFirst(b *testing.B) {
	p := jp.R().D().C("a").N(2).C("c")
	data := buildTree(10, 4, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.First(data)
		//fmt.Printf("*** %v\n", z)
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

const sampleJSON = `[
  [],
  null,
  true,
  false,
  77,
  123.456e7,
  "",
  "a string with \t unicode \u2669 and quotes \".",
  [1, 1.23, -44, "six"],
  [[null,[true,[false,[123,[4.56e7,["abcdef"]]]]]]],
  {
    "abc": 3,
    "def": {
      "ghi": true
    },
    "xyz": "another string",
    "nest": {
      "nest": {
        "nest": {
          "nest": {
            "nest": {
              "egg": [12345678, 87654321]
            }
          }
        }
      }
    }
  }
]
`

const sampleSen = `[
  []
  null
  true
  false
  77
  123.456e7
  ""
  "a string with \t unicode \u2669 and quotes \"."
  [1 1.23 -44 six]
  [[null[true[false[123[4.56e7[abcdef]]]]]]]
  {
    abc: 3
    def: {
      ghi: true
    }
    xyz: "another string"
    nest: {
      nest: {
        nest: {
          nest: {
            nest: {
              egg: [12345678 87654321]
            }
          }
        }
      }
    }
  }
]
`

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
