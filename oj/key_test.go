// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestKey(t *testing.T) {
	k := oj.Key("sample")
	tt.Equal(t, "sample", k.String(), "String()")
	tt.Equal(t, "sample", k.Alter(), "Alter()")
	tt.Equal(t, "sample", k.Simplify(), "Simplify()")
	dup := k.Dup()
	tt.Equal(t, "oj.Key sample", fmt.Sprintf("%T %s", dup, dup), "Dup()")
	tt.Equal(t, false, k.Empty(), "Empty()")
}
