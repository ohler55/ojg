// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gd"
	"github.com/ohler55/ojg/tt"
)

func TestArrayString(t *testing.T) {
	a := gd.Array{gd.Int(3), gd.Array{gd.Int(5)}, gd.Int(7)}

	tt.Equal(t, "[3,[5],7]", a.String())
}

func TestArrayJSON(t *testing.T) {
	a := gd.Array{gd.Int(3), gd.Array{gd.Int(5)}, gd.Int(7)}

	tt.Equal(t, "[3,[5],7]", a.JSON())
}

func TestArrayJSONIndent(t *testing.T) {
	a := gd.Array{gd.Int(3), gd.Array{gd.Int(5)}, gd.Int(7)}

	tt.Equal(t, `[
  3,
  [
    5
  ],
  7
]`, a.JSON(2))
}

func TestArrayNative(t *testing.T) {
	a := gd.Array{gd.Int(3), gd.Int(7)}
	native := a.Native()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", native, native))
}

func TestArrayAlter(t *testing.T) {
	a := gd.Array{gd.Int(3), gd.Int(7)}
	alt := a.Alter()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", alt, alt))

	aa := alt.([]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", aa[0], aa[0], aa[1], aa[1]))
}
