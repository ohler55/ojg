// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestDel(t *testing.T) {
	root := testPlan(t,
		`[
           {one:1 two:2 three:3}
           [del @.one]
           [set $.asm @]
           [del $.asm.three]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{two:2}", sen.String(root["asm"]))
}

func TestDelExprError(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		map[string]interface{}{},
		[]interface{}{"del", jp.D()},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestDelArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"del"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestDelArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"del", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
