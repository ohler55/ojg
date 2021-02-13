// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestPlanSimplify(t *testing.T) {
	p := asm.Plan{
		Fn: asm.Fn{
			Name: "fun",
			Args: []interface{}{
				&asm.Fn{Name: "+", Args: []interface{}{3, 4}},
				&asm.Fn{Name: "list", Args: []interface{}{1, 2, 3}},
			},
		},
	}
	tt.Equal(t, "[fun [+ 3 4] [list 1 2 3]]", sen.String(&p), "plan simplify")
}
