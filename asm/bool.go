// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "bool?",
		Eval: boolEval,
		Desc: `Returns true if the single required argumement is a boolean
otherwise false is returned.`,
	})
}

func boolEval(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("bool? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(bool)

	return ok
}
