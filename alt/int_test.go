// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type intData struct {
	value    any
	defaults []int64
	expect   int64
}

func TestInt(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for _, d := range []intData{
		{value: 3, expect: 3},
		{value: int8(3), expect: 3},
		{value: int16(3), expect: 3},
		{value: int32(3), expect: 3},
		{value: int64(3), expect: 3},
		{value: uint(3), expect: 3},
		{value: uint8(3), expect: 3},
		{value: uint16(3), expect: 3},
		{value: uint32(3), expect: 3},
		{value: uint64(3), expect: 3},
		{value: gen.Int(3), expect: 3},

		{value: nil, expect: 0},
		{value: nil, expect: 5, defaults: []int64{4, 5}},

		{value: true, expect: 0},
		{value: true, expect: 4, defaults: []int64{4}},
		{value: true, expect: 4, defaults: []int64{4, 5}},
		{value: gen.True, expect: 0},
		{value: gen.True, expect: 4, defaults: []int64{4}},
		{value: gen.True, expect: 4, defaults: []int64{4, 5}},

		{value: 3.1, expect: 3},
		{value: 3.2, expect: 3, defaults: []int64{4}},
		{value: 3.3, expect: 5, defaults: []int64{4, 5}},
		{value: float32(3.1), expect: 3},
		{value: float32(3.2), expect: 3, defaults: []int64{4}},
		{value: float32(3.3), expect: 5, defaults: []int64{4, 5}},
		{value: gen.Float(3.1), expect: 3},
		{value: gen.Float(3.2), expect: 3, defaults: []int64{4}},
		{value: gen.Float(3.3), expect: 5, defaults: []int64{4, 5}},

		{value: "3", expect: 3},
		{value: "3", expect: 3, defaults: []int64{4}},
		{value: "3", expect: 5, defaults: []int64{4, 5}},
		{value: "3x", expect: 0},
		{value: "3x", expect: 4, defaults: []int64{4}},
		{value: "3.0", expect: 3},
		{value: "3.1", expect: 3},
		{value: "3.2", expect: 4, defaults: []int64{4}},
		{value: "3.3", expect: 5, defaults: []int64{4, 5}},
		{value: gen.String("3"), expect: 3},
		{value: gen.String("3x"), expect: 5, defaults: []int64{4, 5}},
		{value: gen.Big("3"), expect: 3},

		{value: tm, expect: 1586709244123456789},
		{value: tm, expect: 1586709244123456789, defaults: []int64{4}},
		{value: tm, expect: 5, defaults: []int64{4, 5}},
		{value: gen.Time(tm), expect: 1586709244123456789},
		{value: gen.Time(tm), expect: 5, defaults: []int64{4, 5}},

		{value: []any{}, expect: 0},
		{value: []any{}, expect: 4, defaults: []int64{4, 5}},
	} {
		result := alt.Int(d.value, d.defaults...)
		tt.Equal(t, d.expect, result, "Int(", d.value, d.defaults, ")")
	}
}
