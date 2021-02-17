// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestNeqNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.true [neq 1 2]]
           [set $.asm.false [neq 1 1.0]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{false:false true:true}", sen.String(root["asm"], &sopt))
}
