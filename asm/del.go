// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "del",
		Eval: delEval,
		Desc: `Deletes the first matching value in either the root ($) or
local (@) data. Exactly one argument is required and it must be
a path. The jp.DelOne() function is used to delete the value.
The local (@) value is returned.`,
	})
}

func delEval(root map[string]interface{}, at interface{}, args ...interface{}) (list interface{}) {
	if len(args) != 1 {
		panic(fmt.Errorf("del expects exactly one arguments. %d given", len(args)))
	}
	x, _ := args[0].(jp.Expr)
	if x == nil {
		panic(fmt.Errorf("the first argument to del must be a path not a %T", args[0]))
	}
	var err error
	if 0 < len(x) {
		if _, ok := x[0].(jp.At); ok {
			err = x.DelOne(at)
		} else {
			err = x.DelOne(root)
		}
	}
	if err != nil {
		panic(err)
	}
	return at
}
