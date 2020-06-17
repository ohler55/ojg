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
	"github.com/ohler55/ojg/tt"
)

type testData struct {
	path   string
	data   interface{}
	expect []interface{}
}

func xTestDev(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	data = gen.Object{"x": gen.Int(1)}
	//fmt.Println(oj.JSON(data, 2))
	x, err := jp.ParseString("$")
	tt.Nil(t, err)
	results := x.Get(data)
	sort.Slice(results, func(i, j int) bool {
		iv, _ := results[i].(gen.Int)
		jv, _ := results[j].(gen.Int)
		return iv < jv
	})
	fmt.Printf("*** %s\n", oj.JSON(results))
}

func inpectExpr(x jp.Expr) {
	fmt.Printf("*** Expr %d\n", len(x))
	for i, f := range x {
		fmt.Printf("***   %d: %T %s\n", i, f, f)
	}
}

var getTestData = []*testData{
	{path: "", expect: []interface{}{}},
	{path: "$.a.*.b", expect: []interface{}{112, 122, 132, 142}},
	{path: "@.b[1].c", expect: []interface{}{223}},
	{path: "..[1].b", expect: []interface{}{122, 222, 322, 422}},
	{path: "[-1]", expect: []interface{}{3}, data: []interface{}{0, 1, 2, 3}},
	{path: "..[1].b", expect: []interface{}{122, 222, 322, 422}},
	{path: "[1,'a']['b',2]['c',3]", expect: []interface{}{133}},
	{path: "a[1:-1:2].a", expect: []interface{}{121, 141}},
	{path: "a[?(@.a > 135)].b", expect: []interface{}{142}},
	{path: "$.*[*].a", expect: []interface{}{111, 121, 131, 141, 211, 221, 231, 241, 311, 321, 331, 341, 411, 421, 431, 441}},
	{path: "a[2].*", expect: []interface{}{131, 132, 133, 134}},
	{path: "[*]", expect: []interface{}{1, 2, 3}, data: []interface{}{1, 2, 3}},
	{path: "$", expect: []interface{}{map[string]interface{}{"x": 1}}, data: map[string]interface{}{"x": 1}},
	{path: "@", expect: []interface{}{map[string]interface{}{"x": 1}}, data: map[string]interface{}{"x": 1}},
	{path: "['x',-1]", expect: []interface{}{3}, data: []interface{}{1, 2, 3}},
	{path: "[-4:-1:2]", expect: []interface{}{3, 5}, data: []interface{}{1, 2, 3, 4, 5, 6}},
	{path: "[-4:]", expect: []interface{}{}, data: []interface{}{1, 2, 3}},
	{path: "[-1:1:-2]", expect: []interface{}{4, 6}, data: []interface{}{1, 2, 3, 4, 5, 6}},
	{path: "c[-1:1:-1].a", expect: []interface{}{331, 341}},
	{path: "a[2]..", expect: []interface{}{map[string]interface{}{"a": 131, "b": 132, "c": 133, "d": 134}, 131, 132, 133, 134}},
	{path: "..", expect: []interface{}{[]interface{}{1, 2}, 1, 2}, data: []interface{}{1, 2}},
}

var (
	firstData1 = map[string]interface{}{"a": []interface{}{map[string]interface{}{"b": 2}}}

	firstTestData = []*testData{
		{path: "", expect: []interface{}{nil}, data: map[string]interface{}{"x": 1}},
		{path: "$", expect: []interface{}{map[string]interface{}{"x": 1}}, data: map[string]interface{}{"x": 1}},
		{path: "@", expect: []interface{}{map[string]interface{}{"x": 1}}, data: map[string]interface{}{"x": 1}},
		{path: "$.a.*.b", expect: []interface{}{2}, data: firstData1},
		{path: "@.a[0].b", expect: []interface{}{2}, data: firstData1},
		{path: "..[0].b", expect: []interface{}{2}, data: firstData1},
		{path: "[-1]", expect: []interface{}{2}, data: []interface{}{1, 2}},
		{path: "[1,'a']", expect: []interface{}{2}, data: []interface{}{1, 2}},
		{path: "[:2]", expect: []interface{}{1}, data: []interface{}{1, 2}},
		//{path: "[?(@ > 1)]", expect: []interface{}{2}, data: []interface{}{int64(1), int64(2)}}, // TBD change to simple int
		{path: "[*]", expect: []interface{}{1}, data: []interface{}{1, 2}},
		{path: "a.*.*", expect: []interface{}{2}, data: firstData1},
	}
)

func TestExprGet(t *testing.T) {
	data := buildTree(4, 3, 0)
	for i, d := range getTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var results []interface{}
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
		tt.Equal(t, d.expect, results)
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
		var results []interface{}
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
		var expect []interface{}
		for _, n := range d.expect {
			expect = append(expect, alt.Generify(n))
		}
		tt.Equal(t, expect, results)
	}
}

func TestExprFirst(t *testing.T) {
	for i, d := range firstTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		result := x.First(d.data)
		tt.Equal(t, d.expect[0], result)
	}
}

func TestExprFirstOnNode(t *testing.T) {
	for i, d := range firstTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		result := x.First(alt.Generify(d.data))
		tt.Equal(t, alt.Generify(d.expect[0]), result)
	}
}

func xTestExprFirstNode(t *testing.T) {
	for i, d := range firstTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		result := x.FirstNode(alt.Generify(d.data))
		tt.Equal(t, alt.Generify(d.expect[0]), result)
	}
}

func xTestExprGetNodes(t *testing.T) {
	data := buildNodeTree(4, 3, 0)
	// TBD replace with getTestData
	for i, d := range []*testData{
		{path: "", expect: []interface{}{}},
		{path: "$.a.*.b", expect: []interface{}{112, 122, 132, 142}},
		{path: "$.b[1].c", expect: []interface{}{223}},
		{path: "..[1].b", expect: []interface{}{122, 222, 322, 422}},
		{path: "[1,'a']['b',2]['c',3]", expect: []interface{}{133}},
		//{path: "a[1:-1:2].a", expect: []interface{}{121, 141}},
		//{path: "a[?(@.a > 135)].b", expect: []interface{}{142}},
	} {
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
		tt.Equal(t, alt.Generify(d.expect), ar)
	}
}

func buildTree(size, depth, iv int) interface{} {
	if depth%2 == 0 {
		list := []interface{}{}
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
	obj := map[string]interface{}{}
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

// TBD add object for relfection tests
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
