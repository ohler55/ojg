// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg/oj"
)

func ExampleParseString() {
	v, err := oj.ParseString(`{"a": 1, "b":[2,3,4]}`)
	if err == nil {
		// Sorted output allows for consistent results.
		fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))
	} else {
		fmt.Println(err.Error())
	}
	// Output: {"a":1,"b":[2,3,4]}
}

func ExampleParser_Parse() {
	// The parser can be reused for better performance by reusing buffers.
	var p oj.Parser
	v, err := p.Parse([]byte(`{"a": 1, "b":[2,3,4]}`))
	if err == nil {
		// Sorted output allows for consistent results.
		fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))
	} else {
		fmt.Println(err.Error())
	}
	// Output: {"a":1,"b":[2,3,4]}
}

func ExampleParser_ParseReader() {
	// The parser can be reused for better performance by reusing buffers.
	var p oj.Parser
	v, err := p.ParseReader(strings.NewReader(`{"a": 1, "b":[2,3,4]}`))
	if err == nil {
		// Sorted output allows for consistent results.
		fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))
	} else {
		fmt.Println(err.Error())
	}
	// Output: {"a":1,"b":[2,3,4]}
}

func ExampleParser_Parse_callback() {
	var results []byte
	cb := func(n interface{}) bool {
		if 0 < len(results) {
			results = append(results, ' ')
		}
		results = append(results, oj.JSON(n)...)
		return false
	}
	var p oj.Parser
	_, _ = p.Parse([]byte("[1,2][3,4][5,6]"), cb)
	fmt.Println(string(results))
	// Output: [1,2] [3,4] [5,6]
}

func ExampleNodeParser_Parse() {
	// The parser can be reused for better performance by reusing buffers.
	var p oj.NodeParser
	v, err := p.Parse([]byte(`{"a": 1, "b":[2,3,4]}`))
	if err == nil {
		// Sorted output allows for consistent results.
		fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))
		fmt.Printf("type: %T\n", v)
	} else {
		fmt.Println(err.Error())
	}
	// Output: {"a":1,"b":[2,3,4]}
	// type: oj.Object
}

func ExampleValidateString() {
	err := oj.ValidateString(`{"a": 1, "b":[2,3,4]}`)
	fmt.Println(oj.JSON(err))
	// Output: null
}

func ExampleScript() {
	data := []interface{}{
		map[string]interface{}{"a": 1, "b": 2, "c": 3},
		map[string]interface{}{"a": int64(52), "b": 4, "c": 6},
	}
	// Build an Equation and generate a Script from the Equation.
	s := oj.Or(
		oj.Lt(oj.Get(oj.A().C("a")), oj.ConstInt(52)),
		oj.Eq(oj.Get(oj.A().C("x")), oj.ConstString("cool")),
	).Script()
	fmt.Println(s.String())
	// Normally Scripts are using in Expr (JSON paths).
	result := s.Eval([]interface{}{}, data)
	fmt.Println(oj.JSON(result, &oj.Options{Sort: true}))
	// Output:
	// (@.a < 52 || @.x == 'cool')
	// [{"a":1,"b":2,"c":3}]
}

func ExampleExpr_noparse() {
	data := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{"x": 1, "y": 2, "z": 3},
			map[string]interface{}{"x": 1, "y": 4, "z": 9},
		},
		"b": []interface{}{
			map[string]interface{}{"x": 4, "y": 5, "z": 6},
			map[string]interface{}{"x": 16, "y": 25, "z": 36},
		},
	}
	x := oj.C("b").F(oj.Gt(oj.Get(oj.A().C("y")), oj.ConstInt(10))).C("x")
	fmt.Println(x.String())
	result := x.Get(data)
	fmt.Println(oj.JSON(result, &oj.Options{Sort: true}))
	// Output:
	// b[?(@.y > 10)].x
	// [16]
}

func ExampleBuilder() {
	var b oj.Builder

	_ = b.Object()
	_ = b.Array("a")
	_ = b.Value(true)
	_ = b.Object()
	_ = b.Value(123, "x")
	b.Pop()
	_ = b.Value(nil)
	b.PopAll()
	v := b.Result()
	fmt.Println(oj.JSON(v))
	// Output: {"a":[true,{"x":123},null]}
}

func ExampleParseExprString() {
	data := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{"x": 1, "y": 2, "z": 3},
			map[string]interface{}{"x": 1, "y": 4, "z": 9},
		},
		"b": []interface{}{
			map[string]interface{}{"x": 4, "y": 5, "z": 6},
			map[string]interface{}{"x": 16, "y": 25, "z": 36},
		},
	}
	x, err := oj.ParseExprString("b[?(@.y > 10)].x")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(x.String())
	result := x.Get(data)
	fmt.Println(oj.JSON(result))
	// Output:
	// b[?(@.y > 10)].x
	// [16]
}
