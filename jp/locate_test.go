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
	noSort bool
}

var (
	locateTestData = []*locateData{
		{path: "", expect: []string{}},
		{path: "a.b", data: map[string]any{"a": map[string]any{"b": 2}, "x": 3}, expect: []string{"a.b"}},
		{path: "a[1]", data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[1]"}},
		{path: "a[-1]", data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[2]"}},
		{path: "a[*]", data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[0]", "a[1]", "a[2]"}},
		{path: "$.a.*.b", max: 2, expect: []string{"$.a[0].b", "$.a[1].b"}},
		{path: "$.a[1].*", expect: []string{"$.a[1].a", "$.a[1].b", "$.a[1].c", "$.a[1].d"}},
		{path: "$.*[1].c", expect: []string{"$.a[1].c", "$.b[1].c", "$.c[1].c", "$.d[1].c"}},
		{path: "*[*]", max: 1, data: map[string]any{"a": []any{1, 2, 3}}, expect: []string{"a[0]"}},
		{path: "*", max: 1, data: map[string]any{"a": 1}, expect: []string{"a"}},
		{path: "@.a[?(@.b == 122)].c", max: 1, expect: []string{"@.a[1].c"}},
		{path: "@.a[?(@.b == 122)]", max: 1, expect: []string{"@.a[1]"}},
		{path: "a[1:3].a", noSort: true, expect: []string{"a[1].a", "a[2].a"}},
		{path: "a[1:3].a", max: 1, expect: []string{"a[1].a"}},
		{path: "a[2:0:-1].a", max: 1, expect: []string{"a[2].a"}},
		{path: "a[1:3]", noSort: true, expect: []string{"a[1]", "a[2]"}},
		{path: "a[2:0:-1].a", noSort: true, expect: []string{"a[2].a", "a[1].a"}},
		{path: "a[-2:0:-1]", noSort: true, expect: []string{"a[2]", "a[1]"}},
		{path: "a[5:0:-1]", noSort: true, expect: []string{"a[3]", "a[2]", "a[1]"}},
		{path: "a[1:-7:-1]", noSort: true, expect: []string{"a[1]", "a[0]"}},
		{path: "a[-6:6:]", noSort: true, expect: []string{"a[0]", "a[1]", "a[2]", "a[3]"}},
		{path: "a[2:0:0]", expect: []string{}},
		{path: "a[1,3].a", noSort: true, expect: []string{"a[1].a", "a[3].a"}},
		{path: "a[1,3]", noSort: true, expect: []string{"a[1]", "a[3]"}},
		{path: "a[-1,1]", noSort: true, expect: []string{"a[3]", "a[1]"}},
		{path: "a[-1,1]", max: 1, expect: []string{"a[3]"}},
		{path: "a[1]['a','c']", noSort: true, expect: []string{"a[1].a", "a[1].c"}},
		{
			path:   "$..",
			noSort: true,
			data:   map[string]any{"a": []any{1, map[string]any{"x": 3}}},
			expect: []string{"$", "$.a", "$.a[0]", "$.a[1]", "$.a[1].x"},
		},
		{
			path:   "$..",
			max:    3,
			noSort: true,
			data:   map[string]any{"a": []any{1, map[string]any{"x": 3}}},
			expect: []string{"$", "$.a", "$.a[0]"},
		},
		{
			path: "$..", max: 2, noSort: true, data: []any{map[string]any{"a": 3}}, expect: []string{"$", "$[0]"},
		},
		{path: "$..a", expect: []string{
			"$.a",
			"$.a[0].a",
			"$.a[1].a",
			"$.a[2].a",
			"$.a[3].a",
			"$.b[0].a",
			"$.b[1].a",
			"$.b[2].a",
			"$.b[3].a",
			"$.c[0].a",
			"$.c[1].a",
			"$.c[2].a",
			"$.c[3].a",
			"$.d[0].a",
			"$.d[1].a",
			"$.d[2].a",
			"$.d[3].a",
		}},
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

func TestExprLocateAny(t *testing.T) {
	data := buildTree(4, 3, 0)
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
		if !d.noSort {
			sort.Strings(results)
		}
		diff := alt.Compare(d.expect, results)
		if 0 < len(diff) {
			t.Fatal(testDiffString(d.expect, results, diff))
		}
	}
}

func TestExprLocateNode(t *testing.T) {
	data := alt.Generify(buildTree(4, 3, 0))
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
			locs = x.Locate(alt.Generify(d.data), d.max)
		}
		var results []string
		for _, loc := range locs {
			results = append(results, loc.String())
		}
		if !d.noSort {
			sort.Strings(results)
		}
		diff := alt.Compare(d.expect, results)
		if 0 < len(diff) {
			t.Fatal(testDiffString(d.expect, results, diff))
		}
	}
}

func TestExprLocateOrdered(t *testing.T) {
	data := orderedFromSimple(buildTree(4, 3, 0))
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
			locs = x.Locate(orderedFromSimple(d.data), d.max)
		}
		var results []string
		for _, loc := range locs {
			results = append(results, loc.String())
		}
		if !d.noSort {
			sort.Strings(results)
		}
		diff := alt.Compare(d.expect, results)
		if 0 < len(diff) {
			t.Fatal(testDiffString(d.expect, results, diff))
		}
	}
}

