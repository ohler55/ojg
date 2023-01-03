// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "toupper",
		Eval: toupper,
		Desc: `Convert a string to uppercase. There must be exactly one
string argument.`,
	})
}

func toupper(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("toupper expects exactly one arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("toupper expected a string argument, not a %T", v))
	}
	return strings.ToUpper(s)
}
