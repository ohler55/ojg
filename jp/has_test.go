// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

type hasData struct {
	path   string
	data   any
	expect bool
}

var (
	hasTestData = []*hasData{
		{path: "", expect: false},
		{path: "$.a.*.b", expect: true},
		{path: "$", expect: true, data: map[string]any{"x": 1}},
		{path: "@", expect: true, data: map[string]any{"x": 1}},
		{path: "$.a.*.b", expect: true, data: firstData1},
		{path: "@.a[0].b", expect: true, data: firstData1},
		{path: "..[0].b", expect: true, data: firstData1},
		{path: "[-1]", expect: true, data: []any{1, 2}},
		{path: "[1,'a']", expect: true, data: []any{1, 2}},
		{path: "[:2]", expect: true, data: []any{1, 2}},
		{path: "a[:-3].b", expect: false, data: firstData1},
		{path: "a[:].b", expect: true, data: firstData1},
		{path: "a[-1:0:-1].b", expect: false, data: firstData1},
		{path: "[?(@ > 1)]", expect: true, data: []any{1, 2}},
		{path: "$[?(@ > 1)]", expect: true, data: []any{1, 2}},
		{path: "[*]", expect: true, data: []any{1, 2}},
		{path: "a.*.*", expect: true, data: firstData1},
		{path: "@.*[0].b", expect: true, data: firstData1},
		{path: "@.a[0]..", expect: true, data: firstData1},
		{path: "..", expect: true, data: []any{1, 2}},
		{path: "..a", expect: false, data: []any{1, 2}},
		{path: "..[1]", expect: true, data: []any{1, true}},
		{path: "a..b", expect: true},
		{path: "[0,'a'][-1,'a']['b',1]", expect: true, data: firstData1},
		{path: "a[-1:2].b", expect: true, data: firstData1},
		{path: "a[-2:2].b", expect: true, data: firstData1},
		{path: "x[:2]", expect: true, data: map[string]any{"x": []any{2, 3}}},
		{path: "[1]", expect: true, data: []int{1, 2, 3}},
		{path: "[-1]", expect: true, data: []int{1, 2, 3}},
		{path: "[-1,'a']", expect: true, data: []int{1, 2, 3}},
		{path: "[::0]", expect: false, data: []any{1, 2, 3}},
		{path: "[10:]", expect: false, data: []any{1, 2, 3}},
		{path: "[:-10:-1]", expect: true, data: []any{1, 2, 3}},
		{path: "[-1:0:-1].x", expect: true, data: []any{
			map[string]any{"x": 1},
			map[string]any{"x": 2},
		}},
		{path: "a.b", expect: false, data: map[string]any{"a": nil}},
		{path: "*.*", expect: false, data: map[string]any{"a": nil}},
		{path: "*.*", expect: false, data: []any{nil}},
		{path: "[0][0]", expect: false, data: []any{nil}},
		{path: "['a','b'].c", expect: false, data: map[string]any{"a": nil}},
		{path: "[1:0:-1].c", expect: false, data: []any{nil, nil}},
		{path: "[0:1][0]", expect: false, data: []any{nil}},
	}
	hasTestReflectData = []*hasData{
		{path: "$.a", expect: true, data: &Sample{A: 3, B: "sample"}},
		{path: "x.a", expect: true, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0,'x'].a", expect: true, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0].a", expect: true, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$.*", expect: true, data: &One{A: 3}},
		{path: "[*].*", expect: true, data: []*One{{A: 3}}},
		{path: "[*].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[*].a", expect: true, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$.*.a", expect: true, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: true, data: map[string]any{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: true, data: []any{&Sample{A: 3, B: "sample"}}},
		{path: "$[1:2].a", expect: true, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "$[2:1:-1].a", expect: true, data: []any{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "[0:-1:2].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[-1:0:-2].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "$.*[0]", expect: true, data: &Any{X: []any{3}}},
		{path: "$[1:2]", expect: true, data: []int{1, 2, 3}},
		{path: "$[1:1][0]", expect: true, data: []gen.Array{{gen.Int(1)}, {gen.Int(2)}, {gen.Int(3)}}},
		{path: "$.*", expect: false, data: &one},
		{path: "['a',-1]", expect: true, data: []any{1, 2, 3}},
		{path: "['a','b']", expect: false, data: []any{1, 2, 3}},
		{path: "$.*.x", expect: false, data: &Any{X: 5}},
		{path: "$.*.x", expect: false, data: &Any{X: 5}},
		{path: "[0:1].z", expect: false, data: []*Any{nil, {X: 5}}},
		{path: "[0:1].z", expect: false, data: []int{1}},
	}
)

func TestExprHas(t *testing.T) {
	data := buildTree(4, 3, 0)
	for i, d := range append(hasTestData, hasTestReflectData...) {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result bool
		if d.data == nil {
			result = x.Has(data)
		} else {
			result = x.Has(d.data)
		}
		tt.Equal(t, d.expect, result, i, " : ", x)
	}
}

func TestExprHasNode(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	for i, d := range hasTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result bool
		if d.data == nil {
			result = x.Has(data)
		} else {
			result = x.Has(alt.Generify(d.data))
		}
		tt.Equal(t, d.expect, result, i, " : ", x)
	}
}

func TestHasKeyedIndexed(t *testing.T) {
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
		expect bool
	}{
		{src: "$.b", expect: true},
		{src: "$.c.c2", expect: true},
		{src: "$[1]", expect: true},
		{src: "$[2][1]", expect: true},
		{src: "$.c.*", expect: true},
		{src: "$[*][*]", expect: true},
		{src: "$..", expect: true},
		{src: "$..c2", expect: true},
		{src: "$['a',1]", expect: true},
		{src: "$['c',1].c3", expect: true},
		{src: "$['a',-1][-1]", expect: true},
		{src: "$[0:2]", expect: true},
		{src: "$[-4:][-2:-1]", expect: true},
		{src: "$[-1:0:-1][-2:-1]", expect: true},
		{src: "$[3:]", expect: false},
		{src: "$[2:0:-1][2:-5:-1]", expect: true},
		{src: "$[?(@.c2 == 12)].c2", expect: true},
	} {
		x := jp.MustParseString(d.src)
		tt.Equal(t, d.expect, x.Has(data), d.src)
	}
}

func TestHasIndexed(t *testing.T) {
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
		expect bool
	}{
		{src: "$[1]", expect: true},
		{src: "$[2][1]", expect: true},
		{src: "$[*][*]", expect: true},
		{src: "$..", expect: true},
		{src: "$[1,2][-2,0]", expect: true},
		{src: "$..[1]", expect: true},
	} {
		x := jp.MustParseString(d.src)
		tt.Equal(t, d.expect, x.Has(data), d.src)
	}
}

