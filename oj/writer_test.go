// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

// Used to test Simplifier objects in simple data.
type simon struct {
	x int
}

func (s *simon) Simplify() any {
	return map[string]any{
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

type Mix struct {
	Val   int     `json:"v,omitempty"`
	Str   int     `json:"str,omitempty,string"`
	Title string  `json:",omitempty"`
	Skip  int     `json:"-"`
	Dash  float64 `json:"-,omitempty"`
	Boo   bool    `json:"boo"`
}

type Nest struct {
	List []*Dummy
}

type Dummy struct {
	Val int
}

type Stew []int

func (s Stew) String() string {
	return fmt.Sprintf("%v", []int(s))
}

type Panik struct {
}

func (p *Panik) Simplify() any {
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
	opt := &oj.Options{}
	tm := time.Date(2020, time.May, 7, 19, 29, 19, 123456789, time.UTC)
	tm2 := time.Unix(-10, -100000000)
	for i, d := range []data{
		{value: nil, expect: "null"},
		{value: true, expect: "true"},
		{value: false, expect: "false"},
		{value: "string", expect: `"string"`},
		{value: "\\\t\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122",
			expect:  `"\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029\u0007𝄢"`,
			options: &oj.Options{HTMLUnsafe: false},
		},
		{value: []byte{'a', 'b', 'c'}, expect: `"abc"`, options: &oj.Options{BytesAs: ojg.BytesAsString}},
		{value: []byte{'a', 'b', 'c'}, expect: `"YWJj"`, options: &oj.Options{BytesAs: ojg.BytesAsBase64}},
		{value: []byte{'a', 'b', 'c'}, expect: "[97,98,99]", options: &oj.Options{BytesAs: ojg.BytesAsArray}},
		{value: "&<>", expect: `"&<>"`},
		{value: gen.String("string"), expect: `"string"`},
		{value: []any{true, false}, expect: "[true,false]"},
		{value: gen.Array{gen.Bool(true), nil}, expect: "[true,null]"},
		{value: []any{true, false}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: []any{true, false}, expect: "[\n\ttrue,\n\tfalse\n]", options: &oj.Options{Tab: true}},
		{value: []any{[]any{}, []any{}}, expect: "[[],[]]"},
		{value: []any{map[string]any{}, map[string]any{}}, expect: "[{},{}]"},
		{value: []any{[]any{}, []any{}}, expect: "[\n  [],\n  []\n]", options: &ojg.Options{Indent: 2}},
		{value: []any{map[string]any{}, map[string]any{}}, expect: "[{},{}]", options: &ojg.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, indent: 2, expect: "[\n  true,\n  false\n]"},
		{value: gen.Object{"t": gen.True, "f": gen.False}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: map[string]any{"t": true, "f": false}, expect: `{"f":false,"t":true}`, options: &oj.Options{Sort: true}},
		{value: gen.Array{gen.True, gen.False}, expect: "[true,false]", options: opt},
		{value: gen.Array{gen.False, gen.True}, expect: "[false,true]", options: opt},
		{value: []any{-1, int8(2), int16(-3), int32(4), int64(-5)}, expect: "[-1,2,-3,4,-5]"},
		{value: []any{uint(1), 'A', uint8(2), uint16(3), uint32(4), uint64(5)}, expect: "[1,65,2,3,4,5]"},
		{value: gen.Array{gen.Int(1), gen.Float(1.2)}, expect: "[1,1.2]"},
		{value: []any{float32(1.2), float64(2.1)}, expect: "[1.2,2.1]"},
		{value: []any{tm}, expect: "[1588879759123456789]"},
		{value: []any{tm}, expect: `[{"^":"Time","value":"2020-05-07T19:29:19.123456789Z"}]`,
			options: &oj.Options{TimeMap: true, CreateKey: "^", TimeFormat: time.RFC3339Nano}},
		{value: []any{tm}, expect: `[{"^":"time/Time","value":"2020-05-07T19:29:19.123456789Z"}]`,
			options: &oj.Options{TimeMap: true, CreateKey: "^", TimeFormat: time.RFC3339Nano, FullTypePath: true}},
		{value: tm2, expect: "-10.100000000", options: &oj.Options{TimeFormat: "second"}},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759123456789]"},
		{value: gen.Array{gen.Time(tm)}, expect: `["2020-05-07T19:29:19.123456789Z"]`,
			options: &oj.Options{TimeFormat: time.RFC3339Nano}},
		{value: gen.Array{gen.Time(tm)}, expect: "[1588879759.123456789]", options: &oj.Options{TimeFormat: "second"}},
		{value: gen.Array{gen.Time(tm)}, expect: `[{"@":1588879759123456789}]`, options: &oj.Options{TimeWrap: "@"}},
		{value: map[string]any{"t": true, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: map[string]any{"t": true, "f": false}, expect: "{\n  \"f\": false,\n  \"t\": true\n}",
			options: &oj.Options{Sort: true, Indent: 2}},

		{value: map[string]any{"t": true}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\"f\":false,\"t\":true}",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\n  \"f\": false,\n  \"t\": true\n}",
			options: &oj.Options{OmitNil: true, Sort: true, Indent: 2}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\n  \"f\": false,\n  \"n\": null,\n  \"t\": true\n}",
			options: &oj.Options{OmitNil: false, Sort: true, Indent: 2}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\"f\":false,\"t\":true}",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\"f\":false,\"n\":null,\"t\":true}",
			options: &oj.Options{OmitNil: false, Sort: true}},
		{value: map[string]any{"t": true, "n": nil, "f": false}, expect: "{\n\t\"f\": false,\n\t\"n\": null,\n\t\"t\": true\n}",
			options: &oj.Options{OmitNil: false, Sort: true, Tab: true}},
		{value: map[string]any{"n": nil}, expect: "{\"n\":null}"},
		{value: map[string]any{"n": nil}, expect: "{}", options: &oj.Options{OmitNil: true, Sort: false, Indent: 2}},
		{value: map[string]any{"x": "", "y": []any{}, "z": map[string]any{}}, expect: "{}", options: &oj.Options{OmitEmpty: true}},
		{value: map[string]any{"x": "", "y": []any{}, "z": map[string]any{}}, expect: "{}",
			options: &oj.Options{OmitEmpty: true, Sort: true}},
		{value: map[string]any{"x": "", "y": []any{}, "z": map[string]any{}}, expect: "{}",
			options: &oj.Options{OmitEmpty: true, Indent: 2}},
		{value: map[string]any{"x": "", "y": []any{}, "z": map[string]any{}}, expect: "{}",
			options: &oj.Options{OmitEmpty: true, Sort: true, Indent: 2}},

		{value: gen.Object{"t": gen.True, "x": nil}, expect: "{\"t\":true}", options: &oj.Options{OmitNil: true}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2}},
		{value: gen.Object{"t": gen.True}, expect: "{\n  \"t\": true\n}", options: &oj.Options{Indent: 2, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\"f\":false,\"t\":true}",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\n  \"f\": false,\n  \"t\": true\n}",
			options: &oj.Options{OmitNil: true, Sort: true, Indent: 2}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\n  \"f\": false,\n  \"n\": null,\n  \"t\": true\n}",
			options: &oj.Options{OmitNil: false, Sort: true, Indent: 2}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\"f\":false,\"t\":true}",
			options: &oj.Options{OmitNil: true, Sort: true}},
		{value: gen.Object{"t": gen.True, "n": nil, "f": gen.False}, expect: "{\"f\":false,\"n\":null,\"t\":true}",
			options: &oj.Options{OmitNil: false, Sort: true}},
		{value: gen.Object{"n": nil}, expect: "{\"n\":null}"},
		{value: gen.Object{"n": nil}, expect: "{}", options: &oj.Options{OmitNil: true, Sort: false, Indent: 2}},

		{value: &simon{x: 3}, expect: `{"type":"simon","x":3}`, options: &oj.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `{"type":"genny","val":3}`, options: &oj.Options{Sort: true}},
		{value: &genny{val: 3}, expect: `{
  "type": "genny",
  "val": 3
}`, options: &oj.Options{Sort: true, Indent: 2}},
		{value: &Dummy{Val: 3}, expect: `{"val":3}`},
		{value: Stew{3}, expect: `"[3]"`, options: &oj.Options{NoReflect: true}},
		{value: &Dummy{Val: 3}, expect: `{"^":"Dummy","val":3}`, options: &oj.Options{Sort: true, CreateKey: "^"}},
		{value: &Dummy{Val: 3}, expect: `{"Val":3}`, options: &oj.Options{KeyExact: true}},
		{value: complex(1, 7), expect: `{
  "imag": 7,
  "real": 1
}`, options: &ojg.Options{Indent: 2, Sort: true}},
		{value: complex(1, 7), expect: `{"imag":7,"real":1}`, options: &ojg.Options{Sort: true}},
		{value: complex(1, 7), expect: `"(1+7i)"`, options: &ojg.Options{Indent: 2, NoReflect: true}},
		{value: complex(1, 7), expect: `"(1+7i)"`, options: &ojg.Options{Indent: 0, NoReflect: true}},
		{value: []int{1, 2}, expect: "[1,2]", options: &ojg.Options{Indent: 0}},
		{value: []int{1, 2}, expect: `[
  1,
  2
]`, options: &ojg.Options{Indent: 2}},
	} {
		var s string
		if d.options == nil {
			if 0 < d.indent {
				s = oj.JSON(d.value, d.indent)
			} else {
				s = oj.JSON(d.value)
			}
		} else {
			s = oj.JSON(d.value, d.options)
		}
		tt.Equal(t, d.expect, s, fmt.Sprintf("%d: %v", i, d.value))
	}
}

func TestWrite(t *testing.T) {
	var b strings.Builder

	err := oj.Write(&b, []any{true, false})
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	opt := oj.Options{WriteLimit: 8}
	b.Reset()
	err = oj.Write(&b, []any{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	// A second time.
	b.Reset()
	err = oj.Write(&b, []any{true, false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", b.String())

	b.Reset()
	err = oj.Write(&b, []any{false, true}, 2)
	tt.Nil(t, err)
	tt.Equal(t, "[\n  false,\n  true\n]", b.String())

	b.Reset()
	// Force a realloc of string buffer.
	err = oj.Write(&b, strings.Repeat("Xyz ", 63)+"\U0001D122", 2)
	tt.Nil(t, err)
	tt.Equal(t, 258, len(b.String()))

	// Make sure a comma separator is added in unsorted-unindent mode.
	b.Reset()
	err = oj.Write(&b, map[string]any{"t": true, "f": false})
	tt.Nil(t, err)
	tt.Equal(t, 20, len(b.String()))
	b.Reset()
	err = oj.Write(&b, gen.Object{"t": gen.True, "f": gen.False})
	tt.Nil(t, err)
	tt.Equal(t, 20, len(b.String()))

	b.Reset()
	opt.Sort = true
	err = oj.Write(&b, map[string]any{"t": true, "f": false}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 20, len(b.String()))
	b.Reset()
	err = oj.Write(&b, gen.Object{"t": gen.True, "f": gen.False}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 20, len(b.String()))
}

func TestWriteWide(t *testing.T) {
	var b strings.Builder
	opt := oj.Options{Indent: 300}
	err := oj.Write(&b, []any{[]any{true, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 530, len(b.String()))

	b.Reset()
	err = oj.Write(&b, gen.Array{gen.Array{gen.True, nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 530, len(b.String()))

	b.Reset()
	err = oj.Write(&b, map[string]any{"x": map[string]any{"y": true, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 545, len(b.String()))

	b.Reset()
	err = oj.Write(&b, gen.Object{"x": gen.Object{"y": gen.True, "z": nil}}, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 545, len(b.String()))
}

func TestWriteDeep(t *testing.T) {
	var b strings.Builder
	opt := oj.Options{Tab: true}
	a := []any{map[string]any{"x": true}}
	for i := 40; 0 < i; i-- {
		a = []any{a}
	}
	err := oj.Write(&b, a, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1797, len(b.String()))

	b.Reset()
	g := gen.Array{gen.Object{"x": gen.True}}
	for i := 40; 0 < i; i-- {
		g = gen.Array{g}
	}
	err = oj.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1797, len(b.String()))

	opt.Sort = true
	b.Reset()
	err = oj.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 1797, len(b.String()))

	opt.Tab = false
	opt.Indent = 4
	b.Reset()
	err = oj.Write(&b, g, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 6833, len(b.String()))
}

func TestWriteShort(t *testing.T) {
	opt := oj.Options{Indent: 2, WriteLimit: 2}
	err := oj.Write(&shortWriter{max: 3}, []any{true, nil}, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	opt.Indent = 0
	err = oj.Write(&shortWriter{max: 3}, []any{true, nil}, &opt)
	tt.NotNil(t, err)
	err = oj.Write(&shortWriter{max: 3}, gen.Array{gen.True, nil}, &opt)
	tt.NotNil(t, err)

	obj := map[string]any{"t": true, "n": nil}
	sobj := gen.Object{"t": gen.True, "n": nil}
	opt.Indent = 0
	for i := 2; i < 19; i += 2 {
		err = oj.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = oj.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)

		opt.Sort = true
		err = oj.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = oj.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)
	}
	opt.Indent = 2
	for i := 2; i < 19; i += 2 {
		err = oj.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = oj.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)

		opt.Sort = false
		err = oj.Write(&shortWriter{max: i}, obj, &opt)
		tt.NotNil(t, err)
		err = oj.Write(&shortWriter{max: i}, sobj, &opt)
		tt.NotNil(t, err)
	}
}

func TestMarshal(t *testing.T) {
	b, err := oj.Marshal([]gen.Node{gen.True, gen.False}, 0)
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", string(b))

	b, err = oj.Marshal([]any{true, false}, &ojg.Options{})
	tt.Nil(t, err)
	tt.Equal(t, "[true,false]", string(b))

	_, err = oj.Marshal([]any{true, TestMarshal})
	tt.NotNil(t, err)

	_, err = oj.Marshal([]any{true, &Panik{}})
	tt.NotNil(t, err)

	b, err = oj.Marshal([]any{true, &Dummy{Val: 3}})
	tt.Nil(t, err)
	tt.Equal(t, `[true,{"Val":3}]`, string(b))
	save := b
	b, err = oj.Marshal([]any{true, &Dummy{Val: 5}})
	tt.Nil(t, err)
	tt.Equal(t, `[true,{"Val":5}]`, string(b))
	tt.Equal(t, `[true,{"Val":3}]`, string(save))

	b, err = oj.Marshal([]any{true, &Dummy{Val: 3}}, &oj.Options{UseTags: false})
	tt.Nil(t, err)
	tt.Equal(t, `[true,{"val":3}]`, string(b))

	_, err = oj.Marshal([]any{true, &Dummy{Val: 3}}, &oj.Options{NoReflect: true})
	tt.NotNil(t, err)

	_, err = oj.Marshal([]any{true, &Dummy{Val: 3}}, &oj.Options{NoReflect: true, Indent: 2})
	tt.NotNil(t, err)

	wr := oj.Writer{}
	s := wr.JSON([]any{true, TestMarshal})
	tt.Equal(t, `[true,null]`, s)

	wr.Indent = 2
	s = wr.JSON([]any{true, TestMarshal})
	tt.Equal(t, `[
  true,
  null
]`, s)
}

func TestWriteBad(t *testing.T) {
	var b strings.Builder

	err := oj.Write(&b, []any{true, TestWriteBad})
	tt.Nil(t, err)
	tt.Equal(t, `[true,null]`, b.String())

	b.Reset()
	err = oj.Write(&b, []any{true, &Panik{}})
	tt.NotNil(t, err)
}

func TestJSONBad(t *testing.T) {
	out := oj.JSON([]any{true, &Panik{}})
	tt.Equal(t, 0, len(out))
}

func TestMarshalStruct(t *testing.T) {
	n := Nest{
		List: []*Dummy{
			{Val: 1},
			{Val: 2},
		},
	}
	j, err := oj.Marshal(&n)
	tt.Nil(t, err)
	tt.Equal(t, `{"List":[{"Val":1},{"Val":2}]}`, string(j))

	j, err = oj.Marshal(&n, 2)
	tt.Nil(t, err)
	expect := `{
  "List": [
    {
      "Val": 1
    },
    {
      "Val": 2
    }
  ]
}`
	tt.Equal(t, expect, string(j))

	type empty struct {
		X int `json:"x,omitempty"`
		Y int `json:"y,omitempty"`
	}
	j, err = oj.Marshal(&empty{X: 0, Y: 0}, 2)
	tt.Nil(t, err)
	tt.Equal(t, `{}`, string(j))

	j, err = oj.Marshal(&empty{X: 1, Y: 0}, 2)
	tt.Nil(t, err)
	tt.Equal(t, `{
  "x": 1
}`, string(j))
}

type One struct {
	X int
}

type Two struct {
	Y int
	One
}

type Three struct {
	Two
	Z int
}

func TestMarshalNestedStruct(t *testing.T) {
	obj := &Three{Two: Two{One: One{X: 1}, Y: 2}, Z: 3}
	j, err := oj.Marshal(obj)
	tt.Nil(t, err)
	tt.Equal(t, `{"X":1,"Y":2,"Z":3}`, string(j))
}

func TestMarshalNestedPtr(t *testing.T) {
	type Inner struct {
		X int
	}
	type Wrap struct {
		*Inner
	}
	obj := &Wrap{Inner: &Inner{X: 3}}
	j, err := oj.Marshal(obj)
	tt.Nil(t, err)
	tt.Equal(t, `{"X":3}`, string(j))
}

func TestMarshalMap(t *testing.T) {
	type Dap struct {
		M map[string]*Dummy
	}
	d := Dap{M: map[string]*Dummy{
		"a": {Val: 1},
	}}
	j, err := oj.Marshal(&d)
	tt.Nil(t, err)
	tt.Equal(t, `{"M":{"a":{"Val":1}}}`, string(j))
}

func TestMarshalTypeAlias(t *testing.T) {
	type Stringy string
	d := Stringy("s")
	j, err := oj.Marshal(&d)
	tt.Nil(t, err)
	tt.Equal(t, `"s"`, string(j))
}

func TestMarshalWithWriter(t *testing.T) {
	j, err := oj.Marshal([]any{true}, &oj.Writer{})
	tt.Nil(t, err)
	tt.Equal(t, "[true]", string(j))
}

func TestWriteStructWide(t *testing.T) {
	type Nest struct {
		Dig   *Nest
		Level int
	}
	var b strings.Builder
	opt := oj.Options{Tab: true, OmitNil: true}
	var n *Nest
	for i := 40; 0 < i; i-- {
		n = &Nest{Dig: n, Level: i}
	}
	err := oj.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 3098, len(b.String()))

	b.Reset()
	opt.Tab = false
	opt.Indent = 8
	err = oj.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 13166, len(b.String()))

	b.Reset()
	opt.Indent = 0
	err = oj.Write(&b, n, &opt)
	tt.Nil(t, err)
	tt.Equal(t, 744, len(b.String()))
}

func TestWriteStructCreateKey(t *testing.T) {
	type Sample struct {
		X int
		Y int
		z int
	}
	sample := Sample{X: 1, Y: 2, z: 3}
	opt := oj.Options{Indent: 2, CreateKey: "^"}
	s := oj.JSON(&sample, &opt)
	tt.Equal(t, `{
  "^": "Sample",
  "x": 1,
  "y": 2
}`, s)

	opt.FullTypePath = true
	s = oj.JSON(&sample, &opt)
	tt.Equal(t, `{
  "^": "github.com/ohler55/ojg/oj_test/Sample",
  "x": 1,
  "y": 2
}`, s)

	opt.Indent = 0
	s = oj.JSON(&sample, &opt)
	tt.Equal(t, `{"^":"github.com/ohler55/ojg/oj_test/Sample","x":1,"y":2}`, s)

	opt.FullTypePath = false
	s = oj.JSON(&sample, &opt)
	tt.Equal(t, `{"^":"Sample","x":1,"y":2}`, s)
}

func TestWriteStruct(t *testing.T) {
	type Sample struct {
		X []int
		Y map[string]int
		Z *int
	}
	sample := Sample{X: []int{1}, Y: map[string]int{"y": 2}}
	opt := oj.Options{Indent: 2, OmitNil: true}
	s := oj.JSON(&sample, &opt)
	tt.Equal(t, `{
  "x": [
    1
  ],
  "y": {
    "y": 2
  }
}`, s)
}

func TestWriteStructSkip(t *testing.T) {
	type Skippy struct {
		X int `json:"a,omitempty"`
		Y int `json:"b,omitempty"`
		z int
	}
	opt := oj.Options{Indent: 2, OmitNil: true, UseTags: true}
	skippy := Skippy{X: 0, Y: 1, z: 2}
	s := oj.JSON(&skippy, &opt)
	tt.Equal(t, `{
  "b": 1
}`, s)
	opt.Indent = 0
	s = oj.JSON(&skippy, &opt)
	tt.Equal(t, `{"b":1}`, s)

	opt.Indent = 2
	skippy.X = 1
	skippy.Y = 0
	s = oj.JSON(&skippy, &opt)
	tt.Equal(t, `{
  "a": 1
}`, s)

	opt.Indent = 0
	s = oj.JSON(&skippy, &opt)
	tt.Equal(t, `{"a":1}`, s)
}

func TestWriteStructSkipEmpty(t *testing.T) {
	opt := oj.Options{Indent: 2, OmitNil: true, OmitEmpty: true}
	type empty struct {
		X string
		Y []any
		Z map[string]any
	}
	s := oj.JSON(&empty{X: "", Y: []any{}, Z: map[string]any{}}, &opt)
	tt.Equal(t, `{}`, s)
}

func TestWriteStructSkipReflectEmpty(t *testing.T) {
	opt := oj.Options{Indent: 2, OmitNil: true, OmitEmpty: true}
	type empty struct {
		X string
		Y []int
		Z map[string]int
	}
	s := oj.JSON(&empty{X: "", Y: []int{}, Z: map[string]int{}}, &opt)
	tt.Equal(t, `{}`, s)

	opt.Indent = 0
	s = oj.JSON(&empty{X: "", Y: []int{}, Z: map[string]int{}}, &opt)
	tt.Equal(t, `{}`, s)
}

func TestWriteReflectMapEmpty(t *testing.T) {
	opt := oj.Options{Indent: 2, OmitNil: true, OmitEmpty: true}

	type str string
	s := oj.JSON(map[string]str{"x": str("abc")}, &opt)
	tt.Equal(t, `{
  "x": "abc"
}`, s)
	s = oj.JSON(map[string]str{"x": str("")}, &opt)
	tt.Equal(t, `{}`, s)

	s = oj.JSON(map[string][]int{"x": {}}, &opt)
	tt.Equal(t, `{}`, s)

	s = oj.JSON(map[string]map[string]int{"x": {}}, &opt)
	tt.Equal(t, `{}`, s)

	opt.Indent = 0
	s = oj.JSON(map[string]str{"x": str("abc")}, &opt)
	tt.Equal(t, `{"x":"abc"}`, s)

	s = oj.JSON(map[string]str{"x": str("")}, &opt)
	tt.Equal(t, `{}`, s)

	s = oj.JSON(map[string][]int{"x": {}}, &opt)
	tt.Equal(t, `{}`, s)

	s = oj.JSON(map[string]map[string]int{"x": {}}, &opt)
	tt.Equal(t, `{}`, s)
}

func TestWriteSliceWide(t *testing.T) {
	type Nest struct {
		Dig []Nest
	}
	opt := oj.Options{Tab: true}
	n := &Nest{}
	for i := 20; 0 < i; i-- {
		n = &Nest{Dig: []Nest{*n}}
	}
	s := oj.JSON(n, &opt)
	tt.Equal(t, 1852, len(s))

	opt.Tab = false
	opt.Indent = 4
	s = oj.JSON(n, &opt)
	tt.Equal(t, 6713, len(s))

	opt.Indent = 0
	s = oj.JSON(n, &opt)
	tt.Equal(t, 210, len(s))
}

func TestWriteSliceArray(t *testing.T) {
	type Matrix struct {
		Rows [][4]int
	}
	opt := oj.Options{Indent: 2}
	m := &Matrix{Rows: [][4]int{{1, 2, 3, 4}}}
	s := oj.JSON(m, &opt)
	tt.Equal(t, `{
  "rows": [
    [
      1,
      2,
      3,
      4
    ]
  ]
}`, s)

	opt.Indent = 0
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{"rows":[[1,2,3,4]]}`, s)

	opt.OmitNil = false
	s = oj.JSON(&Matrix{}, &opt)
	tt.Equal(t, `{"rows":[]}`, s)
}

func TestWriteSliceMap(t *testing.T) {
	type SMS struct {
		Maps []map[string]int
	}
	opt := oj.Options{Indent: 2, Sort: true}
	m := &SMS{Maps: []map[string]int{{"x": 1, "y": 2}}}
	s := oj.JSON(m, &opt)
	tt.Equal(t, `{
  "maps": [
    {
      "x": 1,
      "y": 2
    }
  ]
}`, s)

	opt.Indent = 0
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{"maps":[{"x":1,"y":2}]}`, s)

	opt.OmitEmpty = true
	s = oj.JSON(&SMS{}, &opt)
	tt.Equal(t, `{}`, s)
}

func TestWriteMapWide(t *testing.T) {
	type Nest struct {
		Dig map[string]*Nest
	}
	opt := oj.Options{Tab: true, OmitNil: true}
	n := &Nest{map[string]*Nest{"x": nil}}
	for i := 16; 0 < i; i-- {
		n = &Nest{Dig: map[string]*Nest{"x": n}}
	}
	s := oj.JSON(n, &opt)
	tt.Equal(t, 1396, len(s))

	opt.Tab = false
	opt.Indent = 4
	s = oj.JSON(n, &opt)
	tt.Equal(t, 4685, len(s))

	opt.Indent = 0
	s = oj.JSON(n, &opt)
	tt.Equal(t, 234, len(s))

	opt.OmitEmpty = true
	s = oj.JSON(&Nest{}, &opt)
	tt.Equal(t, "{}", s)
}

func TestWriteMapSlice(t *testing.T) {
	m := map[string][]int{"x": {1, 2, 3}, "y": {}}
	opt := oj.Options{Indent: 2, OmitNil: true, Sort: true}
	s := oj.JSON(m, &opt)
	tt.Equal(t, `{
  "x": [
    1,
    2,
    3
  ]
}`, s)

	opt.OmitNil = false
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{
  "x": [
    1,
    2,
    3
  ],
  "y": []
}`, s)

	opt.Indent = 0
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{"x":[1,2,3],"y":[]}`, s)

	opt.OmitNil = true
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{"x":[1,2,3]}`, s)
}

func TestWriteMapMap(t *testing.T) {
	m := map[string]map[string]int{"x": {"y": 3}, "z": {}}
	opt := oj.Options{Indent: 2, OmitNil: true, Sort: true}
	s := oj.JSON(m, &opt)
	tt.Equal(t, `{
  "x": {
    "y": 3
  }
}`, s)

	opt.Indent = 0
	s = oj.JSON(m, &opt)
	tt.Equal(t, `{"x":{"y":3}}`, s)
}

func TestWriteStructOther(t *testing.T) {
	type Sample struct {
		X *int
		Y int
	}
	x := 1
	sample := Sample{X: &x, Y: 2}
	opt := oj.Options{Indent: 2}
	b, err := oj.Marshal(&sample, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{
  "x": 1,
  "y": 2
}`, string(b))

	opt.Indent = 0
	b, err = oj.Marshal(&sample, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{"x":1,"y":2}`, string(b))
}

func TestWriteStructOmit(t *testing.T) {
	type Sample struct {
		X *int
	}
	sample := Sample{X: nil}
	opt := oj.Options{Indent: 2, OmitNil: true}
	b, err := oj.Marshal(&sample, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{}`, string(b))

	opt.Indent = 0
	b, err = oj.Marshal(&sample, &opt)
	tt.Nil(t, err)
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
	opt := oj.Options{Indent: 2}
	b, err := oj.Marshal(&o, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{
  "in": {
    "x": 1
  },
  "y": 2
}`, string(b))

	opt.Indent = 0
	b, err = oj.Marshal(&o, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{"in":{"x":1},"y":2}`, string(b))
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
	opt := oj.Options{Indent: 2, NestEmbed: true}
	b, err := oj.Marshal(&o, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{
  "in": {
    "x": 1
  },
  "y": 2
}`, string(b))

	opt.Indent = 0
	b, err = oj.Marshal(&o, &opt)
	tt.Nil(t, err)
	tt.Equal(t, `{"in":{"x":1},"y":2}`, string(b))
}

func TestWriteSliceNil(t *testing.T) {
	var a []any
	b, err := oj.Marshal(a)
	tt.Nil(t, err)
	tt.Equal(t, "null", string(b))

	a = []any{}
	b, err = oj.Marshal(a)
	tt.Nil(t, err)
	tt.Equal(t, "[]", string(b))
}

func TestWriteStrictPanic(t *testing.T) {
	_, err := oj.Marshal(func() {}, 2)
	tt.NotNil(t, err)

	_, err = oj.Marshal(func() {})
	tt.NotNil(t, err)
}

type Marsha struct {
	val int
}

func (m *Marsha) MarshalJSON() ([]byte, error) {
	if m.val == 5 {
		return nil, fmt.Errorf("oops")
	}
	return []byte(fmt.Sprintf(`{"v":%d}`, m.val)), nil
}

func TestMarshalMarshaler(t *testing.T) {
	j, err := oj.Marshal(&Marsha{val: 3})
	tt.Nil(t, err)
	tt.Equal(t, `{"v":3}`, string(j))

	_, err = oj.Marshal(&Marsha{val: 5})
	tt.NotNil(t, err)
}

type TM struct {
	val int
}

func (tm *TM) MarshalText() ([]byte, error) {
	if tm.val == 5 {
		return nil, fmt.Errorf("oops")
	}
	return []byte(fmt.Sprintf("-- %d --", tm.val)), nil
}

func TestMarshalTextMarshaler(t *testing.T) {
	j, err := oj.Marshal(&TM{val: 3})
	tt.Nil(t, err)
	tt.Equal(t, `"-- 3 --"`, string(j))

	_, err = oj.Marshal(&TM{val: 5})
	tt.NotNil(t, err)
}

func TestMarshalNoAddr(t *testing.T) {
	type Sample struct {
		When time.Time
		Mars Marsha
	}
	testCase := Sample{
		When: time.Unix(0, 0),
		Mars: Marsha{val: 3},
	}
	ojg.ErrorWithStack = true
	_, err := oj.Marshal(testCase)
	tt.Nil(t, err)
}

func TestWriteFloatFormat(t *testing.T) {
	var wr oj.Writer
	wr.FloatFormat = "%05.2f"
	j := wr.MustJSON(1.234)
	tt.Equal(t, `01.23`, string(j))

	j = wr.MustJSON(float32(1.234))
	tt.Equal(t, `01.23`, string(j))
}

func BenchmarkMarshalFlat(b *testing.B) {
	m := Mix{
		Val:   1,
		Str:   2,
		Title: "Mix",
		Skip:  4,
		Dash:  5.5,
		Boo:   true,
	}
	for i := 0; i < b.N; i++ {
		if _, err := oj.Marshal(&m); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNestList(b *testing.B) {
	n := Nest{
		List: []*Dummy{
			{Val: 1},
			{Val: 2},
			{Val: 3},
		},
	}
	for i := 0; i < b.N; i++ {
		if _, err := oj.Marshal(&n); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteNestIndentList(b *testing.B) {
	n := Nest{
		List: []*Dummy{
			{Val: 1},
			{Val: 2},
			{Val: 3},
		},
	}
	wr := oj.Writer{Options: oj.Options{Indent: 2}}
	for i := 0; i < b.N; i++ {
		_ = wr.JSON(&n)
	}
}

func BenchmarkMarshalMap(b *testing.B) {
	type Dap struct {
		M map[string]*Dummy
	}
	d := Dap{M: map[string]*Dummy{
		"a": {Val: 1},
		"b": {Val: 2},
		"c": {Val: 3},
	}}
	for i := 0; i < b.N; i++ {
		if _, err := oj.Marshal(&d); err != nil {
			b.Fatal(err)
		}
	}
}

func TestWriteDev(t *testing.T) {
	data := map[string]any{"x": "", "y": []any{}, "z": map[string]any{}}
	opt := oj.Options{OmitEmpty: true}
	s := oj.JSON(data, &opt)
	tt.Equal(t, `{}`, s)
}
