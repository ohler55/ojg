// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestKey(t *testing.T) {
	k := ojg.Key("sample")
	tt.Equal(t, "sample", k.String(), "String()")
	tt.Equal(t, "sample", k.Alter(), "Alter()")
	tt.Equal(t, "sample", k.Simplify(), "Simplify()")
	dup := k.Dup()
	fmt.Sprintf("%T %s", dup, dup)
	tt.Equal(t, "ojg.Key sample", fmt.Sprintf("%T %s", dup, dup), "Dup()")
	tt.Equal(t, false, k.Empty(), "Empty()")
}
