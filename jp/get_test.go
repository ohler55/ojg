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

type One struct {
	A int
}

type Any struct {
	X any
}

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
		{path: "$.*[*].a", expect: []any{111, 121, 131, 141, 211, 221, 231, 241, 311, 321, 331, 341, 411, 421, 431, 441}},
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
		// filter with map
		{
			path:   "$.x[?(@.b=='sample1')].a",
			expect: []any{3},
			data:   map[string]any{"x": []any{map[string]any{"a": 3, "b": "sample1"}}},
		},
		{
			path:   "$.x[?(@.a==3)].b",
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
		{path: "$.*", expect: []any{}, data: &one},
		{path: "['a',-1]", expect: []any{3}, data: []any{1, 2, 3}},
		{path: "['a','b']", expect: []any{}, data: []any{1, 2, 3}},
		{path: "$.*.x", expect: []any{}, data: &Any{X: 5}},
		{path: "$.*.x", expect: []any{}, data: &Any{X: 5}},
		{path: "[0:1].z", expect: []any{}, data: []*Any{nil, {X: 5}}},
		{path: "[0:1].z", expect: []any{}, data: []int{1}},
	}
)

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
