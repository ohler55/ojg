// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "at",
		Eval: at,
		Desc: `Forms a path starting with @. The remaining string arguments are
joined with a '.' and parsed to form a jp.Expr.`,
	})
}

func at(root map[string]interface{}, at interface{}, args ...interface{}) interface{} {
	var b []byte
	for i, arg := range args {
		v := evalArg(root, at, arg)
		s, ok := v.(string)
		if !ok {
			panic(fmt.Errorf("at expected a string argument, not a %T", v))
		}
		if 0 < i {
			b = append(b, '.')
		}
		b = append(b, s...)
	}
	x, err := jp.Parse(b)
	if err != nil {
		panic(err)
	}
	return append(jp.A(), x...)
}
