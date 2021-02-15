// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestStringCheck(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [string? abc]]
           [set $.asm.b [string? 123]]
           [set $.asm.c [string? true]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:false c:false}", sen.String(root["asm"], &sopt))
}

func TestStringCheckArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"string?", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestStringConv(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [string 1]]
           [set $.asm.b [string 1 "0x%02x"]]
           [set $.asm.c [string 2.2]]
           [set $.asm.d [string 2.2 "%e"]]
           [set $.asm.e [string [time "2021-02-09T01:02:03.123456789Z"]]]
           [set $.asm.f [string [time "2021-02-09T01:02:03.123456789Z"] "02 Jan 2006"]]
           [set $.asm.g [string true]]
           [set $.asm.h [string [1 2]]]
           [set $.asm.i [string {a:1}]]
           [set $.asm.j [string abc]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: 1
  b: 0x01
  c: 2.2
  d: 2.200000e+00
  e: "2021-02-09T01:02:03.123456789Z"
  f: "09 Feb 2021"
  g: true
  h: "[1 2]"
  i: "{a:1}"
  j: abc
}`, sen.String(root["asm"], &opt))
}

func TestStringConvArgCountT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"string", 1, "x", 3},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}

func TestStringConvFormatTypeT(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"string", 1, 2},
	})
	err := p.Execute(map[string]interface{}{})
	tt.NotNil(t, err)
}
