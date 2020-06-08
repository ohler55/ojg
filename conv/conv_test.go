// Copyright (c) 2020, Peter Ohler, All rights reserved.

package conv_test

import (
	"testing"

	"github.com/ohler55/ojg/conv"
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

func TestDecomposeReflectNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := conv.Decompose(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestDecomposeReflectStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := conv.Decompose(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = conv.Decompose(&d, &conv.ConvOptions{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/conv_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/conv_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestDecomposeReflectComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := conv.Decompose(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestDecomposeReflectMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := conv.Decompose(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = conv.Decompose(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestDecomposeReflectArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := conv.Decompose(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestDecomposeReflectOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := conv.Decompose(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestDecomposeReflectSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := conv.Decompose(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}

func TestAlterReflectNumbers(t *testing.T) {
	a := []interface{}{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := conv.Alter(a)
	tt.Equal(t, []interface{}{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestAlterReflectStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := conv.Alter(&d)
	tt.Equal(t, map[string]interface{}{"type": "Dummy", "val": 3, "nest": map[string]interface{}{"type": "Dummy", "val": 2}}, v)

	v = conv.Alter(&d, &conv.ConvOptions{CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]interface{}{
		"^":   "github.com/ohler55/ojg/conv_test/Dummy",
		"val": 3,
		"nest": map[string]interface{}{
			"^":   "github.com/ohler55/ojg/conv_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestAlterReflectComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := conv.Alter(c)
	tt.Equal(t, map[string]interface{}{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestAlterReflectMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := conv.Alter(m)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = conv.Alter(m2)
	tt.Equal(t, map[string]interface{}{"1": 1, "2": 4, "3": 9}, v)
}

func TestAlterReflectArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := conv.Alter(a)
	tt.Equal(t, []interface{}{
		map[string]interface{}{"type": "Dummy", "val": 1},
		map[string]interface{}{"type": "Dummy", "val": 2},
		map[string]interface{}{"type": "Dummy", "val": 3},
	}, v)
}

func TestAlterReflectOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := conv.Alter(odd)
	tt.Equal(t, []interface{}{nil, nil}, v)
}

func TestAlterReflectSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := conv.Alter(&s)
	tt.Equal(t, map[string]interface{}{"type": "silly", "val": 3}, v)
}
