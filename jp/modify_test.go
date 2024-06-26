// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

func TestExprModify(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result, err := x.Modify(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "[[1 2 3 4] [1 2 3 4]]", string(pw.Encode(result)))
}

func TestExprMustModify(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result := x.MustModify(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Equal(t, "[[1 2 3 4] [1 2 3 4]]", string(pw.Encode(result)))
}

func TestExprModifyOne(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result, err := x.ModifyOne(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "[[1 2 3 4] [4 3 2 1]]", string(pw.Encode(result)))
}

func TestExprMustModifyOne(t *testing.T) {
	x := jp.MustParseString("[*]")

	data := []any{[]any{1, 3, 2, 4}, []any{4, 3, 2, 1}}
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if a, ok := element.([]any); ok {
			sort.Slice(a, func(i, j int) bool { return a[i].(int) < a[j].(int) })
			element = a
		}
		return element, true
	})
	tt.Equal(t, "[[1 2 3 4] [4 3 2 1]]", string(pw.Encode(result)))
}

func TestExprModifyOneEmpty(t *testing.T) {
	x := jp.MustParseString("")

	data := []any{}
	_, err := x.ModifyOne(data, func(element any) (any, bool) {
		return element, false
	})
	tt.NotNil(t, err)
}

func TestExprModifyEmpty(t *testing.T) {
	x := jp.MustParseString("")

	data := []any{}
	_, err := x.Modify(data, func(element any) (any, bool) {
		return element, false
	})
	tt.NotNil(t, err)
}

func TestExprModifyDescent(t *testing.T) {
	x := jp.MustParseString("$..key")
	var data any
	data = map[string]any{"a": []any{map[string]any{"key": 1}}}
	result, err := x.Modify(data, func(_ any) (any, bool) {
		return 4, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(result)))
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(data)))

	data = gen.Object{"a": gen.Array{gen.Object{"key": gen.Int(1)}}}
	result, err = x.Modify(data, func(_ any) (any, bool) {
		return gen.Int(4), true
	})
	tt.Nil(t, err)
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(result)))
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(data)))

	data = map[string]any{"a": []any{map[string]int{"key": 1}}}
	result, err = x.Modify(data, func(_ any) (any, bool) {
		return 4, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(result)))
	tt.Equal(t, "{a: [{key: 4}]}", string(pw.Encode(data)))

	data = map[string]any{"a": map[string]int{"key": 1}}
	result, err = x.Modify(data, func(_ any) (any, bool) {
		return 4, true
	})
	tt.Nil(t, err)
	tt.Equal(t, "{a: {key: 4}}", string(pw.Encode(result)))
	tt.Equal(t, "{a: {key: 4}}", string(pw.Encode(data)))
}

func keyedData() any {
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

func TestExprModifyKeyedChild(t *testing.T) {
	data := keyedData()
	x := jp.MustParseString("b")
	result, err := x.Modify(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 1
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, 3, jp.C("b").First(result))

	x = jp.C("c").C("c2")
	result = x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 17, x.First(result))
}

func TestExprModifyKeyedWild(t *testing.T) {
	data := keyedData()
	x := jp.W().W()
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.C("c").C("c1").First(result))
}

func TestExprModifyKeyedUnion(t *testing.T) {
	data := keyedData()
	x := jp.U(1, "c").U(2, "c1")
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.C("c").C("c1").First(result))
}

func TestExprModifyKeyedDescent(t *testing.T) {
	data := keyedData()
	x := jp.D().C("c1")
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.C("c").C("c1").First(result))
}

func indexedData() any {
	return &indexed{
		ordered: ordered{
			entries: []*entry{
				{value: 1},
				{value: 2},
				{value: &indexed{
					ordered: ordered{
						entries: []*entry{
							{value: 11},
							{value: 12},
							{value: 13},
						},
					},
				},
				},
			},
		},
	}
}

func TestExprModifyIndexedNth(t *testing.T) {
	data := indexedData()
	x := jp.N(1)
	result, err := x.Modify(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 1
		}
		return element, true
	})
	tt.Nil(t, err)
	tt.Equal(t, 3, x.First(result))

	x = jp.N(-1).N(0)
	result = x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, x.First(result))
}

func TestExprModifyIndexedWild(t *testing.T) {
	data := indexedData()
	x := jp.W().W()
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.N(2).N(0).First(result))
}

func TestExprModifyIndexedUnion(t *testing.T) {
	data := indexedData()
	x := jp.U("a", -1).U("b", 0)
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.N(2).N(0).First(result))
}

func TestExprModifyIndexedSlice(t *testing.T) {
	data := indexedData()
	x := jp.S(-2, -1).S(0, 4)
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.N(2).N(0).First(result))

	x = jp.S(2, 0, -1).S(2, 0, -1)
	result = x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 18, jp.N(2).N(2).First(result))

	x = jp.S(-5, 2).S(2, 0)
	var changed bool
	_ = x.MustModifyOne(data, func(element any) (any, bool) {
		changed = true
		return element, true
	})
	tt.Equal(t, false, changed)
}

func TestExprModifyIndexedFilter(t *testing.T) {
	data := indexedData()
	x := jp.MustParseString("*[?@ == 11]")
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.N(2).N(0).First(result))
}

func TestExprModifyIndexedDescent(t *testing.T) {
	data := indexedData()
	x := jp.D().N(0)
	result := x.MustModifyOne(data, func(element any) (any, bool) {
		if num, ok := element.(int); ok {
			element = num + 5
		}
		return element, true
	})
	tt.Equal(t, 16, jp.N(2).N(0).First(result))
}
