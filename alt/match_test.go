// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

func TestMatchInt(t *testing.T) {
	tt.Equal(t, true, alt.Match(map[string]any{"x": 1}, map[string]any{"x": 1, "y": 2}))
	tt.Equal(t, true, alt.Match(map[string]any{"x": nil}, map[string]any{"x": nil, "y": 2}))
	tt.Equal(t, true, alt.Match(map[string]any{"x": nil}, map[string]any{"y": 2}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": nil}, map[string]any{"x": 1, "y": 2}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": 1, "z": 3}, map[string]any{"x": 1, "y": 2}))
}

func TestMatchBool(t *testing.T) {
	tt.Equal(t, true, alt.Match(map[string]any{"x": true}, map[string]any{"x": true, "y": 2}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": true}, map[string]any{"x": false, "y": 2}))
}

func TestMatchFloat(t *testing.T) {
	tt.Equal(t, true, alt.Match(map[string]any{"x": 1.5}, map[string]any{"x": 1.5, "y": 2}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": 1.5}, map[string]any{"x": 2.5, "y": 2}))
}

func TestMatchString(t *testing.T) {
	tt.Equal(t, true, alt.Match(map[string]any{"x": "a"}, map[string]any{"x": "a"}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": "a"}, map[string]any{"x": "b"}))
}

func TestMatchTime(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	tt.Equal(t, true, alt.Match(map[string]any{"x": tm}, map[string]any{"x": tm}))
	tt.Equal(t, false, alt.Match(map[string]any{"x": tm}, map[string]any{"x": "b"}))
}

func TestMatchSlice(t *testing.T) {
	tt.Equal(t, true,
		alt.Match(map[string]any{"x": []any{1}}, map[string]any{"x": []any{1}}))
	tt.Equal(t, false,
		alt.Match(map[string]any{"x": []any{1}}, map[string]any{"x": []any{2}}))
	tt.Equal(t, false,
		alt.Match(map[string]any{"x": []any{1}}, map[string]any{"x": []any{1, 2}}))
}

func TestMatchMap(t *testing.T) {
	tt.Equal(t, false, alt.Match(map[string]any{"x": []any{1}}, 7))
}

func TestMatchStruct(t *testing.T) {
	type Sample struct {
		Int int
	}
	type Dample struct {
		Int int
	}
	tt.Equal(t, true, alt.Match(&Sample{Int: 3}, &Sample{Int: 3}))
	tt.Equal(t, false, alt.Match(&Sample{Int: 3}, &Dample{Int: 3}))
}

func TestMatchSimplify(t *testing.T) {
	tt.Equal(t, true, alt.Match(&silly{val: 3}, &silly{val: 3}))
}
