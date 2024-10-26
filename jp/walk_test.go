// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"bytes"
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
		{path: "[-1][-3]", data: "[1 [2 3]]", xpath: "", nodes: ""},
		{path: "[1][1]", data: "[1 [2 3]]", xpath: "[1][1]", nodes: "3"},
		{path: "*", data: "[1 [2 3]]", xpath: "[0] [1]", nodes: "1 [2 3]"},
		{path: "*.*", data: "[1 [2 3]]", xpath: "[1][0] [1][1]", nodes: "2 3"},
		{path: "*", data: "{a:1 b:{c:3}}", xpath: "a b", nodes: "1 {c:3}"},
		{path: "*.*", data: "{a:1 b:{c:3}}", xpath: "b.c", nodes: "3"},
		{path: "@", data: "{a:1}", xpath: "", nodes: "{a:1}"},
		{path: "@.a", data: "{a:1 b:{c:3}}", xpath: "a", nodes: "1"},
		{path: "$", data: "{a:1}", xpath: "", nodes: "{a:1}"},
		{path: "$.a", data: "{a:1 b:{c:3}}", xpath: "a", nodes: "1"},
		{path: "[1,'a']", data: "{a:1 b:{c:3}}", xpath: "a", nodes: "1"},
		{path: "['b','a']['c',2]", data: "{a:1 b:{c:3}}", xpath: "b.c", nodes: "3"},
		{path: "[0,'a']", data: "[1 [2 3]]", xpath: "[0]", nodes: "1"},
		{path: "[1,2][0,4]", data: "[1 [2 3]]", xpath: "[1][0]", nodes: "2"},
		{path: "[0:4:2]", data: "[1 2 3 4 5 6]", xpath: "[0] [2]", nodes: "1 3"},
		{path: "[0:4:0]", data: "[1 2 3 4 5 6]", xpath: "", nodes: ""},
		{path: "[4:0:-2]", data: "[1 2 3 4 5 6]", xpath: "[4] [2]", nodes: "5 3"},
		{path: "[?(@.x == 1)]", data: "[{x:0}{x:1}]", xpath: "[1]", nodes: "{x:1}"},
		{path: "[?(@.x == 1)].x", data: "[{x:0}{x:1}]", xpath: "[1].x", nodes: "1"},
		{path: "[?(@.x == 1)]", data: "{y:{x:0} z:{x:1}}", xpath: "z", nodes: "{x:1}"},
		{path: "[?(@.x == 1)].x", data: "{y:{x:0} z:{x:1}}", xpath: "z.x", nodes: "1"},
		{path: "..", data: "{a:1 b:{c:3}}", xpath: "a b b.c", nodes: "1 {c:3} 3"},
		{path: "..", data: "[1 [2 3]]", xpath: "[0] [1] [1][0] [1][1]", nodes: "1 [2 3] 2 3"},
		{path: "a.b.c", data: "{a:1 b:{c:3}}", xpath: "", nodes: ""},
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
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkBracket(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	data := sen.MustParse([]byte("{a:1}"))
	x := jp.B()
	var (
		ps []byte
		ns []byte
	)
	x.Walk(data, func(path jp.Expr, nodes []any) {
		ps = append(ps, path.String()...)
		ns = append(ns, bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})...)
	})
	tt.Equal(t, "", string(ps))
	tt.Equal(t, "{a:1}", string(ns))

	x = jp.B().C("a")
	ps = ps[:0]
	ns = ns[:0]
	x.Walk(data, func(path jp.Expr, nodes []any) {
		ps = append(ps, path.String()...)
		ns = append(ns, bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})...)
	})
	tt.Equal(t, "a", string(ps))
	tt.Equal(t, "1", string(ns))
}

