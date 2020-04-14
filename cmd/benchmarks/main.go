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
	base := testing.Benchmark(benchmarkBase)

	treeFrom := testing.Benchmark(benchmarkTree)

	fmt.Println()
	treeNs := treeFrom.NsPerOp() - base.NsPerOp()
	treeBytes := treeFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	treeAllocs := treeFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("tree.FromNative:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		treeNs, 1.0, treeBytes, 1.0, treeAllocs, 1.0)

	gdFrom := testing.Benchmark(benchmarkFromNative)
	fromNs := gdFrom.NsPerOp() - base.NsPerOp()
	fromBytes := gdFrom.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	fromAllocs := gdFrom.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.FromNative:  %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		fromNs, float64(treeNs)/float64(fromNs),
		fromBytes, float64(treeBytes)/float64(fromBytes),
		fromAllocs, float64(treeAllocs)/float64(fromAllocs))

	gdAlter := testing.Benchmark(benchmarkAlterNative)
	alterNs := gdAlter.NsPerOp() - base.NsPerOp()
	alterBytes := gdAlter.AllocedBytesPerOp() - base.AllocedBytesPerOp()
	alterAllocs := gdAlter.AllocsPerOp() - base.AllocsPerOp()
	fmt.Printf("  gd.AlterNative: %10d ns/op (%3.1fx)  %10d B/op (%3.1fx)  %10d allocs/op (%3.1fx)\n",
		alterNs, float64(treeNs)/float64(alterNs),
		alterBytes, float64(treeBytes)/float64(alterBytes),
		alterAllocs, float64(treeAllocs)/float64(alterAllocs))
	fmt.Println()
}

func benchmarkBase(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		benchmarkData(tm)
	}
}

func benchmarkAlterNative(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.AlterNative(native)
	}
}

func benchmarkFromNative(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = gd.FromNative(native)
	}
}

func benchmarkTree(b *testing.B) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for n := 0; n < b.N; n++ {
		native := benchmarkData(tm)
		_, _ = tree.FromNative(native)
	}
}

func benchmarkData(tm time.Time) interface{} {
	return map[string]interface{}{
		"a": []interface{}{1, 2, true, tm},
		"b": 2.3,
		"c": map[string]interface{}{
			"x": "xxx",
		},
	}
}
