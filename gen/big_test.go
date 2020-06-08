// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestBigString(t *testing.T) {
	b := gen.Big("1234")

	tt.Equal(t, "1234", b.String())
}

func TestBigSimplify(t *testing.T) {
	b := gen.Big("1234")
	simple := b.Simplify()

	tt.Equal(t, "string 1234", fmt.Sprintf("%T %v", simple, simple))
}

func TestBigAlter(t *testing.T) {
	b := gen.Big("1234")
	alt := b.Alter()

	tt.Equal(t, "string 1234", fmt.Sprintf("%T %v", alt, alt))
}

func TestBigDup(t *testing.T) {
	b := gen.Big("1234")

	dup := b.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "1234", dup.String())
}

func TestBigEmpty(t *testing.T) {
	tt.Equal(t, true, gen.Big("").Empty())
	tt.Equal(t, false, gen.Big("1").Empty())
}
