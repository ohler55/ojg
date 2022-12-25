// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "replace",
		Eval: replace,
		Desc: `Replace an occurrences the second argument with the third
argument. All three arguments must be strings.`,
	})
}

func replace(root map[string]any, at any, args ...any) any {
	if len(args) != 3 {
		panic(fmt.Errorf("replace expects exactly three arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("replace expects a string argument, not a %T", v))
	}
	v = evalArg(root, at, args[1])
	var old string
	if old, ok = v.(string); !ok {
		panic(fmt.Errorf("replace expects a string second argument, not a %T", v))
	}
	v = evalArg(root, at, args[2])
	var rep string
	if rep, ok = v.(string); !ok {
		panic(fmt.Errorf("replace expects a string third argument, not a %T", v))
	}
	return strings.ReplaceAll(s, old, rep)
}
