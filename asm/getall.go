// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "getall",
		Eval: getall,
		Desc: `Gets all matching values in either the root ($), or local (@),
or if present, the second argument. The required first argument
must be a path and the option second argument is the
data to apply the path to. The jp.Get() function is used to get
the results`,
	})
}

func getall(root map[string]interface{}, at interface{}, args ...interface{}) (list interface{}) {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("getall expects one or two arguments. %d given", len(args)))
	}
	x, _ := args[0].(jp.Expr)
	if x == nil {
		panic(fmt.Errorf("the first argument to getall must be a path not a %T", args[0]))
	}
	if 0 < len(x) {
		if 1 < len(args) {
			list = x.Get(evalArg(root, at, args[1]))
		} else if _, ok := x[0].(jp.At); ok {
			list = x.Get(at)
		} else {
			list = x.Get(root)
		}
	}
	return
}
