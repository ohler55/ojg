// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "cond",
		Eval: cond,
		Desc: `A conditional construct modeled after the LISP cond. All
arguments must be array of two elements. The first element must
evaluate to a boolean and the second can be any value. The value
of the first true first argument is returned. If none match nil
is returned.`,
	})
}

func cond(root map[string]any, at any, args ...any) any {
	for _, arg := range args {
		v := evalArg(root, at, arg)
		list, ok := v.([]any)
		if !ok {
			panic(fmt.Errorf("cond expects array arguments, not a %T", v))
		}
		if len(list) != 2 {
			panic(fmt.Errorf("cond array arguments must have two elements, not a %d", len(list)))
		}
		if b, _ := evalArg(root, at, list[0]).(bool); b {
			return list[1]
		}
	}
	return nil
}
