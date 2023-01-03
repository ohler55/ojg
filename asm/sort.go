// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"
	"sort"
	"time"

	"github.com/ohler55/ojg/jp"
)

func init() {
	Define(&Fn{
		Name: "sort",
		Eval: sortEval,
		Desc: `Sort the items in an array and return a copy of the array. Valid
types for comparison are strings, numbers, and times. Any other
type returned or a type mismatch will raise an error.`,
	})
}

func sortEval(root map[string]any, at any, args ...any) any {
	if len(args) != 2 {
		panic(fmt.Errorf("sort expects exactly two argument. %d given", len(args)))
	}
	v := evalArg(root, at, args[0])
	list, ok := v.([]any)
	if !ok {
		panic(fmt.Errorf("sort expects an array argument, not a %T", v))
	}
	// Make a copy so not to change the original.
	list2 := make([]any, len(list))
	_ = copy(list2, list)

	var x jp.Expr
	x, ok = args[1].(jp.Expr)
	if !ok {
		panic(fmt.Errorf("sort expects a path second argument, not a %T", args[1]))
	}
	sort.Slice(list2, func(i, j int) bool {
		vi := x.First(list2[i])
		vj := x.First(list2[j])
		switch ti := vi.(type) {
		case string:
			var sj string
			if sj, ok = vj.(string); !ok {
				panic(fmt.Errorf("sort has mixed key values, string vs %T", vj))
			}
			return ti < sj
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			fi, _ := asFloat(vi)
			var fj float64
			if fj, ok = asFloat(vj); !ok {
				panic(fmt.Errorf("sort has mixed key values, number vs %T", vj))
			}
			return fi < fj
		case time.Time:
			var tj time.Time
			if tj, ok = vj.(time.Time); !ok {
				panic(fmt.Errorf("sort has mixed key values, time vs %T", vj))
			}
			return ti.Before(tj)
		default:
			panic(fmt.Errorf("sort key values must be strings, numbers, or time, not %T", vi))
		}
	})
	return list2
}
