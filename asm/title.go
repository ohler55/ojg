// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"unicode"
)

func init() {
	Define(&Fn{
		Name: "title",
		Eval: title,
		Desc: `Convert a string to capitalized string. There must be exactly
one string argument.`,
	})
}

func title(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("title expects exactly one arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("title expected a string argument, not a %T", v))
	}
	ra := []rune(s)
	if 0 < len(ra) {
		ra[0] = unicode.ToUpper(ra[0])
	}
	return string(ra)
}
