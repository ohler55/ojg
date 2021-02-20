// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSortShallow(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [sort [b a c] @]]
           [set $.asm.b [sort [2 5 1] @]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `{a:[a b c] b:[1 2 5]}`, sen.String(root["asm"], &sopt))
}

func TestSortDeep(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [sort [{x:2}{x:5}{x:1}] @.x]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `[{x:1}{x:2}{x:5}]`, sen.String(root["asm"], &sopt))
}

func TestSortTime(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [sort [list [time "2021-02-09T00:00:00Z"][time "2021-01-05T00:00:00Z"]] @]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `["2021-01-05T00:00:00Z" "2021-02-09T00:00:00Z"]`, sen.String(root["asm"], &sopt))
}

func TestSortArgCount(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{}, "@", 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortArgType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", 1, "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortArgType2(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{}, 1},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortMixedString(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{"x", 1}, "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortMixedNum(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{1, "x"}, "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortMixedTime(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{1, time.Now(), time.Now().Add(time.Hour)}, "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestSortWrongType(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"sort", []interface{}{true, false}, "@"},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
