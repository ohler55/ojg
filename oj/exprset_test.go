// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestExprSet(t *testing.T) {
	data := map[string]interface{}{}
	err := oj.R().C("a").Set(data, map[string]interface{}{})
	tt.Nil(t, err)
	tt.Equal(t, `{"a":{}}`, oj.JSON(data, &oj.Options{Sort: true}))

	err = oj.R().C("b").Set(data, []interface{}{1, 2, 3})
	tt.Nil(t, err)
	tt.Equal(t, `{"a":{},"b":[1,2,3]}`, oj.JSON(data, &oj.Options{Sort: true}))

	err = oj.R().C("b").Nth(1).Set(data, map[string]interface{}{})
	tt.Nil(t, err)
	tt.Equal(t, `{"a":{},"b":[1,{},3]}`, oj.JSON(data, &oj.Options{Sort: true}))

	err = oj.R().C("b").N(1).C("x").Set(data, 7)
	tt.Nil(t, err)
	tt.Equal(t, `{"a":{},"b":[1,{"x":7},3]}`, oj.JSON(data, &oj.Options{Sort: true}))

	err = oj.R().C("b").W().C("x").Set(data, 5)
	tt.Nil(t, err)
	tt.Equal(t, `{"a":{},"b":[1,{"x":5},3]}`, oj.JSON(data, &oj.Options{Sort: true}))

	oj.R().C("b").W().C("x").Del(data)
	tt.Equal(t, `{"a":{},"b":[1,{},3]}`, oj.JSON(data, &oj.Options{Sort: true}))
}
