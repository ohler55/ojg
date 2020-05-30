// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

type xdata struct {
	src    string
	expect string
	err    string
}

func TestExprParse(t *testing.T) {
	for i, d := range []xdata{
		{src: "@", expect: "@"},
		{src: "@.abc", expect: "@.abc"},
		{src: "@.a.b.c", expect: "@.a.b.c"},
		{src: "$", expect: "$"},
		{src: "$.abc", expect: "$.abc"},
		{src: "$.a.b.c", expect: "$.a.b.c"},
		{src: "abc", expect: "abc"},
		{src: "abc.def", expect: "abc.def"},
		{src: "abc.*.def", expect: "abc.*.def"},
		{src: "abc..def", expect: "abc..def"},
		{src: "abc[*].def", expect: "abc[*].def"},
		{src: "abc[0].def", expect: "abc[0].def"},
		{src: "abc[2].def", expect: "abc[2].def"},
		{src: "abc[ -2 ].def", expect: "abc[-2].def"},
		{src: "abc[0:]", expect: "abc[:]"},
		{src: "abc[:]", expect: "abc[:]"},
		{src: "abc[:-1]", expect: "abc[:]"},
		{src: "abc[1:]", expect: "abc[1:]"},
		{src: "abc[1:3]", expect: "abc[1:3]"},
		{src: "abc[1 : -1 : 2]", expect: "abc[1::2]"},
		{src: "$['abc']", expect: "$.abc"},
		{src: "$['a b']", expect: "$['a b']"},
		{src: "$['ぴーたー']", expect: "$.ぴーたー"},
		{src: "$[1,2]", expect: "$[1,2]"},
		{src: "$['a','b']", expect: "$['a','b']"},
		{src: "$[ 1, 'a' , 2 ,'b' ]", expect: "$[1,'a',2,'b']"},

		//{src: "$[?(@.x == 'abc')]", expect: "$[?(@.x == 'abc')]"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		x, err := oj.ParseExprString(d.src)
		if 0 < len(d.err) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.NotNil(t, x)
			tt.Equal(t, d.expect, x.String(), i, ": ", d.src)
		}
	}
}

func xTestExprParseDev(t *testing.T) {
	x, err := oj.ParseExprString("@.abc")
	tt.Nil(t, err)
	tt.NotNil(t, x)
	tt.Equal(t, "@.abc", x.String())
}

func BenchmarkExprParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _ = oj.ParseExpr([]byte("@.abc.*[2,3]..xyz[2]"))
		//fmt.Printf("*** x: %s\n", x)
	}
}
