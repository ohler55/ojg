// copyright (c) 2024, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type matchHandlerData struct {
	target string
	src    string
	expect string
}

func (md *matchHandlerData) runTest(t *testing.T, i int) {
	var buf []byte
	h := jp.NewMatchHandler(jp.MustParseString(md.target), func(path jp.Expr, data any) {
		buf = fmt.Appendf(buf, "%s: %v\n", path, pretty.SEN(data))
	})
	err := sen.TokenizeString(md.src, h)
	tt.Nil(t, err)
	tt.Equal(t, md.expect, string(buf), "%d: %s - %s", i, md.target, md.src)
}

func TestMatchHandlerRoot(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$", src: "123", expect: "$: 123\n"},
		{target: "$", src: "2.5", expect: "$: 2.5\n"},
		{target: "$", src: "abc", expect: "$: abc\n"},
		{target: "$", src: "null", expect: "$: null\n"},
		{target: "$", src: "true", expect: "$: true\n"},
		{target: "$", src: "123456789012345678901234567890", expect: "$: \"123456789012345678901234567890\"\n"},
	} {
		md.runTest(t, i)
	}
}

func TestMatchHandlerChild(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$.a", src: "{a:1 b:2}", expect: "$.a: 1\n"},
		{target: "$.a.b", src: "{a:{b:1} b:2}", expect: "$.a.b: 1\n"},
	} {
		md.runTest(t, i)
	}
}

func TestMatchHandlerNth(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$[1]", src: "[1 2 3 4]", expect: "$[1]: 2\n"},
		{target: "$[1][2]", src: "[1 [2 4 8] 3 4]", expect: "$[1][2]: 8\n"},
	} {
		md.runTest(t, i)
	}
}

func TestMatchHandlerObjectChild(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$.a", src: "{a:{b:2}}", expect: "$.a: {b: 2}\n"},
		{target: "$.a", src: "{a:{b:{c: 2}}}", expect: "$.a: {b: {c: 2}}\n"},
	} {
		md.runTest(t, i)
	}
}

func TestMatchHandlerArrayNth(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$[1]", src: "[1 [2 3 4] 5]", expect: "$[1]: [2 3 4]\n"},
	} {
		md.runTest(t, i)
	}
}

func TestMatchHandlerFilter(t *testing.T) {
	for i, md := range []*matchHandlerData{
		{target: "$[?@.x == 1]", src: "[{x:0 y:0} {x:1 y:1}]", expect: "$[1]: {x: 1 y: 1}\n"},
		{target: "$[?@.x == 2]", src: "[{x:0 y:0} {x:1 y:1}]", expect: ""},
	} {
		md.runTest(t, i)
	}
}
