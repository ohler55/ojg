// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestEachNumbers(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [each [1 2 3] [set @.asm [sum 1 @.src]]]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `[2 3 4]`, sen.String(root["asm"], &sopt))
}

func TestEachFromRoot(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [each [getall "$.src.*"] [set @.xyz [sum 1 @.src.x]] xyz]]
         ]`,
		"{src: {a:{x:1}}}",
	)
	tt.Equal(t, `[2]`, sen.String(root["asm"], &sopt))
}

func TestEachArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"each", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestEachArgList(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"each", 1, []interface{}{"list"}},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestEachArgSecond(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"each", []interface{}{1, 2, 3}, true},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestEachArgThird(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"each", []interface{}{1, 2, 3}, []interface{}{"+"}, true},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
