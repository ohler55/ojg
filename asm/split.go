// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"strings"
)

func init() {
	Define(&Fn{
		Name: "split",
		Eval: split,
		Desc: `Split a string on using a specified separator.`,
	})
}

func split(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("split expects exactly two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	s, ok := v.(string)
	if !ok {
		panic(fmt.Errorf("split expected a string argument, not a %T", v))
	}
	v = evalArg(root, at, args[1])
	var sep string
	if sep, ok = v.(string); !ok {
		panic(fmt.Errorf("split expected a string separator argument, not a %T", v))
	}
	var list []interface{}
	for _, s := range strings.Split(s, sep) {
		list = append(list, s)
	}
	return list
}
