// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

func ExampleExpr_Set() {
	data := []any{
		map[string]any{"a": 1, "b": 2, "c": 3},
	}
	// Set members with a JSONPath expression.
	if err := jp.N(0).C("b").Set(data, 7); err != nil {
		panic(err)
	}
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))

	// Add members with a JSONPath expression.
	if err := jp.N(0).C("d").Set(data, 4); err != nil {
		panic(err)
	}
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))

	// Output:
	// [{a:1 b:7 c:3}]
	// [{a:1 b:7 c:3 d:4}]
}

func ExampleExpr_MustSet() {
	data := []any{
		map[string]any{"a": 1, "b": 2, "c": 3},
	}
	// Set members with a JSONPath expression.
	jp.N(0).C("b").MustSet(data, 7)
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))

	// Add members with a JSONPath expression.
	jp.N(0).C("d").MustSet(data, 4)
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))

	// Output:
	// [{a:1 b:7 c:3}]
	// [{a:1 b:7 c:3 d:4}]
}

func ExampleExpr_Del() {
	data := []any{
		map[string]any{"a": 1, "b": 2, "c": 3},
	}
	if err := jp.N(0).C("b").Del(data); err != nil {
		panic(err)
	}
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))

	// Output:
	// [{a:1 c:3}]
}

func ExampleExpr_MustDel() {
	data := []any{
		map[string]any{"a": 1, "b": 2, "c": 3},
	}
	jp.N(0).C("b").MustDel(data)
	fmt.Println(sen.String(data, &ojg.Options{Sort: true}))
	// Output:
	// [{a:1 c:3}]
}

func ExampleScript() {
	data := []any{
		map[string]any{"a": 1, "b": 2, "c": 3},
		map[string]any{"a": int64(52), "b": 4, "c": 6},
	}
	// Build an Equation and generate a Script from the Equation.
	s := jp.Or(
		jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)),
		jp.Eq(jp.Get(jp.A().C("x")), jp.ConstString("cool")),
	).Script()
	fmt.Println(s.String())
	// Normally Scripts are using in Expr (JSON paths).
	result := s.Eval([]any{}, data)
	fmt.Println(oj.JSON(result, &oj.Options{Sort: true}))
	// Output:
	// (@.a < 52 || @.x == 'cool')
	// [{"a":1,"b":2,"c":3}]
}

func ExampleExpr_noparse() {
	data := map[string]any{
		"a": []any{
			map[string]any{"x": 1, "y": 2, "z": 3},
			map[string]any{"x": 1, "y": 4, "z": 9},
		},
		"b": []any{
			map[string]any{"x": 4, "y": 5, "z": 6},
			map[string]any{"x": 16, "y": 25, "z": 36},
		},
	}
	x := jp.C("b").F(jp.Gt(jp.Get(jp.A().C("y")), jp.ConstInt(10))).C("x")
	fmt.Println(x.String())
	result := x.Get(data)
	fmt.Println(oj.JSON(result, &oj.Options{Sort: true}))
	// Output:
	// b[?(@.y > 10)].x
	// [16]
}

func ExampleParseString() {
	data := map[string]any{
		"a": []any{
			map[string]any{"x": 1, "y": 2, "z": 3},
			map[string]any{"x": 1, "y": 4, "z": 9},
		},
		"b": []any{
			map[string]any{"x": 4, "y": 5, "z": 6},
			map[string]any{"x": 16, "y": 25, "z": 36},
		},
	}
	x, err := jp.ParseString("b[?(@.y > 10)].x")
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

func ExampleMustParseString() {
	data := map[string]any{
		"a": []any{
			map[string]any{"x": 1, "y": 2, "z": 3},
			map[string]any{"x": 1, "y": 4, "z": 9},
		},
		"b": []any{
			map[string]any{"x": 4, "y": 5, "z": 6},
			map[string]any{"x": 16, "y": 25, "z": 36},
		},
	}
	x := jp.MustParseString("b[?(@.y > 10)].x")
	fmt.Println(x.String())
	result := x.Get(data)
	fmt.Println(oj.JSON(result))
	// Output:
	// b[?(@.y > 10)].x
	// [16]
}
