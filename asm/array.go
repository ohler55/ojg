// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "array?",
		Eval: arrayEval,
		Desc: `Returns true if the single required argumement is an array
otherwise false is returned.`,
	})
}

func arrayEval(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("array? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).([]interface{})

	return ok
}
