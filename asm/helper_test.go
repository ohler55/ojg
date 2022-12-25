// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

var sopt = sen.Options{Sort: true, TimeFormat: time.RFC3339Nano}

func testPlan(t *testing.T, plan, root string) map[string]any {
	parser := sen.Parser{}
	val, err := parser.Parse([]byte(plan))
	tt.Nil(t, err)
	list, _ := val.([]any)
	p := asm.NewPlan(list)

	val, err = parser.Parse([]byte(root))
	tt.Nil(t, err)
	r, _ := val.(map[string]any)
	err = p.Execute(r)
	tt.Nil(t, err)

	return r
}
