// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gd"

	"gitlab.com/uhn/core/pkg/tree"
)

// TBD remove tree before going public.

func main() {
	testing.Init()
	flag.Parse()
	tree.Sort = false
	gd.TimeFormat = "nano"

	validateBenchmarks()
	parseBenchmarks()

	base := testing.Benchmark(runBase)

	convBenchmarks(base)
	jsonBenchmarks(base, false, false)
	jsonBenchmarks(base, true, false)
	jsonBenchmarks(base, true, true)

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

	treeFrom := testing.Benchmark(convTree)

	treeNs := treeFrom.NsPerOp() - base.NsPerOp()
	treeBytes := treeFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	treeAllocs := treeFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("tree.FromNative:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdFrom := testing.Benchmark(convFromSimple)
	fromNs := gdFrom.NsPerOp() - base.NsPerOp()
	fromBytes := gdFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	fromAllocs := gdFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.FromSimple:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		fromNs, float64(treeNs)/float64(fromNs),
		fromBytes, float64(treeBytes)/float64(fromBytes),
		fromAllocs, float64(treeAllocs)/float64(fromAllocs))

	gdAlter := testing.Benchmark(convAlterSimple)
	alterNs := gdAlter.NsPerOp() - base.NsPerOp()
	alterBytes := gdAlter.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	alterAllocs := gdAlter.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.AlterSimple: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		alterNs, float64(treeNs)/float64(alterNs),
		alterBytes, float64(treeBytes)/float64(alterBytes),
		alterAllocs, float64(treeAllocs)/float64(alterAllocs))
}

func convAlterSimple(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.AlterSimple(native)
	}
}

func convFromSimple(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.FromSimple(native)
	}
}

func convTree(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = tree.FromNative(native)
	}
}

func jsonBenchmarks(base testing.BenchmarkResult, indent, sort bool) {
	fmt.Println()
	fmt.Printf("JSON() benchmarks, indent: %t, sort: %t\n", indent, sort)

	var treeRes testing.BenchmarkResult
	var ojgSRes testing.BenchmarkResult
	var ojgWRes testing.BenchmarkResult

	if sort {
		treeRes = testing.Benchmark(treeJSONSort)
	} else if indent {
		treeRes = testing.Benchmark(treeJSON2)
	} else {
		treeRes = testing.Benchmark(treeJSON)
	}
	treeNs := treeRes.NsPerOp()
	treeBytes := treeRes.AllocedBytesPerOp()
	treeAllocs := treeRes.AllocsPerOp()
	fmt.Printf("tree.JSON:   %10d ns/op (%3.2fx)  %10d B/op (%3.2fx)  %10d allocs/op (%3.2fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	if sort {
		ojgSRes = testing.Benchmark(ojgStringSort)
	} else if indent {
		ojgSRes = testing.Benchmark(ojgString2)
	} else {
		ojgSRes = testing.Benchmark(ojgString)
	}
	ojgNs := ojgSRes.NsPerOp()
	ojgBytes := ojgSRes.AllocedBytesPerOp()
	ojgAllocs := ojgSRes.AllocsPerOp()
	fmt.Printf(" ojg.String: %10d ns/op (%3.2fx)  %10d B/op (%3.2fx)  %10d allocs/op (%3.2fx)\n",
		ojgNs, float64(treeNs)/float64(ojgNs),
		ojgBytes, float64(treeBytes)/float64(ojgBytes),
		ojgAllocs, float64(treeAllocs)/float64(ojgAllocs))

	if sort {
		ojgWRes = testing.Benchmark(ojgWriteSort)
	} else if indent {
		ojgWRes = testing.Benchmark(ojgWrite2)
	} else {
		ojgWRes = testing.Benchmark(ojgWrite)
	}
	ojgNs = ojgWRes.NsPerOp()
	ojgBytes = ojgWRes.AllocedBytesPerOp()
	ojgAllocs = ojgWRes.AllocsPerOp()
	fmt.Printf(" ojg.Write:  %10d ns/op (%3.2fx)  %10d B/op (%3.2fx)  %10d allocs/op (%3.2fx)\n",
		ojgNs, float64(treeNs)/float64(ojgNs),
		ojgBytes, float64(treeBytes)/float64(ojgBytes),
		ojgAllocs, float64(treeAllocs)/float64(ojgAllocs))
}

func treeJSON(b *testing.B) {
	tree.Sort = false
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON()
	}
}

