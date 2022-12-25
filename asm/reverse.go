// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "reverse",
		Eval: reverse,
		Desc: `Reverse the items in an array and return a copy of it.`,
	})
}

func reverse(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("reverse expects exactly one argument. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]any)
	if !ok {
		panic(fmt.Errorf("reverse expected an array argument, not a %T", v))
	}
	// Make a copy so not to change the original.
	rev := make([]any, len(list))
	var i int
	for i, v = range list {
		rev[len(list)-i-1] = v
	}
	return rev
}
