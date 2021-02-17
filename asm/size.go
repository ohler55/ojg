// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "size",
		Eval: size,
		Desc: `Returns the size or length of a string, array, or object (map).
For all other types zero is returned`,
	})
}

func size(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 1 {
		panic(fmt.Errorf("size expects exactly one arguments. %d given", len(args)))
	}
	var length int
	switch tv := evalArg(root, at, args[0]).(type) {
	case string:
		length = len(tv)
	case []interface{}:
		length = len(tv)
	case map[string]interface{}:
		length = len(tv)
	}
	return length
}
