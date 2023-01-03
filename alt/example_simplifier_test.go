// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"

	"github.com/ohler55/ojg/oj"
)

type simmer struct {
	val int
}

func (s *simmer) Simplify() any {
	return map[string]any{"type": "simmer", "val": s.val}
}

func ExampleSimplifier() {
	// Non public types can be encoded with the Simplifier interface which
	// should decompose into a simple type.
	fmt.Println(oj.JSON(&simmer{val: 3}, &oj.Options{Sort: true}))

	// Output: {"type":"simmer","val":3}
}
