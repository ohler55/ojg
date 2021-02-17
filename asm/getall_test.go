// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestGetall(t *testing.T) {
	root := testPlan(t,
		`[
           {one:1 two:2}
           [set $.at [getall "@.*"]]
           [set $.root [getall "$.src.*"]]
           [set $.arg [getall "@.*" {x:1 y:2}]]
         ]`,
		"{src: {a:1 b:2 c:3}}",
	)
	got, _ := root["at"].([]interface{})
	sort.Slice(got, func(i, j int) bool {
		a, _ := got[i].(int64)
		b, _ := got[j].(int64)
		return a < b
	})
	tt.Equal(t, "[1 2]", sen.String(got))

	got, _ = root["root"].([]interface{})
	sort.Slice(got, func(i, j int) bool {
		a, _ := got[i].(int64)
		b, _ := got[j].(int64)
		return a < b
	})
	tt.Equal(t, "[1 2 3]", sen.String(got))

	got, _ = root["arg"].([]interface{})
	sort.Slice(got, func(i, j int) bool {
		a, _ := got[i].(int64)
		b, _ := got[j].(int64)
		return a < b
	})
	tt.Equal(t, "[1 2]", sen.String(got))
}

func TestGetallArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"getall"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestGetallArgNotExpr(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"getall", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestGetallArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"getall", []interface{}{"sum"}},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
