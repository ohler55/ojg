// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "mod",
		Eval: mod,
		Desc: `Returns the remainer of a modulo operation on the first two
argument. Both arguments must be integers and are both required.
An error is raised if the wrong argument types are given.`,
	})
}

func mod(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("mod expects exactly two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	n0, ok := asInt(v)
	if !ok {
		panic(fmt.Errorf("mod expected only integer arguments, not a %T", v))
	}
	v = evalArg(root, at, args[1])
	var n1 int64
	if n1, ok = asInt(v); !ok {
		panic(fmt.Errorf("mod expected only integer arguments, not a %T", v))
	}
	return n0 % n1
}
