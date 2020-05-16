// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestArrayString(t *testing.T) {
	a := oj.Array{oj.Int(3), oj.Array{oj.Int(5)}, nil}

	tt.Equal(t, "[3,[5],null]", a.String())
}

func TestArraySimplify(t *testing.T) {
	a := oj.Array{oj.Int(3), oj.Int(7)}
	simple := a.Simplify()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", simple, simple))
}

func TestArrayAlter(t *testing.T) {
	a := oj.Array{oj.Int(3), oj.Int(7)}
	alt := a.Alter()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", alt, alt))

	aa := alt.([]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", aa[0], aa[0], aa[1], aa[1]))
}

func TestArrayDup(t *testing.T) {
	a := oj.Array{oj.Int(3), oj.Array{oj.Int(5)}, nil}

	dup := a.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "[3,[5],null]", dup.String())
}

func TestArrayEmpty(t *testing.T) {
	tt.Equal(t, true, oj.EmptyArray.Empty())
}
