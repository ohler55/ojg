// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "include",
		Eval: include,
		Desc: `Returns true if a list first argument includes the second
argument. It will also return true if the first argument is a
string and the second string argument is included in the first.`,
	})
}

func include(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("include expects two arguments. %d given", len(args)))
	}
	v1 := evalArg(root, at, args[1])
	switch v := evalArg(root, at, args[0]).(type) {
	case []interface{}:
		for _, m := range v {
			if m == v1 {
				return true
			}
		}
	case string:
		s, ok := v1.(string)
		if !ok {
			panic(fmt.Errorf("include expects a string second argument if the fist is a string, not a %T", v1))
		}
		return strings.Contains(v, s)
	default:
		panic(fmt.Errorf("include expects an array or string first argument, not a %T", v))
	}
	return false
}
