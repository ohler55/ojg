// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestGteInt(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [gte 4 "$.src[2]" 3]]
           [set $.asm.b [gte 1 2]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestGteFloat(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [">=" 4.1 "$.src[1]" 2.2]]
           [set $.asm.b [">=" 1.0 2.0]]
         ]`,
		"{src: [1.1 2.2]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestGteString(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [gte xyz "$.src[0]" abc]]
           [set $.asm.b [gte abc def]]
         ]`,
		"{src: [abc xyz]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestGteWrongType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", "$.asm.i", []any{"gte", true, false}},
	})
	root := map[string]any{}
	err := p.Execute(root)
	tt.NotNil(t, err)
}

func TestGteWrongType2(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", "$.asm.i", []any{"gte", 1, false}},
	})
	root := map[string]any{}
	err := p.Execute(root)
	tt.NotNil(t, err)
}
