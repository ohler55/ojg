// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/tt"
)

type edata struct {
	src     string
	value   any
	noMatch bool
}

func scriptBenchData(size int64) any {
	list := []any{}
	for i := int64(0); i < size; i++ {
		list = append(list, map[string]any{string([]byte{'a' + byte(i%26)}): i, "x": i})
	}
	return list
}

func TestScriptBasicEval(t *testing.T) {
	data := []any{
		map[string]any{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]any{
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

	stack := s.Eval([]any{}, data)
	tt.Equal(t, `[{"a":1,"b":2,"c":3}]`, oj.JSON(stack, &oj.Options{Sort: true}))
}

func TestScriptParse(t *testing.T) {
	for i, d := range []xdata{
		{src: "($.x == 'abc')", expect: "($.x == 'abc')"},
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
		{src: "(@ has true)", expect: "(@ has true)"},
		{src: "(@ exists true)", expect: "(@ exists true)"},
		{src: "(@ =~ /abc/)", expect: "(@ ~= /abc/)"},
		{src: "(@ ~= /a\\/c/)", expect: "(@ ~= /a\\/c/)"},

		{src: "(length(@.xyz))", expect: "(length(@.xyz))"},
		{src: "(3 == length(@.xyz))", expect: "(3 == length(@.xyz))"},
		{src: "(length(@.xyz) == 3)", expect: "(length(@.xyz) == 3)"},
		{src: "(length(@.xyz) == Nothing)", expect: "(length(@.xyz) == Nothing)"},
		{src: "(length(@.xyz == 3)", err: "not terminated at 15 in (length(@.xyz == 3)"},
		{src: "(leng(@.xyz) == 3)", err: "expected a length function at 2 in (leng(@.xyz) == 3)"},

		{src: "(count(@.xyz))", expect: "(count(@.xyz))"},
		{src: "(3 == count(@.xyz))", expect: "(3 == count(@.xyz))"},
		{src: "(count(@.xyz) == 3)", expect: "(count(@.xyz) == 3)"},
		{src: "(count(7) == 3)", expect: "(count(7) == 3)"},
		{src: "(count(@.xyz == 3)", err: "not terminated at 14 in (count(@.xyz == 3)"},
		{src: "(coun(@.xyz) == 3)", err: "expected a count function at 2 in (coun(@.xyz) == 3)"},

		{src: "(match(@.x, 'xy.'))", expect: "(match(@.x, 'xy.'))"},
		{src: "(match(@.x, 'xy.') == false)", expect: "(match(@.x, 'xy.') == false)"},
		{src: "(false == match(@.x, 'xy.'))", expect: "(false == match(@.x, 'xy.'))"},
		{src: "(matc(@.x, 'xy.'))", err: "expected a match function at 2 in (matc(@.x, 'xy.'))"},

		{src: "(search(@.x, 'xy.'))", expect: "(search(@.x, 'xy.'))"},
		{src: "(search(@.x, 'xy.') == false)", expect: "(search(@.x, 'xy.') == false)"},
		{src: "(false == search(@.x, 'xy.'))", expect: "(false == search(@.x, 'xy.'))"},
		{src: "(sear(@.x, 'xy.'))", err: "expected a search function at 2 in (sear(@.x, 'xy.'))"},

		{src: "@.x == 4", expect: "(@.x == 4)"},
		{src: "(@.x ++ 4)", err: "'++' is not a valid operation at 8 in (@.x ++ 4)"},
		{src: "(@[1:5} == 3)", err: "invalid slice syntax at 8 in (@[1:5} == 3)"},
		{src: "(@ =~ /a[c/)", err: "error parsing regexp: missing closing ]: `[c` at 12 in (@ =~ /a[c/)"},
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
	tt.Equal(t, true, s.Match(map[string]any{"x": 3}))
	tt.Equal(t, true, s.Match(gen.Object{"x": gen.Int(3)}))
}

func TestScriptNormalizeEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	for _, v := range []any{
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
		result, _ := s.Eval([]any{}, []any{v}).([]any)
		tt.Equal(t, 1, len(result), fmt.Sprintf("%T %v", v, v))
	}
	s, err = jp.NewScript("(@ == 'x')")
	tt.Nil(t, err)
	result, _ := s.Eval([]any{}, []any{"x"}).([]any)
	tt.Equal(t, 1, len(result), "string normalize")
	result, _ = s.Eval([]any{}, gen.Array{gen.String("x")}).([]any)
	tt.Equal(t, 1, len(result), "gen.String normalize")

	s, err = jp.NewScript("(@ == true)")
	tt.Nil(t, err)
	result, _ = s.Eval([]any{}, gen.Array{gen.True}).([]any)
	tt.Equal(t, 1, len(result), "bool normalize")
}

func TestScriptNonListEval(t *testing.T) {
	s, err := jp.NewScript("(@ == 3)")
	tt.Nil(t, err)
	result, _ := s.Eval([]any{}, "bad").([]any)
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
		{src: "(2 in @)", value: []any{int64(1), int64(2), int64(3)}},

		{src: "(@ empty false)", value: []any{int64(1)}},
		{src: "(@ empty true)", value: []any{}},
		{src: "(@ empty true)", value: map[string]any{}},
		{src: "(@ empty true)", value: ""},
		{src: "(@ empty true)", value: []any{1}, noMatch: true},
		{src: "(@ empty true)", value: map[string]any{"x": 1}, noMatch: true},
		{src: "(@ empty true)", value: "x", noMatch: true},

		{src: "(@ has true)", value: 5},
		{src: "(@ has false)", value: jp.Nothing},

		{src: "(@ exists true)", value: 5},
		{src: "(@.x exists false)", value: map[string]any{}},

		{src: "(@ ~= /a.c/)", value: "abc"},
		{src: "(@ =~ 'a.c')", value: "abc"},
		{src: "(@ ~= 'a.c')", value: "abb", noMatch: true},
		{src: "(@ =~ 'a.c')", value: int64(3), noMatch: true},

		{src: "(@.x || @.y)", value: map[string]any{"x": false, "y": false}, noMatch: true},
		{src: "(@.x || @.y)", value: map[string]any{"x": false, "y": true}},

		{src: "(@.x && @.y)", value: map[string]any{"x": true, "y": false}, noMatch: true},
		{src: "(@.x && @.y)", value: map[string]any{"x": true, "y": true}},

		{src: "(!@.x)", value: map[string]any{"x": true}, noMatch: true},
		{src: "(!@.x)", value: map[string]any{"x": false}},

		{src: "(@.x + @.y == 0)", value: map[string]any{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x + @.y == 3)", value: map[string]any{"x": 1, "y": 2}},
		{src: "(@.x + @.y == 3.1)", value: map[string]any{"x": 1.1, "y": 2}},
		{src: "(@.x + @.y == 3.2)", value: map[string]any{"x": 1, "y": 2.2}},
		{src: "(@.x + @.y == 3.5)", value: map[string]any{"x": 1.2, "y": 2.3}},
		{src: "(@.x + @.y == null)", value: map[string]any{"x": 1.2, "y": "abc"}},
		{src: "(@.x + @.y == null)", value: map[string]any{"x": 1, "y": "abc"}},
		{src: "(@.x + @.y == 'abcdef')", value: map[string]any{"x": "abc", "y": "def"}},
		{src: "(@.x + @.y == null)", value: map[string]any{"x": "abc", "y": nil}},

		{src: "(@.x - @.y == 0)", value: map[string]any{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x - @.y == 2)", value: map[string]any{"x": 3, "y": 1}},
		{src: "(@.x - @.y == 1.1)", value: map[string]any{"x": 3.1, "y": 2}},
		{src: "(@.x - @.y == 0.5)", value: map[string]any{"x": 3, "y": 2.5}},
		{src: "(@.x - @.y == 1.5)", value: map[string]any{"x": 3.5, "y": 2.0}},
		{src: "(@.x - @.y == null)", value: map[string]any{"x": 1.2, "y": "abc"}},
		{src: "(@.x - @.y == null)", value: map[string]any{"x": 1, "y": "abc"}},
		{src: "(@.x-1 == @.y)", value: map[string]any{"x": 1, "y": 0}},
		{src: `(@["x-1"] == @.y)`, value: map[string]any{"x-1": 1, "y": 1}},

		{src: "(@.x * @.y == 0)", value: map[string]any{"x": 1, "y": 2}, noMatch: true},
		{src: "(@.x * @.y == 2)", value: map[string]any{"x": 1, "y": 2}},
		{src: "(@.x * @.y == 2.2)", value: map[string]any{"x": 1.1, "y": 2}},
		{src: "(@.x * @.y == 2.2)", value: map[string]any{"x": 1, "y": 2.2}},
		{src: "(@.x * @.y == 5.0)", value: map[string]any{"x": 2.0, "y": 2.5}},
		{src: "(@.x * @.y == null)", value: map[string]any{"x": 1.2, "y": "abc"}},
		{src: "(@.x * @.y == null)", value: map[string]any{"x": 1, "y": "abc"}},

		{src: "(@.x / @.y == 0)", value: map[string]any{"x": 2, "y": 1}, noMatch: true},
		{src: "(@.x / @.y == 0)", value: map[string]any{"x": 1, "y": 2}},
		{src: "(@.x / @.y == 1.1)", value: map[string]any{"x": 2.2, "y": 2}},
		{src: "(@.x / @.y == 2.5)", value: map[string]any{"x": 5, "y": 2.0}},
		{src: "(@.x / @.y == 2.0)", value: map[string]any{"x": 5.0, "y": 2.5}},
		{src: "(@.x / @.y == null)", value: map[string]any{"x": 1.2, "y": "abc"}},
		{src: "(@.x / @.y == null)", value: map[string]any{"x": 1, "y": "abc"}},
		{src: "(@.x / @.y == null)", value: map[string]any{"x": 1, "y": 0}},

		{src: "($.x + @.y == 0)", value: map[string]any{"x": 1, "y": 2}, noMatch: true},

		{src: "(length(@.x) == 3)", value: map[string]any{"x": []any{1, 2, 3}}},
		{src: "(length(@.x) == 2)", value: map[string]any{"x": []any{1, 2, 3}}, noMatch: true},
		{src: "(length(@.x) == 3)", value: map[string]any{"x": "abc"}},
		{src: "(length(@.x) == 3)", value: map[string]any{"x": map[string]any{"a": 1, "b": 2, "c": 3}}},
		{src: "(length(@.x) == Nothing)", value: map[string]any{"y": "abc"}},

		{src: "(count(@.x[*]) == 3)", value: map[string]any{"x": []any{1, 2, 3}}},
		{src: "(count(7) == 3)", value: map[string]any{"x": []any{1, 2, 3}}, noMatch: true},
		{src: "(count(@.x[*]) == 2)", value: map[string]any{"x": []any{1, 2, 3}}, noMatch: true},
		{src: "(count(@.x) == 1)", value: map[string]any{"x": "abc"}},
		{src: "(count(@.x.*) == 3)", value: map[string]any{"x": map[string]any{"a": 1, "b": 2, "c": 3}}},

		{src: "(match(@.x, 'ab.'))", value: map[string]any{"x": "abc"}},
		{src: "(match(@.x, 'ab'))", value: map[string]any{"x": "abc"}, noMatch: true},

		{src: "(search(@.x, 'ab'))", value: map[string]any{"x": "abc"}},
		{src: "(search(@.x, 'abx'))", value: map[string]any{"x": "abc"}, noMatch: true},
	} {
		if testing.Verbose() {
			if d.value == nil {
				fmt.Printf("... %d: %s in nil\n", i, d.src)
			} else {
				fmt.Printf("... %d: %s in %s\n", i, d.src, oj.JSON(d.value))
			}
		}
		s, err := jp.NewScript(d.src)
		tt.Nil(t, err)
		result, _ := s.Eval([]any{}, []any{d.value}).([]any)
		if d.noMatch {
			tt.Equal(t, 0, len(result), d.src, " in ", d.value)
		} else {
			tt.Equal(t, 1, len(result), d.src, " in ", d.value)
		}
	}
}

func TestScriptInspect(t *testing.T) {
	type idata struct {
		src    string
		expect string
	}
	for i, d := range []idata{
		{src: "(@ == 3)", expect: `{left: @ op: "==" right: 3}`},
		{src: "(3 == @)", expect: `{left: 3 op: "==" right: @}`},
		{src: "(@.x - @.y == 0)", expect: `{left: {left: @.x op: - right: @.y} op: "==" right: 0}`},
		{src: "(0 == @.x - @.y)", expect: `{left: 0 op: "==" right: {left: @.x op: - right: @.y}}`},
		{src: "(!@.x)", expect: `{left: @.x op: "!" right: null}`},
	} {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.src)
		}
		s, err := jp.NewScript(d.src)
		tt.Nil(t, err)
		f := s.Inspect()
		tt.Equal(t, d.expect, pretty.SEN(f))
	}
}

func BenchmarkOjScriptDev(b *testing.B) {
	s := jp.Lt(jp.Get(jp.A().C("a")), jp.ConstInt(52)).Script()
	data := scriptBenchData(100)
	stack := []any{}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		stack = stack[:0]
		stack, _ = s.Eval(stack, data).([]any)
	}
}
