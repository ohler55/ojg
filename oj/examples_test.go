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

func ExampleValidateString() {
	err := oj.ValidateString(`{"a": 1, "b":[2,3,4]}`)
	fmt.Println(oj.JSON(err))
	// Output: null
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

func ExampleMarshal() {
	type Valley struct {
		Val int `json:"value"`
	}

	b, err := oj.Marshal(&Valley{Val: 3})
	fmt.Printf("%v %s\n", err, b)
	// Output: <nil> {"value":3}
}
