// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "not",
		Eval: not,
		Desc: `Returns the boolean NOT of the argument. Exactly one argument
is expected and it must be a boolean.`,
	})
}

func not(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("not expects exactly one arguments. %d given", len(args)))
	}
	if boo, ok := evalArg(root, at, args[0]).(bool); ok {
		return !boo
	}
	panic(fmt.Errorf("not expects only a single boolean arguments. %T is not a boolean", args[0]))
}
