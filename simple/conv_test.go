// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple_test

import (
	"testing"

	"github.com/ohler55/ojg/simple"
	"github.com/ohler55/ojg/tt"
)

type Dummy struct {
	Val  int
	Nest interface{}
}

type silly struct {
	val int
}

func (s *silly) Simplify() interface{} {
	return map[string]interface{}{"type": "silly", "val": s.val}
}

func TestSimpleFromReflectNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := simple.From(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestSimpleFromReflectStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := simple.From(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = simple.From(&d, &simple.ConvOptions{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/simple_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/simple_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestSimpleFromReflectComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := simple.From(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestSimpleFromReflectMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := simple.From(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = simple.From(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestSimpleFromReflectArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := simple.From(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestSimpleFromReflectOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := simple.From(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestSimpleFromReflectSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := simple.From(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}

func TestSimpleAlterReflectNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := simple.Alter(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestSimpleAlterReflectStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := simple.Alter(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = simple.Alter(&d, &simple.ConvOptions{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/simple_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/simple_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestSimpleAlterReflectComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := simple.Alter(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestSimpleAlterReflectMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := simple.Alter(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = simple.Alter(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestSimpleAlterReflectArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := simple.Alter(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestSimpleAlterReflectOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := simple.Alter(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestSimpleAlterReflectSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := simple.Alter(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}
