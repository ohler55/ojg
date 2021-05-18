// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
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

func ExampleJSON() {
	type Valley struct {
		Val int `json:"value"`
	}
	b := oj.JSON(&Valley{Val: 3})
	fmt.Printf("%s\n", b)
	// Output: {"val":3}
}

func ExampleParse() {
	val, err := oj.Parse([]byte(`[true,false,[3,2,1],{"a":1,"b":2,"c":3,"d":["x","y","z",[]]}]`))
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.JSON(val, 80.3))

	// Output:
	// [
	//   true,
	//   false,
	//   [3, 2, 1],
	//   {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
	// ]
}

func ExampleMustParse() {
	val := oj.MustParse([]byte(`[true,false,[3,2,1],{"a":1,"b":2,"c":3,"d":["x","y","z",[]]}]`))
	fmt.Println(pretty.JSON(val, 80.3))

	// Output:
	// [
	//   true,
	//   false,
	//   [3, 2, 1],
	//   {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
	// ]
}

func ExampleMustParseString() {
	val := oj.MustParseString(`[true,false,[3,2,1],{"a":1,"b":2,"c":3,"d":["x","y","z",[]]}]`)
	fmt.Println(pretty.JSON(val, 80.3))

	// Output:
	// [
	//   true,
	//   false,
	//   [3, 2, 1],
	//   {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
	// ]
}

func ExampleLoad() {
	r := strings.NewReader(`[true,false,[3,2,1],{"a":1,"b":2,"c":3,"d":["x","y","z",[]]}]`)
	val, err := oj.Load(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.JSON(val, 80.3))

	// Output:
	// [
	//   true,
	//   false,
	//   [3, 2, 1],
	//   {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
	// ]
}

func ExampleMustLoad() {
	r := strings.NewReader(`[true,false,[3,2,1],{"a":1,"b":2,"c":3,"d":["x","y","z",[]]}]`)
	val := oj.MustLoad(r)
	fmt.Println(pretty.JSON(val, 80.3))

	// Output:
	// [
	//   true,
	//   false,
	//   [3, 2, 1],
	//   {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
	// ]
}

func ExampleWrite() {
	var b strings.Builder
	data := []interface{}{
		map[string]interface{}{
			"x": 1,
			"y": 2,
		},
	}
	if err := oj.Write(&b, data, &ojg.Options{Sort: true}); err != nil {
		panic(err)
	}
	fmt.Println(b.String())

	// Output: [{"x":1,"y":2}]
}

func ExampleWriter_MustWrite() {
	var b strings.Builder
	data := []interface{}{
		map[string]interface{}{
			"x": 1,
			"y": 2,
		},
	}
	wr := oj.Writer{Options: ojg.Options{Sort: true}}
	wr.MustWrite(&b, data)
	fmt.Println(b.String())

	// Output: [{"x":1,"y":2}]
}

func ExampleWriter_JSON() {
	data := []interface{}{
		map[string]interface{}{
			"x": 1,
			"y": 2,
		},
	}
	wr := oj.Writer{Options: ojg.Options{Sort: true}}
	j := wr.JSON(data)
	fmt.Println(j)

	// Output: [{"x":1,"y":2}]
}
