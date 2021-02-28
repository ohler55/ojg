// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

const sample = `[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]`

var testColor = sen.Options{
	Color:       true,
	SyntaxColor: "s",
	KeyColor:    "k",
	NullColor:   "n",
	BoolColor:   "b",
	NumberColor: "0",
	StringColor: "q",
	NoColor:     "x",
	TimeFormat:  time.RFC3339Nano,
}

type Dummy struct {
	Val int
}

func (d *Dummy) String() string {
	return fmt.Sprintf("{val: %d}", d.Val)
}

type genny struct {
	val int
}

func (g *genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
}

type Pan int

func (p Pan) Simplify() interface{} {
	panic("force fail")
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

func TestJSONDepth(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	s := pretty.JSON(val, &opt, 80.1)
	tt.Equal(t, `[
  true,
  false,
  [
    3,
    2,
    1
  ],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": [
      "x",
      "y",
      "z",
      []
    ]
  }
]`, s)

	s = pretty.JSON(val, &opt, 80.0)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, s)

	s = pretty.JSON(val, &opt, 80.3)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
]`, s)

	s = pretty.JSON(val, &opt, 0.4)
	tt.Equal(t, `[true, false, [3, 2, 1], {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}]`, s)
}

func TestJSONEdge(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	s := pretty.JSON(val, &opt, 60.4)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {"a": 1, "b": 2, "c": 3, "d": ["x", "y", "z", []]}
]`, s)

	s = pretty.JSON(val, &opt, 40.4)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, s)

	s = pretty.JSON(val, &opt, 20.4)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": [
      "x",
      "y",
      "z",
      []
    ]
  }
]`, s)
}

func TestJSONIntArg(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	s := pretty.JSON(val, &opt, 30)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, s)
}

func TestJSONOjOptions(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := oj.DefaultOptions
	s := pretty.JSON(val, &opt, 80.2)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, s)
}

func TestInit(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	s := pretty.JSON(val, &sen.Options{}, 80.2)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, s)
}

func TestTypes(t *testing.T) {
	when := time.Date(2021, 2, 9, 10, 11, 12, 111, time.UTC)
	val := []interface{}{nil, 1.25, float32(1.5), "abc", when, map[string]interface{}{}}
	opt := sen.DefaultOptions
	opt.TimeFormat = time.RFC3339Nano
	s := pretty.JSON(val, &opt)
	tt.Equal(t, `[null, 1.25, 1.5, "abc", "2021-02-09T10:11:12.000000111Z", {}]`, s)

	opt = testColor
	s = pretty.JSON(val, &opt)
	tt.Equal(t, `s[xnnullxs,x 01.25xs,x 01.5xs,x q"abc"xs,x q"2021-02-09T10:11:12.000000111Z"xs,x s{xs}xs]x`, s)
}

func TestQuotedString(t *testing.T) {
	val := []interface{}{"\\\t\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122 ぴーたー"}
	s := pretty.JSON(val, &sen.Options{HTMLSafe: true})
	tt.Equal(t, `["\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029\u0007𝄢 ぴーたー"]`, s)
	s = pretty.JSON(val, &sen.Options{HTMLSafe: false})
	tt.Equal(t, `["\\\t\n\r\b\f\"&<>\u2028\u2029\u0007𝄢 ぴーたー"]`, s)
}

func TestIntTypes(t *testing.T) {
	val := []interface{}{
		[]interface{}{int8(-8), int16(-16), int32(-32), int64(-64), int(-1)},
		[]interface{}{uint8(8), uint16(16), uint32(32), uint64(64), uint(1)},
	}
	s := pretty.JSON(val, 80.2)
	tt.Equal(t, `[
  [-8, -16, -32, -64, -1],
  [8, 16, 32, 64, 1]
]`, s)
}

func TestGen(t *testing.T) {
	when := time.Date(2021, 2, 9, 10, 11, 12, 111, time.UTC)
	val := gen.Array{
		gen.True,
		gen.Int(3),
		gen.Float(1.5),
		gen.String("abc"),
		gen.Object{"x": nil, "y": gen.False},
		gen.Time(when),
	}
	opt := sen.DefaultOptions
	opt.TimeFormat = time.RFC3339Nano
	s := pretty.JSON(val, &opt, 80.3)
	tt.Equal(t, `[
  true,
  3,
  1.5,
  "abc",
  {"x": null, "y": false},
  "2021-02-09T10:11:12.000000111Z"
]`, s)
}

func TestPanic(t *testing.T) {
	s := pretty.JSON(Pan(1), &sen.Options{})
	tt.Equal(t, "", s)
}

func TestSEN(t *testing.T) {
	when := time.Date(2021, 2, 9, 10, 11, 12, 111, time.UTC)
	p := sen.Parser{}
	val, err := p.Parse([]byte(`[true {abc: 123 def: null} 1.25, xyz]`))
	a, _ := val.([]interface{})
	a = append(a, when)
	tt.Nil(t, err)
	opt := testColor
	s := pretty.SEN(a, &opt, 80.2)

	tt.Equal(t, `s[x
  btruex
  s{xkabcxs:x 0123x kdefxs:x nnullxs}x
  01.25x
  qxyzx
  q"2021-02-09T10:11:12.000000111Z"x
s]x`, s)
}

func TestSENGenMap(t *testing.T) {
	val := gen.Object{"a": gen.Int(1), "b": gen.Int(2)}
	opt := testColor
	s := pretty.SEN(val, &opt)
	tt.Equal(t, `s{xkaxs:x 01x kbxs:x 02xs}x`, s)
}

func TestDeep(t *testing.T) {
	val := []interface{}{}
	for i := 0; i < 10; i++ {
		val = []interface{}{val}
	}
	opt := sen.DefaultOptions
	s := pretty.SEN(val, &opt, 10.1)
	tt.Equal(t, `[
 [
  [
   [
    [
     [
      [
       [
        [
         [
          [
          ]
         ]
        ]
       ]
      ]
     ]
    ]
   ]
  ]
 ]
]`, s)

	// Deeper still to hit the max indent.
	for i := 0; i < 120; i++ {
		val = []interface{}{val}
	}
	s = pretty.SEN(val, &opt, 120.1)
	tt.Equal(t, 16902, len(s))

	// Deep map
	m := map[string]interface{}{}
	for i := 0; i < 130; i++ {
		m = map[string]interface{}{"o": m}
	}
	s = pretty.SEN(m, &opt, 120.1)
	tt.Equal(t, 17292, len(s))
}

func TestWriteJSON(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	opt.WriteLimit = 20
	var b strings.Builder
	err = pretty.WriteJSON(&b, val, &opt, 80.2)
	tt.Nil(t, err)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    "a": 1,
    "b": 2,
    "c": 3,
    "d": ["x", "y", "z", []]
  }
]`, b.String())
}

