// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "or",
		Eval: or,
		Desc: `Returns true if any of the argument evaluate to true. Any
arguments that do not evaluate to a boolean or null (false)
raise an error.`,
	})
}

func or(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	val := false
	for _, arg := range args {
		switch tv := evalArg(root, at, arg).(type) {
		case nil:
		case bool:
			val = tv
		default:
			panic(fmt.Errorf("or expects only boolean arguments. %T is not a boolean", tv))
		}
		if val {
			break
		}
	}
	return val
}
