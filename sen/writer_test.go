// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type wdata struct {
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect  string
	value   interface{}
	options *sen.Options
	indent  int
}

// Used to test Simplifier objects in simple data.
type simon struct {
	x int
}

func (s *simon) Simplify() interface{} {
	return map[string]interface{}{
		"type": "simon",
		"x":    s.x,
	}
}

type genny struct {
	val int
}

func (g *genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
}

type Dummy struct {
	Val int
}

func (d *Dummy) String() string {
	return fmt.Sprintf("{val: %d}", d.Val)
}

type Panik struct {
}

func (p *Panik) Simplify() interface{} {
	panic("force panic")
}

type shortWriter struct {
	max int
}

func (w *shortWriter) Write(p []byte) (n int, err error) {
	w.max -= len(p)
	if w.max < 0 {
		return 0, fmt.Errorf("fail now")
	}
	return len(p), nil
}

func TestString(t *testing.T) {
	opt := &sen.Options{}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	tm2 := time.Unix(-10, -100000000)
	for i, d := range []wdata{
		{value: nil, expect: "null"},
		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: "string", expect: "string"},
		{value: "str#ing", expect: `"str#ing"`},
		{value: "\\\tx\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122",
			expect: `"\\	x
\r\b\f\"&<>\u2028\u2029\u0007𝄢"`, options: &sen.Options{HTMLUnsafe: true}},
		{value: "<x>", expect: `<x>`, options: &sen.Options{HTMLUnsafe: true}},
		{value: "<x>", expect: `"\u003cx\u003e"`, options: &sen.Options{HTMLUnsafe: false}},
		{value: "ぴーたー", expect: "ぴーたー", options: &sen.Options{}},
		{value: gen.String("string"), expect: "string"},
		{value: []interface{}{true, nil}, expect: "[true null]"},
		{value: gen.Array{gen.Bool(true), nil}, expect: "[true null]"},
		{value: []interface{}{true, false}, indent: 2, expect: "[\n  true\n  false\n]"},
		{value: []interface{}{}, indent: 2, expect: "[]"},
		{value: []interface{}{[]interface{}{}, []interface{}{}}, expect: "[[][]]"},
		{value: []interface{}{map[string]interface{}{}, map[string]interface{}{}}, expect: "[{}{}]"},
		{value: map[string]interface{}{}, indent: 0, expect: "{}", options: &sen.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, indent: 2, expect: "[\n  true\n  false\n]"},
		{value: map[string]interface{}{"t": true, "f": false}, expect: `{f:false t:true}`, options: &sen.Options{Sort: true}},
		{value: gen.Object{"t": gen.True, "f": gen.False}, expect: `{f:false t:true}`, options: &sen.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "[true false]", options: opt},
		{value: gen.Array{gen.False, gen.True}, expect: "[false true]", options: opt},
		{value: []interface{}{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "[-1 2 -3 4 -5]"},
		{value: []interface{}{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "[1 65 2 3 4 5]"},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "[1 1.2]"},
		{value: []interface{}{float32(1.2), float64(2.1)}, expect: "[1.2 2.1]"},
		{value: []interface{}{tm}, expect: "[1588879759123456789]"},
		{value: []interface{}{tm}, expect: `[{^:Time value:"2020-05-07T19:29:19.123456789Z"}]`,
			options: &sen.Options{TimeMap: true, CreateKey: "^", TimeFormat: time.RFC3339Nano}},
		{value: []interface{}{tm}, expect: `[{^:"time/Time" value:"2020-05-07T19:29:19.123456789Z"}]`,
			options: &sen.Options{TimeMap: true, CreateKey: "^", TimeFormat: time.RFC3339Nano, FullTypePath: true}},
		{value: tm2, expect: "-10.100000000", options: &sen.Options{TimeFormat: "second"}},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759123456789]"},
		{value: gen.Array{gen.Time(tm)}, expect: `["2020-05-07T19:29:19.123456789Z"]`, options: &sen.Options{TimeFormat: time.RFC3339Nano}},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759.123456789]", options: &sen.Options{TimeFormat: "second"}},
		{value: gen.Array{gen.Time(tm)}, expect: `[{@:1588879759123456789}]`, options: &sen.Options{TimeWrap: "@"}},
		{value: map[string]interface{}{"t": true, "x": nil}, expect: "{t:true}", options: &sen.Options{OmitNil: true}},
		{value: map[string]interface{}{"t": true, "f": false}, expect: "{\n  f: false\n  t: true\n}", options: &sen.Options{Sort: true, Indent: 2}},

		{value: map[string]interface{}{"t": true}, expect: "{\n  t: true\n}", options: &sen.Options{Indent: 2}},
		{value: map[string]interface{}{"t": true, "n": nil, "f": false}, expect: "{f:false t:true}",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: map[string]interface{}{"t": true, "n": nil, "f": false}, expect: "{\n  f: false\n  t: true\n}",
			options: &sen.Options{OmitNil: true, Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true, "n": nil, "f": false}, expect: "{\n  f: false\n  n: null\n  t: true\n}",
			options: &sen.Options{OmitNil: false, Sort: true, Indent: 2}},
		{value: map[string]interface{}{"t": true, "n": nil, "f": false}, expect: "{f:false t:true}",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: map[string]interface{}{"t": true, "n": nil, "f": false}, expect: "{f:false n:null t:true}",
			options: &sen.Options{OmitNil: false, Sort: true}},
		{value: map[string]interface{}{"n": nil}, expect: "{n:null}"},
		{value: map[string]interface{}{"n": nil}, expect: "{\n}", options: &sen.Options{OmitNil: true, Sort: false, Indent: 2}},

		{value: gen.Object{"t": gen.True, "x": nil}, expect: "{t:true}", options: &sen.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  t: true\n}", options: &sen.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  t: true\n}", options: &sen.Options{Indent: 2, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{f:false t:true}",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\n  f: false\n  t: true\n}",
			options: &sen.Options{OmitNil: true, Sort: true, Indent: 2}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\n  f: false\n  n: null\n  t: true\n}",
			options: &sen.Options{OmitNil: false, Sort: true, Indent: 2}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{f:false t:true}",
			options: &sen.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{f:false n:null t:true}",
			options: &sen.Options{OmitNil: false, Sort: true}},
		{value: gen.Object{"n": nil}, expect: "{n:null}"},
		{value: gen.Object{"n": nil}, expect: "{\n}", options: &sen.Options{OmitNil: true, Sort: false, Indent: 2}},

		{value: &simon{x: 3}, expect: `{type:simon x:3}`, options: &sen.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `{type:genny val:3}`, options: &sen.Options{Sort: true}},
		{value: &genny{val: 3}, expect: "{\n  type: genny\n  val: 3\n}", options: &sen.Options{Sort: true, Indent: 2}},
		{value: &Dummy{Val: 3}, expect: `{val:3}`, options: &sen.Options{}},
		{value: &Dummy{Val: 3}, expect: "{\n  val: 3\n}", options: &sen.Options{Indent: 2}},
		{value: &Dummy{Val: 3}, expect: `{^:Dummy val:3}`, options: &sen.Options{Sort: true, CreateKey: "^"}},
		{value: []*Dummy{{Val: 3}}, expect: "[\n  {\n    val: 3\n  }\n]", options: &sen.Options{Indent: 2}},
		{value: []*Dummy{{Val: 3}}, expect: "[{val:3}]", options: &sen.Options{Indent: 0}},
		{value: map[string]*Dummy{"d": {Val: 3}}, expect: "{\n  d: {\n    val: 3\n  }\n}", options: &sen.Options{Indent: 2}},
		{value: complex(1, 7), expect: "{\n  imag: 7\n  real: 1\n}", options: &sen.Options{Indent: 2, Sort: true}},
		{value: complex(1, 7), expect: `"(1+7i)"`, options: &sen.Options{Indent: 2, NoReflect: true}},
		{value: complex(1, 7), expect: `"(1+7i)"`, options: &sen.Options{Indent: 0, NoReflect: true}},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, oj.JSON(d.value))
		}
		var s string
		if d.options == nil {
			if 0 < d.indent {
				s = sen.String(d.value, d.indent)
			} else {
				s = sen.String(d.value)
			}
		} else {
			s = sen.String(d.value, d.options)
		}
		tt.Equal(t, d.expect, s, fmt.Sprintf("%d: %v", i, d.value))
	}
}

