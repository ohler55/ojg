// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestGet(t *testing.T) {
	root := testPlan(t,
		`[
           {one:1 two:2}
           [set $.at [get "@.*"]]
           [set $.root [get "$.src.*"]]
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
}

func TestGetArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"get"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestGetArgNotExprT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"get", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
