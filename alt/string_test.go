// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type stringData struct {
	value    any
	defaults []string
	expect   string
}

func TestString(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	for _, d := range []stringData{
		{value: 3, expect: "3"},
		{value: int8(3), expect: "3"},
		{value: int16(3), expect: "3"},
		{value: int32(3), expect: "3"},
		{value: int64(3), expect: "3"},
		{value: uint(3), expect: "3"},
		{value: uint8(3), expect: "3"},
		{value: uint16(3), expect: "3"},
		{value: uint32(3), expect: "3"},
		{value: uint64(3), expect: "3"},
		{value: gen.Int(3), expect: "3"},

		{value: nil, expect: ""},
		{value: nil, expect: "nil", defaults: []string{"x", "nil"}},

		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: true, expect: "true", defaults: []string{"test"}},
		{value: true, expect: "yes", defaults: []string{"test", "yes"}},
		{value: gen.True, expect: "true"},
		{value: gen.False, expect: "false"},

		{value: 3.1, expect: "3.1"},
		{value: 3.2, expect: "d2", defaults: []string{"d1", "d2"}},
		{value: 3.2, expect: "3.2", defaults: []string{"d1"}},
		{value: float32(3.2), expect: "3.2"},
		{value: gen.Float(3.1), expect: "3.1"},
		{value: gen.Float(3.2), expect: "d2", defaults: []string{"d1", "d2"}},
		{value: gen.Float(3.2), expect: "3.2", defaults: []string{"d1"}},

		{value: gen.Big("3.1"), expect: "3.1"},

		{value: "xyz", expect: "xyz"},
		{value: "xyz", expect: "xyz", defaults: []string{"d1"}},
		{value: "xyz", expect: "xyz", defaults: []string{"d1", "d2"}},
		{value: gen.String("xyz"), expect: "xyz"},
		{value: gen.String("xyz"), expect: "xyz", defaults: []string{"d1"}},
		{value: gen.String("xyz"), expect: "xyz", defaults: []string{"d1", "d2"}},

		{value: tm, expect: "2020-04-12T16:34:04.123456789Z"},
		{value: gen.Time(tm), expect: "2020-04-12T16:34:04.123456789Z"},

		{value: []any{}, expect: ""},
		{value: []any{}, expect: "empty", defaults: []string{"empty"}},
	} {
		result := alt.String(d.value, d.defaults...)
		tt.Equal(t, d.expect, result, "String(", d.value, d.defaults, ")")
	}
}
