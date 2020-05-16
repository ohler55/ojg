// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestStringString(t *testing.T) {
	b := oj.String("abc")

	tt.Equal(t, `"abc"`, b.String())
}

func TestStringSimplify(t *testing.T) {
	b := oj.String("abc")
	simple := b.Simplify()

	tt.Equal(t, "string abc", fmt.Sprintf("%T %v", simple, simple))
}

func TestStringAlter(t *testing.T) {
	b := oj.String("abc")
	alt := b.Alter()

	tt.Equal(t, "string abc", fmt.Sprintf("%T %v", alt, alt))
}

func TestStringDup(t *testing.T) {
	b := oj.String("abc")

	dup := b.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, `"abc"`, dup.String())
}

func TestStringEmpty(t *testing.T) {
	tt.Equal(t, true, oj.String("").Empty())
	tt.Equal(t, false, oj.String("1").Empty())
}