func TestWriteBasic(t *testing.T) {
	var b strings.Builder

	err := sen.Write(&b, []interface{}{true, false})
	tt.Nil(t, err)
	tt.Equal(t, "[true false]", b.String())

	opt := sen.Options{WriteLimit: 8}
	b.Reset()
	err = sen.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true false]", b.String())

	// A second time.
	b.Reset()
	err = sen.Write(&b, []interface{}{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true false]", b.String())

	b.Reset()
	err = sen.Write(&b, []interface{}{false, true}, 2)
	tt.Nil(t, err)
	tt.Equal(t, "[\n  false\n  true\n]", b.String())

	b.Reset()
	// Force a realloc of string buffer.
	err = sen.Write(&b, strings.Repeat("Xyz ", 63)+"\U0001D122", 2)
	tt.Nil(t, err)
	tt.Equal(t, 258, len(b.String()))

	// Make sure a comma separator is added in unsorted-unindent mode.
	b.Reset()
	err = sen.Write(&b, map[string]interface{}{"t": true, "f": false})
	tt.Nil(t, err)
	tt.Equal(t, 16, len(b.String()))
	b.Reset()
	err = sen.Write(&b, gen.Object{"t": gen.True, "f": gen.False})
	tt.Nil(t, err)
	tt.Equal(t, 16, len(b.String()))

	b.Reset()
	opt.Sort = true
	err = sen.Write(&b, map[string]interface{}{"t": true, "f": false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 16, len(b.String()))
	b.Reset()
	err = sen.Write(&b, gen.Object{"t": gen.True, "f": gen.False}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 16, len(b.String()))
}

func TestWriteWide(t *testing.T) {
	var b strings.Builder
	opt := ojg.Options{Indent: 300}
	err := sen.Write(&b, []interface{}{[]interface{}{true, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 529, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Array{gen.Array{gen.True, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 529, len(b.String()))

	b.Reset()
	err = sen.Write(&b, map[string]interface{}{"x": map[string]interface{}{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 538, len(b.String()))

	b.Reset()
	err = sen.Write(&b, gen.Object{"x": gen.Object{"y": gen.True, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 538, len(b.String()))

	opt = ojg.Options{Indent: 300, Sort: true}
	b.Reset()
	err = sen.Write(&b, map[string]interface{}{"x": map[string]interface{}{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 538, len(b.String()))
}

func TestWriteDeep(t *testing.T) {
	var b strings.Builder
	opt := sen.Options{Tab: true}
	a := []interface{}{map[string]interface{}{"x": true}}
	for i := 40; 0 < i; i-- {
		a = []interface{}{a}
	}
	err := sen.Write(&b, a, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1795, len(b.String()))

	b.Reset()
	g := gen.Array{gen.Object{"x": gen.True}}
	for i := 40; 0 < i; i-- {
		g = gen.Array{g}
	}
	err = sen.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1795, len(b.String()))

	opt.Sort = true
	b.Reset()
	err = sen.Write(&b, a, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1795, len(b.String()))
}

func TestWriteShort(t *testing.T) {
	opt := sen.Options{Indent: 2, WriteLimit: 2}
	err := sen.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	opt.Indent = 0
	err = sen.Write(&shortWriter{max: 3}, []interface{}{true, nil}, &opt)
	tt.NotNil(t, err)
	err = sen.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	obj := map[string]interface{}{"t": true, "n": nil}
	sobj := gen.Object{"t": gen.True, "n": nil}
	opt.Indent = 0
	for i := 2; i < 15; i += 2 {
		err = sen.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = sen.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)

		opt.Sort = true
		err = sen.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = sen.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)
	}
	opt.Indent = 2
	for i := 2; i < 15; i += 2 {
		err = sen.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = sen.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)

		opt.Sort = false
		err = sen.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = sen.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)
	}
}

func TestWriteBad(t *testing.T) {
	var b strings.Builder
	err := sen.Write(&b, []interface{}{true, &Panik{}})
	tt.NotNil(t, err)
}

func TestStringBad(t *testing.T) {
	out := sen.String([]interface{}{true, &Panik{}})
	tt.Equal(t, 0, len(out))
}

func TestMustWritePanic(t *testing.T) {
	tt.Panic(t, func() { sen.MustWrite(&shortWriter{max: 3}, []interface{}{func() {}}) })
}

func TestBytes(t *testing.T) {
	wr := sen.Writer{Options: ojg.Options{Sort: true}}
	obj := map[string]interface{}{"t": true, "n": nil}
	b := sen.Bytes(obj, &wr)
	tt.Equal(t, "{n:null t:true}", string(b))
}

func TestWriteStructWide(t *testing.T) {
	type Nest struct {
		Dig   *Nest
		Level int
	}
	var b strings.Builder
	opt := sen.Options{Tab: true, OmitNil: true}
	var n *Nest
	for i := 40; 0 < i; i-- {
		n = &Nest{Dig: n, Level: i}
	}
	err := sen.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 2901, len(b.String()))

	b.Reset()
	opt.Tab = false
	opt.Indent = 8
	err = sen.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 12969, len(b.String()))

	b.Reset()
	opt.Indent = 0
	err = sen.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 586, len(b.String()))
}

func TestWriteStructCreateKey(t *testing.T) {
	type Sample struct {
		X int
		Y int
	}
	sample := Sample{X: 1, Y: 2}
	opt := sen.Options{Indent: 2, CreateKey: "^"}
	s := sen.String(&sample, &opt)
	tt.Equal(t, `{
  ^: "Sample"
  x: 1
  y: 2
}`, s)

	opt.FullTypePath = true
	s = sen.String(&sample, &opt)
	tt.Equal(t, `{
  ^: "github.com/ohler55/ojg/sen_test/Sample"
  x: 1
  y: 2
}`, s)

	opt.Indent = 0
	s = sen.String(&sample, &opt)
	tt.Equal(t, `{^:"github.com/ohler55/ojg/sen_test/Sample" x:1 y:2}`, s)

	opt.FullTypePath = false
	s = sen.String(&sample, &opt)
	tt.Equal(t, `{^:Sample x:1 y:2}`, s)
}

func TestWriteStruct(t *testing.T) {
	type Sample struct {
		X []int
		Y map[string]int
		Z *int
	}
	sample := Sample{X: []int{1}, Y: map[string]int{"y": 2}}
	opt := sen.Options{Indent: 2, OmitNil: true}
	s := sen.String(&sample, &opt)
	tt.Equal(t, `{
  x: [
    1
  ]
  y: {
    y: 2
  }
}`, s)
}

func TestWriteStructSkip(t *testing.T) {
	type Skippy struct {
		X int `json:"a,omitempty"`
		Y int `json:"b,omitempty"`
	}
	opt := sen.Options{Indent: 2, OmitNil: true, UseTags: true}
	skippy := Skippy{X: 0, Y: 1}
	s := sen.String(&skippy, &opt)
	tt.Equal(t, `{
  b: 1
}`, s)
	opt.Indent = 0
	s = sen.String(&skippy, &opt)
	tt.Equal(t, `{b:1}`, s)

	opt.Indent = 2
	skippy.X = 1
	skippy.Y = 0
	s = sen.String(&skippy, &opt)
	tt.Equal(t, `{
  a: 1
}`, s)

	opt.Indent = 0
	s = sen.String(&skippy, &opt)
	tt.Equal(t, `{a:1}`, s)
}

func TestWriteSliceWide(t *testing.T) {
	type Nest struct {
		Dig []Nest
	}
	opt := sen.Options{Tab: true}
	n := &Nest{}
	for i := 20; 0 < i; i-- {
		n = &Nest{Dig: []Nest{*n}}
	}
	s := sen.String(n, &opt)
	tt.Equal(t, 1810, len(s))

	opt.Tab = false
	opt.Indent = 4
	s = sen.String(n, &opt)
	tt.Equal(t, 6671, len(s))

	opt.Indent = 0
	s = sen.String(n, &opt)
	tt.Equal(t, 168, len(s))
}

func TestWriteSliceArray(t *testing.T) {
	type Matrix struct {
		Rows [][4]int
	}
	opt := sen.Options{Indent: 2}
	m := &Matrix{Rows: [][4]int{[4]int{1, 2, 3, 4}}}
	s := sen.String(m, &opt)
	tt.Equal(t, `{
  rows: [
    [
      1
      2
      3
      4
    ]
  ]
}`, s)

	opt.Indent = 0
	s = sen.String(m, &opt)
	tt.Equal(t, `{rows:[[1 2 3 4]]}`, s)
}

func TestWriteSliceMap(t *testing.T) {
	type SMS struct {
		Maps []map[string]int
	}
	opt := sen.Options{Indent: 2, Sort: true}
	m := &SMS{Maps: []map[string]int{map[string]int{"x": 1, "y": 2}}}
	s := sen.String(m, &opt)
	tt.Equal(t, `{
  maps: [
    {
      x: 1
      y: 2
    }
  ]
}`, s)

	opt.Indent = 0
	s = sen.String(m, &opt)
	tt.Equal(t, `{maps:[{x:1 y:2}]}`, s)
}

func TestWriteMapWide(t *testing.T) {
	type Nest struct {
		Dig map[string]*Nest
	}
	opt := sen.Options{Tab: true, OmitNil: true}
	n := &Nest{map[string]*Nest{"x": nil}}
	for i := 16; 0 < i; i-- {
		n = &Nest{Dig: map[string]*Nest{"x": n}}
	}
	s := sen.String(n, &opt)
	tt.Equal(t, 1330, len(s))

	opt.Tab = false
	opt.Indent = 4
	s = sen.String(n, &opt)
	tt.Equal(t, 4619, len(s))

	opt.Indent = 0
	s = sen.String(n, &opt)
	tt.Equal(t, 168, len(s))
}

func TestWriteMapSlice(t *testing.T) {
	m := map[string][]int{"x": []int{1, 2, 3}, "y": []int{}}
	opt := sen.Options{Indent: 2, OmitNil: true, Sort: true}
	s := sen.String(m, &opt)
	tt.Equal(t, `{
  x: [
    1
    2
    3
  ]
}`, s)

	opt.OmitNil = false
	s = sen.String(m, &opt)
	tt.Equal(t, `{
  x: [
    1
    2
    3
  ]
  y: []
}`, s)

	opt.OmitNil = false
	opt.Indent = 0
	s = sen.String(m, &opt)
	tt.Equal(t, `{x:[1 2 3] y:[]}`, s)

	opt.OmitNil = true
	s = sen.String(m, &opt)
	tt.Equal(t, `{x:[1 2 3]}`, s)
}

func TestWriteMapMap(t *testing.T) {
	m := map[string]map[string]int{"x": map[string]int{"y": 3}, "z": map[string]int{}}
	opt := sen.Options{Indent: 2, OmitNil: true, Sort: true}
	s := sen.String(m, &opt)
	tt.Equal(t, `{
  x: {
    y: 3
  }
}`, s)

	opt.Indent = 0
	s = sen.String(m, &opt)
	tt.Equal(t, `{x:{y:3}}`, s)
}

func TestWriteStructOther(t *testing.T) {
	type Sample struct {
		X *int
		Y int
	}
	x := 1
	sample := Sample{X: &x, Y: 2}
	opt := sen.Options{Indent: 2}
	b := sen.Bytes(&sample, &opt)
	tt.Equal(t, `{
  x: 1
  y: 2
}`, string(b))

	opt.Indent = 0
	b = sen.Bytes(&sample, &opt)
	tt.Equal(t, `{x:1 y:2}`, string(b))
}

func TestWriteStructOmit(t *testing.T) {
	type Sample struct {
		X *int
	}
	sample := Sample{X: nil}
	opt := sen.Options{Indent: 2, OmitNil: true}
	b := sen.Bytes(&sample, &opt)
	tt.Equal(t, `{}`, string(b))

	opt.Indent = 0
	b = sen.Bytes(&sample, &opt)
	tt.Equal(t, `{}`, string(b))
}

func TestWriteStructEmbed(t *testing.T) {
	type In struct {
		X int
	}
	type Out struct {
		In In
		Y  int
	}
	o := Out{In: In{X: 1}, Y: 2}
	opt := sen.Options{Indent: 2}
	b := sen.Bytes(&o, &opt)
	tt.Equal(t, `{
  in: {
    x: 1
  }
  y: 2
}`, string(b))

	opt.Indent = 0
	b = sen.Bytes(&o, &opt)
	tt.Equal(t, `{in:{x:1} y:2}`, string(b))
}

func TestWriteStructAnonymous(t *testing.T) {
	type In struct {
		X int
	}
	type Out struct {
		In
		Y int
	}
	o := Out{In: In{X: 1}, Y: 2}
	opt := ojg.Options{Indent: 2, NestEmbed: true}
	b := sen.Bytes(&o, &opt)
	tt.Equal(t, `{
  in: {
    x: 1
  }
  y: 2
}`, string(b))

	opt.Indent = 0
	b = sen.Bytes(&o, &opt)
	tt.Equal(t, `{in:{x:1} y:2}`, string(b))
}
