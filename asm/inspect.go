// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/oj"
)

func init() {
	Define(&Fn{
		Name: "inspect",
		Eval: inspect,
		Desc: `Print the arguments as JSON unless the argument is an integer.
Integers are assumed to be the indentation for the arguments
that follow.`,
	})
}

func inspect(root map[string]any, at any, args ...any) any {
	indent := 0
	if 0 < len(args) {
		for _, a := range args {
			val := evalArg(root, at, a)
			switch tv := val.(type) {
			case string:
				if 0 < indent {
					fmt.Printf("%s:\n", tv)
				} else {
					fmt.Printf("%s: ", tv)
				}
			case int:
				indent = tv
			case int64:
				indent = int(tv)
			default:
				fmt.Println(oj.JSON(val, indent))
			}
		}
	}
	return at
}
