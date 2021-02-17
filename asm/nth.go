// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "nth",
		Eval: nth,
		Desc: `Returns a nth element of an array. The second argument must be
an integer that indicates the element of the array to return.
If the index is less than 0 then the index is from the end of
the array.`,
	})
}

func nth(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("nth expects exactly two arguments. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]interface{})
	if !ok {
		panic(fmt.Errorf("nth expected an array argument, not a %T", v))
	}
	v = evalArg(root, at, args[1])
	var index int64
	if index, ok = asInt(v); !ok {
		panic(fmt.Errorf("nth expects an integer second argument, not a %T", v))
	}
	if index < 0 {
		index = int64(len(list)) + index
	}
	if index < 0 || int64(len(list)) <= index {
		return nil
	}
	return list[index]
}
