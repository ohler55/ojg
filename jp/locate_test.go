// Copyright (c) 2023, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

type locateData struct {
	path   string
	max    int
	data   any
	expect []string
}

var (
	locateTestData = []*locateData{
		{path: "", expect: []string{}},
		{path: "a.b", data: map[string]any{"a": map[string]any{"b": nil}, "x": 3}, expect: []string{"a.b"}},
		{path: "a[1]", data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[1]"}},
		{path: "a[*]", data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[0]", "a[1]", "a[2]"}},
		{path: "$.a.*.b", max: 2, expect: []string{"$.a[0].b", "$.a[1].b"}},
	}
)

func testDiffString(expect, actual []string, diff alt.Path) string {
	var b []byte

	b = fmt.Appendf(b, "\n      diff at %s\n", diff)
	b = append(b, "      expect: ["...)
	for _, str := range expect {
		b = append(b, "\n        "...)
		b = append(b, str...)
	}
	b = append(b, "\n      ]\n      actual: ["...)
	for _, str := range actual {
		b = append(b, "\n        "...)
		b = append(b, str...)
	}
	b = append(b, "\n      ]\n"...)

	return string(b)
}

func TestExprLocate(t *testing.T) {
	data := buildTree(4, 3, 0)
	// fmt.Printf("*** %s\n", pretty.SEN(data))
	for i, d := range locateTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var locs []jp.Expr
		if d.data == nil {
			locs = x.Locate(data, d.max)
		} else {
			locs = x.Locate(d.data, d.max)
		}
		var results []string
		for _, loc := range locs {
			results = append(results, loc.String())
		}
		sort.Strings(results)
		diff := alt.Compare(d.expect, results)
		if 0 < len(diff) {
			t.Fatal(testDiffString(d.expect, results, diff))
		}
	}
}

// func TestLocateDev(t *testing.T) {
// 	data := []any{map[string]any{"b": 1, "c": 2}, []any{1, 2, 3}}
// 	x := jp.MustParseString("$[?(@[1] == 2)].*")
// 	for _, ep := range x.Locate(data, 0) {
// 		fmt.Printf("*** %s\n", ep.BracketString())
// 	}
// }
