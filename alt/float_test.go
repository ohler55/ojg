// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type floatData struct {
	value    interface{}
	defaults []float64
	expect   float64
}

func TestFloat(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for _, d := range []floatData{
		{value: 3, expect: 3.0},
		{value: int8(3), expect: 3.0},
		{value: int16(3), expect: 3.0},
		{value: int32(3), expect: 3.0},
		{value: int64(3), expect: 3.0},
		{value: uint(3), expect: 3.0},
		{value: uint8(3), expect: 3.0},
		{value: uint16(3), expect: 3.0},
		{value: uint32(3), expect: 3.0},
		{value: uint64(3), expect: 3.0},
		{value: gen.Int(3), expect: 3.0},

		{value: nil, expect: 0.0},
		{value: nil, expect: 5.5, defaults: []float64{4.4, 5.5}},

		{value: true, expect: 0.0},
		{value: true, expect: 4.4, defaults: []float64{4.4}},
		{value: true, expect: 5.5, defaults: []float64{4.4, 5.5}},
		{value: gen.True, expect: 0.0},
		{value: gen.True, expect: 4.4, defaults: []float64{4.4}},
		{value: gen.True, expect: 5.5, defaults: []float64{4.4, 5.5}},

		{value: 3.2, expect: 3.2},
		{value: 3.2, expect: 3.2, defaults: []float64{4.4}},
		{value: 3.3, expect: 3.3, defaults: []float64{4.4, 5.5}},

		{value: float32(3.2), expect: float64(float32(3.2))},
		{value: float32(3.2), expect: float64(float32(3.2)), defaults: []float64{4.4}},
		{value: float32(3.2), expect: float64(float32(3.2)), defaults: []float64{4.4, 5.5}},
		{value: gen.Float(3.1), expect: 3.1},
		{value: gen.Float(3.2), expect: 3.2, defaults: []float64{4.4}},
		{value: gen.Float(3.3), expect: 3.3, defaults: []float64{4.4, 5.5}},

		{value: "3.3", expect: 3.3},
		{value: "3.3", expect: 3.3, defaults: []float64{4.4}},
		{value: "3.3", expect: 5.5, defaults: []float64{4.4, 5.5}},

		{value: "3x", expect: 0.0},
		{value: "3x", expect: 4.4, defaults: []float64{4.4}},
		{value: "3.0", expect: 3.0},
		{value: "3.1", expect: 3.1},
		{value: gen.String("3.3"), expect: 3.3},
		{value: gen.String("3x"), expect: 5.5, defaults: []float64{4.4, 5.5}},
		{value: gen.Big("3.3"), expect: 3.3},

		{value: tm, expect: 1586709244.123456789},
		{value: tm, expect: 5.5, defaults: []float64{4.4, 5.5}},
		{value: gen.Time(tm), expect: 1586709244.123456789},
		{value: gen.Time(tm), expect: 5.5, defaults: []float64{4.4, 5.5}},

		{value: []interface{}{}, expect: 0.0},
		{value: []interface{}{}, expect: 4.4, defaults: []float64{4.4}},
		{value: []interface{}{}, expect: 5.5, defaults: []float64{4.4, 5.5}},
	} {
		result := alt.Float(d.value, d.defaults...)
		tt.Equal(t, d.expect, result, "Float(", d.value, d.defaults, ")")
	}
}
