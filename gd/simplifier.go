// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gd

type Simplifier interface {

	// Simplify should return either a gd.Node or one of the simple type which
	// are: nil, bool, int64, float64, string, time.Time, []interface{}, or
	// map[string]interface{}.
	Simplify() interface{}
}
