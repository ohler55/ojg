// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "string?",
		Eval: stringEval,
		Desc: `Returns true if the single required argumement is a string
otherwise false is returned.`,
	})
}

func stringEval(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("string? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(string)

	return ok
}
