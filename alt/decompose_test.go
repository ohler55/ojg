// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

type Dummy struct {
	Val    int
	Nest   any
	hidden int
}

type Bummy struct {
	Dummy
	Num int
}

type Rummy struct {
	Bummy
}

type Anno struct {
	Val   int    `json:"v,omitempty"`
	Str   int    `json:"str,omitempty,string"`
	Title int    `json:",omitempty"`
	Skip  int    `json:"-"`
	Dash  int    `json:"-,omitempty"`
	Buf   []byte `json:"buf,omitempty"`
}

type Pointy struct {
	X   *int   `json:"x,omitempty"`
	Buf []byte `json:"buf,omitempty"`
}

type silly struct {
	val int
}

func (s *silly) Simplify() any {
	return map[string]any{"type": "silly", "val": s.val}
}

func TestDecomposeNumbers(t *testing.T) {
	a := []any{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := alt.Dup(a)
	tt.Equal(t, []any{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestDecomposeStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}, hidden: 3}
	v := alt.Decompose(&d)
	tt.Equal(t, map[string]any{"type": "Dummy", "val": 3, "nest": map[string]any{"type": "Dummy", "val": 2}}, v)

	v = alt.Decompose(&d, &alt.Options{CreateKey: "^", FullTypePath: true, OmitNil: true})
	tt.Equal(t, map[string]any{
		"^":   "github.com/ohler55/ojg/alt_test/Dummy",
		"val": 3,
		"nest": map[string]any{
			"^":   "github.com/ohler55/ojg/alt_test/Dummy",
			"val": 2,
		},
	}, v)

	a := Anno{Val: 3}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{"v": 3}, v)

	a = Anno{Val: 0, Str: 1}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{"str": "1"}, v)

	a = Anno{Val: 0, Str: 0}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{}, v)

	a = Anno{Dash: 1, Skip: 2, Title: 3}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{"-": 1, "Title": 3}, v)

	a = Anno{Buf: []byte{}, Val: 3}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{"v": 3}, v)

	a = Anno{Buf: []byte{'a', 'b'}, Val: 3}
	v = alt.Decompose(&a, &alt.Options{UseTags: true})
	tt.Equal(t, map[string]any{"v": 3, "buf": "ab"}, v)

	v = alt.Decompose(&a, &alt.Options{UseTags: true, BytesAs: ojg.BytesAsArray})
	tt.Equal(t, map[string]any{"v": 3, "buf": []any{97, 98}}, v)

	v = alt.Decompose(&a, &alt.Options{UseTags: true, BytesAs: ojg.BytesAsBase64})
	tt.Equal(t, map[string]any{"v": 3, "buf": "YWI="}, v)
}

func TestDecomposeStructWithPointers(t *testing.T) {
	x := 3
	p := Pointy{X: &x, Buf: []byte("byte me")}
	v := alt.Decompose(&p)
	tt.Equal(t, map[string]any{"type": "Pointy", "x": 3, "buf": "byte me"}, v)

	v = alt.Decompose(&p, &alt.Options{CreateKey: "", UseTags: true, BytesAs: ojg.BytesAsBase64})
	tt.Equal(t, map[string]any{"x": 3, "buf": "Ynl0ZSBtZQ=="}, v)

	v = alt.Decompose(&p, &alt.Options{CreateKey: "", UseTags: true, BytesAs: ojg.BytesAsArray})
	tt.Equal(t, map[string]any{"x": 3, "buf": []any{98, 121, 116, 101, 32, 109, 101}}, v)

	// TBD try empty for each
}

