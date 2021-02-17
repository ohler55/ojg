// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "append",
		Eval: appendEval,
		Desc: `Appends the second argument to the first argument which must be
an array.`,
	})
}

func appendEval(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("append expects exactly two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]interface{})
	if !ok {
		panic(fmt.Errorf("append expected an array argument, not a %T", v))
	}
	v = evalArg(root, at, args[1])

	return append(list, v)
}
