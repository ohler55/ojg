// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestQuotient(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [quotient 42 3 2]]
           [set $.asm.b ["/" 5.0 2]]
           [set $.asm.c [quotient 5 2]]
           [set $.asm.d [quotient]]
           [set $.asm.e [quotient 1]]
           [set $.asm.f [quotient 1.2]]
           [set $.asm.g ["/" 2.5 0.5]]
           [set $.asm.h ["/" 5 2.0]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 7
  b: 2.5
  c: 2
  d: 0
  e: 1
  f: 1.2
  g: 5
  h: 2.5
}`, sen.String(root["asm"], &opt))
}

func TestQuotientArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"quotient", 1, true},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestQuotientZero(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"quotient", 1, 0},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
