// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "delall",
		Eval: delall,
		Desc: `Deletes the all matching values in either the root ($) or
local (@) data. Exactly one argument is required and it must be
a path. The jp.DelOne() function is used to delete the value.
The local (@) value is returned.`,
	})
}

func delall(root map[string]any, at any, args ...any) (list any) {
	if len(args) != 1 {
		panic(fmt.Errorf("delall expects exactly one arguments. %d given", len(args)))
	}
	x, _ := args[0].(jp.Expr)
	if x == nil {
		panic(fmt.Errorf("the first argument to delall must be a path not a %T", args[0]))
	}
	var err error
	if 0 < len(x) {
		if _, ok := x[0].(jp.At); ok {
			err = x.Del(at)
		} else {
			err = x.Del(root)
		}
	}
	if err != nil {
		panic(err)
	}
	return at
}
