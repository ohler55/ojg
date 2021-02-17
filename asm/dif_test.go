// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestDif(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [dif 5 3 1]]
           [set $.asm.b ["-" 2.5 1]]
           [set $.asm.c [dif 5 2.5 1]]
           [set $.asm.d [dif]]
           [set $.asm.e [dif 1]]
           [set $.asm.f [dif 1.2]]
           [set $.asm.g ["-" 2.5 1.5]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 1
  b: 1.5
  c: 1.5
  d: 0
  e: 1
  f: 1.2
  g: 1
}`, sen.String(root["asm"], &opt))
}

func TestDifArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"dif", 1, true},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
