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

type setData struct {
	path     string
	data     string // JSON
	value    any
	expect   string // JSON
	err      string
	noNode   bool
	noSimple bool
}

type setReflectData struct {
	path   string
	data   any
	value  any
	expect string // JSON
	err    string
}

var (
	setTestData = []*setData{
		{path: "@.a", data: `{}`, value: 3, expect: `{"a":3}`},
		{path: "a.b", data: `{}`, value: 3, expect: `{"a":{"b":3}}`},
		{path: "a.b", data: `{"a":{}}`, value: 3, expect: `{"a":{"b":3}}`},
		{path: "[1]", data: `[1,2,3]`, value: 5, expect: `[1,5,3]`},
		{path: "[-1]", data: `[1,2,3]`, value: 5, expect: `[1,2,5]`},
		{path: "[1].a", data: `[1,{},3]`, value: 5, expect: `[1,{"a":5},3]`},
		{path: "[*]", data: `[1,2,3]`, value: 5, expect: `[5,5,5]`},
		{path: ".*", data: `{"a":1,"b":2}`, value: 5, expect: `{"a":5,"b":5}`},
		{path: "$.*.a", data: `{"a":{"a":1,"b":2},"b":{"a":2}}`, value: 5, expect: `{"a":{"a":5,"b":2},"b":{"a":5}}`},
		{path: "[*].a", data: `[{"a":1,"b":2},{"a":2}]`, value: 5, expect: `[{"a":5,"b":2},{"a":5}]`},
		{path: "..a", data: `{"a":{"a":1,"b":2},"b":{"a":2}}`, value: 5, expect: `{"a":5,"b":{"a":5}}`},
		{path: "..a", data: `[{"a":1,"b":2},{"a":2}]`, value: 5, expect: `[{"a":5,"b":2},{"a":5}]`},
		{path: "[-1,'x'].a", data: `[{"a":1,"b":2},{"a":2}]`, value: 5, expect: `[{"a":1,"b":2},{"a":5}]`},
		{path: "[1,'a'].a", data: `{"a":{"a":1,"b":2},"b":{"a":2}}`, value: 5, expect: `{"a":{"a":5,"b":2},"b":{"a":2}}`},
		{path: "[:-1:2].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":5}]`},
		{path: "[-1:0:-2].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":5}]`},
		{path: "[:5].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":5},{"a":5}]`},
		{path: "[?(@.b == 2)].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":3}]`},
		{path: "a[0]", data: `{}`, value: 3, expect: `{"a":[3]}`},
		{path: "*.x", data: `{"a":null}`, value: 3, expect: `{"a":null}`},
		{path: "a", data: `{"a":null}`, value: nil, expect: `{"a":null}`},
		{path: "[*].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "...x", data: "[null]", value: 3, expect: `[null]`},
		{path: "[0,1].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "['a','b'].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "['a','b'].x", data: `{"a":null}`, value: 3, expect: `{"a":null}`},
		{path: "a[1,2]", data: `{"a":[0,1,2,3]}`, value: 5, expect: `{"a":[0,5,5,3]}`},
		{path: "['a','b']", data: `{"a":1,"b":2,"c":3}`, value: 5, expect: `{"a":5,"b":5,"c":3}`},

		{path: "", data: `{}`, value: 3, err: "can not set with an empty expression"},
		{path: "$", data: `{}`, value: 3, err: "can not set with an expression ending with a Root"},
		{path: "@", data: `{}`, value: 3, err: "can not set with an expression ending with a At"},
		{path: "a", data: `{}`, value: func() {}, err: "can not set a func() in a gen.Object", noSimple: true},
		{path: "a.b", data: `{"a":4}`, value: 3, err: "/can not follow a .+ at 'a'/"},
		{path: "a[-1]", data: `{}`, value: 3, err: "can not deduce the length of the array to add at 'a'"},
		{path: "a[1,2].x", data: `{}`, value: 3, err: "can not deduce what element to add at 'a'"},
		{path: "[0].1", data: `[1]`, value: 3, err: "/can not follow a .+ at '\\[0\\]'/"},
		{path: "[1]", data: `[1]`, value: 3, err: "can not follow out of bounds array index at '[1]'"},
	}
	setOneTestData = []*setData{
		{path: "@.a", data: `{}`, value: 3, expect: `{"a":3}`},
		{path: "a.b", data: `{}`, value: 3, expect: `{"a":{"b":3}}`},
		{path: "a.b", data: `{"a":{}}`, value: 3, expect: `{"a":{"b":3}}`},
		{path: "a.b", data: `{"a":{"b":2}}`, value: 3, expect: `{"a":{"b":3}}`},
		{path: "[1]", data: `[1,2,3]`, value: 5, expect: `[1,5,3]`},
		{path: "[-1]", data: `[1,2,3]`, value: 5, expect: `[1,2,5]`},
		{path: "[1].a", data: `[1,{},3]`, value: 5, expect: `[1,{"a":5},3]`},
		{path: "[*]", data: `[1,2,3]`, value: 5, expect: `[5,2,3]`},
		{path: "*", data: `{"a":1}`, value: 5, expect: `{"a":5}`},
		{path: "$.*.a", data: `{"a":{"a":1}}`, value: 5, expect: `{"a":{"a":5}}`},
		{path: "[*].a", data: `[{"a":1,"b":2}]`, value: 5, expect: `[{"a":5,"b":2}]`},
		{path: "..a", data: `{"a":1,"b":2}`, value: 5, expect: `{"a":5,"b":2}`},
		{path: "..a", data: `[{"a":1,"b":2}]`, value: 5, expect: `[{"a":5,"b":2}]`},
		{path: "..b", data: `{"x":{"a":{}}}`, value: 5, expect: `{"x":{"a":{"b":5}}}`},
		{path: "[-1,'x'].a", data: `[{"a":1,"b":2},{"a":2}]`, value: 5, expect: `[{"a":1,"b":2},{"a":5}]`},
		{path: "[1,'a'].a", data: `{"a":{"a":1,"b":2},"b":{"a":2}}`, value: 5, expect: `{"a":{"a":5,"b":2},"b":{"a":2}}`},
		{path: "[:-1:2].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":3}]`},
		{path: "[-1:0:-2].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":1,"b":2},{"a":2},{"a":5}]`},
		{path: "[:5].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":3}]`},
		{path: "[?(@.b == 2)].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":5,"b":2},{"a":2},{"a":3}]`},
		{path: "[-5:3].a", data: `[{"a":1,"b":2},{"a":2},{"a":3}]`, value: 5, expect: `[{"a":1,"b":2},{"a":2},{"a":3}]`},
		{path: "a[0]", data: `{}`, value: 3, expect: `{"a":[3]}`},
		{path: "*.x", data: `{"a":null}`, value: 3, expect: `{"a":null}`},
		{path: "[*].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "...x", data: "[null]", value: 3, expect: `[null]`},
		{path: "[0,1].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "['a','b'].x", data: "[null]", value: 3, expect: `[null]`},
		{path: "['a','b'].x", data: `{"a":null}`, value: 3, expect: `{"a":null}`},
		{path: "a[1,2]", data: `{"a":[0,1,2,3]}`, value: 5, expect: `{"a":[0,5,2,3]}`},
		{path: "['a','b']", data: `{"a":1,"b":2,"c":3}`, value: 5, expect: `{"a":5,"b":2,"c":3}`},

		{path: "", data: `{}`, value: 3, err: "can not set with an empty expression"},
		{path: "$", data: `{}`, value: 3, err: "can not set with an expression ending with a Root"},
		{path: "@", data: `{}`, value: 3, err: "can not set with an expression ending with a At"},
		{path: "a", data: `{}`, value: func() {}, err: "can not set a func() in a gen.Object", noSimple: true},
		{path: "a.b", data: `{"a":4}`, value: 3, err: "/can not follow a .+ at 'a'/"},
		{path: "a[-1]", data: `{}`, value: 3, err: "can not deduce the length of the array to add at 'a'"},
		{path: "a[1,2].x", data: `{}`, value: 3, err: "can not deduce what element to add at 'a'"},
		{path: "[0].1", data: `[1]`, value: 3, err: "/can not follow a .+ at '\\[0\\]'/"},
		{path: "[1]", data: `[1]`, value: 3, err: "can not follow out of bounds array index at '[1]'"},
	}
	setReflectTestData = []*setReflectData{
		{path: "a", data: &Sample{A: 1, B: "a string"}, value: 3, expect: `{"^":"Sample","a":3,"b":"a string"}`},
		{path: "x.a", data: &Any{X: map[string]any{"a": 1}}, value: 3, expect: `{"^":"Any","x":{"a":3}}`},
		{path: "x.a", data: &Any{X: &Sample{A: 1}}, value: 3, expect: `{"^":"Any","x":{"^":"Sample","a":3,"b":""}}`},
		{path: "x.a", data: map[string]any{"x": &Sample{A: 1}}, value: 3, expect: `{"x":{"^":"Sample","a":3,"b":""}}`},
		{path: "[1]", data: []int{1, 2, 3}, value: 5, expect: `[1,5,3]`},
		{path: "[-2]", data: []int{1, 2, 3}, value: 5, expect: `[1,5,3]`},
		{path: "[1].x", data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5, expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "$.*.x", data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5, expect: `[{"^":"Any","x":5},{"^":"Any","x":5},{"^":"Any","x":5}]`},
		{path: "[1].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":1},{"x":5},{"x":3}]`},
		{path: "[*].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":5},{"x":5},{"x":5}]`},
		{path: "[1,'a'].x",
			data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[1,'a'].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":1},{"x":5},{"x":3}]`},
		{path: "[1,'a'].x",
			data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[1,'a'].x",
			data:  map[string]any{"a": &Any{X: 1}, "b": &Any{X: 2}, "c": &Any{X: 3}},
			value: 5, expect: `{"a":{"^":"Any","x":5},"b":{"^":"Any","x":2},"c":{"^":"Any","x":3}}`},
		{path: "[1,'x'].x", data: &Any{X: &Any{X: 1}}, value: 5, expect: `{"^":"Any","x":{"^":"Any","x":5}}`},
		{path: "[1,'x'].x", data: &Any{X: map[string]any{"x": 1}}, value: 5, expect: `{"^":"Any","x":{"x":5}}`},

		{path: "[1].x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[*].x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":5},{"^":"Any","x":5}]`},
		{path: "$.*.x",
			data:  map[string]any{"a": &Any{X: 1}, "b": &Any{X: 2}, "c": &Any{X: 3}},
			value: 5, expect: `{"a":{"^":"Any","x":5},"b":{"^":"Any","x":5},"c":{"^":"Any","x":5}}`},
		{path: "..x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":5},{"^":"Any","x":5}]`},
		{path: "..x",
			data:  map[string]any{"a": &Any{X: 1}, "b": &Any{X: 2}, "c": &Any{X: 3}},
			value: 5, expect: `{"a":{"^":"Any","x":5},"b":{"^":"Any","x":5},"c":{"^":"Any","x":5},"x":5}`},
		{path: "[:2:2].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":5},{"x":2},{"x":3}]`},
		{path: "[:2:2].x",
			data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "*.a", data: &Any{X: 1}, value: 3, expect: `{"^":"Any","x":1}`},
		{path: "[0,1].x", data: []int{1}, value: 3, expect: `[1]`},
		{path: "[0:1].x", data: []int{1}, value: 3, expect: `[1]`},
		{path: "['x','y'].a", data: &Any{X: 1}, value: 3, expect: `{"^":"Any","x":1}`},

		{path: "x.a", data: map[string]any{"x": func() {}}, value: 3, err: "can not follow a func() at 'x'"},
		{path: "x.a", data: &Any{X: 1}, value: 3, err: "can not follow a int at 'x'"},
		{path: "x.a", data: &Any{X: func() {}}, value: 3, err: "can not follow a func() at 'x'"},
		{path: "[0].x", data: []any{func() {}}, value: 5, err: "can not follow a func() at '[0]'"},
		{path: "[0].x", data: []int{1, 2, 3}, value: 5, err: "can not follow a int at '[0]'"},
		{path: "[0].x", data: []func(){func() {}}, value: 5, err: "can not follow a func() at '[0]'"},
	}
	setOneReflectTestData = []*setReflectData{
		{path: "a", data: &Sample{A: 1, B: "a string"}, value: 3, expect: `{"^":"Sample","a":3,"b":"a string"}`},
		{path: "x.a", data: &Any{X: map[string]any{"a": 1}}, value: 3, expect: `{"^":"Any","x":{"a":3}}`},
		{path: "x.a", data: &Any{X: &Sample{A: 1}}, value: 3, expect: `{"^":"Any","x":{"^":"Sample","a":3,"b":""}}`},
		{path: "x.a", data: map[string]any{"x": &Sample{A: 1}}, value: 3, expect: `{"x":{"^":"Sample","a":3,"b":""}}`},
		{path: "[1]", data: []int{1, 2, 3}, value: 5, expect: `[1,5,3]`},
		{path: "[-2]", data: []int{1, 2, 3}, value: 5, expect: `[1,5,3]`},
		{path: "[1].x", data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5, expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "$.*.x", data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5, expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "[1].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":1},{"x":5},{"x":3}]`},
		{path: "[*].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":5},{"x":2},{"x":3}]`},
		{path: "[1,'a'].x",
			data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[1,'a']",
			data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: &Any{X: 5},
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[1,'a'].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":1},{"x":5},{"x":3}]`},
		{path: "[1,'a'].x",
			data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[1,'a'].x",
			data:  map[string]any{"a": &Any{X: 1}, "b": &Any{X: 2}, "c": &Any{X: 3}},
			value: 5, expect: `{"a":{"^":"Any","x":5},"b":{"^":"Any","x":2},"c":{"^":"Any","x":3}}`},
		{path: "[1,'x'].x", data: &Any{X: &Any{X: 1}}, value: 5, expect: `{"^":"Any","x":{"^":"Any","x":5}}`},
		{path: "[1,'x'].x", data: &Any{X: map[string]any{"x": 1}}, value: 5, expect: `{"^":"Any","x":{"x":5}}`},
		{path: "[1,'x']", data: &Any{X: 1}, value: 5, expect: `{"^":"Any","x":5}`},

		{path: "[1].x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":5},{"^":"Any","x":3}]`},
		{path: "[*].x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "$.*.x",
			data:  map[string]any{"a": &Any{X: 1}},
			value: 5, expect: `{"a":{"^":"Any","x":5}}`},
		{path: "..x", data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "..x",
			data:  map[string]any{"a": &Any{X: 1}},
			value: 5, expect: `{"a":{"^":"Any","x":5}}`},
		{path: "[:2:2].x",
			data: []map[string]any{
				{"x": 1},
				{"x": 2},
				{"x": 3},
			},
			value: 5, expect: `[{"x":5},{"x":2},{"x":3}]`},
		{path: "[:2:2].x",
			data: []*Any{{X: 1}, {X: 2}, {X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "[:2:2].x",
			data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":5},{"^":"Any","x":2},{"^":"Any","x":3}]`},
		{path: "[2:0:-2].x",
			data: []any{&Any{X: 1}, &Any{X: 2}, &Any{X: 3}}, value: 5,
			expect: `[{"^":"Any","x":1},{"^":"Any","x":2},{"^":"Any","x":5}]`},
		{path: "*.a", data: &Any{X: 1}, value: 3, expect: `{"^":"Any","x":1}`},
		{path: "['x','y'].a", data: &Any{X: 1}, value: 3, expect: `{"^":"Any","x":1}`},
		{path: "[0,1].x", data: []int{1}, value: 3, expect: `[1]`},
		{path: "[0:1].x", data: []int{1}, value: 3, expect: `[1]`},
		{path: "[0:1].x", data: []any{1}, value: 3, expect: `[1]`},
		{path: "[1:0:-1].x", data: []any{1, 2}, value: 3, expect: `[1,2]`},

		{path: "x.a", data: map[string]any{"x": func() {}}, value: 3, err: "can not follow a func() at 'x'"},
		{path: "x.a", data: &Any{X: 1}, value: 3, err: "can not follow a int at 'x'"},
		{path: "x.a", data: &Any{X: func() {}}, value: 3, err: "can not follow a func() at 'x'"},
		{path: "[0].x", data: []any{func() {}}, value: 5, err: "can not follow a func() at '[0]'"},
		{path: "[0].x", data: []int{1, 2, 3}, value: 5, err: "can not follow a int at '[0]'"},
		{path: "[0].x", data: []func(){func() {}}, value: 5, err: "can not follow a func() at '[0]'"},
	}
)

