// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

type RemObj struct {
	Field any
}

var (
	pw = pretty.Writer{
		Options:  ojg.Options{Sort: true},
		Width:    80,
		MaxDepth: 5,
		SEN:      true,
	}

	remTestData = []*delData{
		{path: "key[2]", data: `{key:[1,2,3,4]}`, expect: `{key: [1 2 4]}`},
		{path: "@.key[2]", data: `{key:[1,2,3,4]}`, expect: `{key: [1 2 4]}`},
		{path: "$[1][2]", data: `[1 [1,2,3,4]]`, expect: `[1 [1 2 4]]`},
		{path: "$[1]", data: `[1,2,3,4]`, expect: `[1 3 4]`},
		{path: "@[1]", data: `[1,2,3,4]`, expect: `[1 3 4]`},
		{path: "key", data: `{key:[1,2,3,4]}`, expect: `{}`},
		{path: "key[*]", data: `{key:[1,2,3,4]}`, expect: `{key: []}`},
		{path: "key.gee", data: `{key:{gee:3}}`, expect: `{key: {}}`},
		{path: "*[0]", data: `{one:[1],two:[1,2]}`, expect: `{one: [] two: [2]}`},
		{path: "*[0]", data: `[[1],[1,2]]`, expect: `[[] [2]]`},
		{path: "key[0]", data: `[[1],[1,2]]`, expect: `[[1] [1 2]]`},
		{path: "one.two[2]", data: `{one:{two:[1,2,3,4]}}`, expect: `{one: {two: [1 2 4]}}`},
		{path: "one.two[2]", data: `{one:{two:2}}`, expect: `{one: {two: 2}}`},
		{path: "key[-2]", data: `{key:[1,2,3,4]}`, expect: `{key: [1 2 4]}`},
		{path: "[-1][-2]", data: `[1,2,[1,2,3,4]]`, expect: `[1 2 [1 2 4]]`},
		{path: "[0][-1][-2]", data: `[[1,2,[1,2,3,4]]]`, expect: `[[1 2 [1 2 4]]]`},
		{path: "*.two[*]", data: `{one:{two:[1,2,3,4]}}`, expect: `{one: {two: []}}`},
		{path: "*.two[*]", data: `{one:1}`, expect: `{one: 1}`},
		{path: "*.two.*", data: `{one:{two:{x:1 y:2}}}`, expect: `{one: {two: {}}}`},
		{path: "*[*][*]", data: `[[1,2,[1,2,3,4]]]`, expect: `[[1 2 []]]`},
		{path: "['a','b']['x','y'][1]", data: `{a:[] b:{x:[1,2,3]}}`, expect: `{a: [] b: {x: [1 3]}}`},
		{path: "[0,1][0,-1][1]", data: `[[[][1,2,3]]]`, expect: `[[[] [1 3]]]`},
		{path: "[1:3:2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 3] [4 5] [6] [8 9]]"},
		{path: "[3:1:-2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 3] [4 5] [6] [8 9]]"},
		{path: "[-4:-2:2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 3] [4 5] [6] [8 9]]"},
		{path: "[-6:-2:2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 2 3] [4 5] [6 7] [8 9]]"},
		{path: "[:3][1:3:2][1]", data: `[[[][1,2,3][4,5][6,7][8,9]]]`, expect: "[[[] [1 3] [4 5] [6] [8 9]]]"},
		{path: "[-1:0:-1][1:3:2][1]", data: `[[[][1,2,3][4,5][6,7][8,9]]]`, expect: "[[[] [1 3] [4 5] [6] [8 9]]]"},
		{path: "[?(@.x == 1)].y", data: `[{x:1 y:2}]`, expect: `[{x: 1}]`},
		{path: "[?(@[0] != 0)][?(@.x == 1)].y", data: `[[{x:1 y:2}]]`, expect: `[[{x: 1}]]`},
		{path: "['a','b']", data: `{a:1 b:2 c:3}`, expect: `{c: 3}`},
		{path: "[1,2]", data: `[1,2,3]`, expect: `[1]`},
		// {path: "[?(@.x == 1)]", data: `[{x:1}{y:2}{x:3}]`, expect: `[{x: 1}]`},
	}
	remOneTestData = []*delData{
		{path: "key[2]", data: "{key:[1,2,3,4]}", expect: "{key: [1 2 4]}"},
		{path: "*[1]", data: "[[0,1,2][3,2,1]]", expect: "[[0 2] [3 2 1]]"},
		{path: "*[*]", data: "[[0,1,2][3,2,1]]", expect: "[[1 2] [3 2 1]]"},
		{path: "@[*][1]", data: "[[0,1,2][3,2,1]]", expect: "[[0 2] [3 2 1]]"},
		{path: "*[*]", data: "{one:[1,2]}", expect: "{one: [2]}"},
		{path: "*.*", data: "{one:{two: 2 three: 3}}", expect: "{one: {two: 2}}"},
		{path: "[1:3:2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 3] [4 5] [6 7] [8 9]]"},
		{path: "[3:1:-2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 2 3] [4 5] [6] [8 9]]"},
		{path: "[-4:-2:2][1]", data: `[[][1,2,3][4,5][6,7][8,9]]`, expect: "[[] [1 3] [4 5] [6 7] [8 9]]"},
		{path: "@[2]", data: `[1,2,3,4]`, expect: `[1 2 4]`},
		{path: "$[2]", data: `[1,2,3,4]`, expect: `[1 2 4]`},
		{path: "$[1][2]", data: `[1 [1,2,3,4]]`, expect: `[1 [1 2 4]]`},
		{path: "['a','b']['x','y'][1]", data: `{a:[] b:{x:[1,2,3]}}`, expect: `{a: [] b: {x: [1 3]}}`},
		{path: "[0,1][0,-1][1]", data: `[[[][1,2,3]]]`, expect: `[[[] [1 3]]]`},
		{path: "[?(@.x == 1)].y", data: `[{x:1 y:2}]`, expect: `[{x: 1}]`},
		{path: "['a','b']", data: `{a:1 b:2 c:3}`, expect: `{b: 2 c: 3}`},
		{path: "[1,2]", data: `[1,2,3]`, expect: `[1 3]`},
	}
)

