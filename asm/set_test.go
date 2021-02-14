// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSet(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		map[string]interface{}{}, // Sets @
		[]interface{}{"set", "@.one", 1},
		[]interface{}{"set", "$.asm", "@"},
	})
	tt.Equal(t, "[asm {} [set @.one 1] [set $.asm @]]", sen.String(p), "plan string")

	root := map[string]interface{}{
		"src": []interface{}{1, 2, 3},
	}
	err := p.Execute(root)
	tt.Nil(t, err)
	tt.Equal(t, "{one:1}", sen.String(root["asm"]))
}

func TestSetExprError(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		map[string]interface{}{}, // Sets @
		[]interface{}{"set", jp.D(), 1},
		[]interface{}{"set", "$.asm", "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSetArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"set", "@.x"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSetArgNotExprT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"set", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