func TestExprSet(t *testing.T) {
	for i, d := range setTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		var data any
		if !d.noSimple {
			data, err = oj.ParseString(d.data)
			tt.Nil(t, err, i, " : ", x)
			err = x.Set(data, d.value)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := oj.JSON(data, &oj.Options{Sort: true})
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
		if !d.noNode {
			var p gen.Parser
			data, err = p.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			err = x.Set(data, d.value)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := oj.JSON(data, &oj.Options{Sort: true})
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
	}
}

func TestExprSetOne(t *testing.T) {
	for i, d := range setOneTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		var data any
		if !d.noSimple {
			data, err = oj.ParseString(d.data)
			tt.Nil(t, err, i, " : ", x)
			err = x.SetOne(data, d.value)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := oj.JSON(data, &oj.Options{Sort: true})
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
		if !d.noNode {
			var p gen.Parser
			data, err = p.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			err = x.SetOne(data, d.value)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := oj.JSON(data, &oj.Options{Sort: true})
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
	}
}

func TestExprSetReflect(t *testing.T) {
	for i, d := range setReflectTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		err = x.Set(d.data, d.value)
		if 0 < len(d.err) {
			tt.NotNil(t, err, i, " : ", x)
			tt.Equal(t, d.err, err.Error(), i, " : ", x)
		} else {
			result := oj.JSON(d.data, &oj.Options{Sort: true, CreateKey: "^"})
			tt.Equal(t, d.expect, result, i, " : ", x)
		}
	}
}

func TestExprSetOneReflect(t *testing.T) {
	for i, d := range setOneReflectTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		err = x.SetOne(d.data, d.value)
		if 0 < len(d.err) {
			tt.NotNil(t, err, i, " : ", x)
			tt.Equal(t, d.err, err.Error(), i, " : ", x)
		} else {
			result := oj.JSON(d.data, &oj.Options{Sort: true, CreateKey: "^"})
			tt.Equal(t, d.expect, result, i, " : ", x)
		}
	}
}

