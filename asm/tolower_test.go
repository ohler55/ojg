// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestTolower(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [tolower low]]
           [set $.asm.b [tolower UP]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: low
  b: up
}`, sen.String(root["asm"], &opt))
}

func TestTolowerArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"tolower", "x", "y"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestTolowerArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"tolower", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
