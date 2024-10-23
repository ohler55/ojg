// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type simple int

type walkData struct {
	path  string
	data  string
	xpath string
	nodes string
}

func (s simple) Simplify() any {
	return map[string]any{"x": int64(s)}
}

type other int

func TestWalk(t *testing.T) {
	data := map[string]any{"a": []any{1, 2, 3}, "b": nil, "c": simple(4), "d": other(5)}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value any) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	tt.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b $.c $.c.x $.d]`, string(sen.Bytes(paths)))

	leaves := map[string]any{}
	jp.Walk(data, func(path jp.Expr, value any) {
		leaves[path[1:].String()] = value
	}, true)
	tt.Equal(t, `{"a[0]": 1 "a[1]": 2 "a[2]": 3 b: null c.x: 4 d: 5}`, pretty.SEN(leaves))
}

func TestWalkNode(t *testing.T) {
	data := gen.Object{"a": gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}, "b": nil}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value any) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	tt.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b]`, string(sen.Bytes(paths)))
}

var (
	walkTestData = []*walkData{
		{path: "a", data: "{a:1 b:{c:3}}", xpath: "a", nodes: "1"},
		{path: "b.c", data: "{a:1 b:{c:3}}", xpath: "b.c", nodes: "3"},
		{path: "[0]", data: "[1 [2 3]]", xpath: "[0]", nodes: "1"},
		{path: "[1][1]", data: "[1 [2 3]]", xpath: "[1][1]", nodes: "3"},
		{path: "*", data: "[1 [2 3]]", xpath: "[0] [1]", nodes: "1 [2 3]"},
		{path: "*.*", data: "[1 [2 3]]", xpath: "[1][0] [1][1]", nodes: "2 3"},
		{path: "*", data: "{a:1 b:{c:3}}", xpath: "a b", nodes: "1 {c:3}"},
		{path: "*.*", data: "{a:1 b:{c:3}}", xpath: "b.c", nodes: "3"},
	}
)

func TestExprWalkAny(t *testing.T) {
	testExprWalk(t, false)
}

func TestExprWalkGen(t *testing.T) {
	testExprWalk(t, true)
}

func testExprWalk(t *testing.T, generic bool) {
	opt := ojg.Options{Sort: true, Indent: 0}
	for i, wd := range walkTestData {
		x := jp.MustParseString(wd.path)
		var (
			ps   []string
			ns   []string
			data any
		)
		data = sen.MustParse([]byte(wd.data))
		if generic {
			data = alt.Generify(data)
		}
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		sort.Strings(ps)
		sort.Strings(ns)
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkDev(t *testing.T) {
	data := sen.MustParse([]byte("{a:1 b:{c:3}}"))
	x := jp.W()
	x.Walk(data, func(path jp.Expr, nodes []any) {
		fmt.Printf("*** %s %v\n", path, nodes)
	})
}
