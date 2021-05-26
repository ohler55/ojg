// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

const (
	blocks    = " ▏▎▍▌▋▊▉█"
	darkBlock = "▓"
	//mediumBlock = "▒"
	//lightBlock  = "░"
)

var (
	filename    = "test/patient.json"
	patFilename = "test/patient.json"
	catFilename = "test/citm_catalog.json"

	useCat     = flag.Bool("cat", false, "marshal catalog instead of patient")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
)

type specs struct {
	os        string
	model     string
	processor string
	cores     string
	speed     string
}

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
	if len([]interface{}{&Patient{}, &Catalog{}}) == 3 {
		// Dummy to avoid linter complaints when not using one set of types.
		fmt.Println("how did we get here?")
	}
	testing.Init()
	flag.Parse()
	if 0 < len(flag.Args()) {
		filename = flag.Args()[0]
	} else if *useCat {
		filename = catFilename
	}
	gen.TimeFormat = "nano"
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if *useCat {
		benchSuite("Marshal Struct", []*bench{
			{pkg: "json", name: "Marshal", fun: goMarshalCatalog},
			{pkg: "oj", name: "Marshal", fun: ojMarshalCatalog},
		})
	} else {
		benchSuite("Marshal Struct", []*bench{
			{pkg: "json", name: "Marshal", fun: goMarshalPatient},
			{pkg: "oj", name: "Marshal", fun: ojMarshalPatient},
		})
	}
	return
	benchSuite("Parse string/[]byte", []*bench{
		{pkg: "json", name: "Unmarshal", fun: goParse},
		{pkg: "oj", name: "Parse", fun: ojParse},
		{pkg: "oj-reuse", name: "Parse", fun: ojParseReuse},
		{pkg: "gen", name: "Parse", fun: genParse},
		{pkg: "gen-reuse", name: "Parse", fun: genParseReuse},
		{pkg: "sen", name: "Parse", fun: senParse},
		{pkg: "sen-reuse", name: "Parse", fun: senParseReuse},
	})

	if *useCat {
		benchSuite("Unmarshal []byte to type", []*bench{
			{pkg: "json", name: "Unmarshal", fun: goUnmarshalCatalog},
			{pkg: "oj", name: "Unmarshal", fun: ojUnmarshalCatalog},
			{pkg: "sen", name: "Unmarshal", fun: senUnmarshalCatalog},
		})
	} else {
		benchSuite("Unmarshal []byte to type", []*bench{
			{pkg: "json", name: "Unmarshal", fun: goUnmarshalPatient},
			{pkg: "oj", name: "Unmarshal", fun: ojUnmarshalPatient},
			{pkg: "sen", name: "Unmarshal", fun: senUnmarshalPatient},
		})
	}
	benchSuite("Tokenize", []*bench{
		{pkg: "json", name: "Decode", fun: goDecode},
		{pkg: "oj", name: "Tokenize", fun: ojTokenize},
		{pkg: "sen", name: "Tokenize", fun: senTokenize},
	})
	benchSuite("Parse io.Reader", []*bench{
		{pkg: "json", name: "Decode", fun: goDecodeReader},
		{pkg: "oj", name: "ParseReader", fun: ojParseReader},
		{pkg: "oj-reuse", name: "ParseReader", fun: ojParseReaderReuse},
		{pkg: "gen", name: "ParseReder", fun: genParseReader},
		{pkg: "gen-reuse", name: "ParseReder", fun: genParseReaderReuse},
		{pkg: "sen", name: "ParseReader", fun: senParseReader},
		{pkg: "sen-reuse", name: "ParseReader", fun: senParseReaderReuse},
		{pkg: "oj", name: "TokenizeLoad", fun: ojTokenizeLoad},
		{pkg: "sen", name: "TokenizeLoad", fun: senTokenizeLoad},
	})
	benchSuite("Parse chan interface{}", []*bench{
		{pkg: "json", name: "Parse-chan", fun: goParseChan},
		{pkg: "oj", name: "Parse", fun: ojParseChan},
		{pkg: "gen", name: "Parse", fun: genParseChan},
		{pkg: "sen", name: "Parse", fun: senParseChan},
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
		{pkg: "sen", name: "SEN", fun: senSEN},
	})
	benchSuite("to JSON with indentation", []*bench{
		{pkg: "json", name: "Marshal", fun: marshalJSONIndent},
		{pkg: "oj", name: "JSON", fun: ojJSONIndent},
		{pkg: "sen", name: "Bytes", fun: senSENIndent},
		{pkg: "pretty", name: "JSON", fun: prettyJSON},
		{pkg: "pretty", name: "SEN", fun: prettySEN},
	})
	benchSuite("to JSON with indentation and sorted keys", []*bench{
		{pkg: "oj", name: "JSON", fun: ojJSONSort},
		{pkg: "sen", name: "Bytes", fun: senSENSort},
		{pkg: "pretty", name: "JSON", fun: prettyJSON},
		{pkg: "pretty", name: "SEN", fun: prettySEN},
	})
	benchSuite("Write indented JSON", []*bench{
		{pkg: "json", name: "Encode", fun: jsonEncodeIndent},
		{pkg: "oj", name: "Write", fun: ojWriteIndent},
		{pkg: "sen", name: "Write", fun: senWriteIndent},
		{pkg: "pretty", name: "WriteJSON", fun: prettyWriteJSON},
		{pkg: "pretty", name: "WriteSEN", fun: prettyWriteSEN},
	})

	if *useCat {
		benchSuite("Marshal Struct", []*bench{
			{pkg: "json", name: "Marshal", fun: goMarshalCatalog},
			{pkg: "oj", name: "Marshal", fun: ojMarshalCatalog},
		})
	} else {
		benchSuite("Marshal Struct", []*bench{
			{pkg: "json", name: "Marshal", fun: goMarshalPatient},
			{pkg: "oj", name: "Marshal", fun: ojMarshalPatient},
		})
	}

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
	if s := getSpecs(); s != nil {
		fmt.Println("Tests run on:")
		if 0 < len(s.model) {
			fmt.Printf(" Machine:         %s\n", s.model)
		}
		fmt.Printf(" OS:              %s\n", s.os)
		fmt.Printf(" Processor:       %s\n", s.processor)
		fmt.Printf(" Cores:           %s\n", s.cores)
		fmt.Printf(" Processor Speed: %s\n", s.speed)
	}
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
		fmt.Printf(" %10s.%-16s %8d ns/op %8d B/op %6d allocs/op\n",
			b.pkg, b.name, b.ns, b.bytes, b.allocs)
	}
	fmt.Println()

	scale := 8 // TBD adjust to fit screen better?
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
		fmt.Printf(" %10s.%-12s %s %3.2f\n", b.pkg, b.name, bar, x)
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