func TestHasKeyedIndexedReflect(t *testing.T) {
	data := &keydex{
		keyed: keyed{
			ordered: ordered{
				entries: []*entry{
					{key: "a", value: Any{X: []any{1}}},
					{key: "b", value: Any{X: 2}},
					{key: "c", value: Any{X: 3}},
				},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect bool
	}{
		{src: "$.b.x", expect: true},
		{src: "$[1].x", expect: true},
		{src: "$.c.*", expect: true},
		{src: "$[*][*]", expect: true},
		{src: "$.*.x", expect: true},
		{src: "$..", expect: true},
		{src: "$..x", expect: true},
		{src: "$['a',1].x", expect: true},
		{src: "$[-1,2].x", expect: true},
		{src: "$['a',-1].x", expect: true},
		{src: "$[0:2].x", expect: true},
		{src: "$[-4:].x", expect: true},
		{src: "$[-1:0:-1].x", expect: true},
		{src: "$[3:].x", expect: false},
		{src: "$[2:0:-1].x", expect: true},
	} {
		x := jp.MustParseString(d.src)
		tt.Equal(t, d.expect, x.Has(data), d.src)
	}
}

func TestHasIndexedReflect(t *testing.T) {
	data := &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: Any{X: []any{1}}},
				{key: "b", value: Any{X: 2}},
				{key: "c", value: Any{X: 3}},
			},
		},
	}
	for _, d := range []*struct {
		src    string
		expect bool
	}{
		{src: "$..", expect: true},
		{src: "$..x", expect: true},
		{src: "$.*.x", expect: true},
		{src: "$.*", expect: true},
	} {
		x := jp.MustParseString(d.src)
		tt.Equal(t, d.expect, x.Has(data), d.src)
	}
}
