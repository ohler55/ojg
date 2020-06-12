// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/tt"
)

type genny struct {
	val int
}

func (g *genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
}

func TestGenerifyBase(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	a := []interface{}{
		true,
		gen.False,
		int8(-8),
		int16(-16),
		int32(-32),
		int64(-64),
		uint(0),
		uint8(8),
		uint16(16),
		uint32(32),
		uint64(64),
		gen.Int(77),
		float32(3.5),
		3.7,
		gen.Float(3.7),
		"string",
		gen.String("string"),
		tm,
		gen.Time(tm),
	}
	v := alt.Generify(a)
	tt.Equal(t, gen.Array{
		gen.True,
		gen.False,
		gen.Int(-8),
		gen.Int(-16),
		gen.Int(-32),
		gen.Int(-64),
		gen.Int(0),
		gen.Int(8),
		gen.Int(16),
		gen.Int(32),
		gen.Int(64),
		gen.Int(77),
		gen.Float(3.5),
		gen.Float(3.7),
		gen.Float(3.7),
		gen.String("string"),
		gen.String("string"),
		gen.Time(tm),
		gen.Time(tm),
	}, v)
}

func TestGenerifyStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := alt.Generify(&d, &alt.Options{OmitNil: true, CreateKey: "type"})
	tt.Equal(t, gen.Object{
		"type": gen.String("Dummy"),
		"val":  gen.Int(3),
		"nest": gen.Object{"type": gen.String("Dummy"), "val": gen.Int(2)},
	}, v)
	v = alt.Generify(&d, &alt.Options{CreateKey: "^", FullTypePath: true, OmitNil: true})
	tt.Equal(t, gen.Object{
		"^":   gen.String("github.com/ohler55/ojg/alt_test/Dummy"),
		"val": gen.Int(3),
		"nest": gen.Object{
			"^":   gen.String("github.com/ohler55/ojg/alt_test/Dummy"),
			"val": gen.Int(2),
		},
	}, v)
}

func TestGenerifyComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.Generify(c)
	tt.Equal(t, gen.Object{"type": gen.String("complex"), "real": gen.Float(1.2), "imag": gen.Float(3.4)}, v)
}

func TestGenerifyMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.Generify(m)
	tt.Equal(t, gen.Object{"1": gen.Int(1), "2": gen.Int(4), "3": gen.Int(9)}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.Generify(m2)
	tt.Equal(t, gen.Object{"1": gen.Int(1), "2": gen.Int(4), "3": gen.Int(9)}, v)
}

func TestGenerifyArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.Generify(a, &alt.Options{OmitNil: true, CreateKey: "type"})
	tt.Equal(t, gen.Array{
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(1)},
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(2)},
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(3)},
	}, v)
}

func TestGenerifyOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := alt.Generify(odd)
	tt.Equal(t, gen.Array{nil, nil}, v)
}

func TestGenerifySimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.Generify(&s)
	tt.Equal(t, gen.Object{"type": gen.String("silly"), "val": gen.Int(3)}, v)
}

func TestGenerifyGeneric(t *testing.T) {
	g := genny{val: 3}
	v := alt.Generify(&g)
	tt.Equal(t, gen.Object{"type": gen.String("genny"), "val": gen.Int(3)}, v)
}

func TestGenerifyNode(t *testing.T) {
	n := gen.Array{gen.Int(1)}
	v := alt.Generify(n)
	tt.Equal(t, gen.Array{gen.Int(1)}, v)
}

func TestGenAlterBase(t *testing.T) {
	tm := time.Date(2020, time.April, 12, 16, 34, 04, 123456789, time.UTC)
	a := []interface{}{
		true,
		gen.False,
		int8(-8),
		int16(-16),
		int32(-32),
		int64(-64),
		uint(0),
		uint8(8),
		uint16(16),
		uint32(32),
		uint64(64),
		gen.Int(77),
		float32(3.5),
		3.7,
		gen.Float(3.7),
		"string",
		gen.String("string"),
		tm,
		gen.Time(tm),
	}
	v := alt.GenAlter(a)
	tt.Equal(t, gen.Array{
		gen.True,
		gen.False,
		gen.Int(-8),
		gen.Int(-16),
		gen.Int(-32),
		gen.Int(-64),
		gen.Int(0),
		gen.Int(8),
		gen.Int(16),
		gen.Int(32),
		gen.Int(64),
		gen.Int(77),
		gen.Float(3.5),
		gen.Float(3.7),
		gen.Float(3.7),
		gen.String("string"),
		gen.String("string"),
		gen.Time(tm),
		gen.Time(tm),
	}, v)
}

func TestGenAlterStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := alt.GenAlter(&d, &alt.Options{OmitNil: true, CreateKey: "type"})
	tt.Equal(t, gen.Object{
		"type": gen.String("Dummy"),
		"val":  gen.Int(3),
		"nest": gen.Object{"type": gen.String("Dummy"), "val": gen.Int(2)}}, v)

	v = alt.GenAlter(&d, &alt.Options{CreateKey: "^", FullTypePath: true, OmitNil: true})
	tt.Equal(t, gen.Object{
		"^":   gen.String("github.com/ohler55/ojg/alt_test/Dummy"),
		"val": gen.Int(3),
		"nest": gen.Object{
			"^":   gen.String("github.com/ohler55/ojg/alt_test/Dummy"),
			"val": gen.Int(2),
		},
	}, v)
}

func TestGenAlterComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.GenAlter(c)
	tt.Equal(t, gen.Object{"type": gen.String("complex"), "real": gen.Float(1.2), "imag": gen.Float(3.4)}, v)
}

func TestGenAlterMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.GenAlter(m)
	tt.Equal(t, gen.Object{"1": gen.Int(1), "2": gen.Int(4), "3": gen.Int(9)}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.GenAlter(m2)
	tt.Equal(t, gen.Object{"1": gen.Int(1), "2": gen.Int(4), "3": gen.Int(9)}, v)
}

func TestGenAlterArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.GenAlter(a)
	tt.Equal(t, gen.Array{
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(1)},
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(2)},
		gen.Object{"type": gen.String("Dummy"), "val": gen.Int(3)},
	}, v)
}

func TestGenAlterOdd(t *testing.T) {
	odd := []interface{}{func() {}, nil}
	v := alt.GenAlter(odd)
	tt.Equal(t, gen.Array{nil, nil}, v)
}

func TestGenAlterSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.GenAlter(&s)
	tt.Equal(t, gen.Object{"type": gen.String("silly"), "val": gen.Int(3)}, v)
}

func TestGenAlterGeneric(t *testing.T) {
	g := genny{val: 3}
	v := alt.GenAlter(&g)
	tt.Equal(t, gen.Object{"type": gen.String("genny"), "val": gen.Int(3)}, v)
}

func TestGenAlterNode(t *testing.T) {
	n := gen.Array{gen.Int(1)}
	v := alt.GenAlter(n)
	tt.Equal(t, gen.Array{gen.Int(1)}, v)
}
