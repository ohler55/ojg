// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestUnmarshal(t *testing.T) {
	var obj map[string]interface{}
	src := `{"x":3}`
	err := oj.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, oj.JSON(obj))
	tt.Equal(t, 3.0, obj["x"])

	obj = nil
	p := oj.Parser{}
	err = p.Unmarshal([]byte(src), &obj)
	tt.Nil(t, err)
	tt.Equal(t, src, oj.JSON(obj))

	obj = nil
	err = oj.Unmarshal([]byte(src), &obj, &alt.Recomposer{})
	tt.Nil(t, err)
	tt.Equal(t, src, oj.JSON(obj))
}

func TestUnmarshalError(t *testing.T) {
	type Query struct {
		Level  string
		Query  map[string]interface{}
		Expand bool
		Limit  int
	}

	queryJSON := `{
	"Level": "Series",
	"Query": {},
	"Expand": false,
	"Limit": true
}`

	var query Query
	err := oj.Unmarshal([]byte(queryJSON), &query)
	tt.Equal(t, true, strings.Contains(err.Error(), "value of type bool cannot be converted to type int"))
}
