// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func scriptBenchData(size int64) interface{} {
	list := []interface{}{}
	for i := int64(0); i < size; i++ {
		list = append(list, map[string]interface{}{string([]byte{'a' + byte(i%26)}): i, "x": i})
	}
	return list
}

func TestScriptEval(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]interface{}{
			"a": int64(52),
			"b": 4,
			"c": 6,
		},
	}
	e := oj.Or(
		oj.Lt(oj.Get(oj.A().C("a")), oj.ConstInt(52)),
		oj.Eq(oj.Get(oj.A().C("x")), oj.ConstString("cool")),
	)
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", e.String())
	s := e.Script()
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", s.String())
	f := e.Filter()
	tt.Equal(t, "[?(@.a < 52 || @.x == 'cool')]", f.String())

	//fmt.Printf("*** data: %s\n", oj.JSON(data))
	stack := s.Eval([]interface{}{}, data)
	tt.Equal(t, `[{"a":1,"b":2,"c":3}]`, oj.JSON(stack, &oj.Options{Sort: true}))
}

func TestScriptParse(t *testing.T) {
	for i, d := range []xdata{
		{src: "(@.x == 'abc')", expect: "(@.x == 'abc')"},
		{src: "(@.x<5)", expect: "(@.x < 5)"},
		{src: "(@.x<123)", expect: "(@.x < 123)"},
		{src: "(@.x == 3)", expect: "(@.x == 3)"},
		{src: "(@.*.xyz==true)", expect: "(@.*.xyz == true)"},
		{src: "(@.x.* == 3)", expect: "(@.x.* == 3)"},
		{src: "(@.. == 3)", expect: "(@.. == 3)"},
		{src: "(@[3] == 3)", expect: "(@[3] == 3)"},
		{src: "(@[3,4] == 3)", expect: "(@[3,4] == 3)"},
		{src: "(@[3,'four'] == 3)", expect: "(@[3,'four'] == 3)"},
		{src: "(@[1:5] == 3)", expect: "(@[1:5] == 3)"},
		{src: "(@ == 3)", expect: "(@ == 3)"},
		{src: "($ == 3.4e-3)", expect: "($ == 0.0034)"},
		{src: "($ == 3e-3)", expect: "($ == 0.003)"},
		{src: "(3 == @.x)", expect: "(3 == @.x)"},
		{src: "(@ == $)", expect: "(@ == $)"},
		{src: "(@.x[?(@.a == true)].b == false)", expect: "(@.x[?(@.a == true)].b == false)"},
		{src: "(@.x[?(@.a == 5)] == 11)", expect: "(@.x[?(@.a == 5)] == 11)"},
		{src: "((@.x == 3) || (@.y > 5))", expect: "(@.x == 3 || @.y > 5)"},
		{src: "(@.x < 3 && @.x > 1 || @.z == 3)", expect: "(@.x < 3 && @.x > 1 || @.z == 3)"},
		{src: "(!(3 == @.x))", expect: "(!(3 == @.x))"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		s, err := oj.NewScript(d.src)
		if 0 < len(d.err) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.expect, s.String(), i, ": ", d.src)
		}
	}
}

func TestScriptDev(t *testing.T) {
	src := "(!(@.x == 2))"
	s, err := oj.NewScript(src)
	tt.Nil(t, err)
	tt.Equal(t, src, s.String())
}

func BenchmarkOjScriptDev(b *testing.B) {
	s := oj.Lt(oj.Get(oj.A().C("a")), oj.ConstInt(52)).Script()
	data := scriptBenchData(100)
	stack := []interface{}{}
	b.ReportAllocs()
	b.ResetTimer()
	//fmt.Printf("*** data: %s\n", oj.JSON(data))
	for n := 0; n < b.N; n++ {
		stack = stack[:0]
		stack = s.Eval(stack, data)
		//fmt.Printf("*** stack: %s\n", oj.JSON(stack))
	}
}
