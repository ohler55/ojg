// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "set",
		Eval: set,
		Desc: `Sets a single value in either the root ($) or local (@) data. Two
arguments are required, the first must be a path and the second
argument is evaluate to a value and inserted using the
jp.SetOne() function.`,
	})
}

func set(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	if len(args) != 2 {
		panic(fmt.Errorf("set expects exactly two arguments. %d given", len(args)))
	}
	var x jp.Expr
	switch v := args[0].(type) {
	case jp.Expr:
		x = v
	case *Fn:
		if x, _ = evalArg(root, at, v).(jp.Expr); x == nil {
			panic(fmt.Errorf("the first argument to set must be a path not a %T", v))
		}
	default:
		panic(fmt.Errorf("the first argument to set must be a path not a %T", v))
	}
	arg := evalArg(root, at, args[1])
	var err error
	if 0 < len(x) {
		if _, ok := x[0].(jp.At); ok {
			err = x.SetOne(at, arg)
		} else {
			err = x.SetOne(root, arg)
		}
	}
	if err != nil {
		panic(err)
	}
	return at
}
