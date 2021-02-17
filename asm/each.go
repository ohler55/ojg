// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "each",
		Eval: each,
		Desc: `Each .`,
	})
}

func each(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) < 2 || 3 < len(args) {
		panic(fmt.Errorf("each expects two or three argument. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]interface{})
	if !ok {
		panic(fmt.Errorf("each expects an array argument, not a %T", v))
	}
	fn, _ := args[1].(*Fn)
	if fn == nil {
		panic(fmt.Errorf("each expects function as the second argument, not a %T", args[1]))
	}
	key := "asm"
	if 2 < len(args) {
		v = evalArg(root, at, args[2])
		var s string
		if s, ok = v.(string); !ok {
			panic(fmt.Errorf("each expects a string for the optional third argument, not a %T", v))
		}
		key = s
	}
	var result []interface{}
	for _, src := range list {
		at := map[string]interface{}{"src": src}
		fn.Eval(root, at, fn.Args...)
		result = append(result, at[key])
	}
	return result
}
