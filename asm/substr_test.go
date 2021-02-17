// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSubstr(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [substr abcdef 2 2]]
           [set $.asm.b [substr abcdef 2]]
           [set $.asm.c [substr abcdef 2 -2]]
           [set $.asm.d [substr abcdef -3 2]]
           [set $.asm.e [substr abcdef -7 2]]
           [set $.asm.f [substr abcdef 3 4]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: cd
  b: cdef
  c: ""
  d: de
  e: ab
  f: def
}`, sen.String(root["asm"], &opt))
}

func TestSubstrArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"substr", "x", 1, 1, 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSubstrArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"substr", 1, 1, 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSubstrArgType2(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"substr", "x", true, 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSubstrArgType3(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"substr", "x", 1, true},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
