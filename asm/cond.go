// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "cond",
		Eval: cond,
		Desc: `A conditional construct modeled after the LISP cond. All
arguments must be array of two elements. The first element must
evaluate to a boolean and the second can be any value. The value
of the first true first argument is returned. If none match nil
is returned.`,
	})
}

func cond(root map[string]any, at any, args ...any) any {
	for _, arg := range args {
		list, ok := arg.([]any)
		if !ok {
			panic(fmt.Errorf("cond expects array arguments, not a %T", arg))
		}
		if len(list) != 2 {
			panic(fmt.Errorf("cond array arguments must have two elements, not a %d", len(list)))
		}
		if bv, _ := evalValue(root, at, list[0]).(bool); bv {
			return evalValue(root, at, list[1])
		}
	}
	return nil
}

func evalValue(root map[string]any, at any, value any) (result any) {
top:
	switch tv := value.(type) {
	case *Fn:
		result = tv.Eval(root, at, tv.Args...)
	case []any:
		if 0 < len(tv) {
			if name, _ := tv[0].(string); 0 < len(name) {
				if af := NewFn(name); af != nil {
					af.Args = tv[1:]
					af.compile()
					value = af
					goto top
				}
			}
		}
	case jp.Expr:
		if 0 < len(tv) {
			if _, ok := tv[0].(jp.At); ok {
				result = tv.First(at)
			} else {
				result = tv.First(root)
			}
		}
	case string:
		if 0 < len(tv) && (tv[0] == '$' || tv[0] == '@') {
			if x, err := jp.Parse([]byte(tv)); err == nil {
				value = x
				goto top
			}
		}
		result = tv
	default:
		result = value
	}
	return
}
