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

func reverse(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("reverse expects exactly one argument. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]interface{})
	if !ok {
		panic(fmt.Errorf("reverse expected an array argument, not a %T", v))
	}
	// Make a copy so not to change the original.
	rev := make([]interface{}, len(list))
	var i int
	for i, v = range list {
		rev[len(list)-i-1] = v
	}
	return rev
}
