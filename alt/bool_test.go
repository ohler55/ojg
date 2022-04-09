// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type boolData struct {
	value    interface{}
	defaults []bool
	expect   bool
}

func TestBool(t *testing.T) {
	// tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for _, d := range []boolData{
		{value: nil, expect: false},
		{value: nil, expect: true, defaults: []bool{false, true}},
		{value: true, expect: true},
		{value: gen.True, expect: true},
		{value: 1, expect: false},
		{value: 0, expect: true, defaults: []bool{true}},
		{value: 0, expect: true, defaults: []bool{true, false}},
		{value: "true", expect: true, defaults: []bool{false}},
		{value: "true", expect: false, defaults: []bool{true, false}},
		{value: "false", expect: false},
		{value: "yes", expect: false},
		{value: "no", expect: true, defaults: []bool{true}},
		{value: gen.String("true"), expect: true, defaults: []bool{false}},
		{value: gen.String("true"), expect: false, defaults: []bool{true, false}},
		{value: gen.String("false"), expect: false},
		{value: gen.String("yes"), expect: false},
		{value: gen.String("no"), expect: true, defaults: []bool{true}},
	} {
		result := alt.Bool(d.value, d.defaults...)
		tt.Equal(t, d.expect, result, "Bool(", d.value, d.defaults, ")")
	}
}
