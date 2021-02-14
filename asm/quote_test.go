// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestQuote(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [quote @.src]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "@.src", sen.String(root["asm"]))
}