func TestExprRemoveAll(t *testing.T) {
	for i, d := range remTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		var data any
		var out any
		if !d.noSimple {
			data, err = sen.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			out, err = x.Remove(data)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := string(pw.Encode(out))
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
		if !d.noNode {
			data, err = sen.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			data = alt.Generify(data)
			out = x.MustRemove(data)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := string(pw.Encode(out))
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
	}
}

func TestExprRemoveOne(t *testing.T) {
	for i, d := range remOneTestData {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err, i, " : ", x)

		var data any
		var out any
		if !d.noSimple {
			data, err = sen.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			out, err = x.RemoveOne(data)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := string(pw.Encode(out))
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
		if !d.noNode {
			data, err = sen.Parse([]byte(d.data))
			tt.Nil(t, err, i, " : ", x)
			data = alt.Generify(data)
			out = x.MustRemoveOne(data)
			if 0 < len(d.err) {
				tt.NotNil(t, err, i, " : ", x)
				tt.Equal(t, d.err, err.Error(), i, " : ", x)
			} else {
				result := string(pw.Encode(out))
				tt.Equal(t, d.expect, result, i, " : ", x)
			}
		}
	}
}

func TestExprRemoveReflect(t *testing.T) {
	x, err := jp.ParseString("field[1]")
	tt.Nil(t, err)
	obj := &RemObj{Field: []any{1, 2, 3, 4}}
	result := x.MustRemove(obj)
	tt.Equal(t, "{field: [1 3 4]}", string(pw.Encode(result)))
	tt.Equal(t, "{field: [1 3 4]}", string(pw.Encode(obj)))

	obj = &RemObj{Field: []any{1, 2, 3, 4}}
	result = x.MustRemoveOne(obj)
	tt.Equal(t, "{field: [1 3 4]}", string(pw.Encode(result)))
	tt.Equal(t, "{field: [1 3 4]}", string(pw.Encode(obj)))

	x, err = jp.ParseString("field[1][0]")
	tt.Nil(t, err)
	obj = &RemObj{Field: []any{1, []any{2, 3}, 4}}
	result = x.MustRemove(obj)
	tt.Equal(t, "{field: [1 [3] 4]}", string(pw.Encode(result)))
	tt.Equal(t, "{field: [1 [3] 4]}", string(pw.Encode(obj)))

	x, err = jp.ParseString("field[1][0]")
	tt.Nil(t, err)
	obj = &RemObj{Field: 3}
	result = x.MustRemove(obj)
	tt.Equal(t, "{field: 3}", string(pw.Encode(result)))
	tt.Equal(t, "{field: 3}", string(pw.Encode(obj)))
}

func TestExprRemoveOneMap(t *testing.T) {
	x, err := jp.ParseString("*.two")
	tt.Nil(t, err)
	data := sen.MustParse([]byte("{one:{two:2} two:{two:4}}"))
	var result any
	result, err = x.RemoveOne(data)
	tt.Nil(t, err)
	// No telling what order items are removed from a map nor can the print
	// order be counted on to be consistent so just verify there is only one
	// match on the result and the original data.
	tt.Equal(t, 1, len(x.Get(result)))
	tt.Equal(t, 1, len(x.Get(data)))
}

func TestExprRemoveOneFail(t *testing.T) {
	x, err := jp.ParseString("..")
	tt.Nil(t, err)
	data := sen.MustParse([]byte("{one:[0,1,2] two:[3,2,1]}"))
	_, err = x.RemoveOne(data)
	tt.NotNil(t, err)
}

func TestExprRemoveFail(t *testing.T) {
	x, err := jp.ParseString("")
	tt.Nil(t, err)
	data := sen.MustParse([]byte("{one:[0,1,2] two:[3,2,1]}"))
	_, err = x.Remove(data)
	tt.NotNil(t, err)
}

func TestExprRemoveChildReflect(t *testing.T) {
	x, err := jp.ParseString("obj.field[1]")
	tt.Nil(t, err)
	data := map[string]any{"obj": &RemObj{Field: []any{1, 2, 3, 4}}}
	result := x.MustRemove(data)
	tt.Equal(t, "{obj: {field: [1 3 4]}}", string(pw.Encode(result)))
	tt.Equal(t, "{obj: {field: [1 3 4]}}", string(pw.Encode(data)))
}

func TestExprRemoveChildReflectReflect(t *testing.T) {
	x, err := jp.ParseString("field.field[1]")
	tt.Nil(t, err)
	data := &RemObj{Field: &RemObj{Field: []any{1, 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "{field: {field: [1]}}", string(pw.Encode(result)))
	tt.Equal(t, "{field: {field: [1]}}", string(pw.Encode(data)))
}

func TestExprRemoveNthReflect(t *testing.T) {
	x, err := jp.ParseString("[1].field[1]")
	tt.Nil(t, err)
	data := []*RemObj{{}, {Field: []any{1, 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(data)))

	data = []*RemObj{{}, {Field: []any{1, 2}}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(data)))
}

func TestExprRemoveNthReflectLast(t *testing.T) {
	x, err := jp.ParseString("[1].one.two")
	tt.Nil(t, err)
	data := []map[string]any{{}, {"one": map[string]any{"two": 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{} {one: {}}]", string(pw.Encode(result)))
	tt.Equal(t, "[{} {one: {}}]", string(pw.Encode(data)))
}

func TestExprRemoveNthReflectLastMap(t *testing.T) {
	x, err := jp.ParseString("@[1].one")
	tt.Nil(t, err)
	data := []map[string]int{{}, {"one": 1, "two": 2}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{} {two: 2}]", string(pw.Encode(result)))
	tt.Equal(t, "[{} {two: 2}]", string(pw.Encode(data)))
}

func TestExprRemoveNthReflectLastSlice(t *testing.T) {
	x, err := jp.ParseString("@[1][-2]")
	tt.Nil(t, err)
	data := [][]int{{}, {1, 2, 3}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(data)))

	data = [][]int{{}, {1, 2, 3}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(data)))
}

func TestExprRemoveWildReflectMap(t *testing.T) {
	x, err := jp.ParseString("*.field[1]")
	tt.Nil(t, err)
	data := map[string]any{"one": &RemObj{Field: []any{1, 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "{one: {field: [1]}}", string(pw.Encode(result)))
	tt.Equal(t, "{one: {field: [1]}}", string(pw.Encode(data)))

	x, err = jp.ParseString("one.*")
	tt.Nil(t, err)
	data = map[string]any{"one": map[string]int{"x": 1, "y": 2}}
	result = x.MustRemove(data)
	tt.Equal(t, "{one: {}}", string(pw.Encode(result)))
	tt.Equal(t, "{one: {}}", string(pw.Encode(data)))

	data = map[string]any{"one": map[string]int{"x": 1, "y": 2}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "{one: {y: 2}}", string(pw.Encode(result)))
	tt.Equal(t, "{one: {y: 2}}", string(pw.Encode(data)))
}

func TestExprRemoveWildReflectSlice(t *testing.T) {
	x, err := jp.ParseString("*.field[1]")
	tt.Nil(t, err)
	data := []any{&RemObj{Field: []any{1, 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{field: [1]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: [1]}]", string(pw.Encode(data)))

	x, err = jp.ParseString("[0][*]")
	tt.Nil(t, err)
	data = []any{[]int{1, 2}}
	result = x.MustRemove(data)
	tt.Equal(t, "[[]]", string(pw.Encode(result)))
	tt.Equal(t, "[[]]", string(pw.Encode(data)))

	data = []any{[]int{1, 2}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[[2]]", string(pw.Encode(result)))
	tt.Equal(t, "[[2]]", string(pw.Encode(data)))
}

func TestExprRemoveNthNodeInSimple(t *testing.T) {
	x, err := jp.ParseString("@[1][-2]")
	tt.Nil(t, err)
	data := []any{[]any{}, gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(data)))
}

func TestExprRemoveWildReflectWildSlice(t *testing.T) {
	x, err := jp.ParseString("[*][1]")
	tt.Nil(t, err)
	data := [][]int{{}, {1, 2, 3}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(data)))

	data = [][]int{{}, {1, 2, 3}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [1 3]]", string(pw.Encode(data)))
}

func TestExprRemoveWildReflectWildMap(t *testing.T) {
	x, err := jp.ParseString("*[1]")
	tt.Nil(t, err)
	data := map[string][]int{"one": {}, "two": {1, 2, 3}}
	result := x.MustRemove(data)
	tt.Equal(t, "{one: [] two: [1 3]}", string(pw.Encode(result)))
	tt.Equal(t, "{one: [] two: [1 3]}", string(pw.Encode(data)))

	data = map[string][]int{"one": {4, 5}, "two": {1, 2, 3}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "{one: [4] two: [1 2 3]}", string(pw.Encode(result)))
	tt.Equal(t, "{one: [4] two: [1 2 3]}", string(pw.Encode(data)))

	x, err = jp.ParseString("[*].*[1]")
	tt.Nil(t, err)
	data2 := []map[string][]int{{"one": {4, 5}, "two": {1, 2, 3}}}
	result = x.MustRemoveOne(data2)
	tt.Equal(t, "[{one: [4] two: [1 2 3]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{one: [4] two: [1 2 3]}]", string(pw.Encode(data2)))

	x, err = jp.ParseString("[*][0][1]")
	tt.Nil(t, err)
	data3 := [][]any{{[]any{1, 2, 3}}}
	result = x.MustRemoveOne(data3)
	tt.Equal(t, "[[[1 3]]]", string(pw.Encode(result)))
	tt.Equal(t, "[[[1 3]]]", string(pw.Encode(data3)))

	data4 := []int{1, 2, 3}
	result = x.MustRemoveOne(data4)
	tt.Equal(t, "[1 2 3]", string(pw.Encode(result)))
	tt.Equal(t, "[1 2 3]", string(pw.Encode(data4)))
}

func TestExprRemoveWildArrayInSimple(t *testing.T) {
	x, err := jp.ParseString("[1][*]")
	tt.Nil(t, err)
	data := []any{[]any{}, gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] []]", string(pw.Encode(result)))
	tt.Equal(t, "[[] []]", string(pw.Encode(data)))

	data = []any{[]any{}, gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[[] [2 3]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [2 3]]", string(pw.Encode(data)))
}

func TestExprRemoveWildObjectInSimple(t *testing.T) {
	x, err := jp.ParseString("[1].*")
	tt.Nil(t, err)
	data := []any{[]any{}, gen.Object{"one": gen.Int(1), "two": gen.Int(2), "three": gen.Int(3)}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] {}]", string(pw.Encode(result)))
	tt.Equal(t, "[[] {}]", string(pw.Encode(data)))

	data = []any{[]any{}, gen.Object{"one": gen.Int(1), "two": gen.Int(2), "three": gen.Int(3)}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[[] {three: 3 two: 2}]", string(pw.Encode(result)))
	tt.Equal(t, "[[] {three: 3 two: 2}]", string(pw.Encode(data)))
}

func TestExprRemoveDescent(t *testing.T) {
	x, err := jp.ParseString("..[1]")
	tt.Nil(t, err)
	data := sen.MustParse([]byte(`[[1,2,[1,2,3,4]]]`))
	tt.Panic(t, func() { _ = x.MustRemove(data) })
}

func TestExprRemoveUnionReflectMap(t *testing.T) {
	x, err := jp.ParseString("['a','b'].field[0]")
	tt.Nil(t, err)
	data := map[string]any{"a": &RemObj{Field: []any{1, 2, 3}}}
	result := x.MustRemove(data)
	tt.Equal(t, "{a: {field: [2 3]}}", string(pw.Encode(result)))
	tt.Equal(t, "{a: {field: [2 3]}}", string(pw.Encode(data)))

	x, err = jp.ParseString("['field','x'][0]")
	tt.Nil(t, err)
	data2 := &RemObj{Field: []any{1, 2, 3}}
	result = x.MustRemove(data2)
	tt.Equal(t, "{field: [2 3]}", string(pw.Encode(result)))
	tt.Equal(t, "{field: [2 3]}", string(pw.Encode(data2)))

	x, err = jp.ParseString("['field','x'][0][1]")
	tt.Nil(t, err)
	data2 = &RemObj{Field: []any{[]any{1, 2, 3}}}
	result = x.MustRemove(data2)
	tt.Equal(t, "{field: [[1 3]]}", string(pw.Encode(result)))
	tt.Equal(t, "{field: [[1 3]]}", string(pw.Encode(data2)))

	x, err = jp.ParseString("['x','field']['field','x'][1]")
	tt.Nil(t, err)
	data2 = &RemObj{Field: &RemObj{Field: []any{1, 2, 3}}}
	result = x.MustRemove(data2)
	tt.Equal(t, "{field: {field: [1 3]}}", string(pw.Encode(result)))
	tt.Equal(t, "{field: {field: [1 3]}}", string(pw.Encode(data2)))

	data2 = &RemObj{Field: &RemObj{Field: []any{1, 2, 3}}}
	result = x.MustRemoveOne(data2)
	tt.Equal(t, "{field: {field: [1 3]}}", string(pw.Encode(result)))
	tt.Equal(t, "{field: {field: [1 3]}}", string(pw.Encode(data2)))

	x, err = jp.ParseString("['a','c']")
	tt.Nil(t, err)
	data3 := map[string][]any{"a": {1}, "b": {2}, "c": {3}}
	result = x.MustRemove(data3)
	tt.Equal(t, "{b: [2]}", string(pw.Encode(result)))
	tt.Equal(t, "{b: [2]}", string(pw.Encode(data3)))

	data3 = map[string][]any{"a": {1}, "b": {2}, "c": {3}}
	result = x.MustRemoveOne(data3)
	tt.Equal(t, "{b: [2] c: [3]}", string(pw.Encode(result)))
	tt.Equal(t, "{b: [2] c: [3]}", string(pw.Encode(data3)))
}

func TestExprRemoveUnionReflectSlice(t *testing.T) {
	x, err := jp.ParseString("[0,2].field[0]")
	tt.Nil(t, err)
	data := []any{&RemObj{Field: []any{1, 2, 3}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(data)))

	data2 := []*RemObj{{Field: []any{1, 2, 3}}}
	result = x.MustRemove(data2)
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(data2)))

	x, err = jp.ParseString("[0,-1][1]")
	tt.Nil(t, err)
	data3 := [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemove(data3)
	tt.Equal(t, "[[1 3] [4 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[1 3] [4 6]]", string(pw.Encode(data3)))

	data3 = [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemoveOne(data3)
	tt.Equal(t, "[[1 3] [4 5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[1 3] [4 5 6]]", string(pw.Encode(data3)))

	x, err = jp.ParseString("[0,2][0,1][1]")
	tt.Nil(t, err)
	data4 := [][]any{{[]any{1, 2, 3}}}
	result = x.MustRemove(data4)
	tt.Equal(t, "[[[1 3]]]", string(pw.Encode(result)))
	tt.Equal(t, "[[[1 3]]]", string(pw.Encode(data4)))

	x, err = jp.ParseString("[0,2]")
	tt.Nil(t, err)
	data4 = [][]any{{1}, {2}, {3}}
	result = x.MustRemove(data4)
	tt.Equal(t, "[[2]]", string(pw.Encode(result)))
	tt.Equal(t, "[[1] [2] [3]]", string(pw.Encode(data4))) // should be unchanged
}

func TestExprRemoveSliceReflect(t *testing.T) {
	x, err := jp.ParseString("[0:3].field[0]")
	tt.Nil(t, err)
	data := []*RemObj{{Field: []any{1, 2, 3}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(data)))

	data2 := []map[string]any{{"field": []any{1, 2, 3}}}
	result = x.MustRemove(data2)
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: [2 3]}]", string(pw.Encode(data2)))

	x, err = jp.ParseString("[0:3][0]")
	tt.Nil(t, err)
	data3 := [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemove(data3)
	tt.Equal(t, "[[2 3] [5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[2 3] [5 6]]", string(pw.Encode(data3)))

	data3 = [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemoveOne(data3)
	tt.Equal(t, "[[2 3] [4 5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[2 3] [4 5 6]]", string(pw.Encode(data3)))

	x, err = jp.ParseString("[-1:-2:-1][0]")
	tt.Nil(t, err)
	data3 = [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemove(data3)
	tt.Equal(t, "[[2 3] [5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[2 3] [5 6]]", string(pw.Encode(data3)))

	data3 = [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemoveOne(data3)
	tt.Equal(t, "[[1 2 3] [5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[1 2 3] [5 6]]", string(pw.Encode(data3)))

	x, err = jp.ParseString("[3:][0]")
	tt.Nil(t, err)
	data3 = [][]any{{1, 2, 3}, {4, 5, 6}}
	result = x.MustRemove(data3)
	tt.Equal(t, "[[1 2 3] [4 5 6]]", string(pw.Encode(result)))
	tt.Equal(t, "[[1 2 3] [4 5 6]]", string(pw.Encode(data3)))
}

func TestExprRemoveFilterSlice(t *testing.T) {
	x, err := jp.ParseString("[?(@.x == 1)].y")
	tt.Nil(t, err)
	data := []map[string]any{{"x": 1, "y": 2}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{x: 1}]", string(pw.Encode(result)))
	tt.Equal(t, "[{x: 1}]", string(pw.Encode(data)))

	data = []map[string]any{{"x": 1, "y": 2}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "[{x: 1}]", string(pw.Encode(result)))
	tt.Equal(t, "[{x: 1}]", string(pw.Encode(data)))
}

func TestExprRemoveFilterMap(t *testing.T) {
	x, err := jp.ParseString("[?(@.x == 1)].y")
	tt.Nil(t, err)
	data := map[string]map[string]any{"a": {"x": 1, "y": 2}, "b": {"x": 1, "y": 3}}
	result := x.MustRemove(data)
	tt.Equal(t, "{a: {x: 1} b: {x: 1}}", string(pw.Encode(result)))
	tt.Equal(t, "{a: {x: 1} b: {x: 1}}", string(pw.Encode(data)))

	data = map[string]map[string]any{"a": {"x": 1, "y": 2}, "b": {"x": 1, "y": 3}}
	result = x.MustRemoveOne(data)
	tt.Equal(t, "{a: {x: 1} b: {x: 1 y: 3}}", string(pw.Encode(result)))
	tt.Equal(t, "{a: {x: 1} b: {x: 1 y: 3}}", string(pw.Encode(data)))
}

func xTestExprRemoveDev(t *testing.T) {
	x, err := jp.ParseString("[?(@.x == 1)].y")
	tt.Nil(t, err)
	data := []map[string]any{{"x": 1, "y": 2}}
	result := x.MustRemove(data)
	fmt.Printf("*** %s\n", pw.Encode(result))
	fmt.Printf("*** %s\n", pw.Encode(data))
}
