// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
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

func TestScriptBasicEval(t *testing.T) {
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
	e := jp.Or(
		jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)),
		jp.Eq(jp.Get(jp.A().C("x")), jp.ConstString("cool")),
	)
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", e.String())
	s := e.Script()
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", s.String())
	f := e.Filter()
	tt.Equal(t, "[?(@.a < 52 || @.x == 'cool')]", f.String())

	//fmt.Printf("*** data: %s\n", jp.JSON(data))
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

		{src: "@.x == 4", err: "a script must start with a '('"},
		{src: "(@.x ++ 4)", err: "'++' is not a valid operation at 7 in (@.x ++ 4)"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		s, err := jp.NewScript(d.src)
		if 0 < len(d.err) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.expect, s.String(), i, ": ", d.src)
		}
	}
}

func TestScriptMatch(t *testing.T) {
	s, err := jp.NewScript("(@.x == 3)")
	tt.Nil(t, err)
	tt.Equal(t, true, s.Match(map[string]interface{}{"x": 3}))
	tt.Equal(t, true, s.Match(gen.Object{"x": gen.Int(3)}))
}

func TestScriptNormalizeEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	for _, v := range []interface{}{
		int8(3),
		int16(3),
		int32(3),
		int64(3),
		uint(3),
		uint8(3),
		uint16(3),
		uint32(3),
		uint64(3),
		float32(3.0),
		3.0,
		gen.Int(3),
		gen.Float(3.0),
	} {
		result := s.Eval([]interface{}{}, []interface{}{v})
		tt.Equal(t, 1, len(result), fmt.Sprintf("%T %v", v, v))
	}
	s, err = jp.NewScript("(@ == 'x')")
	tt.Nil(t, err)
	result := s.Eval([]interface{}{}, []interface{}{"x"})
	tt.Equal(t, 1, len(result), "string normalize")
	result = s.Eval([]interface{}{}, gen.Array{gen.String("x")})
	tt.Equal(t, 1, len(result), "gen.String normalize")

	s, err = jp.NewScript("(@ == true)")
	tt.Nil(t, err)
	result = s.Eval([]interface{}{}, gen.Array{gen.True})
	tt.Equal(t, 1, len(result), "bool normalize")
}

func TestScriptNonListEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	result := s.Eval([]interface{}{}, "bad")
	tt.Equal(t, 0, len(result))
}

type edata struct {
	src     string
	value   interface{}
	noMatch bool
}

func TestScriptEval(t *testing.T) {
	for i, d := range []edata{
		{src: "(@ == 3)", value: int64(3)},
		{src: "(@ == 3)", value: 3.0},
		{src: "(@ == 3.0)", value: int64(3)},
		{src: "(@ == 'abc')", value: "abc"},
		{src: "(@ != 3)", value: int64(3), noMatch: true},
		{src: "(@ != 3)", value: int64(4)},
		{src: "(@ != 3)", value: 3.1},
		{src: "(@ != 3.0)", value: int64(4)},
		{src: "(@ != 'abc')", value: "xyz"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s in %s\n", i, d.src, oj.JSON(d.value))
		}
		s, err := jp.NewScript(d.src)
		tt.Nil(t, err)
		result := s.Eval([]interface{}{}, []interface{}{d.value})
		if d.noMatch {
			tt.Equal(t, 0, len(result), d.src, " in ", d.value)
		} else {
			tt.Equal(t, 1, len(result), d.src, " in ", d.value)
		}
	}
}

func BenchmarkOjScriptDev(b *testing.B) {
	s := jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)).Script()
	data := scriptBenchData(100)
	stack := []interface{}{}
	b.ReportAllocs()
	b.ResetTimer()
	//fmt.Printf("*** data: %s\n", jp.JSON(data))
	for n := 0; n < b.N; n++ {
		stack = stack[:0]
		stack = s.Eval(stack, data)
		//fmt.Printf("*** stack: %s\n", jp.JSON(stack))
	}
}
