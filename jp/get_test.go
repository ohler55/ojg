// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

type getData struct {
	path   string
	data   any
	expect []any
}

type Sample struct {
	A int
	B string
}

type NestedSample struct {
	A *bool
	B *string
	C *[]NestedSample
}

type One struct {
	A int
}

type Any struct {
	X any
}

type triple [3]int

var (
	getTestData = []*getData{
		{path: "", expect: []any{}},
		{path: "$.a.*.b", expect: []any{112, 122, 132, 142}},
		{path: "@.b[1].c", expect: []any{223}},
		{path: "..[1].b", expect: []any{122, 222, 322, 422}},
		{path: "[-1]", expect: []any{3}, data: []any{0, 1, 2, 3}},
		{path: "[1,'a']['b',2]['c',3]", expect: []any{133}},
		{path: "a[1::2].a", expect: []any{121, 141}},
		{path: "a[?(@.a > 135)].b", expect: []any{142}},
		{path: "[?(@[1].a > 230)][1].b", expect: []any{322, 422}},
		{path: "[?(@ > 1)]", expect: []any{2, 3}, data: []any{1, 2, 3}},
		{path: "$[?(1==1)]", expect: []any{1, 2, 3}, data: []any{1, 2, 3}},
		{
			path:   "$.*[*].a",
			expect: []any{111, 121, 131, 141, 211, 221, 231, 241, 311, 321, 331, 341, 411, 421, 431, 441},
		},
		{path: `$['\\']`, expect: []any{3}, data: map[string]any{`\`: 3}},
		{path: `$['\x41']`, expect: []any{3}, data: map[string]any{"A": 3}},
		{path: `$['\x4A']`, expect: []any{3}, data: map[string]any{"J": 3}},
		{path: `$['\x4a']`, expect: []any{3}, data: map[string]any{"J": 3}},
		{path: `$['\u03A0']`, expect: []any{3}, data: map[string]any{"Π": 3}},
		{path: `$['\b']`, expect: []any{3}, data: map[string]any{"\b": 3}},
		{path: `$['\t']`, expect: []any{3}, data: map[string]any{"\t": 3}},
		{path: `$['\n']`, expect: []any{3}, data: map[string]any{"\n": 3}},
		{path: `$['\r']`, expect: []any{3}, data: map[string]any{"\r": 3}},
		{path: `$['\f']`, expect: []any{3}, data: map[string]any{"\f": 3}},
		{path: `$["\'"]`, expect: []any{3}, data: map[string]any{"'": 3}},
		{path: `$['\"']`, expect: []any{3}, data: map[string]any{`"`: 3}},

		{path: "$.a[*].y",
			expect: []any{2, 4},
			data: map[string]any{
				"a": []any{
					map[string]any{"x": 1, "y": 2, "z": 3},
					map[string]any{"x": 2, "y": 4, "z": 6},
				},
			},
		},
		{path: "$..x",
			expect: []any{map[string]any{"x": 2}, 1, 2, 3, 4},
			data: map[string]any{
				"o": map[string]any{
					"a": []any{
						map[string]any{"x": 1},
						map[string]any{
							"x": map[string]any{
								"x": 2,
							},
						},
					},
					"x": 3,
				},
				"x": 4,
			},
		},
		{path: "$..[1].x",
			expect: []any{42, 200, 500},
			data: map[string]any{
				"x": []any{0, 1},
				"y": []any{
					map[string]any{"x": 0},
					map[string]any{"x": 42},
				},
				"z": []any{
					[]any{
						map[string]any{"x": 100},
						map[string]any{"x": 200},
						map[string]any{"x": 300},
					},
					[]any{
						map[string]any{"x": 400},
						map[string]any{"x": 500},
						map[string]any{"x": 600},
					},
				},
			},
		},
		{path: "$['a-b']",
			expect: []any{1},
			data:   map[string]any{"a-b": 1, "c-d": 2},
		},
		{path: "$.a..x",
			expect: []any{1, 2, 3, 4, 5},
			data: map[string]any{
				"a": map[string]any{
					"b": []any{
						map[string]any{"x": 1, "y": true},
						map[string]any{"x": 2, "y": false},
						map[string]any{"x": 3, "y": true},
						map[string]any{"x": 4, "y": false},
					},
					"c": map[string]any{"x": 5, "y": nil},
				},
			},
		},
		{path: "a[2].*", expect: []any{131, 132, 133, 134}},
		{path: "[*]", expect: []any{1, 2, 3}, data: []any{1, 2, 3}},
		{path: "$", expect: []any{map[string]any{"x": 1}}, data: map[string]any{"x": 1}},
		{path: "@", expect: []any{map[string]any{"x": 1}}, data: map[string]any{"x": 1}},
		{path: "['x',-1]", expect: []any{3}, data: []any{1, 2, 3}},
		{path: "$[1:3]", expect: []any{2, 3}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[::0]", expect: []any{}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[10:]", expect: []any{}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[:-10:-1]", expect: []any{1}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[1:10]", expect: []any{2, 3, 4, 5, 6}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[-4:-4]", expect: []any{}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[-4:-3]", expect: []any{3}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[-4:2]", expect: []any{}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[-4:3]", expect: []any{3}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[:2]", expect: []any{1, 2}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "$[-4:]", expect: []any{1, 2, 3}, data: []any{1, 2, 3}},
		{path: "$[0:3:1]", expect: []any{1, 2, 3}, data: []any{1, 2, 3, 4, 5}},
		{path: "$[0:4:2]", expect: []any{1, 3}, data: []any{1, 2, 3, 4, 5}},
		{path: "[-4:-1:2]", expect: []any{3, 5}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "[-4:]", expect: []any{1, 2, 3}, data: []any{1, 2, 3}},
		{path: "[-1:1:-2]", expect: []any{4, 6}, data: []any{1, 2, 3, 4, 5, 6}},
		{path: "c[-1:1:-1].a", expect: []any{331, 341}},
		{path: "a[2]..", expect: []any{map[string]any{"a": 131, "b": 132, "c": 133, "d": 134}, 131, 132, 133, 134}},
		{path: "..", expect: []any{[]any{1, 2}, 1, 2}, data: []any{1, 2}},
		{path: "..a", expect: []any{}, data: []any{1, 2}},
		{path: "a..b", expect: []any{112, 122, 132, 142}},
		{path: "[1]", expect: []any{2}, data: []int{1, 2, 3}},
		{path: "[-1]", expect: []any{3}, data: []int{1, 2, 3}},
		{path: "[1]", expect: []any{2}, data: triple{1, 2, 3}},
		{path: "[-1,'a']", expect: []any{3}, data: []int{1, 2, 3}},
		{path: "$[::]", expect: []any{1, 2, 3}, data: []int{1, 2, 3}},
		{path: "[-1,'a'].x",
			expect: []any{2},
			data: []any{
				map[string]any{"x": 1, "y": 2, "z": 3},
				map[string]any{"x": 2, "y": 4, "z": 6},
			},
		},
		{path: "$[1:3:]", expect: []any{2, 3}, data: []any{1, 2, 3, 4, 5}},
		{path: "$[01:03:01]", expect: []any{2, 3}, data: []any{1, 2, 3, 4, 5}},
		{path: "$[:]['x','y']",
			expect: []any{1, 2, 4, 5},
			data: []any{
				map[string]any{"x": 1, "y": 2, "z": 3},
				map[string]any{"x": 4, "y": 5, "z": 6},
			},
		},
		{path: "a.b", expect: []any{}, data: map[string]any{"a": nil}},
		{path: "*.*", expect: []any{}, data: map[string]any{"a": nil}},
		{path: "*.*", expect: []any{}, data: []any{nil}},
		{path: "[0][0]", expect: []any{}, data: []any{nil}},
		{path: "['a','b'].c", expect: []any{}, data: map[string]any{"a": nil}},
		{path: "[1:0:-1].c", expect: []any{}, data: []any{nil, nil}},
		{path: "[0:1][0]", expect: []any{}, data: []any{nil}},
	}
	getTestReflectData = []*getData{
		{path: "['a','b']", expect: []any{"sample", 3}, data: &Sample{A: 3, B: "sample"}},
		{path: "$.*", expect: []any{"sample", 3}, data: &Sample{A: 3, B: "sample"}},
		{path: "$.a", expect: []any{3}, data: &Sample{A: 3, B: "sample"}},
		{path: "x.a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0,'x'].a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0].a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "[*].*", expect: []any{"sample", 3}, data: []*Sample{{A: 3, B: "sample"}}},
		{path: "[*].a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$.*.a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$[1:2].a", expect: []any{2}, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "$[2:1:-1].a", expect: []any{3}, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "[0::2].a", expect: []any{1, 3}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[-1:0:-2].a", expect: []any{3}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[4:0:-2].a", expect: []any{}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "$.*[0]", expect: []any{3}, data: &Any{X: []any{3}}},
		{path: "$[1:2]", expect: []any{2}, data: []int{1, 2, 3}},
		{path: "$[1:2][0]", expect: []any{gen.Int(2)},
			data: []gen.Array{{gen.Int(1)}, {gen.Int(2)}, {gen.Int(3)}}},
		{path: "$[-10:]", expect: []any{1, 2, 3}, data: []int{1, 2, 3}},
		{path: "$[1:-10:-1]", expect: []any{1, 2}, data: []int{1, 2, 3}},
		{path: "$[2:10]", expect: []any{3}, data: []int{1, 2, 3}},
		{
			path:   "$.x[*].c[0].b",
			expect: []any{ptr("1")},
			data:   Any{X: []*NestedSample{{C: &[]NestedSample{{A: ptr(false), B: ptr("1")}, {A: ptr(true), B: ptr("2")}}}}},
		},
		// filter with map
		{
			path:   "$.x[?(@.b=='sample1')].a",
			expect: []any{3},
			data:   map[string]any{"x": []any{map[string]any{"a": 3, "b": "sample1"}}},
		},
		{
			path:   "$.x[?@.b=='sample1'].a",
			expect: []any{3},
			data:   map[string]any{"x": []any{map[string]any{"a": 3, "b": "sample1"}}},
		},
		{
			path:   "$.x[?(@.a==3)].b",
			expect: []any{"sample1"},
			data:   map[string]any{"x": []any{map[string]any{"a": 3, "b": "sample1"}}},
		},
		{
			path:   "$.x[?@.a==3].b",
			expect: []any{"sample1"},
			data:   map[string]any{"x": []any{map[string]any{"a": 3, "b": "sample1"}}},
		},
		// filter with struct
		{
			path:   "$.x[?(@.b=='sample2')].a",
			expect: []any{3},
			data:   Any{X: []*Sample{{A: 3, B: "sample2"}}},
		},
		{
			path:   "$.x[?(@.a==4)].b",
			expect: []any{"sample2"},
			data:   Any{X: []*Sample{{A: 4, B: "sample2"}}},
		},
		{
			path:   "$.x[*].c[?(@.a==true)].b",
			expect: []any{ptr("2")},
			data:   Any{X: []*NestedSample{{C: &[]NestedSample{{A: ptr(false), B: ptr("1")}, {A: ptr(true), B: ptr("2")}}}}},
		},
		{path: "$.*", expect: []any{}, data: &one},
		{path: "['a',-1]", expect: []any{3}, data: []any{1, 2, 3}},
		{path: "['a','b']", expect: []any{}, data: []any{1, 2, 3}},
		{path: "$.*.x", expect: []any{}, data: &Any{X: 5}},
		{path: "$.*.x", expect: []any{}, data: &Any{X: 5}},
		{path: "[0:1].z", expect: []any{}, data: []*Any{nil, {X: 5}}},
		{path: "[0:1].z", expect: []any{}, data: []int{1}},
	}
)

func ptr[T any](v T) *T {
	return &v
}

var (
	firstData1    = map[string]any{"a": []any{map[string]any{"b": 2}}}
	one           = &One{A: 3}
	firstTestData = []*getData{
		{path: "", expect: []any{nil}, data: map[string]any{"x": 1}},
		{path: "$", expect: []any{map[string]any{"x": 1}}, data: map[string]any{"x": 1}},
		{path: "@", expect: []any{map[string]any{"x": 1}}, data: map[string]any{"x": 1}},
		{path: "$.a.*.b", expect: []any{2}, data: firstData1},
		{path: "@.a[0].b", expect: []any{2}, data: firstData1},
		{path: "..[0].b", expect: []any{2}, data: firstData1},
		{path: "[-1]", expect: []any{2}, data: []any{1, 2}},
		{path: "[1,'a']", expect: []any{2}, data: []any{1, 2}},
		{path: "[:2]", expect: []any{1}, data: []any{1, 2}},
		{path: "a[:-3].b", expect: []any{nil}, data: firstData1},
		{path: "a[:].b", expect: []any{2}, data: firstData1},
		{path: "a[-1:0:-1].b", expect: []any{nil}, data: firstData1},
		{path: "[?(@ > 1)]", expect: []any{2}, data: []any{1, 2}},
		{path: "$[?(@ > 1)]", expect: []any{2}, data: []any{1, 2}},
		{path: "[*]", expect: []any{1}, data: []any{1, 2}},
		{path: "a.*.*", expect: []any{2}, data: firstData1},
		{path: "@.*[0].b", expect: []any{2}, data: firstData1},
		{path: "@.a[0]..", expect: []any{2}, data: firstData1},
		{path: "..", expect: []any{1}, data: []any{1, 2}},
		{path: "..a", expect: []any{nil}, data: []any{1, 2}},
		{path: "..[1]", expect: []any{[]any{2}}, data: []any{1, []any{2}}},
		{path: "a..b", expect: []any{112}},
		{path: "[0,'a'][-1,'a']['b',1]", expect: []any{2}, data: firstData1},
		{path: "a[-1:2].b", expect: []any{2}, data: firstData1},
		{path: "a[-2:2].b", expect: []any{2}, data: firstData1},
		{path: "x[:2]", expect: []any{2}, data: map[string]any{"x": []any{2, 3}}},
		{path: "[1]", expect: []any{2}, data: []int{1, 2, 3}},
		{path: "[-1]", expect: []any{3}, data: []int{1, 2, 3}},
		{path: "[-1,'a']", expect: []any{3}, data: []int{1, 2, 3}},
		{path: "[::0]", expect: []any{nil}, data: []any{1, 2, 3}},
		{path: "[10:]", expect: []any{nil}, data: []any{1, 2, 3}},
		{path: "[:-10:-1]", expect: []any{1}, data: []any{1, 2, 3}},
		{path: "[-1:0:-1].x", expect: []any{2}, data: []any{
			map[string]any{"x": 1},
			map[string]any{"x": 2},
		}},
		{path: "a.b", expect: []any{nil}, data: map[string]any{"a": nil}},
		{path: "*.*", expect: []any{nil}, data: map[string]any{"a": nil}},
		{path: "*.*", expect: []any{nil}, data: []any{nil}},
		{path: "[0][0]", expect: []any{nil}, data: []any{nil}},
		{path: "['a','b'].c", expect: []any{nil}, data: map[string]any{"a": nil}},
		{path: "[1:0:-1].c", expect: []any{nil}, data: []any{nil, nil}},
		{path: "[0:1][0]", expect: []any{nil}, data: []any{nil}},
	}
	firstTestReflectData = []*getData{
		{path: "$.a", expect: []any{3}, data: &Sample{A: 3, B: "sample"}},
		{path: "x.a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0,'x'].a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0].a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$.*", expect: []any{3}, data: &One{A: 3}},
		{path: "[*].*", expect: []any{3}, data: []*One{{A: 3}}},
		{path: "[*].a", expect: []any{1}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[*].a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$.*.a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: []any{3}, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: []any{3}, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$[1:2].a", expect: []any{2}, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "$[2:1:-1].a", expect: []any{3}, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "[0:-1:2].a", expect: []any{1}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[-1:0:-2].a", expect: []any{3}, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "$.*[0]", expect: []any{3}, data: &Any{X: []any{3}}},
		{path: "$[1:2]", expect: []any{2}, data: []int{1, 2, 3}},
		{path: "$[1:1][0]", expect: []any{gen.Int(2)},
			data: []gen.Array{{gen.Int(1)}, {gen.Int(2)}, {gen.Int(3)}}},
		{path: "$.*", expect: []any{nil}, data: &one},
		{path: "['a',-1]", expect: []any{3}, data: []any{1, 2, 3}},
		{path: "$.*", expect: []any{nil}, data: &one},
		{path: "['a',-1]", expect: []any{3}, data: []any{1, 2, 3}},
		{path: "['a','b']", expect: []any{nil}, data: []any{1, 2, 3}},
		{path: "$.*.x", expect: []any{nil}, data: &Any{X: 5}},
		{path: "$.*.x", expect: []any{nil}, data: &Any{X: 5}},
		{path: "[0:1].z", expect: []any{nil}, data: []*Any{nil, {X: 5}}},
		{path: "[0:1].z", expect: []any{nil}, data: []int{1}},
	}
)

func TestExprGet(t *testing.T) {
	data := buildTree(4, 3, 0)
	for i, d := range append(getTestData, getTestReflectData...) {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var results []any
		if d.data == nil {
			results = x.Get(data)
		} else {
			results = x.Get(d.data)
		}
		sort.Slice(results, func(i, j int) bool {
			iv, _ := results[i].(int)
			jv, _ := results[j].(int)
			return iv < jv
		})
		tt.Equal(t, d.expect, results, i, " : ", x)
	}
}

func TestExprGetOnNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	for i, d := range getTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var results []any
		if d.data == nil {
			results = x.Get(data)
		} else {
			results = x.Get(alt.Generify(d.data))
		}
		sort.Slice(results, func(i, j int) bool {
			iv, _ := results[i].(gen.Int)
			jv, _ := results[j].(gen.Int)
			return iv < jv
		})
		var expect []any
		for _, n := range d.expect {
			expect = append(expect, alt.Generify(n))
		}
		tt.Equal(t, expect, results, i, " : ", x)
	}
}

func TestExprFirst(t *testing.T) {
	data := buildTree(4, 3, 0)
	for i, d := range append(firstTestData, firstTestReflectData...) {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result any
		if d.data == nil {
			result = x.First(data)
		} else {
			result = x.First(d.data)
		}
		tt.Equal(t, d.expect[0], result, i, " : ", x)
	}
}

func TestExprFirstOnNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	for i, d := range firstTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result any
		if d.data == nil {
			result = x.First(data)
		} else {
			result = x.First(alt.Generify(d.data))
		}
		tt.Equal(t, alt.Generify(d.expect[0]), result, i, " : ", x)
	}
}

