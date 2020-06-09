// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestIntString(t *testing.T) {
	b := gen.Int(1234)

	tt.Equal(t, "1234", b.String())
}

func TestIntSimplify(t *testing.T) {
	b := gen.Int(1234)
	simple := b.Simplify()

	tt.Equal(t, "int64 1234", fmt.Sprintf("%T %v", simple, simple))
}

func TestIntAlter(t *testing.T) {
	b := gen.Int(1234)
	alt := b.Alter()

	tt.Equal(t, "int64 1234", fmt.Sprintf("%T %v", alt, alt))
}

func TestIntDup(t *testing.T) {
	b := gen.Int(1234)

	dup := b.Dup()
	tt.NotNil(t, dup)
	tt.Equal(t, "1234", dup.String())
}

func TestIntEmpty(t *testing.T) {
	tt.Equal(t, false, gen.Int(1234).Empty())
}
