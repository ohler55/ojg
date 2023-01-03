// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestSizeString(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [size a_string]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "8", sen.String(root["asm"]))
}

func TestSizeArray(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [size [1 2 3]]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "3", sen.String(root["asm"]))
}

func TestSizeMap(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [size {a:1 b:2 c:3}]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "3", sen.String(root["asm"]))
}

func TestSizeOther(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [size true]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "0", sen.String(root["asm"]))
}

func TestSizeArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"size", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
