// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
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
		{path: "*[0]", data: `{one:[1],two:[1,2]}`, expect: `{one: [] two: [2]}`},
		{path: "*[0]", data: `[[1],[1,2]]`, expect: `[[] [2]]`},
		{path: "key[0]", data: `[[1],[1,2]]`, expect: `[[1] [1 2]]`},
		{path: "one.two[2]", data: `{one:{two:[1,2,3,4]}}`, expect: `{one: {two: [1 2 4]}}`},
		{path: "one.two[2]", data: `{one:{two:2}}`, expect: `{one: {two: 2}}`},
		{path: "key[-2]", data: `{key:[1,2,3,4]}`, expect: `{key: [1 2 4]}`},
		{path: "[-1][-2]", data: `[1,2,[1,2,3,4]]`, expect: `[1 2 [1 2 4]]`},
		{path: "[0][-1][-2]", data: `[[1,2,[1,2,3,4]]]`, expect: `[[1 2 [1 2 4]]]`},
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

func TestExprRemoveReflect(t *testing.T) {
	x, err := jp.ParseString("field[1]")
	tt.Nil(t, err)
	obj := &RemObj{Field: []any{1, 2, 3, 4}}
	result := x.MustRemove(obj)
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

	// TBD try other variations
}

func TestExprRemoveOne(t *testing.T) {
	x, err := jp.ParseString("*[1]")
	tt.Nil(t, err)
	data := sen.MustParse([]byte("[[0,1,2][3,2,1]]"))
	var result any
	result, err = x.RemoveOne(data)
	tt.Nil(t, err)
	tt.Equal(t, "[[0 2] [3 2 1]]", string(pw.Encode(result)))
	tt.Equal(t, "[[0 2] [3 2 1]]", string(pw.Encode(data)))
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

func TestExprMustRemoveOne(t *testing.T) {
	x, err := jp.ParseString("@[*][1]")
	tt.Nil(t, err)
	data := sen.MustParse([]byte("[[0,1,2][3,2,1]]"))
	result := x.MustRemoveOne(data)
	tt.Equal(t, "[[0 2] [3 2 1]]", string(pw.Encode(result)))
	tt.Equal(t, "[[0 2] [3 2 1]]", string(pw.Encode(data)))
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
	data := []*RemObj{&RemObj{}, &RemObj{Field: []any{1, 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(result)))
	tt.Equal(t, "[{field: null} {field: [1]}]", string(pw.Encode(data)))
}

func TestExprRemoveNthReflectLast(t *testing.T) {
	x, err := jp.ParseString("[1].one.two")
	tt.Nil(t, err)
	data := []map[string]any{map[string]any{}, map[string]any{"one": map[string]any{"two": 2}}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{} {one: {}}]", string(pw.Encode(result)))
	tt.Equal(t, "[{} {one: {}}]", string(pw.Encode(data)))
}

func TestExprRemoveNthReflectLastMap(t *testing.T) {
	x, err := jp.ParseString("@[1].one")
	tt.Nil(t, err)
	data := []map[string]int{map[string]int{}, map[string]int{"one": 1, "two": 2}}
	result := x.MustRemove(data)
	tt.Equal(t, "[{} {two: 2}]", string(pw.Encode(result)))
	tt.Equal(t, "[{} {two: 2}]", string(pw.Encode(data)))
}

// TBD
func xTestExprRemoveNthReflectLastSlice(t *testing.T) {
	x, err := jp.ParseString("@[1][0]")
	tt.Nil(t, err)
	data := [][]int{[]int{}, []int{1, 2}}
	result := x.MustRemove(data)
	tt.Equal(t, "[[] [2]]", string(pw.Encode(result)))
	tt.Equal(t, "[[] [2]]", string(pw.Encode(data)))
}

func TestExprRemoveDev(t *testing.T) {
	x, err := jp.ParseString("field[1]")
	tt.Nil(t, err)

	obj := RemObj{Field: []any{1, 2, 3, 4}}

	result := x.MustRemove(obj)
	fmt.Printf("*** %s\n", pw.Encode(result))
	fmt.Printf("*** %s\n", pw.Encode(obj))

}
