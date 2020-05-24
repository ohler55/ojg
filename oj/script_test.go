// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
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

func TestOjScriptDev(t *testing.T) {
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
