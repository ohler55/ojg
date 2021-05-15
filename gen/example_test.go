// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

func ExampleParser_Parse() {
	// The parser can be reused for better performance by reusing buffers.
	var p gen.Parser
	v, err := p.Parse([]byte(`{"a": 1, "b":[2,3,4]}`))
	if err == nil {
		// Sorted output allows for consistent results.
		fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))
		fmt.Printf("type: %T\n", v)
	} else {
		fmt.Println(err.Error())
	}
	// Output: {"a":1,"b":[2,3,4]}
	// type: gen.Object
}

func ExampleBuilder() {
	var b gen.Builder

	b.MustObject()
	b.MustArray("a")
	b.MustValue(gen.True)
	b.MustObject()
	b.MustValue(gen.Int(123), "x")
	b.Pop()
	b.MustValue(nil)
	b.PopAll()
	v := b.Result()

	fmt.Println(sen.String(v))

	// Output: {a:[true {x:123}null]}
}
