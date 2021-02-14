// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "setall",
		Eval: setall,
		Desc: `Sets multiple values in either the root ($) or local (@) data.
Two arguments are required, the first must be a path and the
second argument is evaluate to a value and inserted using the
jp.Set() function.`,
	})
}

func setall(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("setall expects exactly two arguments. %d given", len(args)))
	}
	x, _ := args[0].(jp.Expr)
	if x == nil {
		panic(fmt.Errorf("the first argument to setall must be a path not a %T", args[0]))
	}
	arg := evalArg(root, at, args[1])
	var err error
	if 0 < len(x) {
		if _, ok := x[0].(jp.At); ok {
			err = x.Set(at, arg)
		} else {
			err = x.Set(root, arg)
		}
	}
	if err != nil {
		panic(err)
	}
	return at
}
