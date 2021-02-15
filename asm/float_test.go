// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestFloat(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [float 1]]
           [set $.asm.b [float 2.2]]
           [set $.asm.c [float "1.23"]]
           [set $.asm.d [float [time "2021-02-09T01:02:03.123456Z"]]]
           [set $.asm.e [float true]]
           [set $.asm.f [float abc]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 1
  b: 2.2
  c: 1.23
  d: 1.612832523123456e+09
  e: null
  f: null
}`, sen.String(root["asm"], &opt))
}

func TestFloatArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"float", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
