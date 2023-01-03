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

func size(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("size expects exactly one arguments. %d given", len(args)))
	}
	var length int
	switch tv := evalArg(root, at, args[0]).(type) {
	case string:
		length = len(tv)
	case []any:
		length = len(tv)
	case map[string]any:
		length = len(tv)
	}
	return length
}
