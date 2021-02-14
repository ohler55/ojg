// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestQuote(t *testing.T) {
	p := asm.NewPlan([]interface{}{
		[]interface{}{"quote", "@.src"},
	})
	// TBD use set and then verify the output
	tt.Equal(t, "[asm [quote @.src]]", sen.String(p), "quote string")

	root := map[string]interface{}{
		"src": []interface{}{1, 2, 3},
	}
	err := p.Execute(root)
	tt.Nil(t, err)
}
