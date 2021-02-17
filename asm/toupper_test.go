// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestToupper(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [toupper low]]
           [set $.asm.b [toupper UP]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: LOW
  b: UP
}`, sen.String(root["asm"], &opt))
}

func TestToupperArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"toupper", "x", "y"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestToupperArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"toupper", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
