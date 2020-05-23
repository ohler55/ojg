// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
)

func scriptBenchData(size int64) interface{} {
	list := []interface{}{}
	for i := int64(0); i < size; i++ {
		list = append(list, map[string]interface{}{string([]byte{'a' + byte(i%26)}): i, "x": i})
	}
	return list
}

func TestOjScriptDev(t *testing.T) {
	s := oj.Foo()
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
	//fmt.Printf("*** data: %s\n", oj.JSON(data))
	stack := []interface{}{}
	stack = s.Eval(stack, data)
	fmt.Printf("*** stack after: %s\n", oj.JSON(stack))
	//fmt.Printf("*** script string: %s\n", s.String())
	s = oj.Foo2()
	fmt.Printf("*** script string: %s\n", s.String())
}

func BenchmarkOjScriptDev(b *testing.B) {

	s := oj.Foo()
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