func TestExprWalkIndexed(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	data := &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{
					key: "b",
					value: &indexed{
						ordered: ordered{
							entries: []*entry{
								{key: "b2", value: 2},
								{key: "b3", value: 3},
							},
						},
					},
				},
			},
		},
	}
	for i, wd := range []*walkData{
		{path: "[0]", xpath: "[0]", nodes: "1"},
		{path: "[1][1]", xpath: "[1][1]", nodes: "3"},
		{path: "[-1][-3]", xpath: "", nodes: ""},
		{path: "*", xpath: "[0] [1]", nodes: "1 [{key:b2 value:2}{key:b3 value:3}]"},
		{path: "*.*", xpath: "[1][0] [1][1]", nodes: "2 3"},
		{path: "[0:4:2]", xpath: "[0]", nodes: "1"},
		{path: "[?(@ == 1)]", xpath: "[0]", nodes: "1"},
		{path: "[?(@[0] == 2)][1]", xpath: "[1][1]", nodes: "3"},
		{path: "..", xpath: "[0] [1] [1][0] [1][1]", nodes: "1 [{key:b2 value:2}{key:b3 value:3}] 2 3"},
	} {
		x := jp.MustParseString(wd.path)
		var (
			ps []string
			ns []string
		)
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkKeyed(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	data := &keyed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{
					key: "b",
					value: &keyed{
						ordered: ordered{
							entries: []*entry{
								{key: "c", value: 3},
							},
						},
					},
				},
			},
		},
	}
	for i, wd := range []*walkData{
		{path: "a", xpath: "a", nodes: "1"},
		{path: "b.c", xpath: "b.c", nodes: "3"},
		{path: "*", xpath: "a b", nodes: "1 [{key:c value:3}]"},
		{path: "*.*", xpath: "b.c", nodes: "3"},
		{path: "[?(@ == 1)]", xpath: "a", nodes: "1"},
		{path: "[?(@.c == 3)].c", xpath: "b.c", nodes: "3"},
		{path: "..", xpath: "a b b.c", nodes: "1 [{key:c value:3}] 3"},
	} {
		x := jp.MustParseString(wd.path)
		var (
			ps []string
			ns []string
		)
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkSliceReflect(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	type AA []any
	data := AA{1, AA{2, 3}}
	for i, wd := range []*walkData{
		{path: "[0]", xpath: "[0]", nodes: "1"},
		{path: "[1][1]", xpath: "[1][1]", nodes: "3"},
		{path: "[-1][-3]", xpath: "", nodes: ""},
		{path: "*", xpath: "[0] [1]", nodes: "1 [2 3]"},
		{path: "*.*", xpath: "[1][0] [1][1]", nodes: "2 3"},
		{path: "[0:4:2]", xpath: "[0]", nodes: "1"},
		{path: "[?(@ == 1)]", xpath: "[0]", nodes: "1"},
		{path: "[?(@[0] == 2)][1]", xpath: "[1][1]", nodes: "3"},
		{path: "..", xpath: "[0] [1] [1][0] [1][1]", nodes: "1 [2 3] 2 3"},
	} {
		x := jp.MustParseString(wd.path)
		var (
			ps []string
			ns []string
		)
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkStruct(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	type B struct {
		C int
	}
	type top struct {
		A int
		B B
	}
	data := &top{A: 1, B: B{C: 3}}
	for i, wd := range []*walkData{
		{path: "a", xpath: "a", nodes: "1"},
		{path: "b.c", xpath: "b.c", nodes: "3"},
		{path: "*", xpath: "A B", nodes: "1 {c:3}"},
		{path: "*.*", xpath: "B.C", nodes: "3"},
		{path: "..", xpath: "A B B.C", nodes: "1 {c:3} 3"},
	} {
		x := jp.MustParseString(wd.path)
		var (
			ps []string
			ns []string
		)
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}

func TestExprWalkMap(t *testing.T) {
	opt := ojg.Options{Sort: true, Indent: 0}
	type name string
	type MM map[name]int

	data := map[name]any{"a": 1, "b": MM{"c": 3}}
	for i, wd := range []*walkData{
		{path: "b", xpath: "b", nodes: "{c:3}"},
		{path: "b.c", xpath: "b.c", nodes: "3"},
		{path: "*", xpath: "a b", nodes: "1 {c:3}"},
		{path: "*.*", xpath: "b.c", nodes: "3"},
		{path: "[?(@.c == 3)]", xpath: "b", nodes: "{c:3}"},
		{path: "[?(@.c == 3)].c", xpath: "b.c", nodes: "3"},
		{path: "..", xpath: "a b b.c", nodes: "1 {c:3} 3"},
	} {
		x := jp.MustParseString(wd.path)
		var (
			ps []string
			ns []string
		)
		x.Walk(data, func(path jp.Expr, nodes []any) {
			ps = append(ps, path.String())
			ns = append(ns, string(bytes.ReplaceAll(sen.Bytes(nodes[len(nodes)-1], &opt), []byte{'\n'}, []byte{})))
		})
		tt.Equal(t, wd.xpath, strings.Join(ps, " "), "%d: path mismatch for %s", i, wd.path)
		tt.Equal(t, wd.nodes, strings.Join(ns, " "), "%d: nodes mismatch for %s", i, wd.path)
	}
}
