// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty_test

import (
	"testing"

	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

func TestWriteAlignArrayNumbers(t *testing.T) {
	w := pretty.Writer{
		Width:    20,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		[]interface{}{1, 2, 3},
		[]interface{}{10, 20, 30},
		[]interface{}{100, 200, 300},
	}
	out, err := w.Marshal(data)
	tt.Nil(t, err)
	tt.Equal(t, `[
  [  1,   2,   3],
  [ 10,  20,  30],
  [100, 200, 300]
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  [  1   2   3]
  [ 10  20  30]
  [100 200 300]
]`, string(out))
}

func TestWriteAlignArrayStrings(t *testing.T) {
	w := pretty.Writer{
		Width:    30,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		[]interface{}{"alpha", "bravo", "charlie"},
		[]interface{}{"a", "b", "c"},
		[]interface{}{"andy", "betty"},
	}
	out, err := w.Marshal(data)
	tt.Nil(t, err)
	tt.Equal(t, `[
  ["alpha", "bravo", "charlie"],
  ["a"    , "b"    , "c"      ],
  ["andy" , "betty"]
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  [alpha bravo charlie]
  [a     b     c      ]
  [andy  betty]
]`, string(out))
}

func TestWriteAlignArrayNested(t *testing.T) {
	w := pretty.Writer{
		Width:    40,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		[]interface{}{1, 2, 3, []interface{}{100, 200, 300}},
		[]interface{}{1, 2, 3, "fourth"},
		[]interface{}{10, 20, 30, []interface{}{1, 20, 300}},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  [ 1,  2,  3, [100, 200, 300]],
  [ 1,  2,  3, "fourth"       ],
  [10, 20, 30, [  1,  20, 300]]
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  [ 1  2  3 [100 200 300]]
  [ 1  2  3 fourth       ]
  [10 20 30 [  1  20 300]]
]`, string(out))
}

func TestWriteAlignMixed(t *testing.T) {
	w := pretty.Writer{
		Width:    20,
		MaxDepth: 3,
		Align:    true,
		SEN:      true,
	}
	out, err := w.Marshal([]interface{}{
		[]interface{}{1, 2, 3},
		map[string]interface{}{"x": 1, "y": 2},
	})
	tt.Nil(t, err)
	tt.Equal(t, `[
  [1 2 3]
  {x: 1 y: 2}
]`, string(out))
}

func TestWriteAlignMapNumber(t *testing.T) {
	w := pretty.Writer{
		Width:    50,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		map[string]interface{}{"x": 1, "y": 2},
		map[string]interface{}{"z": 3, "y": 2},
		map[string]interface{}{"x": 100, "y": 200, "z": 300},
		map[string]interface{}{"x": 10, "z": 30},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  {"x":   1, "y":   2,         },
  {          "y":   2, "z":   3},
  {"x": 100, "y": 200, "z": 300},
  {"x":  10,           "z":  30}
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  {x:   1 y:   2       }
  {       y:   2 z:   3}
  {x: 100 y: 200 z: 300}
  {x:  10        z:  30}
]`, string(out))
}

func TestWriteAlignMapString(t *testing.T) {
	w := pretty.Writer{
		Width:    50,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		map[string]interface{}{"x": true, "y": false},
		map[string]interface{}{"z": nil, "y": "yoda"},
		map[string]interface{}{"x": "x-ray", "y": "yellow", "z": "zoo"},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  {"x": true   , "y": false   ,           },
  {              "y": "yoda"  , "z": null },
  {"x": "x-ray", "y": "yellow", "z": "zoo"}
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  {x: true  y: false         }
  {         y: yoda   z: null}
  {x: x-ray y: yellow z: zoo }
]`, string(out))
}

func TestWriteAlignMapNested(t *testing.T) {
	w := pretty.Writer{
		Width:    60,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		map[string]interface{}{"x": 1, "y": 2, "z": map[string]interface{}{"a": 1, "b": 2, "c": 3}},
		map[string]interface{}{"x": 100, "y": 200, "z": map[string]interface{}{"a": 10, "b": 20, "c": 30}},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  {"x":   1, "y":   2, "z": {"a":  1, "b":  2, "c":  3}},
  {"x": 100, "y": 200, "z": {"a": 10, "b": 20, "c": 30}}
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  {x:   1 y:   2 z: {a:  1 b:  2 c:  3}}
  {x: 100 y: 200 z: {a: 10 b: 20 c: 30}}
]`, string(out))
}

func TestWriteAlignMapArray(t *testing.T) {
	w := pretty.Writer{
		Width:    60,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		map[string]interface{}{"x": 1, "y": 2, "z": []interface{}{1, 2, 3}},
		map[string]interface{}{"x": 10, "y": 20, "z": []interface{}{10, 200, 3000}},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  {"x":  1, "y":  2, "z": [ 1,   2,    3]},
  {"x": 10, "y": 20, "z": [10, 200, 3000]}
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  {x:  1 y:  2 z: [ 1   2    3]}
  {x: 10 y: 20 z: [10 200 3000]}
]`, string(out))
}

func TestWriteAlignArrayMap(t *testing.T) {
	w := pretty.Writer{
		Width:    60,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		[]interface{}{1, 2, 3, map[string]interface{}{"x": 1, "y": 2, "z": 3}},
		[]interface{}{100, 200, 300, map[string]interface{}{"x": 1, "y": 20, "z": 300}},
	}
	out := w.Encode(data)
	tt.Equal(t, `[
  [  1,   2,   3, {"x": 1, "y":  2, "z":   3}],
  [100, 200, 300, {"x": 1, "y": 20, "z": 300}]
]`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `[
  [  1   2   3 {x: 1 y:  2 z:   3}]
  [100 200 300 {x: 1 y: 20 z: 300}]
]`, string(out))
}

func TestWriteAlignColor(t *testing.T) {
	w := pretty.Writer{
		Options:  testColor,
		Width:    60,
		MaxDepth: 3,
		Align:    true,
	}
	data := []interface{}{
		[]interface{}{1, 2, 3, map[string]interface{}{"x": 1, "y": 2, "z": 3}},
		[]interface{}{100, 200, 300, map[string]interface{}{"x": 1, "y": 20, "z": 300}},
	}
	out := w.Encode(data)
	tt.Equal(t, `s[x
  s[x  01xs,x   02xs,x   03xs,x s{xk"x"x: 01xs,x k"y"x:  02xs,x k"z"x:   03xs}xs]xs,x
  s[x0100xs,x 0200xs,x 0300xs,x s{xk"x"x: 01xs,x k"y"x: 020xs,x k"z"x: 0300xs}xs]x
s]x`, string(out))

	w.SEN = true
	out = w.Encode(data)
	tt.Equal(t, `s[x
  s[x  01x   02x   03x s{xkxx: 01x kyx:  02x kzx:   03xs}xs]x
  s[x0100x 0200x 0300x s{xkxx: 01x kyx: 020x kzx: 0300xs}xs]x
s]x`, string(out))
}

type simplyPanic int

func (sp simplyPanic) Simplify() interface{} {
	panic("no can do")
}

func TestMarshalError(t *testing.T) {
	w := pretty.Writer{}
	_, err := w.Marshal(simplyPanic(0))
	tt.NotNil(t, err)
}
