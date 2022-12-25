// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "trim",
		Eval: trim,
		Desc: `Trim white space from both ends of a string unless a second
argument provides an alternative cut set.`,
	})
}

func trim(root map[string]any, at any, args ...any) any {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("trim expects one or two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("trim expected a string argument, not a %T", v))
	}
	if 1 < len(args) {
		v = evalArg(root, at, args[1])
		var cut string
		cut, ok = v.(string)
		if !ok {
			panic(fmt.Errorf("trim expected a string cut set argument, not a %T", v))
		}
		return strings.Trim(s, cut)
	}
	return strings.TrimSpace(s)
}
