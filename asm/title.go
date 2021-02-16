// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "title",
		Eval: title,
		Desc: `Convert a string to capitalized string. There must be exactly
one string argument.`,
	})
}

func title(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("title expects exactly one arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("title expected a string argument, not a %T", v))
	}
	return strings.Title(s)
}
