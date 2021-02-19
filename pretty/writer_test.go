// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

const sample = `[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]`

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
	s := pretty.JSON(val, &opt)
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
	s := pretty.JSON(val, &sen.Options{})
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
	val := []interface{}{nil, 1.25, float32(1.5), "abc", when}
	opt := sen.DefaultOptions
	opt.TimeFormat = time.RFC3339Nano
	s := pretty.JSON(val, &opt)
	tt.Equal(t, `[null, 1.25, 1.5, "abc", "2021-02-09T10:11:12.000000111Z"]`, s)
}

func TestQuotedString(t *testing.T) {
	val := []interface{}{"\\\t\n\r\b\f\"&<>\u2028\u2029\x07\U0001D122 ぴーたー"}
	s := pretty.JSON(val, &sen.Options{})
	tt.Equal(t, `["\\\t\n\r\b\f\"\u0026\u003c\u003e\u2028\u2029\u0007𝄢 ぴーたー"]`, s)
}

func TestIntTypes(t *testing.T) {
	val := []interface{}{
		[]interface{}{int8(-8), int16(-16), int32(-32), int64(-64), int(-1)},
		[]interface{}{uint8(8), uint16(16), uint32(32), uint64(64), uint(1)},
	}
	s := pretty.JSON(val)
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

type Pan int

func (p Pan) Simplify() interface{} {
	panic("force fail")
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
	opt := sen.DefaultOptions
	opt.Color = true
	opt.TimeFormat = time.RFC3339Nano
	s := pretty.SEN(a, &opt)

	// Uncomment the next two lines to see the colored string and the
	// corresponding hex.
	// fmt.Printf("*** %s\n", s)
	// fmt.Printf("*** % 2x\n", s)
	tt.Equal(t,
		"1b 5b 6d 5b 0a 20 20 1b 5b 33 33 6d 74 72 75 65 1b 5b 6d 0a 20 20 1b 5b 6d 7b 1b 5b 33 34 6d 1b "+
			"5b 33 34 6d 61 62 63 1b 5b 6d 1b 5b 6d 3a 20 1b 5b 33 36 6d 31 32 33 1b 5b 6d 20 1b 5b 33 34 6d "+
			"1b 5b 33 34 6d 64 65 66 1b 5b 6d 1b 5b 6d 3a 20 1b 5b 33 31 6d 6e 75 6c 6c 1b 5b 6d 1b 5b 6d 7d "+
			"0a 20 20 1b 5b 33 36 6d 31 2e 32 35 1b 5b 6d 0a 20 20 1b 5b 33 32 6d 78 79 7a 1b 5b 6d 0a 20 20 "+
			"1b 5b 33 32 6d 22 32 30 32 31 2d 30 32 2d 30 39 54 31 30 3a 31 31 3a 31 32 2e 30 30 30 30 30 30 "+
			"31 31 31 5a 22 1b 5b 6d 0a 1b 5b 6d 5d 1b 5b 6d",
		fmt.Sprintf("% 2x", s))
}
