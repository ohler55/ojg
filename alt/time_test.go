// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type timeData struct {
	value    interface{}
	defaults []time.Time
	expect   time.Time
}

func TestTime(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	tm2 := time.Date(2020, time.April, 12, 16, 34, 04, 111112000, time.UTC)
	tm3 := time.Date(2020, time.April, 12, 16, 34, 0, 0, time.UTC)
	var zero time.Time
	for _, d := range []timeData{
		{value: tm, expect: tm},
		{value: gen.Time(tm), expect: tm},

		{value: int8(3), expect: zero},
		{value: int64(1586709244123456789), expect: tm},
		{value: uint64(1586709244123456789), expect: tm},
		{value: int(1586709244123456789), expect: tm},
		{value: uint(1586709244123456789), expect: tm},
		{value: gen.Int(1586709244123456789), expect: tm},

		{value: nil, expect: zero},
		{value: nil, expect: tm2, defaults: []time.Time{tm, tm2}},

		{value: 1586709244.111112, expect: tm2},
		{value: 1586709244.111112, expect: tm, defaults: []time.Time{zero, tm}},
		{value: gen.Float(1586709244.111112), expect: tm2},
		{value: gen.Float(1586709244.111112), expect: tm, defaults: []time.Time{zero, tm}},

		{value: float32(1586709244.0), expect: tm3},

		{value: "2020-04-12T16:34:04.123456789Z", expect: tm},
		{value: gen.String("2020-04-12T16:34:04.123456789Z"), expect: tm},
		{value: "x", expect: tm, defaults: []time.Time{tm}},
		{value: gen.String("x"), expect: tm, defaults: []time.Time{tm}},

		{value: []interface{}{}, expect: zero},
		{value: []interface{}{}, expect: tm, defaults: []time.Time{tm}},
	} {
		result := alt.Time(d.value, d.defaults...)
		tt.Equal(t, d.expect, result, "Time(", d.value, d.defaults, ")")
	}
}
