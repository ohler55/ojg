// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestArrayString(t *testing.T) {
	a := gen.Array{gen.Int(3), gen.Array{gen.Int(5)}, gen.Int(7)}

	tt.Equal(t, "[3,[5],7]", a.String())
}

func TestArraySimple(t *testing.T) {
	a := gen.Array{gen.Int(3), gen.Int(7)}
	simple := a.Simplify()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", simple, simple))
}

func TestArrayAlter(t *testing.T) {
	a := gen.Array{gen.Int(3), gen.Int(7)}
	alt := a.Alter()

	tt.Equal(t, "[]interface {} [3 7]", fmt.Sprintf("%T %v", alt, alt))

	aa := alt.([]interface{})
	tt.Equal(t, "int64 3  int64 7", fmt.Sprintf("%T %v  %T %v", aa[0], aa[0], aa[1], aa[1]))
}