func treeJSON2(b *testing.B) {
	tree.Sort = false
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON(2)
	}
}

func treeJSONSort(b *testing.B) {
	tree.Sort = true
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON(2)
	}
}

func ojgString(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = ojg.String(data, &opt)
	}
}

func ojgString2(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = ojg.String(data, &opt)
	}
}

func ojgStringSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true, Indent: 2, Sort: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = ojg.String(data, &opt)
	}
}

func ojgWrite(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = ojg.Write(&buf, data, &opt)
		_ = buf.String()
	}
}

func ojgWrite2(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true, Indent: 2}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = ojg.Write(&buf, data, &opt)
		_ = buf.String()
	}
}

func ojgWriteSort(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterSimple(benchmarkData(tm))
	opt := ojg.Options{SkipNil: true, Indent: 2, Sort: true}
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_ = ojg.Write(&buf, data, &opt)
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
	fmt.Printf("json.Decode:       %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojgRes := testing.Benchmark(ojgValidate)
	ojgNs := ojgRes.NsPerOp()
	ojgBytes := ojgRes.AllocedBytesPerOp()
	ojgAllocs := ojgRes.AllocsPerOp()
	fmt.Printf(" ojg.Validate:     %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		ojgNs, float64(goNs)/float64(ojgNs),
		ojgBytes, float64(goBytes)/float64(ojgBytes),
		ojgAllocs, float64(goAllocs)/float64(ojgAllocs))

	treeRes := testing.Benchmark(treeValidate)
	treeNs := treeRes.NsPerOp()
	treeBytes := treeRes.AllocedBytesPerOp()
	treeAllocs := treeRes.AllocsPerOp()
	fmt.Printf("tree.ParseString:  %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		treeNs, float64(goNs)/float64(treeNs),
		treeBytes, float64(goBytes)/float64(treeBytes),
		treeAllocs, float64(goAllocs)/float64(treeAllocs))
}

func goValidate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//_ = json.Valid([]byte(patientJSON))
		_ = json.Valid([]byte(sampleJSON))
	}
}

func ojgValidate(b *testing.B) {
	v := &ojg.Validator{}
	for n := 0; n < b.N; n++ {
		//_ = v.Validate([]byte(patientJSON))
		_ = v.Validate([]byte(sampleJSON))
		//err := v.Validate([]byte(patientJSON))
		//fmt.Println(err)
	}
}

func treeValidate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//_, _ = tree.ParseString(patientJSON)
		_, _ = tree.ParseString(sampleJSON)
	}
}

