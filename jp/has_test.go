// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/tt"
)

type hasData struct {
	path   string
	data   interface{}
	expect bool
}

var (
	hasTestData = []*hasData{
		{path: "", expect: false},
		{path: "$.a.*.b", expect: true},
	}
	hasTestReflectData = []*hasData{
		{path: "$.a", expect: true, data: &Sample{A: 3, B: "sample"}},
	}
)

func TestExprHas(t *testing.T) {
	data := buildTree(4, 3, 0)
	for i, d := range append(hasTestData, hasTestReflectData...) {
		if testing.Verbose() {
			fmt.Printf("... %d: %s\n", i, d.path)
		}
		x, err := jp.ParseString(d.path)
		tt.Nil(t, err)
		var result bool
		if d.data == nil {
			result = x.Has(data)
		} else {
			result = x.Has(d.data)
		}
		tt.Equal(t, d.expect, result, i, " : ", x)
	}
}
