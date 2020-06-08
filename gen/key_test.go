// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestKey(t *testing.T) {
	k := gen.Key("sample")
	tt.Equal(t, "sample", k.String(), "String()")
	tt.Equal(t, "sample", k.Alter(), "Alter()")
	tt.Equal(t, "sample", k.Simplify(), "Simplify()")
	dup := k.Dup()
	tt.Equal(t, "gen.Key sample", fmt.Sprintf("%T %s", dup, dup), "Dup()")
	tt.Equal(t, false, k.Empty(), "Empty()")
}
