// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSetall(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		map[string]interface{}{"x": 2, "y": 3}, // Sets @
		[]interface{}{"setall", "@.*", 1},
		[]interface{}{"setall", "$.asm", "@"},
	})
	tt.Equal(t, `[asm {x:2 y:3} [setall "@.*" 1] [setall $.asm @]]`, sen.String(p, &sopt), "plan string")

	root := map[string]interface{}{
		"src": []interface{}{1, 2, 3},
	}
	err := p.Execute(root)
	tt.Nil(t, err)
	tt.Equal(t, "{x:1 y:1}", sen.String(root["asm"], &sopt))
}

func TestSetallExprError(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		map[string]interface{}{}, // Sets @
		[]interface{}{"setall", jp.D(), 1},
		[]interface{}{"setall", "$.asm", "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSetallArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"setall", "@.x"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSetallArgNotExprT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"setall", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
