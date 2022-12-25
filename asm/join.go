// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "join",
		Eval: join,
		Desc: `Join an array of strings with the provided separator. If a
separator is not provided as the second argument then an empty
string is used.`,
	})
}

func join(root map[string]any, at any, args ...any) any {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("join expects one or two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]any)
	if !ok {
		panic(fmt.Errorf("join expected an array of string arguments, not a %T", v))
	}
	slist := make([]string, 0, len(list))
	for _, v = range list {
		var s string
		if s, ok = v.(string); !ok {
			panic(fmt.Errorf("join expected an array of string arguments, not a %T", v))
		}
		slist = append(slist, s)
	}
	var sep string
	if 1 < len(args) {
		v = evalArg(root, at, args[1])
		if sep, ok = v.(string); !ok {
			panic(fmt.Errorf("join expected a string separator argument, not a %T", v))
		}
	}
	return strings.Join(slist, sep)
}