func TestExprMustSet(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	tt.Panic(t, func() { jp.C("b").N(0).MustSet(data, 7) })
}

func TestExprMustSetOne(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	tt.Panic(t, func() { jp.C("b").N(0).MustSetOne(data, 7) })
}

func flatKeyed() any {
	return &keyed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{key: "c", value: 3},
			},
		},
	}
}

func deepKeyed() any {
	return &keyed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{
					key: "c",
					value: &keyed{
						ordered: ordered{
							entries: []*entry{
								{key: "c1", value: 11},
								{key: "c2", value: 12},
								{key: "c3", value: 13},
							},
						},
					},
				},
			},
		},
	}
}

func flatIndexed() any {
	return &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{key: "c", value: 3},
			},
		},
	}
}

func deepIndexed() any {
	return &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: 1},
				{key: "b", value: 2},
				{
					key: "c",
					value: &indexed{
						ordered: ordered{
							entries: []*entry{
								{key: "c1", value: 11},
								{key: "c2", value: 12},
								{key: "c3", value: 13},
							},
						},
					},
				},
			},
		},
	}
}

func TestSetKeyedIndexed(t *testing.T) {
	for _, d := range []*struct {
		src   string
		data  any
		after string
		value any
		del   bool
		err   bool
		one   bool
	}{
		{src: "$.b", value: 5, data: flatKeyed(), after: "[{key: a value: 1} {key: b value: 5} {key: c value: 3}]"},
		{src: "$.b", value: 5, data: flatKeyed(), after: "[{key: a value: 1} {key: c value: 3}]", del: true},
		{src: "$.b", value: 5, data: flatKeyed(), after: "[{key: a value: 1} {key: c value: 3}]", del: true, one: true},
		{
			src:   "$.c.c2",
			value: 5,
			data:  deepKeyed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
		},
		{
			src:   "$.b.c",
			value: 5,
			data:  &keyed{},
			after: `[
  {key: b value: {c: 5}}
]`,
			one: true,
		},
		{
			src:   "$.b[1]",
			value: 5,
			data:  &keyed{},
			after: `[
  {key: b value: [null 5]}
]`,
			one: true,
		},
		{src: "$.b[-1]", value: 5, data: &keyed{}, err: true},
		{src: "$.b.*", value: 5, data: &keyed{}, err: true},
		{src: "$.b.c2", value: 5, data: flatKeyed(), err: true},

		{src: "$[1]", value: 5, data: flatIndexed(), after: "[{key: a value: 1} {key: b value: 5} {key: c value: 3}]"},
		{src: "$[-2]", value: 5, data: flatIndexed(), after: "[{key: a value: 1} {key: b value: null} {key: c value: 3}]", del: true, one: true},
		{
			src:   "$[2][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`},
		{src: "$[4]", value: 5, data: flatIndexed(), err: true},
		{src: "$[1][1]", value: 5, data: flatIndexed(), err: true},

		{src: "$.*", value: 5, data: flatKeyed(), after: "[{key: a value: 5} {key: b value: 5} {key: c value: 5}]"},
		{
			src:   "$.*.c2",
			value: 5,
			data:  deepKeyed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
		},
		{src: "$.*", value: 5, data: flatKeyed(), after: `[{key: b value: 2} {key: c value: 3}]`, del: true, one: true},
		{
			src:   "$.*",
			value: 5,
			data:  flatKeyed(),
			after: `[{key: a value: 5} {key: b value: 2} {key: c value: 3}]`,
			one:   true,
		},
		{src: "$.*", data: flatIndexed(), after: `[{key: a value: null} {key: b value: 2} {key: c value: 3}]`, del: true, one: true},
		{src: "$.*", value: 5, data: flatIndexed(), after: `[{key: a value: 5} {key: b value: 2} {key: c value: 3}]`, one: true},
		{
			src:   "$.*[1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
			one: true,
		},
		{
			src:   "$..x",
			value: 5,
			data:  deepKeyed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [
      {key: c1 value: 11}
      {key: c2 value: 12}
      {key: c3 value: 13}
      {key: x value: 5}
    ]
  }
  {key: x value: 5}
]`,
		},
		{
			src:   "$..[1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 5}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
		},

		{
			src:   "$['a','c'].x",
			value: 5,
			data:  deepKeyed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [
      {key: c1 value: 11}
      {key: c2 value: 12}
      {key: c3 value: 13}
      {key: x value: 5}
    ]
  }
]`,
		},
		{
			src:   "$['a','b']",
			value: 5,
			data:  flatKeyed(),
			after: `[{key: a value: 5} {key: b value: 2} {key: c value: 3}]`,
			one:   true,
		},
		{src: "$['a','b']", data: flatKeyed(), after: `[{key: c value: 3}]`, del: true},
		{
			src:   "$[0,2][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
		},
		{
			src:   "$[0,2]",
			value: 5,
			data:  flatIndexed(),
			after: `[{key: a value: 5} {key: b value: 2} {key: c value: 3}]`,
			one:   true,
		},
		{
			src:   "$[0,-1]",
			data:  flatIndexed(),
			after: `[{key: a value: null} {key: b value: 2} {key: c value: null}]`,
			del:   true,
		},
		{
			src:   "$[0:-1][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
			one: true,
		},
		{
			src:   "$[2:0:-1][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
			one: true,
		},
		{
			src:   "$[-3:4][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 5} {key: c3 value: 13}]
  }
]`,
			one: true,
		},
		{
			src:   "$[0:1:0][1]",
			value: 5,
			data:  deepIndexed(),
			after: `[
  {key: a value: 1}
  {key: b value: 2}
  {
    key: c
    value: [{key: c1 value: 11} {key: c2 value: 12} {key: c3 value: 13}]
  }
]`,
			one: true,
		},
	} {
		x := jp.MustParseString(d.src)
		var (
			err error
		)
		if d.one {
			if d.del {
				err = x.DelOne(d.data)
			} else {
				err = x.SetOne(d.data, d.value)
			}
		} else {
			if d.del {
				err = x.Del(d.data)
			} else {
				err = x.Set(d.data, d.value)
			}
		}
		if d.err {
			tt.NotNil(t, err, "%s del: %t one: %t", d.src, d.del, d.one)
		} else {
			tt.Nil(t, err, "%s del: %t one: %t", d.src, d.del, d.one)
			tt.Equal(t, d.after, pretty.SEN(d.data), "%s del: %t one: %t", d.src, d.del, d.one)
		}
	}
}

func flatReflectKeyed() any {
	return &keyed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: &Any{X: 1}},
				{key: "b", value: &Any{X: 2}},
				{key: "c", value: &Any{X: 3}},
			},
		},
	}
}

func flatReflectIndexed() any {
	return &indexed{
		ordered: ordered{
			entries: []*entry{
				{key: "a", value: &Any{X: 1}},
				{key: "b", value: &Any{X: 2}},
				{key: "c", value: &Any{X: 3}},
			},
		},
	}
}

func TestSetKeyedIndexedReflect(t *testing.T) {
	for _, d := range []*struct {
		src   string
		data  any
		after string
		value any
		err   bool
	}{
		{
			src:   "$.b.x",
			value: 5,
			data:  flatReflectKeyed(),
			after: `[
  {key: a value: {type: Any x: 1}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
		{src: "$.a.x", value: 5, data: &keyed{ordered: ordered{entries: []*entry{{key: "a", value: func() {}}}}}, err: true},
		{
			src:   "$[1].x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 1}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
		{src: "$[0].x", value: 5, data: &indexed{ordered: ordered{entries: []*entry{{key: "a", value: func() {}}}}}, err: true},
		{
			src:   "$.*.x",
			value: 5,
			data:  flatReflectKeyed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 5}}
]`,
		},
		{
			src:   "$.*.x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 5}}
]`,
		},
		{
			src:   "$..x",
			value: 5,
			data:  flatReflectKeyed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 5}}
  {key: x value: 5}
]`,
		},
		{
			src:   "$..x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 5}}
]`,
		},
		{
			src:   "$['a','b'].x",
			value: 5,
			data:  flatReflectKeyed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
		{
			src:   "$[0,1].x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
		{
			src:   "$[0:1].x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
		{
			src:   "$[1:0:-1].x",
			value: 5,
			data:  flatReflectIndexed(),
			after: `[
  {key: a value: {type: Any x: 5}}
  {key: b value: {type: Any x: 5}}
  {key: c value: {type: Any x: 3}}
]`,
		},
	} {
		x := jp.MustParseString(d.src)
		err := x.Set(d.data, d.value)
		if d.err {
			tt.NotNil(t, err, "%s", d.src)
		} else {
			tt.Nil(t, err, "%s", d.src)
			tt.Equal(t, d.after, pretty.SEN(d.data), d.src)
		}
	}
}

func TestSetMapPtrNil(t *testing.T) {
	data := map[string]any{
		"a": nil,
	}
	_ = jp.C("a").Set(&data, nil)
	tt.Nil(t, data["a"])
}

func TestSetStructPtrNil(t *testing.T) {
	type A struct {
		B any
	}
	data := A{B: 1}
	_ = jp.C("b").Set(&data, nil)
	tt.Nil(t, data.B)
}