func parseBenchmarks() {
	fmt.Println()
	fmt.Println("Parse JSON")

	goRes := testing.Benchmark(goParse)
	goNs := goRes.NsPerOp()
	goBytes := goRes.AllocedBytesPerOp()
	goAllocs := goRes.AllocsPerOp()
	fmt.Printf("json.Unmarshal:    %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		goNs, 1.0, goBytes, 1.0, goAllocs, 1.0)

	ojgRes := testing.Benchmark(ojgParse)
	ojgNs := ojgRes.NsPerOp()
	ojgBytes := ojgRes.AllocedBytesPerOp()
	ojgAllocs := ojgRes.AllocsPerOp()
	fmt.Printf(" ojg.Parse:        %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		ojgNs, float64(goNs)/float64(ojgNs),
		ojgBytes, float64(goBytes)/float64(ojgBytes),
		ojgAllocs, float64(goAllocs)/float64(ojgAllocs))

	treeRes := testing.Benchmark(treeParse)
	treeNs := treeRes.NsPerOp()
	treeBytes := treeRes.AllocedBytesPerOp()
	treeAllocs := treeRes.AllocsPerOp()
	fmt.Printf("tree.ParseString:  %10d ns/op (%3.2fx)  %10d B/op (%4.2fx)  %10d allocs/op (%4.2fx)\n",
		treeNs, float64(goNs)/float64(treeNs),
		treeBytes, float64(goBytes)/float64(treeBytes),
		treeAllocs, float64(goAllocs)/float64(treeAllocs))
}

func goParse(b *testing.B) {
	var result interface{}
	for n := 0; n < b.N; n++ {
		_ = json.Unmarshal([]byte(sampleJSON), &result)
	}
}

func ojgParse(b *testing.B) {
	p := &ojg.Parser{}
	for n := 0; n < b.N; n++ {
		_, _ = p.Parse([]byte(sampleJSON))
		//_, err := p.Parse([]byte(sampleJSON))
		//fmt.Println(err)
	}
}

func treeParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = tree.ParseString(sampleJSON)
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
  "a string",
  [1, 1.23, -44, "six"],
  [[null,[true,[false,[123,[4.56e7,["abcdef"]]]]]]],
  {
    "abc": 3,
    "def": {
      "ghi": true
    },
    "xyz": "another string"
  }
]
`

const patientJSON = `{
  "resourceType": "Patient",
  "id": "example",
  "text": {
    "status": "generated",
    "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">\n\t\t\t<table>\n\t\t\t\t<tbody>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Name</td>\n\t\t\t\t\t\t<td>Peter James \n              <b>Chalmers</b> (&quot;Jim&quot;)\n            </td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Address</td>\n\t\t\t\t\t\t<td>534 Erewhon, Pleasantville, Vic, 3999</td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Contacts</td>\n\t\t\t\t\t\t<td>Home: unknown. Work: (03) 5555 6473</td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Id</td>\n\t\t\t\t\t\t<td>MRN: 12345 (Acme Healthcare)</td>\n\t\t\t\t\t</tr>\n\t\t\t\t</tbody>\n\t\t\t</table>\n\t\t</div>"
  },
  "identifier": [
    {
      "use": "usual",
      "system": "urn:oid:1.2.36.146.595.217.0.3",
      "value": "54321"
    },
    {
      "use": "official",
      "type": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v2-0203",
            "code": "MR"
          }
        ]
      },
      "system": "urn:oid:1.2.36.146.595.217.0.1",
      "value": "12345",
      "period": {
        "start": "2001-05-06"
      },
      "assigner": {
        "display": "Acme Healthcare"
      }
    }
  ],
  "active": true,
  "name": [
    {
      "use": "official",
      "family": "Chalmers",
      "given": [
        "Peter",
        "James"
      ]
    },
    {
      "use": "usual",
      "given": [
        "Jim"
      ]
    },
    {
      "use": "maiden",
      "family": "Windsor",
      "given": [
        "Peter",
        "James"
      ],
      "period": {
        "end": "2002"
      }
    }
  ],
  "telecom": [
    {
      "use": "home"
    },
    {
      "system": "phone",
      "value": "(03) 5555 6473",
      "use": "work",
      "rank": 1
    },
    {
      "system": "phone",
      "value": "(03) 3410 5613",
      "use": "mobile",
      "rank": 2
    },
    {
      "system": "phone",
      "value": "(03) 5555 8834",
      "use": "old",
      "period": {
        "end": "2014"
      }
    }
  ],
  "gender": "male",
  "birthDate": "1974-12-25",
  "_birthDate": {
    "extension": [
      {
        "url": "http://hl7.org/fhir/StructureDefinition/patient-birthTime",
        "valueDateTime": "1974-12-25T14:35:45-05:00"
      }
    ]
  },
  "deceasedBoolean": false,
  "address": [
    {
      "use": "home",
      "type": "both",
      "text": "534 Erewhon St PeasantVille, Rainbow, Vic  3999",
      "line": [
        "534 Erewhon St"
      ],
      "city": "PleasantVille",
      "district": "Rainbow",
      "state": "Vic",
      "postalCode": "3999",
      "period": {
        "start": "1974-12-25"
      }
    }
  ],
  "contact": [
    {
      "relationship": [
        {
          "coding": [
            {
              "system": "http://terminology.hl7.org/CodeSystem/v2-0131",
              "code": "N"
            }
          ]
        }
      ],
      "name": {
        "family": "du Marché",
        "_family": {
          "extension": [
            {
              "url": "http://hl7.org/fhir/StructureDefinition/humanname-own-prefix",
              "valueString": "VV"
            }
          ]
        },
        "given": [
          "Bénédicte"
        ]
      },
      "telecom": [
        {
          "system": "phone",
          "value": "+33 (237) 998327"
        }
      ],
      "address": {
        "use": "home",
        "type": "both",
        "line": [
          "534 Erewhon St"
        ],
        "city": "PleasantVille",
        "district": "Rainbow",
        "state": "Vic",
        "postalCode": "3999",
        "period": {
          "start": "1974-12-25"
        }
      },
      "gender": "female",
      "period": {
        "start": "2012"
      }
    }
  ],
  "managingOrganization": {
    "reference": "Organization/1"
  }
}
`
