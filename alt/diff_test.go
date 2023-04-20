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
		v0     any
		v1     any
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

		{v0: "abc", v1: "abc", expect: true},
		{v0: "abc", v1: "abx", expect: false},
	} {
		diffs := alt.Diff(p.v0, p.v1)
		tt.Equal(t, p.expect, len(diffs) == 0, "Diff(", p.v0, p.v1, ")")
	}
}

func TestCompare(t *testing.T) {
	dif := alt.Compare([]any{1, 2}, []any{1, 2, 3})
	tt.Equal(t, 1, len(dif))
	tt.Equal(t, 2, dif[0])

	dif = alt.Compare([]any{1, 2}, []any{1, 2})
	tt.Equal(t, 0, len(dif))
}

func TestCompareTime(t *testing.T) {
	t0 := time.Date(2021, time.March, 3, 16, 34, 04, 0, time.UTC)
	t1 := time.Date(2021, time.March, 3, 16, 34, 04, 499, time.UTC)

	alt.TimeTolerance = time.Microsecond
	dif := alt.Compare(t0, t1)
	tt.Equal(t, 0, len(dif))

	alt.TimeTolerance = time.Nanosecond
	dif = alt.Compare(t0, t1)
	tt.Equal(t, 1, len(dif))
}

func TestDiffSlice(t *testing.T) {
	diffs := alt.Diff(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{4, 4}},
	)
	tt.Equal(t, 1, len(diffs))
	tt.Equal(t, alt.Path{2, 0}, diffs[0])

	diffs = alt.Diff([]any{1, 2}, 5)
	tt.Equal(t, 1, len(diffs))
	tt.Equal(t, alt.Path{nil}, diffs[0])

	diffs = alt.Diff(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{3, 4, 5}},
		alt.Path{2, 2},
	)
	tt.Equal(t, 0, len(diffs))

	dif := alt.Compare(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{3, 5}},
	)
	tt.Equal(t, alt.Path{2, 1}, dif)

	diffs = alt.Diff(
		[]any{1, 2, []any{3, 4, 5}},
		[]any{1, 2, []any{3, 4}},
	)
	tt.Equal(t, 1, len(diffs))
	tt.Equal(t, alt.Path{2, 2}, diffs[0])
}

func TestDiffSliceIgnores(t *testing.T) {
	diffs := alt.Diff(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{3, 4, 5}},
		alt.Path{2, 2},
	)
	tt.Equal(t, 0, len(diffs))

	diffs = alt.Diff(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{3, 5}},
		alt.Path{2, 1},
	)
	tt.Equal(t, 0, len(diffs))

	diffs = alt.Diff(
		[]any{1, 2, []any{3, 4}},
		[]any{1, 2, []any{3, 5}},
		alt.Path{2, nil},
	)
	tt.Equal(t, 0, len(diffs))
}

func TestDiffMap(t *testing.T) {
	diffs := alt.Diff(
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 3, "b": 4}},
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 4, "b": 4}},
	)
	tt.Equal(t, 1, len(diffs))
	tt.Equal(t, alt.Path{"z", "a"}, diffs[0])

	dif := alt.Compare(
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 3, "b": 4}},
		map[string]any{"x": 1, "y": 2, "z": true},
	)
	tt.Equal(t, alt.Path{"z"}, dif)
}

func TestDiffMapIgnores(t *testing.T) {
	diffs := alt.Diff(
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 3, "b": 4}},
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 4, "b": 4}},
		alt.Path{"z"},
	)
	tt.Equal(t, 0, len(diffs))

	diffs = alt.Diff(
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 3, "b": 4}},
		map[string]any{"x": 1, "y": 2, "z": map[string]any{"a": 4, "b": 4}},
		alt.Path{"z", nil},
	)
	tt.Equal(t, 0, len(diffs))
}

func TestDiffSimplifier(t *testing.T) {
	diffs := alt.Diff(
		&silly{val: 3},
		&silly{val: 3},
	)
	tt.Equal(t, 0, len(diffs))
}

func TestDiffTypes(t *testing.T) {
	diffs := alt.Diff(
		&silly{val: 3},
		&Dummy{Val: 3},
	)
	tt.Equal(t, 1, len(diffs))
}

func TestDiffReflect(t *testing.T) {
	diffs := alt.Diff(
		&Dummy{Val: 3},
		&Dummy{Val: 3},
	)
	tt.Equal(t, 0, len(diffs))
}

func TestDiffArrayIgnores(t *testing.T) {
	v1 := map[string]any{"x": []any{map[string]any{"a": 3}, map[string]any{"a": 3}}}
	v2 := map[string]any{"x": []any{map[string]any{"a": 3}, map[string]any{"a": 4}}}

	diffs := alt.Diff(v1, v2, alt.Path{"x", 1, "a"})
	tt.Equal(t, 0, len(diffs))

	diffs = alt.Diff(v1, v2, alt.Path{"x", 0, "a"})
	tt.Equal(t, 1, len(diffs))

	diffs = alt.Diff(v1, v2, alt.Path{"x", nil, "a"})
	tt.Equal(t, 0, len(diffs))
}
