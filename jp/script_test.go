// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

type edata struct {
	src     string
	value   interface{}
	noMatch bool
}

func scriptBenchData(size int64) interface{} {
	list := []interface{}{}
	for i := int64(0); i < size; i++ {
		list = append(list, map[string]interface{}{string([]byte{'a' + byte(i%26)}): i, "x": i})
	}
	return list
}

func TestScriptBasicEval(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]interface{}{
			"a": int64(52),
			"b": 4,
			"c": 6,
		},
	}
	e := jp.Or(
		jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)),
		jp.Eq(jp.Get(jp.A().C("x")), jp.ConstString("cool")),
	)
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", e.String())
	s := e.Script()
	tt.Equal(t, "(@.a < 52 || @.x == 'cool')", s.String())
	f := e.Filter()
	tt.Equal(t, "[?(@.a < 52 || @.x == 'cool')]", f.String())

	stack := s.Eval([]interface{}{}, data)
	tt.Equal(t, `[{"a":1,"b":2,"c":3}]`, oj.JSON(stack, &oj.Options{Sort: true}))
}

func TestScriptParse(t *testing.T) {
	for i, d := range []xdata{
		{src: "(@.x == 'abc')", expect: "(@.x == 'abc')"},
		{src: "(@.x<5)", expect: "(@.x < 5)"},
		{src: "(@.x<123)", expect: "(@.x < 123)"},
		{src: "(@.x == 3)", expect: "(@.x == 3)"},
		{src: "(@.*.xyz==true)", expect: "(@.*.xyz == true)"},
		{src: "(@.x.* == 3)", expect: "(@.x.* == 3)"},
		{src: "(@.. == 3)", expect: "(@.. == 3)"},
		{src: "(@[3] == 3)", expect: "(@[3] == 3)"},
		{src: "(@[3,4] == 3)", expect: "(@[3,4] == 3)"},
		{src: "(@[3,'four'] == 3)", expect: "(@[3,'four'] == 3)"},
		{src: "(@[1:5] == 3)", expect: "(@[1:5] == 3)"},
		{src: "(@ == 3)", expect: "(@ == 3)"},
		{src: "($ == 3.4e-3)", expect: "($ == 0.0034)"},
		{src: "($ == 3e-3)", expect: "($ == 0.003)"},
		{src: "(3 == @.x)", expect: "(3 == @.x)"},
		{src: "(@ == $)", expect: "(@ == $)"},
		{src: "(@.x[?(@.a == true)].b == false)", expect: "(@.x[?(@.a == true)].b == false)"},
		{src: "(@.x[?(@.a == 5)] == 11)", expect: "(@.x[?(@.a == 5)] == 11)"},
		{src: "((@.x == 3) || (@.y > 5))", expect: "(@.x == 3 || @.y > 5)"},
		{src: "(@.x < 3 && @.x > 1 || @.z == 3)", expect: "(@.x < 3 && @.x > 1 || @.z == 3)"},
		{src: "(!(3 == @.x))", expect: "(!(3 == @.x))"},
		{src: "(@.x in [1,2,3])", expect: "(@.x in [1,2,3])"},
		{src: "(@.x in ['a' , 'b', 'c'])", expect: "(@.x in ['a','b','c'])"},
		{src: "(@ empty true)", expect: "(@ empty true)"},

		{src: "@.x == 4", err: "a script must start with a '('"},
		{src: "(@.x ++ 4)", err: "'++' is not a valid operation at 8 in (@.x ++ 4)"},
		{src: "(@[1:5} == 3)", err: "invalid slice syntax at 8 in (@[1:5} == 3)"},
	} {
		if testing.Verbose() {
			fmt.Printf("... %s\n", d.src)
		}
		s, err := jp.NewScript(d.src)
		if 0 < len(d.err) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.err, err.Error(), i, ": ", d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.expect, s.String(), i, ": ", d.src)
		}
	}
}

func TestScriptMatch(t *testing.T) {
	s, err := jp.NewScript("(@.x == 3)")
	tt.Nil(t, err)
	tt.Equal(t, true, s.Match(map[string]interface{}{"x": 3}))
	tt.Equal(t, true, s.Match(gen.Object{"x": gen.Int(3)}))
}

func TestScriptNormalizeEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	for _, v := range []interface{}{
		int8(3),
		int16(3),
		int32(3),
		int64(3),
		uint(3),
		uint8(3),
		uint16(3),
		uint32(3),
		uint64(3),
		float32(3.0),
		3.0,
		gen.Int(3),
		gen.Float(3.0),
	} {
		result, _ := s.Eval([]interface{}{}, []interface{}{v}).([]interface{})
		tt.Equal(t, 1, len(result), fmt.Sprintf("%T %v", v, v))
	}
	s, err = jp.NewScript("(@ == 'x')")
	tt.Nil(t, err)
	result, _ := s.Eval([]interface{}{}, []interface{}{"x"}).([]interface{})
	tt.Equal(t, 1, len(result), "string normalize")
	result, _ = s.Eval([]interface{}{}, gen.Array{gen.String("x")}).([]interface{})
	tt.Equal(t, 1, len(result), "gen.String normalize")

	s, err = jp.NewScript("(@ == true)")
	tt.Nil(t, err)
	result, _ = s.Eval([]interface{}{}, gen.Array{gen.True}).([]interface{})
	tt.Equal(t, 1, len(result), "bool normalize")
}

func TestScriptNonListEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	result, _ := s.Eval([]interface{}{}, "bad").([]interface{})
	tt.Equal(t, 0, len(result))
}

