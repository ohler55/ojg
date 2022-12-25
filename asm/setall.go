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

func setall(root map[string]any, at any, args ...any) any {
	if len(args) != 2 {
		panic(fmt.Errorf("setall expects exactly two arguments. %d given", len(args)))
	}
	var x jp.Expr
	switch v := args[0].(type) {
	case jp.Expr:
		x = v
	case *Fn:
		if x, _ = evalArg(root, at, v).(jp.Expr); x == nil {
			panic(fmt.Errorf("the first argument to setall must be a path not a %T", v))
		}
	default:
		panic(fmt.Errorf("the first argument to setall must be a path not a %T", v))
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
