// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "tolower",
		Eval: tolower,
		Desc: `Convert a string to lowercase. There must be exactly one
string argument.`,
	})
}

func tolower(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("tolower expects exactly one arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("tolower expected a string argument, not a %T", v))
	}
	return strings.ToLower(s)
}
