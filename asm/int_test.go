// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestInt(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [int 1]]
           [set $.asm.b [int 2.2]]
           [set $.asm.c [int "123"]]
           [set $.asm.d [int [time "2021-02-09T01:02:03.123456789Z"]]]
           [set $.asm.e [int true]]
           [set $.asm.f [int abc]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 1
  b: 2
  c: 123
  d: 1612832523123456789
  e: null
  f: null
}`, sen.String(root["asm"], &opt))
}

func TestIntArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"int", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