func getSpecs() (s *specs) {
	// Assume MacOS and try system_profiler. If that fails assume linux and check /proc.
	out, err := exec.Command("system_profiler", "-json", "SPHardwareDataType").Output()
	if err == nil {
		var js interface{}
		if js, err = oj.Parse(out); err == nil {
			s = &specs{
				model:     alt.String(jp.C("SPHardwareDataType").N(0).C("machine_model").First(js)),
				processor: alt.String(jp.C("SPHardwareDataType").N(0).C("cpu_type").First(js)),
				cores:     alt.String(jp.C("SPHardwareDataType").N(0).C("number_processors").First(js)),
				speed:     alt.String(jp.C("SPHardwareDataType").N(0).C("current_processor_speed").First(js)),
			}
			var b []byte
			if out, err = exec.Command("sw_vers", "-productName").Output(); err == nil {
				b = append(b, bytes.TrimSpace(out)...)
				b = append(b, ' ')
			}
			if out, err = exec.Command("sw_vers", "-productVersion").Output(); err == nil {
				b = append(b, bytes.TrimSpace(out)...)
			}
			s.os = string(b)
		}
		return
	}
	// Try Ubuntu next.
	if out, err = exec.Command("lsb_release", "-d").Output(); err == nil {
		s = &specs{}
		parts := strings.Split(string(out), ":")
		if 1 < len(parts) {
			s.os = string(strings.TrimSpace(parts[1]))
		}
		if out, err = ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			cnt := 0
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "processor") {
					cnt++
				} else if strings.Contains(line, "model name") {
					parts := strings.Split(line, ":")
					if 1 < len(parts) {
						parts = strings.Split(parts[1], "@")
						s.processor = strings.TrimSpace(parts[0])
						if 1 < len(parts) {
							s.speed = strings.TrimSpace(parts[1])
						}
					}
				}
				s.cores = fmt.Sprintf("%d", cnt)
			}
		}
	}
	return
}
