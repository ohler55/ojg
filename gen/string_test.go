// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestStringString(t *testing.T) {
	b := gen.String("abc")

	tt.Equal(t, `"abc"`, b.String())
}

func TestStringSimplify(t *testing.T) {
	b := gen.String("abc")
	simple := b.Simplify()

	tt.Equal(t, "string abc", fmt.Sprintf("%T %v", simple, simple))
}

func TestStringAlter(t *testing.T) {
	b := gen.String("abc")
	alt := b.Alter()

	tt.Equal(t, "string abc", fmt.Sprintf("%T %v", alt, alt))
}

func TestStringDup(t *testing.T) {
	b := gen.String("abc")

	dup := b.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, `"abc"`, dup.String())
}

func TestStringEmpty(t *testing.T) {
	tt.Equal(t, true, gen.String("").Empty())
	tt.Equal(t, false, gen.String("1").Empty())
}
