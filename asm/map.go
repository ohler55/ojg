// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "map?",
		Eval: mapEval,
		Desc: `Returns true if the single required argumement is a map
otherwise false is returned.`,
	})
}

func mapEval(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("map? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(map[string]any)

	return ok
}
