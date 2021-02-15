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

func mapEval(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("map? expects exactly one arguments. %d given", len(args)))
	}
	_, ok := evalArg(root, at, args[0]).(map[string]interface{})

	return ok
}
