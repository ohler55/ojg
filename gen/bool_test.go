// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestBoolString(t *testing.T) {
	b := gen.Bool(true)

	tt.Equal(t, "true", b.String())
}

func TestBoolSimplify(t *testing.T) {
	simple := gen.True.Simplify()

	tt.Equal(t, "bool true", fmt.Sprintf("%T %v", simple, simple))
}

func TestBoolAlter(t *testing.T) {
	b := gen.False
	alt := b.Alter()

	tt.Equal(t, "bool false", fmt.Sprintf("%T %v", alt, alt))
}

func TestBoolDup(t *testing.T) {
	dup := gen.True.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "true", dup.String())
}

func TestBoolEmpty(t *testing.T) {
	tt.Equal(t, false, gen.Bool(true).Empty())
}
