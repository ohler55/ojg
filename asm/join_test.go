// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestJoin(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [join [a b c] +]]
           [set $.asm.b [join [a b c]]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: a+b+c
  b: abc
}`, sen.String(root["asm"], &opt))
}

func TestJoinArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"join", []interface{}{}, "x", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestJoinArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"join", 1, "x"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestJoinArgType2(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"join", []interface{}{}, 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestJoinArgType3(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"join", []interface{}{"x", 3}},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
