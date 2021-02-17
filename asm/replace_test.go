// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestReplace(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [replace hello l x]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t,
		`{a:hexxo}`, sen.String(root["asm"], &sopt))
}

func TestReplaceArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"replace", "x", "y", "z", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestReplaceArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"replace", 1, "x", "y"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestReplaceArgType2(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"replace", "x", 1, "y"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestReplaceArgType3(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"replace", "x", "y", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
