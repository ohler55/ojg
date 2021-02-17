// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestList(t *testing.T) {
	root := testPlan(t,
		`[
           [list 1 "$.src[1]" true $.bad]
           [set $.asm @]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "[1 2 true null]", sen.String(root["asm"]))
}

func TestListEmpty(t *testing.T) {
	root := testPlan(t,
		`[
           [list]
           [set $.asm @]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, "[]", sen.String(root["asm"]))
}
