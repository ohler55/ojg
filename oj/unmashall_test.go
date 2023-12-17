// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestUnmarshal(t *testing.T) {
	var obj map[string]any
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
		Query  map[string]any
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

type TagMap map[string]any

func (tm *TagMap) UnmarshalJSON(data []byte) error {
	*tm = map[string]any{}
	simple, err := oj.Parse(data)
	if err != nil {
		return err
	}
	for _, kv := range simple.([]any) {
		(*tm)[jp.C("key").First(kv).(string)] = jp.C("value").First(kv)
	}
	return nil
}

func TestUnmarshaler(t *testing.T) {
	var tags TagMap
	src := []byte(`[{"key": "k1", "value": 1}]`)
	err := oj.Unmarshal(src, &tags)
	tt.Nil(t, err)
	tt.Equal(t, 1, len(tags))
	tt.Equal(t, 1, tags["k1"])
}

type Triple [3]float64

func TestUnmarshalArray(t *testing.T) {
	var tri Triple
	src := []byte(`[1.0, 2.0, 3.0]`)
	err := oj.Unmarshal(src, &tri)
	tt.Nil(t, err)
	tt.Equal(t, 1.0, tri[0])
	tt.Equal(t, 2.0, tri[1])
	tt.Equal(t, 3.0, tri[2])
}
