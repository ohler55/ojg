// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
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
	tree.Sort = false
	gd.Sort = false
	gd.TimeFormat = "nano"

	fmt.Println()
	fmt.Println("Converting from native to canonical types benchmarks")
	base := testing.Benchmark(convBase)
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

	fmt.Println()
	fmt.Println("JSON() benchmarks")
	treeJSON := testing.Benchmark(treeJSON)

	treeNs = treeJSON.NsPerOp()
	treeBytes = treeJSON.AllocedBytesPerOp()
	treeAllocs = treeJSON.AllocsPerOp()
	fmt.Printf("tree.JSON:   %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON := testing.Benchmark(ojgString)
	gdNs := gdJSON.NsPerOp()
	gdBytes := gdJSON.AllocedBytesPerOp()
	gdAllocs := gdJSON.AllocsPerOp()
	fmt.Printf(" ojg.String: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		gdNs, float64(treeNs)/float64(gdNs),
		gdBytes, float64(treeBytes)/float64(gdBytes),
		gdAllocs, float64(treeAllocs)/float64(gdAllocs))

	gdJSON = testing.Benchmark(ojgWrite)
	gdNs = gdJSON.NsPerOp()
	gdBytes = gdJSON.AllocedBytesPerOp()
	gdAllocs = gdJSON.AllocsPerOp()
	fmt.Printf(" ojg.Write:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		gdNs, float64(treeNs)/float64(gdNs),
		gdBytes, float64(treeBytes)/float64(gdBytes),
		gdAllocs, float64(treeAllocs)/float64(gdAllocs))

	fmt.Println()
	fmt.Println("JSON(2) benchmarks")
	treeJSON = testing.Benchmark(json2Tree)

	treeNs = treeJSON.NsPerOp()
	treeBytes = treeJSON.AllocedBytesPerOp()
	treeAllocs = treeJSON.AllocsPerOp()
	fmt.Printf("tree.JSON:   %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON = testing.Benchmark(ojgString2)
	gdNs = gdJSON.NsPerOp()
	gdBytes = gdJSON.AllocedBytesPerOp()
	gdAllocs = gdJSON.AllocsPerOp()
	fmt.Printf(" ojg.String: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		gdNs, float64(treeNs)/float64(gdNs),
		gdBytes, float64(treeBytes)/float64(gdBytes),
		gdAllocs, float64(treeAllocs)/float64(gdAllocs))

	fmt.Println()
	fmt.Println("JSON(2) sorted benchmarks")
	tree.Sort = true
	gd.Sort = true
	treeJSON = testing.Benchmark(json2Tree)

	treeNs = treeJSON.NsPerOp()
	treeBytes = treeJSON.AllocedBytesPerOp()
	treeAllocs = treeJSON.AllocsPerOp()
	fmt.Printf("tree.JSON:   %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON = testing.Benchmark(ojgStringSort)
	gdNs = gdJSON.NsPerOp()
	gdBytes = gdJSON.AllocedBytesPerOp()
	gdAllocs = gdJSON.AllocsPerOp()
	fmt.Printf(" ojg.String: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		gdNs, float64(treeNs)/float64(gdNs),
		gdBytes, float64(treeBytes)/float64(gdBytes),
		gdAllocs, float64(treeAllocs)/float64(gdAllocs))

	fmt.Println()
}

func convBase(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		benchmarkData(tm)
	}
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

func treeJSON(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON()
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

func json2Tree(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON(2)
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
