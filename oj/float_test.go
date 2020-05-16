// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestFloatString(t *testing.T) {
	b := oj.Float(12.34)

	tt.Equal(t, "12.34", b.String())
}

func TestFloatSimplify(t *testing.T) {
	b := oj.Float(12.34)
	simple := b.Simplify()

	tt.Equal(t, "float64 12.34", fmt.Sprintf("%T %v", simple, simple))
}

func TestFloatAlter(t *testing.T) {
	b := oj.Float(12.34)
	alt := b.Alter()

	tt.Equal(t, "float64 12.34", fmt.Sprintf("%T %v", alt, alt))
}

func TestFloatDup(t *testing.T) {
	b := oj.Float(12.34)

	dup := b.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "12.34", dup.String())
}

func TestFloatEmpty(t *testing.T) {
	tt.Equal(t, false, oj.Float(12.34).Empty())
}