func TestExprGetNodes(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	for i, d := range getTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var results []gen.Node
		if d.data == nil {
			results = x.GetNodes(data)
		} else {
			results = x.GetNodes(alt.Generify(d.data))
		}
		sort.Slice(results, func(i, j int) bool {
			iv, _ := results[i].(gen.Int)
			jv, _ := results[j].(gen.Int)
			return iv < jv
		})
		ar := gen.Array{}
		for _, r := range results {
			ar = append(ar, r)
		}
		tt.Equal(t, alt.Generify(d.expect), ar, i, " : ", x, " on ", oj.JSON(d.data), " - ", oj.JSON(results))
	}
}

func TestExprFirstNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	for i, d := range firstTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result gen.Node
		if d.data == nil {
			result = x.FirstNode(data)
		} else {
			result = x.FirstNode(alt.Generify(d.data))
		}
		tt.Equal(t, alt.Generify(d.expect[0]), result, i, " : ", x)
	}
}

func buildTree(size, depth, iv int) any {
	if depth%2 == 0 {
		list := []any{}
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
	obj := map[string]any{}
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

func buildNodeTree(size, depth, iv int) gen.Node {
	if depth%2 == 0 {
		list := gen.Array{}
		for i := 0; i < size; i++ {
			nv := iv*10 + i + 1
			if 1 < depth {
				list = append(list, buildNodeTree(size, depth-1, nv))
			} else {
				list = append(list, gen.Int(nv))
			}
		}
		return list
	}
	obj := gen.Object{}
	for i := 0; i < size; i++ {
		k := string([]byte{'a' + byte(i)})
		nv := iv*10 + i + 1
		if 1 < depth {
			obj[k] = buildNodeTree(size, depth-1, nv)
		} else {
			obj[k] = gen.Int(nv)
		}
	}
	return obj
}

func TestExprGetWildArray(t *testing.T) {
	obj, err := oj.ParseString(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`)
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[*].y")
	ys := x.Get(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y := x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = oj.ParseString(`{"a":[2,4]}`)
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[*]")
	ys = x.Get(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y = x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetWildGenArray(t *testing.T) {
	p := gen.Parser{}
	obj, err := p.Parse([]byte(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`))
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[*].y")
	ys := x.GetNodes(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y := x.FirstNode(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = p.Parse([]byte(`{"a":[2,4]}`))
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[*]")
	ys = x.GetNodes(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y = x.FirstNode(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetUnionArray(t *testing.T) {
	obj, err := oj.ParseString(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`)
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[0,1].y")
	ys := x.Get(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y := x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = oj.ParseString(`{"a":[2,4]}`)
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[0,1]")
	ys = x.Get(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y = x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetUnionGenArray(t *testing.T) {
	p := gen.Parser{}
	obj, err := p.Parse([]byte(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`))
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[0,1].y")
	ys := x.GetNodes(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y := x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = p.Parse([]byte(`{"a":[2,4,6]}`))
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[0,1]")
	ys = x.GetNodes(obj)
	tt.Equal(t, "[2,4]", oj.JSON(ys))
	y = x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetSlice(t *testing.T) {
	obj, err := oj.ParseString(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`)
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[0:1].y")
	ys := x.Get(obj)
	tt.Equal(t, "[2]", oj.JSON(ys))
	y := x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = oj.ParseString(`{"a":[2,4]}`)
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[0:1]")
	ys = x.Get(obj)
	tt.Equal(t, "[2]", oj.JSON(ys))
	y = x.First(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetGenSlice(t *testing.T) {
	p := gen.Parser{}
	obj, err := p.Parse([]byte(`{
  "a":[
    {"x":1,"y":2,"z":3},
    {"x":2,"y":4,"z":6}
  ]
}`))
	tt.Nil(t, err)
	x, _ := jp.ParseString("$.a[0:1].y")
	ys := x.GetNodes(obj)
	tt.Equal(t, "[2]", oj.JSON(ys))
	y := x.FirstNode(obj)
	tt.Equal(t, "2", oj.JSON(y))

	obj, err = p.Parse([]byte(`{"a":[2,4]}`))
	tt.Nil(t, err)
	x, _ = jp.ParseString("$.a[0:1]")
	ys = x.GetNodes(obj)
	tt.Equal(t, "[2]", oj.JSON(ys))
	y = x.FirstNode(obj)
	tt.Equal(t, "2", oj.JSON(y))
}

func TestExprGetBadPath(t *testing.T) {
	type Instance struct {
		ID          string
		BackendName string
		Name        string
		Type        string
		Status      string
		PrivateIP   string
		PublicIP    string
	}
	items := []*Instance{
		{
			Name:        "sds-sds",
			BackendName: "sd",
			ID:          "i-0sdsd44c0",
			PublicIP:    "23.23.23.23",
			PrivateIP:   "12.12.2.2",
			Status:      "Running",
			Type:        "r5d.large",
		},
	}
	expr, err := jp.ParseString("$[*].{}")
	tt.Nil(t, err)
	tt.Equal(t, []any{}, expr.Get(items))
}

func TestFilterAt(t *testing.T) {
	jsondoc := `{
			"item1": {
				"id": "item1",
				"type": "type1",
				"@type": "attype1"
			},
			"item2": {
				"id": "item2",
				"type": "type2",
				"@type": "attype2"
			}
		}`
	store := oj.MustParseString(jsondoc)
	x := jp.MustParseString(`$[?(@['@type']=="attype1")]`)
	result := x.Get(store)
	tt.Equal(t, 1, len(result))
}

func TestAncesterFilter(t *testing.T) {
	json := `{
  "list": {
    "x": "a",
    "y": "b",
    "sublist": [
      {
        "x": "a",
        "y": "d",
        "subs": [
          {
            "x": "a",
            "y": "c"
          }
        ]
      }
    ]
  }
}`
	doc := oj.MustParseString(json)
	x := jp.MustParseString(`$..subs[?(@.x == 'a')].y`)
	result := x.Get(doc)
	tt.Equal(t, []any{"c"}, result)
}

func TestGetFilterOrder(t *testing.T) {
	jsondoc := `[
			{
				"id": "item1",
				"type": "good"
			},
			{
				"id": "item2",
				"type": "good"
			}
		]`
	store := oj.MustParseString(jsondoc)

	x := jp.MustParseString(`$[?(@.type == 'good')]`)
	result := x.Get(store)
	tt.Equal(t, "[{id: item1 type: good} {id: item2 type: good}]", pretty.SEN(result))
	tt.Equal(t, "{id: item1 type: good}", pretty.SEN(x.First(store)))

	x = jp.MustParseString(`$[?(@.type == 'good')].id`)
	result = x.Get(store)
	tt.Equal(t, "[item1 item2]", pretty.SEN(result))
	tt.Equal(t, "item1", pretty.SEN(x.First(store)))
}

func TestGetFilterRoot(t *testing.T) {
	jsondoc := `{
  "key": "item2",
  "data": [
    {"id": "item1", "type": "good1"},
    {"id": "item2", "type": "good2"}
  ]
}`
	store := oj.MustParseString(jsondoc)

	x := jp.MustParseString(`$.data[?(@.id == $.key)]`)
	result := x.Get(store)
	tt.Equal(t, "[{id: item2 type: good2}]", pretty.SEN(result))
}

func TestGetWildReflectOrder(t *testing.T) {
	type Element struct {
		Value string
	}
	type Root struct {
		Elements []Element
	}
	data := Root{
		Elements: []Element{
			{Value: "e1"},
			{Value: "e2"},
			{Value: "e3"},
		},
	}
	path := jp.MustParseString("$.elements[*]")
	tt.Equal(t, "[{value: e1} {value: e2} {value: e3}]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "{value: e1}", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.elements[*].value")
	tt.Equal(t, "[e1 e2 e3]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "e1", pretty.SEN(path.First(data)))
}

func TestGetChildReflectByJsonTag(t *testing.T) {
	type Element struct {
		Value string
	}
	type Root struct {
		Elements    []Element
		AltElements []Element `json:"anyOtherAttributeName"`
	}
	data := Root{
		Elements: []Element{
			{Value: "e1"},
			{Value: "e2"},
			{Value: "e3"},
		},
		AltElements: []Element{
			{Value: "e4"},
			{Value: "e5"},
			{Value: "e6"},
		},
	}
	path := jp.MustParseString("$.elements[*]")
	tt.Equal(t, "[{value: e1} {value: e2} {value: e3}]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "{value: e1}", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.elements[*].value")
	tt.Equal(t, "[e1 e2 e3]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "e1", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.altElements[*]")
	tt.Equal(t, "[{value: e4} {value: e5} {value: e6}]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "{value: e4}", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.altElements[*].value")
	tt.Equal(t, "[e4 e5 e6]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "e4", pretty.SEN(path.First(data)))

	// Get by "json" tag
	path = jp.MustParseString("$.anyOtherAttributeName[*]")
	tt.Equal(t, "[{value: e4} {value: e5} {value: e6}]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "{value: e4}", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.anyOtherAttributeName[*].value")
	tt.Equal(t, "[e4 e5 e6]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "e4", pretty.SEN(path.First(data)))

	// a non-existent attribute in the struct and also in the json (which populated the struct)
	// would still be non-existent as expected
	path = jp.MustParseString("$.nonExistentAttributeName[*].value")
	tt.Equal(t, "[]", pretty.SEN(path.Get(data)))
}

func TestGetChildReflectByJsonTagInEmbeddedStruct(t *testing.T) {
	type Base struct {
		BaseVal string `json:"anyOtherAttributeName"`
	}
	type Element struct {
		Base  // Embedded Struct
		Value string
	}
	type Root struct {
		Elements []Element
	}
	data := Root{
		Elements: []Element{
			{Base: Base{BaseVal: "b1"}, Value: "e1"},
			{Base: Base{BaseVal: "b2"}, Value: "e2"},
			{Base: Base{BaseVal: "b3"}, Value: "e3"},
		},
	}
	path := jp.MustParseString("$.elements[*]")
	tt.Equal(t, "[{baseVal: b1 value: e1} {baseVal: b2 value: e2} {baseVal: b3 value: e3}]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "{baseVal: b1 value: e1}", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.elements[*].value")
	tt.Equal(t, "[e1 e2 e3]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "e1", pretty.SEN(path.First(data)))

	path = jp.MustParseString("$.elements[*].baseVal")
	tt.Equal(t, "[b1 b2 b3]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "b1", pretty.SEN(path.First(data)))

	// Get by "json" tag
	path = jp.MustParseString("$.elements[*].anyOtherAttributeName")
	tt.Equal(t, "[b1 b2 b3]", pretty.SEN(path.Get(data)))
	tt.Equal(t, "b1", pretty.SEN(path.First(data)))

	// a non-existent attribute in the struct and also in the json (which populated the struct)
	// would still be non-existent as expected
	path = jp.MustParseString("$.elements[*].nonExistentAttributeName")
	tt.Equal(t, "[]", pretty.SEN(path.Get(data)))
}

func TestGetChildReflectInEmbeddedStructsResultsInNothing(t *testing.T) {

	// given embedded structures with fields of the same name
	// when trying to find the field by its name
	// then the search will be invalidated by not knowing which of the fields to get

	type A struct {
		attr string
	}
	type B struct {
		attr string
	}
	type C struct {
		A
		B
	}

	a := A{attr: "_a"}
	b := B{attr: "_b"}
	c := C{
		A: a,
		B: b,
	}

	path := jp.MustParseString("$.attr")
	tt.Equal(t, "[]", pretty.SEN(path.Get(c)))
}

func TestGetChildReflectInCyclicGraphEmbeddedStructsResultsInNothing(t *testing.T) {

	// Given a cyclic graph of embedded structs (because of the interface)
	// when trying to find the field by its name (since the name does not exist)
	// then the search will be invalidated and it will not go into infinite loop

	type I interface{}
	type A struct {
		I
		attr string
	}
	type B struct {
		I
		attr string
	}
	type C struct {
		*A
		*B
		I
	}

	a := &A{attr: "_a"}
	b := &B{attr: "_b"}
	a.I = b
	b.I = a
	c := C{
		A: a,
		B: b,
		I: b, // this member will be completely ignored because it is repeated
	}

	path := jp.MustParseString("$.aNonExistentAttribute")
	tt.Equal(t, "[]", pretty.SEN(path.Get(c)))
}

func TestGetDescentReflect(t *testing.T) {
	type A struct {
		X string
	}
	type B struct {
		X string
		A *A
	}
	type C struct {
		A *A
		B *B
	}
	c := C{
		A: &A{X: "A"},
		B: &B{X: "B", A: &A{X: "BA"}},
	}
	path := jp.MustParseString("$..x")
	tt.Equal(t, "[A BA B]", pretty.SEN(path.Get(c)))
	tt.Equal(t, "A", pretty.SEN(path.First(c)))

	tt.Equal(t, "{a: {x: BA} x: B}", pretty.SEN(jp.R().D().First(c)))
}

func TestGetSliceReflect(t *testing.T) {
	src := "$.vals[-3:]"
	data := map[string]any{"vals": []int{10, 20, 30, 40, 50, 60}}
	x := jp.MustParseString(src)
	tt.Equal(t, "[40 50 60]", pretty.SEN(x.Get(data)))

	src = "$.vals[-3:].x"
	data = map[string]any{"vals": []map[string]int{{"x": 10}, {"x": 20}, {"x": 30}, {"x": 40}, {"x": 50}, {"x": 60}}}
	x = jp.MustParseString(src)
	tt.Equal(t, "[40 50 60]", pretty.SEN(x.Get(data)))
}

func TestGetExists(t *testing.T) {
	src := "[?(@.x exists false)]"
	data := []any{
		map[string]any{"x": 1, "y": 2},
		map[string]any{"y": 4},
		map[string]any{"x": 5},
	}
	x := jp.MustParseString(src)
	tt.Equal(t, "[{y: 4}]", pretty.SEN(x.Get(data)))
}

func TestGetKeyedIndexed(t *testing.T) {
	data := &keydex{
		keyed: keyed{
			ordered: ordered{
				entries: []*entry{
					{key: "a", value: 1},
					{key: "b", value: 2},
					{
						key: "c",
						value: &keydex{
							keyed: keyed{
								ordered: ordered{
									entries: []*entry{
										{key: "c1", value: 11},
										{key: "c2", value: 12},
										{key: "c3", value: 13},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect string
		first  bool
	}{
		{src: "$.b", expect: "[2]"},
		{src: "$.c.c2", expect: "[12]"},
		{src: "$[1]", expect: "[2]"},
		{src: "$[2][1]", expect: "[12]"},
		{src: "$.c.*", expect: "[11 12 13]"},
		{src: "$[*][*]", expect: "[11 12 13]"},
		{src: "$..", expect: `[
  1
  2
  [{key: c1 value: 11} {key: c2 value: 12} {key: c3 value: 13}]
  11
  12
  13
  [
    {key: a value: 1}
    {key: b value: 2}
    {
      key: c
      value: [{key: c1 value: 11} {key: c2 value: 12} {key: c3 value: 13}]
    }
  ]
]`},
		{src: "$['a',1]", expect: "[1 2]"},
		{src: "$['c',1].c3", expect: "[13]"},
		{src: "$['a',-1][-1]", expect: "[13]"},
		{src: "$[0:2]", expect: "[1 2]"},
		{src: "$[-4:][-2:-1]", expect: "[12]"},
		{src: "$[-1:0:-1][-2:-1]", expect: "[12]"},
		{src: "$[3:]", expect: "[]"},
		{src: "$[2:0:-1][2:-5:-1]", expect: "[13 12 11]"},
		{src: "$[?(@.c2 == 12)].c2", expect: "[12]"},

		{src: "$.b", expect: "2", first: true},
		{src: "$.c.c2", expect: "12", first: true},
		{src: "$[1]", expect: "2", first: true},
		{src: "$[2][1]", expect: "12", first: true},
		{src: "$.c.*", expect: "11", first: true},
		{src: "$[*][*]", expect: "11", first: true},
		{src: "$..", expect: "1", first: true},
		{src: "$..c2", expect: "12", first: true},
		{src: "$['a',1]", expect: "1", first: true},
		{src: "$['c',1].c3", expect: "13", first: true},
		{src: "$['a',-1][-1]", expect: "13", first: true},
		{src: "$[0:2]", expect: "1", first: true},
		{src: "$[-4:][-2:-1]", expect: "12", first: true},
		{src: "$[-1:0:-1][-2:-1]", expect: "12", first: true},
		{src: "$[3:]", expect: "null", first: true},
		{src: "$[2:0:-1][2:-5:-1]", expect: "13", first: true},
	} {
		x := jp.MustParseString(d.src)
		if d.first {
			tt.Equal(t, d.expect, pretty.SEN(x.First(data)), d.src)
		} else {
			tt.Equal(t, d.expect, pretty.SEN(x.Get(data)), d.src)
		}
	}
}

func TestGetIndexed(t *testing.T) {
	data := &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{
					key: "c",
					value: &indexed{
						ordered: ordered{
							entries: []*entry{
								{key: "c1", value: 11},
								{key: "c2", value: 12},
								{key: "c3", value: 13},
							},
						},
					},
				},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect string
		first  bool
	}{
		{src: "$[1]", expect: "[2]"},
		{src: "$[2][1]", expect: "[12]"},
		{src: "$[*][*]", expect: "[11 12 13]"},
		{src: "$..", expect: `[
  1
  2
  [{key: c1 value: 11} {key: c2 value: 12} {key: c3 value: 13}]
  11
  12
  13
  [
    {key: a value: 1}
    {key: b value: 2}
    {
      key: c
      value: [{key: c1 value: 11} {key: c2 value: 12} {key: c3 value: 13}]
    }
  ]
]`},
		{src: "$[1,2][-2,0]", expect: "[12 11]"},

		{src: "$[1]", expect: "2", first: true},
		{src: "$[2][1]", expect: "12", first: true},
		{src: "$[*][*]", expect: "11", first: true},
		{src: "$..", expect: "1", first: true},
		{src: "$..[1]", expect: "12", first: true},
		{src: "$[1,2][-2,0]", expect: "12", first: true},
	} {
		x := jp.MustParseString(d.src)
		if d.first {
			tt.Equal(t, d.expect, pretty.SEN(x.First(data)), d.src)
		} else {
			tt.Equal(t, d.expect, pretty.SEN(x.Get(data)), d.src)
		}
	}
}

func TestGetKeyed(t *testing.T) {
	data := &keyed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{
					key: "c",
					value: &keyed{
						ordered: ordered{
							entries: []*entry{
								{key: "c1", value: 11},
								{key: "c2", value: 12},
								{key: "c3", value: 13},
							},
						},
					},
				},
			},
		},
	}
	x := jp.MustParseString("$[?(@.c2 == 12)].c2")
	tt.Equal(t, "[12]", pretty.SEN(x.Get(data)))
}

func TestGetKeyedIndexedReflect(t *testing.T) {
	data := &keydex{
		keyed: keyed{
			ordered: ordered{
				entries: []*entry{
					{key: "a", value: Any{X: 1}},
					{key: "b", value: Any{X: 2}},
					{key: "c", value: Any{X: 3}},
				},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect string
		first  bool
	}{
		{src: "$.b.x", expect: "[2]"},
		{src: "$[1].x", expect: "[2]"},
		{src: "$.c.*", expect: "[3]"},
		{src: "$[*][*]", expect: "[1 2 3]"},
		{src: "$.*.x", expect: "[1 2 3]"},
		{src: "$..", expect: `[
  {x: 1}
  {x: 2}
  {x: 3}
  1
  2
  3
  [
    {key: a value: {type: Any x: 1}}
    {key: b value: {type: Any x: 2}}
    {key: c value: {type: Any x: 3}}
  ]
]`},
		{src: "$..x", expect: "[1 2 3]"},
		{src: "$['a',1].x", expect: "[1 2]"},
		{src: "$['a',-1].x", expect: "[1 3]"},
		{src: "$[0:2].x", expect: "[1 2]"},
		{src: "$[-4:].x", expect: "[1 2 3]"},
		{src: "$[-1:0:-1].x", expect: "[3 2]"},
		{src: "$[3:].x", expect: "[]"},
		{src: "$[2:0:-1].x", expect: "[3 2]"},

		{src: "$.b.x", expect: "2", first: true},
		{src: "$[1].x", expect: "2", first: true},
		{src: "$.c.*", expect: "3", first: true},
		{src: "$[*][*]", expect: "1", first: true},
		{src: "$.*.x", expect: "1", first: true},
		{src: "$..", expect: "{x: 1}", first: true},
		{src: "$..x", expect: "1", first: true},
		{src: "$['a',1].x", expect: "1", first: true},
		{src: "$['a',-1].x", expect: "1", first: true},
		{src: "$[0:2].x", expect: "1", first: true},
		{src: "$[-4:].x", expect: "1", first: true},
		{src: "$[-1:0:-1].x", expect: "3", first: true},
		{src: "$[3:].x", expect: "null", first: true},
		{src: "$[2:0:-1].x", expect: "3", first: true},
	} {
		x := jp.MustParseString(d.src)
		if d.first {
			tt.Equal(t, d.expect, pretty.SEN(x.First(data)), d.src)
		} else {
			tt.Equal(t, d.expect, pretty.SEN(x.Get(data)), d.src)
		}
	}
}

func TestGetIndexedReflect(t *testing.T) {
	data := &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: Any{X: 1}},
				{key: "b", value: Any{X: 2}},
				{key: "c", value: Any{X: 3}},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect string
		first  bool
	}{
		{src: "$..", expect: `[
  {x: 1}
  {x: 2}
  {x: 3}
  1
  2
  3
  [
    {key: a value: {type: Any x: 1}}
    {key: b value: {type: Any x: 2}}
    {key: c value: {type: Any x: 3}}
  ]
]`},
		{src: "$..x", expect: "[1 2 3]"},
		{src: "$.*.x", expect: "[1 2 3]"},

		{src: "$..", expect: "{x: 1}", first: true},
		{src: "$..x", expect: "1", first: true},
		{src: "$.*.x", expect: "1", first: true},
	} {
		x := jp.MustParseString(d.src)
		if d.first {
			tt.Equal(t, d.expect, pretty.SEN(x.First(data)), d.src)
		} else {
			tt.Equal(t, d.expect, pretty.SEN(x.Get(data)), d.src)
		}
	}
}

func TestGetStructMap(t *testing.T) {
	type A struct {
		X map[string]any
	}
	data := A{X: map[string]any{"a": 1}}
	x := jp.MustParseString("$..a")
	tt.Equal(t, "[1]", pretty.SEN(x.Get(data)))
	tt.Equal(t, "1", pretty.SEN(x.First(data)))
}

func TestGetMultiFilter(t *testing.T) {
	data := map[string]any{
		"a": []any{
			map[string]any{
				"b": []any{
					map[string]any{"c": 1},
					map[string]any{"c": 2},
					map[string]any{"c": 3},
				},
			},
		},
	}
	// A match on any c value should return non nil.
	x := jp.MustParseString("a[?(@.b[*].c == 1)].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))

	x = jp.MustParseString("a[?(@.b[*].c == 2)].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))

	x = jp.MustParseString("a[?(@.b[*].c == 3)].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))

	x = jp.MustParseString("a[?(@.b[*].c == 4)].b[0]")
	tt.Equal(t, "null", pretty.SEN(x.First(data)))

	data = map[string]any{
		"a": []any{
			map[string]any{
				"b": []any{
					map[string]any{"c": 1},
				},
			},
		},
	}
	x = jp.MustParseString("a[?(@.b[*].c == 1)].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))

	data = map[string]any{
		"a": []any{
			map[string]any{
				"b": []any{},
			},
		},
	}
	x = jp.MustParseString("a[?(@.b[*].c == 1)]")
	tt.Equal(t, "null", pretty.SEN(x.First(data)))

	// Make sure normalization works.
	m := map[string]any{"c": 2}
	data = map[string]any{
		"a": []any{
			map[string]any{
				"b": []any{
					map[string]any{"c": 1},
					m,
				},
			},
		},
	}
	x = jp.MustParseString("a[?(@.b[*].c == 2)].b[0]")
	for _, v := range []any{
		int8(2),
		int16(2),
		int32(2),
		int64(2),
		uint(2),
		uint8(2),
		uint16(2),
		uint32(2),
		uint64(2),
		gen.Int(2),
	} {
		m["c"] = v
		tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))
	}
	x = jp.MustParseString("a[?(@.b[*].c == 2.5)].b[0]")
	for _, v := range []any{
		float32(2.5),
		float64(2.5),
		gen.Float(2.5),
	} {
		m["c"] = v
		tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))
	}
	m["c"] = gen.Bool(true)
	x = jp.MustParseString("a[?(@.b[*].c == true)].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))

	m["c"] = gen.String("xyz")
	x = jp.MustParseString("a[?(@.b[*].c == 'xyz')].b[0]")
	tt.Equal(t, "{c: 1}", pretty.SEN(x.First(data)))
}

type Key string
type X struct {
	Y Key
}

func TestGetNonStringKey(t *testing.T) {
	data := map[Key]*X{Key("a"): {Y: Key("xyz")}}
	x := jp.MustParseString("$.a.y")
	tt.Equal(t, "[xyz]", pretty.SEN(x.Get(data)))
}

func TestGetKeyMismatch(t *testing.T) {
	data := map[int]*X{1: {Y: Key("xyz")}}
	x := jp.MustParseString("$.a.y")
	tt.Equal(t, "[]", pretty.SEN(x.Get(data)))
}

func TestGetAncestorFilter(t *testing.T) {
	data := map[string]any{
		"array": []any{
			map[string]any{"x": 1, "y": 1},
			map[string]any{"y": 2},
			map[string]any{"child": map[string]any{"x": 1, "y": 3}},
		},
	}
	x := jp.MustParseString("$..[?(@.x)]")
	// x := jp.MustParseString("$..[?(@.x exists true)]")
	tt.Equal(t, "[{x: 1 y: 3} {x: 1 y: 1}]", pretty.SEN(x.Get(data)))
}