func TestScriptEval(t *testing.T) {
	for i, d := range []edata{
		{src: "(@ == 3)", value: int64(3)},
		{src: "(@ == 3)", value: 3.0},
		{src: "(@ == 3.0)", value: int64(3)},
		{src: "(@ == 'abc')", value: "abc"},

		{src: "(@ != 3)", value: int64(3), noMatch: true},
		{src: "(@ != 3)", value: int64(4)},
		{src: "(@ != 3)", value: 3.1},
		{src: "(@ != 3.0)", value: int64(4)},
		{src: "(@ != 'abc')", value: "xyz"},

		{src: "(@ < 3)", value: int64(3), noMatch: true},
		{src: "(@ < 3)", value: int64(2)},
		{src: "(@ < 3)", value: 2.1},
		{src: "(@ < 3.0)", value: int64(2)},
		{src: "(@ < 3.0)", value: 2.9},
		{src: "(@ < 'abc')", value: "aaa"},

		{src: "(@ <= 3)", value: int64(4), noMatch: true},
		{src: "(@ <= 3)", value: int64(3)},
		{src: "(@ <= 3)", value: 2.1},
		{src: "(@ <= 3.0)", value: int64(2)},
		{src: "(@ <= 3.0)", value: 2.9},
		{src: "(@ <= 'abc')", value: "abc"},

		{src: "(@ > 3)", value: int64(3), noMatch: true},
		{src: "(@ > 3)", value: int64(4)},
		{src: "(@ > 3)", value: 3.1},
		{src: "(@ > 3.0)", value: int64(4)},
		{src: "(@ > 3.0)", value: 3.1},
		{src: "(@ > 'abc')", value: "abd"},

		{src: "(@ >= 3)", value: int64(2), noMatch: true},
		{src: "(@ >= 3)", value: int64(3)},
		{src: "(@ >= 3)", value: 3.1},
		{src: "(@ >= 3.0)", value: int64(4)},
		{src: "(@ >= 3.0)", value: 3.1},
		{src: "(@ >= 'abc')", value: "abd"},

		{src: "(@ in [1,2,3])", value: int64(2)},
		{src: "(@ in ['a','b','c'])", value: "b"},
		{src: "(2 in @)", value: []interface{}{int64(1), int64(2), int64(3)}},

		{src: "(@ empty false)", value: []interface{}{int64(1)}},
		{src: "(@ empty true)", value: []interface{}{}},
		{src: "(@ empty true)", value: map[string]interface{}{}},
		{src: "(@ empty true)", value: ""},
		{src: "(@ empty true)", value: []interface{}{1}, noMatch: true},
		{src: "(@ empty true)", value: map[string]interface{}{"x": 1}, noMatch: true},
		{src: "(@ empty true)", value: "x", noMatch: true},

		{src: "(@.x || @.y)", value: map[string]interface{}{"x": false, "y": false}, noMatch: true},
		{src: "(@.x || @.y)", value: map[string]interface{}{"x": false, "y": true}},

		{src: "(@.x && @.y)", value: map[string]interface{}{"x": true, "y": false}, noMatch: true},
		{src: "(@.x && @.y)", value: map[string]interface{}{"x": true, "y": true}},

		{src: "(!@.x)", value: map[string]interface{}{"x": true}, noMatch: true},
		{src: "(!@.x)", value: map[string]interface{}{"x": false}},

		{src: "(@.x + @.y == 0)", value: map[string]interface{}{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x + @.y == 3)", value: map[string]interface{}{"x": 1, "y": 2}},
		{src: "(@.x + @.y == 3.1)", value: map[string]interface{}{"x": 1.1, "y": 2}},
		{src: "(@.x + @.y == 3.2)", value: map[string]interface{}{"x": 1, "y": 2.2}},
		{src: "(@.x + @.y == 3.5)", value: map[string]interface{}{"x": 1.2, "y": 2.3}},
		{src: "(@.x + @.y == null)", value: map[string]interface{}{"x": 1.2, "y": "abc"}},
		{src: "(@.x + @.y == null)", value: map[string]interface{}{"x": 1, "y": "abc"}},
		{src: "(@.x + @.y == 'abcdef')", value: map[string]interface{}{"x": "abc", "y": "def"}},
		{src: "(@.x + @.y == null)", value: map[string]interface{}{"x": "abc", "y": nil}},

		{src: "(@.x - @.y == 0)", value: map[string]interface{}{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x - @.y == 2)", value: map[string]interface{}{"x": 3, "y": 1}},
		{src: "(@.x - @.y == 1.1)", value: map[string]interface{}{"x": 3.1, "y": 2}},
		{src: "(@.x - @.y == 0.5)", value: map[string]interface{}{"x": 3, "y": 2.5}},
		{src: "(@.x - @.y == 1.5)", value: map[string]interface{}{"x": 3.5, "y": 2.0}},
		{src: "(@.x - @.y == null)", value: map[string]interface{}{"x": 1.2, "y": "abc"}},
		{src: "(@.x - @.y == null)", value: map[string]interface{}{"x": 1, "y": "abc"}},

		{src: "(@.x * @.y == 0)", value: map[string]interface{}{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x * @.y == 2)", value: map[string]interface{}{"x": 1, "y": 2}},
		{src: "(@.x * @.y == 2.2)", value: map[string]interface{}{"x": 1.1, "y": 2}},
		{src: "(@.x * @.y == 2.2)", value: map[string]interface{}{"x": 1, "y": 2.2}},
		{src: "(@.x * @.y == 5.0)", value: map[string]interface{}{"x": 2.0, "y": 2.5}},
		{src: "(@.x * @.y == null)", value: map[string]interface{}{"x": 1.2, "y": "abc"}},
		{src: "(@.x * @.y == null)", value: map[string]interface{}{"x": 1, "y": "abc"}},

		{src: "(@.x / @.y == 0)", value: map[string]interface{}{"x": 2, "y": 1}, noMatch: true},
		{src: "(@.x / @.y == 0)", value: map[string]interface{}{"x": 1, "y": 2}},
		{src: "(@.x / @.y == 1.1)", value: map[string]interface{}{"x": 2.2, "y": 2}},
		{src: "(@.x / @.y == 2.5)", value: map[string]interface{}{"x": 5, "y": 2.0}},
		{src: "(@.x / @.y == 2.0)", value: map[string]interface{}{"x": 5.0, "y": 2.5}},
		{src: "(@.x / @.y == null)", value: map[string]interface{}{"x": 1.2, "y": "abc"}},
		{src: "(@.x / @.y == null)", value: map[string]interface{}{"x": 1, "y": "abc"}},
		{src: "(@.x / @.y == null)", value: map[string]interface{}{"x": 1, "y": 0}},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s in %s\n", i, d.src, oj.JSON(d.value))
		}
		s, err := jp.NewScript(d.src)
		tt.Nil(t, err)
		result, _ := s.Eval([]interface{}{}, []interface{}{d.value}).([]interface{})
		if d.noMatch {
			tt.Equal(t, 0, len(result), d.src, " in ", d.value)
		} else {
			tt.Equal(t, 1, len(result), d.src, " in ", d.value)
		}
	}
}

func BenchmarkOjScriptDev(b *testing.B) {
	s := jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)).Script()
	data := scriptBenchData(100)
	stack := []interface{}{}
	b.ReportAllocs()
	b.ResetTimer()
	//fmt.Printf("*** data: %s\n", jp.JSON(data))
	for n := 0; n < b.N; n++ {
		stack = stack[:0]
		stack, _ = s.Eval(stack, data).([]interface{})
		//fmt.Printf("*** stack: %s\n", jp.JSON(stack))
	}
}
