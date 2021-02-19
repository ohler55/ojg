// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"fmt"
	"testing"

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
    a: 1,
    b: 2,
    c: 3,
    d: [
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
    a: 1,
    b: 2,
    c: 3,
    d: ["x", "y", "z", []]
  }
]`, s)

	s = pretty.JSON(val, &opt, 80.3)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {a: 1, b: 2, c: 3, d: ["x", "y", "z", []]}
]`, s)

	s = pretty.JSON(val, &opt, 0.4)
	tt.Equal(t, `[true, false, [3, 2, 1], {a: 1, b: 2, c: 3, d: ["x", "y", "z", []]}]`, s)
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
  {a: 1, b: 2, c: 3, d: ["x", "y", "z", []]}
]`, s)

	s = pretty.JSON(val, &opt, 40.4)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    a: 1,
    b: 2,
    c: 3,
    d: ["x", "y", "z", []]
  }
]`, s)

	s = pretty.JSON(val, &opt, 20.4)
	tt.Equal(t, `[
  true,
  false,
  [3, 2, 1],
  {
    a: 1,
    b: 2,
    c: 3,
    d: [
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
    a: 1,
    b: 2,
    c: 3,
    d: ["x", "y", "z", []]
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
    a: 1,
    b: 2,
    c: 3,
    d: ["x", "y", "z", []]
  }
]`, s)
}

func TestSEN2(t *testing.T) {
	p := sen.Parser{}
	val, err := p.Parse([]byte(`[true {abc: 123 def: true}]`))
	tt.Nil(t, err)
	opt := sen.DefaultOptions
	opt.Color = true
	s := pretty.JSON(val, 25, &opt)

	fmt.Printf("*** %s\n", s)
	//fmt.Printf("*** % x\n", s)

}
