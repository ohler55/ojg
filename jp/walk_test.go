// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type simple int

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
