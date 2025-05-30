// copyright (c) 2024, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type matchData struct {
	target string
	path   string
	expect bool
}

func TestPathMatchCheck(t *testing.T) {
	for i, md := range []*matchData{
		{target: "$.a", path: "a", expect: true},
		{target: "@.a", path: "a", expect: true},
		{target: "a", path: "a", expect: true},
		{target: "a", path: "$.a", expect: true},
		{target: "a", path: "@.a", expect: true},
		{target: "[1]", path: "[1]", expect: true},
		{target: "[1]", path: "[0]", expect: false},
		{target: "*", path: "[1]", expect: true},
		{target: "[*]", path: "[1]", expect: true},
		{target: "*", path: "a", expect: true},
		{target: "[1,'a']", path: "a", expect: true},
		{target: "[1,'a']", path: "[1]", expect: true},
		{target: "[1,'a']", path: "b", expect: false},
		{target: "[1,'a']", path: "[0]", expect: false},
		{target: "$.x[1,'a']", path: "x[1]", expect: true},
		{target: "..x", path: "a.b.x", expect: true},
		{target: "..x", path: "a.b.c", expect: false},
		{target: "x[1:5:2]", path: "x[2]", expect: true},
		{target: "x[1:5:2]", path: "x.y", expect: false},
		{target: "x[?@.a == 2]", path: "x[2]", expect: true},
		{target: "x.y.z", path: "x.y", expect: false},
	} {
		tt.Equal(t, md.expect, jp.PathMatch(jp.MustParseString(md.target), jp.MustParseString(md.path)),
			"%d: %s %s", i, md.target, md.path)
	}
}

func TestPathMatchDoubleRoot(t *testing.T) {
	tt.Equal(t, false, jp.PathMatch(jp.R().R().C("a"), jp.C("a")))
	tt.Equal(t, false, jp.PathMatch(jp.A().A().C("a"), jp.C("a")))
	tt.Equal(t, false, jp.PathMatch(jp.C("a"), jp.R().R().C("a")))
	tt.Equal(t, false, jp.PathMatch(jp.C("a"), jp.A().A().C("a")))
}

func TestPathMatchSkipBracket(t *testing.T) {
	tt.Equal(t, true, jp.PathMatch(jp.B().C("a"), jp.C("a")))
}

func TestMatchNestedMapArray(t *testing.T) {
	var buf []byte
	err := oj.Match([]byte(`[{"a":[1, 2]}]`), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %s", path, pretty.SEN(data))
	}, jp.MustParseString("$[*]"))
	tt.Nil(t, err)
	tt.Equal(t, "$[0]: {a: [1 2]}", string(buf))
}

func TestMatchNestedArrays(t *testing.T) {
	var buf []byte
	err := oj.Match([]byte(`[[[1, [2]]]]`), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %s", path, sen.Bytes(data))
	}, jp.MustParseString("$[*]"))
	tt.Nil(t, err)
	tt.Equal(t, "$[0]: [[1 [2]]]", string(buf))
}

func TestMatchNestedMap(t *testing.T) {
	var buf []byte
	err := oj.Match([]byte(`{"a": {"b": 1, "c": {"d": 2}}}`), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %s", path, pretty.SEN(data))
	}, jp.MustParseString("$[*]"))
	tt.Nil(t, err)
	tt.Equal(t, "$.a: {b: 1 c: {d: 2}}", string(buf))
}