func TestExprLocateReflect(t *testing.T) {
	for i, d := range []*locateData{
		{path: "a", data: &Sample{A: 3, B: "sample"}, expect: []string{"a"}},
		{path: "[1]", data: []int{1, 2, 3}, expect: []string{"[1]"}},
		{path: "['a','b']", data: &Sample{A: 3, B: "sample"}, expect: []string{"a", "b"}},
		{path: "[1,2]", data: []int{1, 2, 3}, expect: []string{"[1]", "[2]"}},
		{path: "[1:2]", data: nil, expect: []string{}},
		{path: "[1:3]", data: []int{1, 2, 3}, expect: []string{"[1]", "[2]"}},
		{path: "[1:3]", max: 1, data: []int{1, 2, 3}, expect: []string{"[1]"}},
		{path: "[2:0:-1]", max: 1, data: []int{1, 2, 3}, expect: []string{"[2]"}},
		{path: "[0:3].b", max: 1, data: []map[string]any{{"a": 1}, {"b": 1}}, expect: []string{"[1].b"}},
		{path: "[2:0:-1].b", max: 1, data: []map[string]any{{"a": 1}, {"b": 1}}, expect: []string{"[1].b"}},
		{path: "$.*", data: nil, expect: []string{}},
		{path: "$.*", data: &Sample{A: 3, B: "sample"}, expect: []string{"$.A", "$.B"}},
		{path: "$.*", max: 1, data: &Sample{A: 3, B: "sample"}, expect: []string{"$.B"}},
		{path: "$.*.a", data: &Any{X: map[string]any{"a": 1}}, expect: []string{"$.X.a"}},
		{path: "$.*.a", max: 1, data: &Any{X: map[string]any{"a": 1}}, expect: []string{"$.X.a"}},
		{path: "$.*", max: 2, data: []int{1, 2, 3}, expect: []string{"$[0]", "$[1]"}},
		{path: "$.*.a", max: 1, data: []map[string]any{{"a": 1}}, expect: []string{"$[0].a"}},
		{path: "$..", data: nil, expect: []string{"$"}},
		{path: "$..", data: &Sample{A: 3, B: "sample"}, expect: []string{"$", "$.A", "$.B"}},
		{path: "$..", max: 2, data: &Sample{A: 3, B: "sample"}, expect: []string{"$", "$.B"}},
		{path: "$..", max: 2, data: []int{1, 2, 3}, expect: []string{"$", "$[0]"}},
		{path: "[0][1]", data: []any{[]int{1, 2, 3}}, expect: []string{"[0][1]"}},
		{path: "[0:2][1]", data: []any{[]int{1, 2, 3}}, expect: []string{"[0][1]"}},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		locs := x.Locate(d.data, d.max)
		var results []string
		for _, loc := range locs {
			results = append(results, loc.String())
		}
		if !d.noSort {
			sort.Strings(results)
		}
		diff := alt.Compare(d.expect, results)
		if 0 < len(diff) {
			t.Fatal(testDiffString(d.expect, results, diff))
		}
	}
}

func TestExprLocateBracket(t *testing.T) {
	data := []any{map[string]any{"b": 1, "c": 2}, []any{1, 2, 3}}
	x := jp.B().N(0).C("b")
	tt.Equal(t, "[0]['b']", x.Locate(data, 0)[0].BracketString())
}
