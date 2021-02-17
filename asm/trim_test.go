// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestTrim(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [trim " a string "]]
           [set $.asm.b [trim str]]
           [set $.asm.c [trim "-- a string ---" " -"]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: "a string"
  b: str
  c: "a string"
}`, sen.String(root["asm"], &opt))
}

func TestTrimArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"trim", "x", "y", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestTrimArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"trim", 1, "x"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestTrimArgType2(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"trim", "x", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
