// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"time"

	"github.com/ohler55/ojg/gen"
)

// Walk data and call the cb callback for each node in the data. The path is
// reused in each call so if the path needs to be save it should be copied.
func Walk(data any, cb func(path Expr, value any)) {
	path := Expr{Root('$')}
	walk(path, data, cb)
}

func walk(path Expr, data any, cb func(path Expr, value any)) {
	cb(path, data)
	switch td := data.(type) {
	case nil, bool, int64, float64, string,
		int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32,
		[]byte, time.Time:
		// leaf node
	case []any:
		for i, v := range td {
			walk(append(path, Nth(i)), v, cb)
		}
	case map[string]any:
		for k, v := range td {
			walk(append(path, Child(k)), v, cb)
		}
	case gen.Array:
		for i, v := range td {
			walk(append(path, Nth(i)), v, cb)
		}
	case gen.Object:
		for k, v := range td {
			walk(append(path, Child(k)), v, cb)
		}
	default:
		// TBD use reflection
	}
}
