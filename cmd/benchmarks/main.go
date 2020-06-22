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
)

// TBD remove tree before going public.

func main() {
	testing.Init()
	flag.Parse()
	gen.TimeFormat = "nano"

	jsonPathGetBenchmarks()
	jsonPathFirstBenchmarks()

	parseBenchmarks()
	parseReaderBenchmarks()
	validateBenchmarks()
	validateReaderBenchmarks()

	base := testing.Benchmark(runBase)

	jsonBenchmarks(base, false)
	jsonBenchmarks(base, true)
	jsonSortBenchmarks(base)
	convBenchmarks(base)

	fmt.Println()
}

func runBase(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		benchmarkData(tm)
	}
}

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

func convBenchmarks(base testing.BenchmarkResult) {
	fmt.Println()
	fmt.Println("Converting from simple to canonical types benchmarks")

	genFrom := testing.Benchmark(ojGenerify)
	fromNs := genFrom.NsPerOp() - base.NsPerOp()
	fromBytes := genFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	fromAllocs := genFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  oj.Generify:          %6d ns/op (%3.1fx)  %6d B/op (%3.1fx)  %6d allocs/op (%3.1fx)\n",
		fromNs, 1.0,
		fromBytes, 1.0,
		fromAllocs, 1.0)

	genAlter := testing.Benchmark(ojGenAlter)
	alterNs := genAlter.NsPerOp() - base.NsPerOp()
	alterBytes := genAlter.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	alterAllocs := genAlter.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  oj.GenAlter:          %6d ns/op (%3.1fx)  %6d B/op (%3.1fx)  %6d allocs/op (%3.1fx)\n",
		alterNs, float64(fromNs)/float64(alterNs),
		alterBytes, float64(fromBytes)/float64(alterBytes),
		alterAllocs, float64(fromAllocs)/float64(alterAllocs))
}

func ojGenAlter(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_ = alt.GenAlter(native)
	}
}

func ojGenerify(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_ = alt.Generify(native)
	}
}

func jsonBenchmarks(base testing.BenchmarkResult, indent bool) {
	fmt.Println()
	fmt.Printf("JSON() benchmarks, indent: %t, sort: false\n", indent)

	var marshalRes testing.BenchmarkResult
	var ojSRes testing.BenchmarkResult
	var ojWRes testing.BenchmarkResult

	if indent {
		marshalRes = testing.Benchmark(marshalJSON2)
	} else {
		marshalRes = testing.Benchmark(marshalJSON)
	}
	marshalNs := marshalRes.NsPerOp()
	marshalBytes := marshalRes.AllocedBytesPerOp()
	marshalAllocs := marshalRes.AllocsPerOp()
	fmt.Printf("json.Marshal:           %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
		marshalNs, 1.0, marshalBytes, 1.0, marshalAllocs, 1.0)

	if indent {
		ojSRes = testing.Benchmark(ojJSON2)
	} else {
		ojSRes = testing.Benchmark(ojJSON)
	}
	ojNs := ojSRes.NsPerOp()
	ojBytes := ojSRes.AllocedBytesPerOp()
	ojAllocs := ojSRes.AllocsPerOp()
	fmt.Printf("  oj.JSON:              %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
		ojNs, float64(marshalNs)/float64(ojNs),
		ojBytes, float64(marshalBytes)/float64(ojBytes),
		ojAllocs, float64(marshalAllocs)/float64(ojAllocs))

	if indent {
		ojWRes = testing.Benchmark(ojWrite2)
	} else {
		ojWRes = testing.Benchmark(ojWrite)
	}
	ojNs = ojWRes.NsPerOp()
	ojBytes = ojWRes.AllocedBytesPerOp()
	ojAllocs = ojWRes.AllocsPerOp()
	fmt.Printf("  oj.Write:             %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
		ojNs, float64(marshalNs)/float64(ojNs),
		ojBytes, float64(marshalBytes)/float64(ojBytes),
		ojAllocs, float64(marshalAllocs)/float64(ojAllocs))
}

func jsonSortBenchmarks(base testing.BenchmarkResult) {
	fmt.Println()
	fmt.Printf("JSON() benchmarks, sort: true\n")

	var ojSRes testing.BenchmarkResult
	var ojWRes testing.BenchmarkResult

	ojSRes = testing.Benchmark(ojJSONSort)
	ns := ojSRes.NsPerOp()
	bytes := ojSRes.AllocedBytesPerOp()
	allocs := ojSRes.AllocsPerOp()
	fmt.Printf("  oj.JSON:              %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
		ns, 1.0,
		bytes, 10.,
		allocs, 1.0)

	ojWRes = testing.Benchmark(ojWriteSort)
	ojNs := ojWRes.NsPerOp()
	ojBytes := ojWRes.AllocedBytesPerOp()
	ojAllocs := ojWRes.AllocsPerOp()
	fmt.Printf("  oj.Write:             %6d ns/op (%3.2fx)  %6d B/op (%3.2fx)  %6d allocs/op (%3.2fx)\n",
		ojNs, float64(ns)/float64(ojNs),
		ojBytes, float64(bytes)/float64(ojBytes),
		ojAllocs, float64(allocs)/float64(ojAllocs))
}

func marshalJSONSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(data)
	}
}

func marshalJSON2(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = json.MarshalIndent(data, "", "  ")
	}
}

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