func TestDecomposeEmbeddedStruct(t *testing.T) {
	b := Bummy{Dummy: Dummy{Val: 3}, Num: 5}
	v := alt.Decompose(&b, &alt.Options{NestEmbed: false, OmitNil: true})
	tt.Equal(t, map[string]any{"val": 3, "num": 5}, v)

	v = alt.Decompose(&b, &alt.Options{NestEmbed: true, OmitNil: true})
	tt.Equal(t, map[string]any{"dummy": map[string]any{"val": 3}, "num": 5}, v)

	v = alt.Decompose(&b, &alt.Options{NestEmbed: true, OmitNil: true, CreateKey: "^"})
	tt.Equal(t, map[string]any{
		"^":     "Bummy",
		"dummy": map[string]any{"^": "Dummy", "val": 3}, "num": 5}, v)

	v = alt.Decompose(&b, &alt.Options{NestEmbed: true, OmitNil: true, CreateKey: "^", FullTypePath: true})
	tt.Equal(t, map[string]any{
		"^":     "github.com/ohler55/ojg/alt_test/Bummy",
		"dummy": map[string]any{"^": "github.com/ohler55/ojg/alt_test/Dummy", "val": 3}, "num": 5}, v)

	r := Rummy{Bummy: Bummy{Dummy: Dummy{Val: 3}, Num: 5}}
	v = alt.Decompose(&r, &alt.Options{NestEmbed: true, OmitNil: true, CreateKey: "^"})
	tt.Equal(t, map[string]any{
		"^": "Rummy",
		"bummy": map[string]any{
			"^":     "Bummy",
			"dummy": map[string]any{"^": "Dummy", "val": 3}, "num": 5}}, v)
	d := Dummy{Nest: &silly{val: 3}}
	v = alt.Decompose(&d, &alt.Options{OmitNil: true})
	tt.Equal(t, map[string]any{"nest": map[string]any{"type": "silly", "val": 3}, "val": 0}, v)
}

func TestDecomposeNestedPtr(t *testing.T) {
	type Inner struct {
		X int
	}
	type Wrap struct {
		*Inner
	}
	obj := &Wrap{Inner: &Inner{X: 3}}
	v := alt.Decompose(obj, &alt.Options{CreateKey: ""})
	tt.Equal(t, map[string]any{"x": 3}, v)
}

func TestDecomposeComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.Decompose(c)
	tt.Equal(t, map[string]any{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestDecomposeMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.Decompose(m)
	tt.Equal(t, map[string]any{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.Decompose(m2)
	tt.Equal(t, map[string]any{"1": 1, "2": 4, "3": 9}, v)
}

func TestDecomposeArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.Decompose(a)
	tt.Equal(t, []any{
		map[string]any{"type": "Dummy", "val": 1},
		map[string]any{"type": "Dummy", "val": 2},
		map[string]any{"type": "Dummy", "val": 3},
	}, v)
}

func TestDecomposeOdd(t *testing.T) {
	odd := []any{func() {}, nil}
	v := alt.Decompose(odd)
	tt.Equal(t, []any{nil, nil}, v)
}

func TestDecomposeSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.Decompose(&s)
	tt.Equal(t, map[string]any{"type": "silly", "val": 3}, v)

	v = alt.Alter([]any{[]byte("abc")}, &alt.Options{UseTags: true})
	tt.Equal(t, []any{"abc"}, v)

	v = alt.Alter([]any{[]byte("abc")}, &alt.Options{UseTags: true, BytesAs: ojg.BytesAsArray})
	tt.Equal(t, []any{[]any{97, 98, 99}}, v)

	v = alt.Alter([]any{[]byte("abc")}, &alt.Options{UseTags: true, BytesAs: ojg.BytesAsBase64})
	tt.Equal(t, []any{"YWJj"}, v)
}

func TestDecomposeTime(t *testing.T) {
	tm := time.Date(2021, time.February, 9, 12, 13, 14, 0, time.UTC)
	a := []any{tm}
	v := alt.Decompose(a, &ojg.Options{TimeFormat: "nano"})
	tt.Equal(t, []any{1612872794000000000}, v)
}

func TestDecomposeAliasedTypes(t *testing.T) {
	type Stringy string
	type Inty int
	type Uinty uint
	type Floaty float64
	type Booly bool
	tcs := [][]any{
		{Stringy("stringy"), "stringy"},
		{Inty(1), 1},
		{Uinty(1), uint(1)},
		{Floaty(1.3), 1.3},
		{Booly(true), true},
	}
	for i, tc := range tcs {
		input := tc[0]
		expect := tc[1]
		dec := alt.Decompose(input)
		tt.Equal(t, expect, dec, fmt.Sprintf("case %d failed", i))
	}
}

func TestAlterNumbers(t *testing.T) {
	a := []any{int8(-8), int16(-16), int32(-32), uint(0), uint8(8), uint16(16), uint32(32), uint64(64), float32(3.2)}
	v := alt.Alter(a)
	tt.Equal(t, []any{-8, -16, -32, 0, 8, 16, 32, 64, 3.2}, v)
}

func TestAlterStruct(t *testing.T) {
	d := Dummy{Val: 3, Nest: &Dummy{Val: 2}}
	v := alt.Alter(&d)
	tt.Equal(t, map[string]any{"type": "Dummy", "val": 3, "nest": map[string]any{"type": "Dummy", "val": 2}}, v)

	v = alt.Alter(&d, &alt.Options{CreateKey: "^", FullTypePath: true, OmitNil: true})
	tt.Equal(t, map[string]any{
		"^":   "github.com/ohler55/ojg/alt_test/Dummy",
		"val": 3,
		"nest": map[string]any{
			"^":   "github.com/ohler55/ojg/alt_test/Dummy",
			"val": 2,
		},
	}, v)
}

func TestAlterComplex(t *testing.T) {
	c := complex(1.2, 3.4)
	v := alt.Alter(c)
	tt.Equal(t, map[string]any{"type": "complex", "real": 1.2, "imag": 3.4}, v)
}

func TestAlterMap(t *testing.T) {
	m := map[int]int{1: 1, 2: 4, 3: 9}
	v := alt.Alter(m)
	tt.Equal(t, map[string]any{"1": 1, "2": 4, "3": 9}, v)

	m2 := map[string]int{"1": 1, "2": 4, "3": 9}
	v = alt.Alter(m2)
	tt.Equal(t, map[string]any{"1": 1, "2": 4, "3": 9}, v)
}

func TestAlterArray(t *testing.T) {
	a := []*Dummy{{Val: 1}, {Val: 2}, {Val: 3}}
	v := alt.Alter(a)
	tt.Equal(t, []any{
		map[string]any{"type": "Dummy", "val": 1},
		map[string]any{"type": "Dummy", "val": 2},
		map[string]any{"type": "Dummy", "val": 3},
	}, v)
}

func TestAlterOdd(t *testing.T) {
	odd := []any{func() {}, nil}
	v := alt.Alter(odd)
	tt.Equal(t, []any{nil, nil}, v)
}

func TestAlterSimplifier(t *testing.T) {
	s := silly{val: 3}
	v := alt.Alter(&s)
	tt.Equal(t, map[string]any{"type": "silly", "val": 3}, v)
}

func TestDecomposeConverter(t *testing.T) {
	c := ojg.Converter{
		Int: []func(val int64) (any, bool){
			func(val int64) (any, bool) { return val + 1, true },
		},
	}
	val := []any{1, true}
	v := alt.Decompose(val, &alt.Options{Converter: &c})
	tt.Equal(t, []any{2, true}, v)
}

func TestAlterConverter(t *testing.T) {
	c := alt.Converter{
		Int: []func(val int64) (any, bool){
			func(val int64) (any, bool) { return val + 1, true },
		},
	}
	val := []any{1, true}
	v := alt.Alter(val, &alt.Options{Converter: &c})
	tt.Equal(t, []any{2, true}, v)
}

func BenchmarkDecompose(b *testing.B) {
	a := Anno{
		Val:   1,
		Str:   2,
		Title: 3,
		Skip:  4,
		Dash:  5,
		Buf:   []byte("abcd"),
	}
	_ = alt.Decompose(&a, &alt.Options{UseTags: true})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = alt.Decompose(&a, &alt.Options{UseTags: true})
	}
}
