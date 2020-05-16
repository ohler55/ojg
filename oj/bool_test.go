// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestBoolString(t *testing.T) {
	b := oj.Bool(true)

	tt.Equal(t, "true", b.String())
}

func TestBoolSimplify(t *testing.T) {
	simple := oj.True.Simplify()

	tt.Equal(t, "bool true", fmt.Sprintf("%T %v", simple, simple))
}

func TestBoolAlter(t *testing.T) {
	b := oj.False
	alt := b.Alter()

	tt.Equal(t, "bool false", fmt.Sprintf("%T %v", alt, alt))
}

func TestBoolDup(t *testing.T) {
	dup := oj.True.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "true", dup.String())
}

func TestBoolEmpty(t *testing.T) {
	tt.Equal(t, false, oj.Bool(true).Empty())
}
