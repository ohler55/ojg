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

func getall(root map[string]any, at any, args ...any) (list any) {
	if len(args) < 1 || 2 < len(args) {
		panic(fmt.Errorf("getall expects one or two arguments. %d given", len(args)))
	}
	var x jp.Expr
	switch v := args[0].(type) {
	case jp.Expr:
		x = v
	case *Fn:
		if x, _ = evalArg(root, at, v).(jp.Expr); x == nil {
			panic(fmt.Errorf("the first argument to getall must be a path not a %T", v))
		}
	default:
		panic(fmt.Errorf("the first argument to getall must be a path not a %T", v))
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
