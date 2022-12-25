// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestTimeCheck(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [time? [time "2021-02-09T01:02:03Z"]]]
           [set $.asm.b [time? 123]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestTimeCheckArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"time?", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestTimeConv(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [time "2021-02-09T01:02:03Z"]]
           [set $.asm.b [time "2021-02-09T01:02:03.123456789Z"]]
           [set $.asm.c [time 1612832523]]
           [set $.asm.d [time 1612832523123456789]]
           [set $.asm.e [time 1612832523.123456789]]
           [set $.asm.f [time "05 Jan 2021 -0400" "02 Jan 2006 -0700"]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	// Note the golang float64 does not have enough precision to represent a
	// time with nonoseconds.
	tt.Equal(t,
		`{
  a: "2021-02-09T01:02:03Z"
  b: "2021-02-09T01:02:03.123456789Z"
  c: "2021-02-09T01:02:03Z"
  d: "2021-02-09T01:02:03.123456789Z"
  e: "2021-02-09T01:02:03.123456716Z"
  f: "2021-01-05T00:00:00-04:00"
}`, sen.String(root["asm"], &opt))
}

func TestTimeConvArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"time", 1, 2, 3},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestTimeConvFormatType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"time", "2021", 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestTimeConvParseErr(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"time", "Jan 05 2021"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
