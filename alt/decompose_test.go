// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

type Dummy struct {
	Val    int
	Nest   interface{}
	hidden int
}

type Anno struct {
	Val int `json:"v"`
}

type silly struct {
	val int
}

func (s *silly) Simplify() interface{} {
	return map[string]interface{}{"type": "silly", "val": s.val}
}

func TestDecomposeNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := alt.Dup(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestDecomposeStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}, hidden: 3}
	v := alt.Decompose(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = alt.Decompose(&d, &alt.Options{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/alt_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/alt_test/Dummy",
			"val": 2,
		},
	}, v)

	a := Anno{Val: 3}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]interface{}{"v": 3}, v)
}

func TestDecomposeComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.Decompose(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestDecomposeMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.Decompose(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.Decompose(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestDecomposeArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.Decompose(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestDecomposeOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := alt.Decompose(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestDecomposeSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.Decompose(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}

func TestAlterNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := alt.Alter(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestAlterStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := alt.Alter(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = alt.Alter(&d, &alt.Options{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/alt_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/alt_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestAlterComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.Alter(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestAlterMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.Alter(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.Alter(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestAlterArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.Alter(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestAlterOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := alt.Alter(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestAlterSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.Alter(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}

func TestDecomposeConverter(t *testing.T) {
	c := alt.Converter{
		Int: []func(val int64) (interface{}, bool){
			func(val int64) (interface{}, bool) { return val + 1, true },
		},
	}
	val := []interface{}{1, true}
	v := alt.Decompose(val, &alt.Options{Converter: &c})
	tt.Equal(t, []interface{}{2, true}, v)
}

func TestAlterConverter(t *testing.T) {
	c := alt.Converter{
		Int: []func(val int64) (interface{}, bool){
			func(val int64) (interface{}, bool) { return val + 1, true },
		},
	}
	val := []interface{}{1, true}
	v := alt.Alter(val, &alt.Options{Converter: &c})
	tt.Equal(t, []interface{}{2, true}, v)
}
