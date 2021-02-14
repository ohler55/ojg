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
		Desc: `Gets all matching values in either the root ($) or local (@)
data. Exactly one argument is required and it must be a path.
The jp.Get() function is used to get the results`,
	})
}

func get(root map[string]interface{}, at interface{}, args ...interface{}) (list interface{}) {
	if len(args) != 1 {
		panic(fmt.Errorf("get expects exactly one arguments. %d given", len(args)))
	}
	x, _ := args[0].(jp.Expr)
	if x == nil {
		panic(fmt.Errorf("the first argument to get must be a path not a %T", args[0]))
	}
	if 0 < len(x) {
		if _, ok := x[0].(jp.At); ok {
			list = x.Get(at)
		} else {
			list = x.Get(root)
		}
	}
	return list
}
