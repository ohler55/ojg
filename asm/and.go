// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "and",
		Eval: and,
		Desc: `Returns true if all argument evaluate to true. Any arguments
that do not evaluate to a boolean or null (false) raise an error.`,
	})
}

func and(root map[string]any, at any, args ...any) any {
	val := true
	for _, arg := range args {
		switch tv := evalArg(root, at, arg).(type) {
		case nil:
			val = false
		case bool:
			val = tv
		default:
			panic(fmt.Errorf("and expects only boolean arguments. %T is not a boolean", tv))
		}
		if !val {
			break
		}
	}
	return val
}
