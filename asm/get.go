// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "get",
		Eval: get,
		Desc: `Gets the first matching value in either the root ($), local (@),
or if present, the second argument. The required first argument
must be a path and the option second argument is the
data to apply the path to. The jp.First() function is used to
get the results`,
	})
}

func get(root map[string]interface{}, at interface{}, args ...interface{}) (val interface{}) {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("get expects one or two arguments. %d given", len(args)))
	}
	var x jp.Expr
	switch v := args[0].(type) {
	case jp.Expr:
		x = v
	case *Fn:
		if x, _ = evalArg(root, at, v).(jp.Expr); x == nil {
			panic(fmt.Errorf("the first argument to get must be a path not a %T", v))
		}
	default:
		panic(fmt.Errorf("the first argument to get must be a path not a %T", v))
	}
	if 0 < len(x) {
		if 1 < len(args) {
			val = x.First(evalArg(root, at, args[1]))
		} else if _, ok := x[0].(jp.At); ok {
			val = x.First(at)
		} else {
			val = x.First(root)
		}
	}
	return
}
