// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

func TestDiffPrimitive(t *testing.T) {
	type prim struct {
		v0     interface{}
		v1     interface{}
		expect bool
	}
	for _, p := range []prim{
		{v0: nil, v1: nil, expect: true},
		{v0: nil, v1: true, expect: false},
		{v0: false, v1: true, expect: false},
		{v0: false, v1: false, expect: true},

		{v0: 3, v1: int64(3), expect: true},
		{v0: 3, v1: int32(3), expect: true},
		{v0: 3, v1: int16(3), expect: true},
		{v0: 3, v1: int8(3), expect: true},
		{v0: 3, v1: uint64(3), expect: true},
		{v0: 3, v1: uint32(3), expect: true},
		{v0: 3, v1: uint16(3), expect: true},
		{v0: 3, v1: uint8(3), expect: true},
		{v0: 3, v1: uint(3), expect: true},
		{v0: 3, v1: gen.Int(3), expect: true},
		{v0: 3, v1: gen.Float(3.0), expect: true},
		{v0: 3, v1: 3.0, expect: true},
		{v0: 3, v1: float32(3.0), expect: true},
		{v0: 3, v1: float32(3.1), expect: false},
		{v0: 3, v1: gen.Float(3.1), expect: false},
		{v0: 3, v1: 3.1, expect: false},
		{v0: 3, v1: 4, expect: false},
		{v0: 3, v1: true, expect: false},

		{v0: 3.0, v1: float32(3.0), expect: true},
		{v0: 3.0, v1: 3, expect: true},
		{v0: 3.0, v1: uint(3), expect: true},
		{v0: 3.0, v1: uint8(3), expect: true},
		{v0: 3.0, v1: uint16(3), expect: true},
		{v0: 3.0, v1: uint32(3), expect: true},
		{v0: 3.0, v1: uint64(3), expect: true},
		{v0: 3.0, v1: int64(3), expect: true},
		{v0: 3.0, v1: int32(3), expect: true},
		{v0: 3.0, v1: int16(3), expect: true},
		{v0: 3.0, v1: int8(3), expect: true},
		{v0: 3.0, v1: gen.Float(3.0), expect: true},
		{v0: 3.0, v1: gen.Int(3), expect: true},
		{v0: 3.0, v1: true, expect: false},
	} {
		diffs := alt.Diff(p.v0, p.v1)
		tt.Equal(t, p.expect, len(diffs) == 0, "Diff(", p.v0, p.v1, ")")
	}
}

func TestCompare(t *testing.T) {
	dif := alt.Compare([]interface{}{1, 2}, []interface{}{1, 2, 3})
	tt.Equal(t, 1, len(dif))
	tt.Equal(t, 2, dif[0])

	dif = alt.Compare([]interface{}{1, 2}, []interface{}{1, 2})
	tt.Equal(t, 0, len(dif))
}

func TestDiffTime(t *testing.T) {
	t0 := time.Date(2021, time.March, 3, 16, 34, 04, 0, time.UTC)
	t1 := time.Date(2021, time.March, 3, 16, 34, 04, 499, time.UTC)

	alt.TimeTolerance = time.Microsecond
	dif := alt.Compare(t0, t1)
	tt.Equal(t, 0, len(dif))

	alt.TimeTolerance = time.Nanosecond
	dif = alt.Compare(t0, t1)
	tt.Equal(t, 1, len(dif))
}

// TBD slices and maps
// ignores

// TBD support simplifier
