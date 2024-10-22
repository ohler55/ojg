// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

// Walk data and call the cb callback for each node in the data. The path is
// reused in each call so if the path needs to be save it should be copied.
func Walk(data any, cb func(path Expr, value any), justLeaves ...bool) {
	path := Expr{Root('$')}
	walk(path, data, cb, 0 < len(justLeaves) && justLeaves[0])
}

func walk(path Expr, data any, cb func(path Expr, value any), justLeaves bool) {
top:
	switch td := data.(type) {
	case nil, bool, int64, float64, string,
		int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32,
		[]byte, time.Time:
		// leaf node
		cb(path, data)
	case []any:
		if !justLeaves {
			cb(path, data)
		}
		pi := len(path)
		path = append(path, nil)
		for i, v := range td {
			path[pi] = Nth(i)
			walk(path, v, cb, justLeaves)
		}
	case map[string]any:
		if !justLeaves {
			cb(path, data)
		}
		pi := len(path)
		path = append(path, nil)
		for k, v := range td {
			path[pi] = Child(k)
			walk(path, v, cb, justLeaves)
		}
	case gen.Array:
		if !justLeaves {
			cb(path, data)
		}
		pi := len(path)
		path = append(path, nil)
		for i, v := range td {
			path[pi] = Nth(i)
			walk(path, v, cb, justLeaves)
		}
	case gen.Object:
		if !justLeaves {
			cb(path, data)
		}
		pi := len(path)
		path = append(path, nil)
		for k, v := range td {
			path[pi] = Child(k)
			walk(path, v, cb, justLeaves)
		}
	case alt.Simplifier:
		data = td.Simplify()
		goto top
	default:
		cb(path, data)
	}
}

// Walk the matching elements in the data and call cb on the matches or follow
// on to the matching if not the last fragment in an expression. The rest
// argument is the rest of the expression after this fragment. The path is the
// normalized path up to this point. Data is the data element to act on.
func (x Expr) Walk(data any, cb func(path Expr, data, parent any)) {
	if 0 < len(x) {
		x[0].Walk(x[1:], Expr{}, data, cb)
	}
}

// Walk continues with the next in rest.
func (f At) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	if 0 < len(rest) {
		rest[0].Walk(rest[1:], path, data, cb)
	} else {
		cb(path, data, data)
	}
}

func (f Bracket) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Child) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Descent) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Filter) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Nth) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD

	// TBD if last ...
	//     else cb

	// TBD to delete in data, mod data at nth and update in parent (slice or map)
	// how to deal with changes in index for unions and ranges?
	//  keep data and just update parent each time.

}

// Walk continues with the next in rest.
func (f Root) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	if 0 < len(rest) {
		rest[0].Walk(rest[1:], path, data, cb)
	} else {
		cb(path, data, data)
	}
}

func (f Slice) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Union) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}

func (f Wildcard) Walk(rest, path Expr, data any, cb func(path Expr, data, parent any)) {
	// TBD
}
