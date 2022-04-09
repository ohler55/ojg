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
	data   interface{}
	expect bool
}

var (
	hasTestData = []*hasData{
		{path: "", expect: false},
		{path: "$.a.*.b", expect: true},
		{path: "$", expect: true, data: map[string]interface{}{"x": 1}},
		{path: "@", expect: true, data: map[string]interface{}{"x": 1}},
		{path: "$.a.*.b", expect: true, data: firstData1},
		{path: "@.a[0].b", expect: true, data: firstData1},
		{path: "..[0].b", expect: true, data: firstData1},
		{path: "[-1]", expect: true, data: []interface{}{1, 2}},
		{path: "[1,'a']", expect: true, data: []interface{}{1, 2}},
		{path: "[:2]", expect: true, data: []interface{}{1, 2}},
		{path: "a[:-3].b", expect: false, data: firstData1},
		{path: "a[:].b", expect: true, data: firstData1},
		{path: "a[-1:0:-1].b", expect: false, data: firstData1},
		{path: "[?(@ > 1)]", expect: true, data: []interface{}{1, 2}},
		{path: "$[?(@ > 1)]", expect: true, data: []interface{}{1, 2}},
		{path: "[*]", expect: true, data: []interface{}{1, 2}},
		{path: "a.*.*", expect: true, data: firstData1},
		{path: "@.*[0].b", expect: true, data: firstData1},
		{path: "@.a[0]..", expect: true, data: firstData1},
		{path: "..", expect: true, data: []interface{}{1, 2}},
		{path: "..a", expect: false, data: []interface{}{1, 2}},
		{path: "..[1]", expect: true, data: []interface{}{1, true}},
		{path: "a..b", expect: true},
		{path: "[0,'a'][-1,'a']['b',1]", expect: true, data: firstData1},
		{path: "a[-1:2].b", expect: true, data: firstData1},
		{path: "a[-2:2].b", expect: true, data: firstData1},
		{path: "x[:2]", expect: true, data: map[string]interface{}{"x": []interface{}{2, 3}}},
		{path: "[1]", expect: true, data: []int{1, 2, 3}},
		{path: "[-1]", expect: true, data: []int{1, 2, 3}},
		{path: "[-1,'a']", expect: true, data: []int{1, 2, 3}},
		{path: "[::0]", expect: false, data: []interface{}{1, 2, 3}},
		{path: "[10:]", expect: false, data: []interface{}{1, 2, 3}},
		{path: "[:-10:-1]", expect: true, data: []interface{}{1, 2, 3}},
		{path: "[-1:0:-1].x", expect: true, data: []interface{}{
			map[string]interface{}{"x": 1},
			map[string]interface{}{"x": 2},
		}},
		{path: "a.b", expect: false, data: map[string]interface{}{"a": nil}},
		{path: "*.*", expect: false, data: map[string]interface{}{"a": nil}},
		{path: "*.*", expect: false, data: []interface{}{nil}},
		{path: "[0][0]", expect: false, data: []interface{}{nil}},
		{path: "['a','b'].c", expect: false, data: map[string]interface{}{"a": nil}},
		{path: "[1:0:-1].c", expect: false, data: []interface{}{nil, nil}},
		{path: "[0:1][0]", expect: false, data: []interface{}{nil}},
	}
	hasTestReflectData = []*hasData{
		{path: "$.a", expect: true, data: &Sample{A: 3, B: "sample"}},
		{path: "x.a", expect: true, data: map[string]interface{}{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0,'x'].a", expect: true, data: map[string]interface{}{"x": &Sample{A: 3, B: "sample"}}},
		{path: "[0].a", expect: true, data: []interface{}{&Sample{A: 3, B: "sample"}}},
		{path: "$.*", expect: true, data: &One{A: 3}},
		{path: "[*].*", expect: true, data: []*One{{A: 3}}},
		{path: "[*].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[*].a", expect: true, data: []interface{}{&Sample{A: 3, B: "sample"}}},
		{path: "$.*.a", expect: true, data: map[string]interface{}{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: true, data: map[string]interface{}{"x": &Sample{A: 3, B: "sample"}}},
		{path: "$..a", expect: true, data: []interface{}{&Sample{A: 3, B: "sample"}}},
		{path: "$[1:2].a", expect: true, data: []interface{}{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "$[2:1:-1].a", expect: true, data: []interface{}{&One{A: 1}, &One{A: 2}, &One{A: 3}}},
		{path: "[0:-1:2].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "[-1:0:-2].a", expect: true, data: []*One{{A: 1}, {A: 2}, {A: 3}}},
		{path: "$.*[0]", expect: true, data: &Any{X: []interface{}{3}}},
		{path: "$[1:2]", expect: true, data: []int{1, 2, 3}},
		{path: "$[1:1][0]", expect: true, data: []gen.Array{{gen.Int(1)}, {gen.Int(2)}, {gen.Int(3)}}},
		{path: "$.*", expect: false, data: &one},
		{path: "['a',-1]", expect: true, data: []interface{}{1, 2, 3}},
		{path: "['a','b']", expect: false, data: []interface{}{1, 2, 3}},
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
