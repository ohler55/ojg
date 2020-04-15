// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"testing"
	"time"

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

	gdFrom := testing.Benchmark(convFromNative)
	fromNs := gdFrom.NsPerOp() - base.NsPerOp()
	fromBytes := gdFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	fromAllocs := gdFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.FromNative:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		fromNs, float64(treeNs)/float64(fromNs),
		fromBytes, float64(treeBytes)/float64(fromBytes),
		fromAllocs, float64(treeAllocs)/float64(fromAllocs))

	gdAlter := testing.Benchmark(convAlterNative)
	alterNs := gdAlter.NsPerOp() - base.NsPerOp()
	alterBytes := gdAlter.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	alterAllocs := gdAlter.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.AlterNative: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		alterNs, float64(treeNs)/float64(alterNs),
		alterBytes, float64(treeBytes)/float64(alterBytes),
		alterAllocs, float64(treeAllocs)/float64(alterAllocs))

	fmt.Println()
	fmt.Println("JSON() benchmarks")
	treeJSON := testing.Benchmark(jsonTree)

	treeNs = treeJSON.NsPerOp()
	treeBytes = treeJSON.AllocedBytesPerOp()
	treeAllocs = treeJSON.AllocsPerOp()
	fmt.Printf("tree.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON := testing.Benchmark(jsonGd)
	gdNs := gdJSON.NsPerOp()
	gdBytes := gdJSON.AllocedBytesPerOp()
	gdAllocs := gdJSON.AllocsPerOp()
	fmt.Printf("  gd.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		gdNs, float64(treeNs)/float64(gdNs),
		gdBytes, float64(treeBytes)/float64(gdBytes),
		gdAllocs, float64(treeAllocs)/float64(gdAllocs))

	fmt.Println()
	fmt.Println("JSON(2) benchmarks")
	treeJSON = testing.Benchmark(json2Tree)

	treeNs = treeJSON.NsPerOp()
	treeBytes = treeJSON.AllocedBytesPerOp()
	treeAllocs = treeJSON.AllocsPerOp()
	fmt.Printf("tree.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON = testing.Benchmark(json2Gd)
	gdNs = gdJSON.NsPerOp()
	gdBytes = gdJSON.AllocedBytesPerOp()
	gdAllocs = gdJSON.AllocsPerOp()
	fmt.Printf("  gd.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
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
	fmt.Printf("tree.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdJSON = testing.Benchmark(json2Gd)
	gdNs = gdJSON.NsPerOp()
	gdBytes = gdJSON.AllocedBytesPerOp()
	gdAllocs = gdJSON.AllocsPerOp()
	fmt.Printf("  gd.JSON:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
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

func convAlterNative(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.AlterNative(native)
	}
}

func convFromNative(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.FromNative(native)
	}
}

func convTree(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = tree.FromNative(native)
	}
}

func jsonTree(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := tree.FromNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON()
	}
}

func jsonGd(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON()
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

func json2Gd(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	data, _ := gd.AlterNative(benchmarkData(tm))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = data.JSON(2)
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