func ojJSON2(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojJSONSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data := alt.Alter(benchmarkData(tm))
	opt := oj.Options{OmitNil: true, Indent: 2, Sort: true}
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

func ojWrite2(b *testing.B) {
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

func validateBenchmarks() {
	fmt.Println()
	fmt.Println("Validate JSON")

	goRes := testing.Benchmark(goValidate)
	goNs := goRes.NsPerOp()
	goBytes := goRes.AllocedBytesPerOp()
	goAllocs := goRes.AllocsPerOp()
	fmt.Printf("json.Valid:             %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojRes := testing.Benchmark(ojValidate)
	ojNs := ojRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp()
	fmt.Printf("  oj.Validate:          %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))
}

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

func validateReaderBenchmarks() {
	fmt.Println()
	fmt.Println("Validate io.Reader JSON")

	baseRes := testing.Benchmark(baseValidateReader)

	goRes := testing.Benchmark(goParseReader)
	goNs := goRes.NsPerOp() - baseRes.NsPerOp()
	goBytes := goRes.AllocedBytesPerOp() - baseRes.AllocedBytesPerOp()
	goAllocs := goRes.AllocsPerOp() - baseRes.AllocsPerOp()
	fmt.Printf("json.Decoder:           %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojRes := testing.Benchmark(ojValidateReader)
	ojNs := ojRes.NsPerOp() - baseRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp() - baseRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp() - baseRes.AllocsPerOp()
	fmt.Printf("  oj.ValidateReader:    %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))
}

func baseValidateReader(b *testing.B) {
	f, err := os.Open("test/sample.json")
	if err != nil {
		fmt.Printf("Failed to read test/sample.json. %s\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		_, _ = f.Seek(0, 2)
	}
}

func goParseReader(b *testing.B) {
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

func parseBenchmarks() {
	fmt.Println()
	fmt.Println("Parse JSON")

	goRes := testing.Benchmark(goParse)
	goNs := goRes.NsPerOp()
	goBytes := goRes.AllocedBytesPerOp()
	goAllocs := goRes.AllocsPerOp()
	fmt.Printf("json.Unmarshal:         %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojRes := testing.Benchmark(ojParse)
	ojNs := ojRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp()
	fmt.Printf("  oj.Parse:             %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))

	ojRes = testing.Benchmark(genParse)
	ojNs = ojRes.NsPerOp()
	ojBytes = ojRes.AllocedBytesPerOp()
	ojAllocs = ojRes.AllocsPerOp()
	fmt.Printf("  oj.GenParse:          %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))
}

func parseReaderBenchmarks() {
	fmt.Println()
	fmt.Println("Parse io.Reader JSON")

	baseRes := testing.Benchmark(baseValidateReader)

	goRes := testing.Benchmark(goParseReader)
	goNs := goRes.NsPerOp() - baseRes.NsPerOp()
	goBytes := goRes.AllocedBytesPerOp() - baseRes.AllocedBytesPerOp()
	goAllocs := goRes.AllocsPerOp() - baseRes.AllocsPerOp()
	fmt.Printf("json.Decoder:           %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojRes := testing.Benchmark(ojParseReader)
	ojNs := ojRes.NsPerOp() - baseRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp() - baseRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp() - baseRes.AllocsPerOp()
	fmt.Printf("  oj.ParseReader:       %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))

	ojRes = testing.Benchmark(genParseReader)
	ojNs = ojRes.NsPerOp() - baseRes.NsPerOp()
	ojBytes = ojRes.AllocedBytesPerOp() - baseRes.AllocedBytesPerOp()
	ojAllocs = ojRes.AllocsPerOp() - baseRes.AllocsPerOp()
	fmt.Printf("  oj.GenParseReader:    %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, float64(goNs)/float64(ojNs),
		ojBytes, float64(goBytes)/float64(ojBytes),
		ojAllocs, float64(goAllocs)/float64(ojAllocs))
}

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

func jsonPathGetBenchmarks() {
	fmt.Println()
	fmt.Println("JSON Path Get")

	ojRes := testing.Benchmark(ojGet)
	ojNs := ojRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp()
	fmt.Printf("  oj.Expr.Get:          %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, 1.0, ojBytes, 1.0, ojAllocs, 1.0)
}

func ojGet(b *testing.B) {
	p := jp.D().C("a").W().C("c")
	data := buildTree(10, 4, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.Get(data)
		//x := p.Get(data)
		//fmt.Printf("*** %s\n", oj.JSON(x))
	}
}

func jsonPathFirstBenchmarks() {
	fmt.Println()
	fmt.Println("JSON Path First")

	ojRes := testing.Benchmark(ojFirst)
	ojNs := ojRes.NsPerOp()
	ojBytes := ojRes.AllocedBytesPerOp()
	ojAllocs := ojRes.AllocsPerOp()
	fmt.Printf("  oj.Expr.First:        %6d ns/op (%3.2fx)  %6d B/op (%4.2fx)  %6d allocs/op (%4.2fx)\n",
		ojNs, 1.0, ojBytes, 1.0, ojAllocs, 1.0)
}

func ojFirst(b *testing.B) {
	p := jp.X().D().C("a").W().C("c").C("d")
	data := buildTree(10, 3, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = p.First(data)
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
