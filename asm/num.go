// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
)

func init() {
	Define(&Fn{
		Name: "num?",
		Eval: num,
		Desc: `Returns true if the single required argumement is number
otherwise false is returned.`,
	})
}

func num(root map[string]any, at any, args ...any) any {
	if len(args) != 1 {
		panic(fmt.Errorf("num? expects exactly one arguments. %d given", len(args)))
	}
	ok := false
	switch evalArg(root, at, args[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		ok = true
	case float32, float64:
		ok = true
	}
	return ok
}
