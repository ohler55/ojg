// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSum(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [sum 1 2 3]]
           [set $.asm.b ['+' 2.2 1]]
           [set $.asm.c ['+' x 1]]
           [set $.asm.d [sum x y]]
           [set $.asm.e [sum 1.1 x 2.2]]
           [set $.asm.f [sum]]
           [set $.asm.g [sum 1 2.3]]
           [set $.asm.h [sum 1 x 2]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 6
  b: 3.2
  c: x1
  d: xy
  e: "1.1x2.2"
  f: 0
  g: 3.3
  h: "1x2"
}`, sen.String(root["asm"], &opt))
}

func TestSumArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"sum", 1, true},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