func TestWriteSEN(t *testing.T) {
	val, err := sen.Parse([]byte(sample))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	opt.WriteLimit = 20
	var b strings.Builder
	err = pretty.WriteSEN(&b, val, &opt, 80.2)
	tt.Nil(t, err)
	tt.Equal(t, `[
  true
  false
  [3 2 1]
  {
    a: 1
    b: 2
    c: 3
    d: [x y z []]
  }
]`, b.String())
}

func TestWritePanic(t *testing.T) {
	opt := testColor
	var b strings.Builder
	err := pretty.WriteSEN(&b, Pan(1), &opt)
	tt.Nil(t, err)
	tt.Equal(t, "x", b.String())
}

func TestWriteShort(t *testing.T) {
	opt := sen.DefaultOptions
	opt.WriteLimit = 2
	err := pretty.WriteJSON(&shortWriter{max: 3}, []interface{}{"abcdef"}, &opt)
	tt.NotNil(t, err)
}

func TestGenericer(t *testing.T) {
	s := pretty.JSON(&genny{val: 3})
	tt.Equal(t, `{"type": "genny", "val": 3}`, s)
}

func TestCreateKey(t *testing.T) {
	opt := sen.DefaultOptions
	opt.CreateKey = "^"
	s := pretty.JSON(&Dummy{Val: 3}, &opt)
	tt.Equal(t, `{"^": "Dummy", "val": 3}`, s)
}

func TestAsString(t *testing.T) {
	s := pretty.JSON(&Dummy{Val: 3})
	tt.Equal(t, `"{val: 3}"`, s)
}
