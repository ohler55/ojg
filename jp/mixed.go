// Copyright (c) 2026, Peter Ohler, All rights reserved.

package jp

import (
	"strconv"

	"github.com/ohler55/ojg/gen"
)

// MixedUnion is a bracketed list of selectors as described in RFC 9535
// section 2.5.1 where at least one member is a slice, a wildcard, or a
// filter. Each member is a Child, Nth, Slice, Wildcard, or *Filter. The
// values selected by each member are returned in member order and duplicates
// are kept. A list of only names and indexes is a Union.
type MixedUnion []Frag

// NewMixedUnion creates a new MixedUnion with the provided members. A string
// member becomes a Child and an int or int64 member becomes an Nth. Members
// of any other type than a Child, Nth, Slice, Wildcard, or *Filter are
// ignored.
func NewMixedUnion(members ...any) (mu MixedUnion) {
	for _, m := range members {
		switch tm := m.(type) {
		case string:
			mu = append(mu, Child(tm))
		case int:
			mu = append(mu, Nth(tm))
		case int64:
			mu = append(mu, Nth(tm))
		case Child, Nth, Slice, Wildcard, *Filter:
			mu = append(mu, tm.(Frag))
		}
	}
	return
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f MixedUnion) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, '[')
	for i, m := range f {
		if 0 < i {
			buf = append(buf, ',')
		}
		switch tm := m.(type) {
		case Child:
			buf = AppendString(buf, string(tm), '\'')
		case Nth:
			buf = append(buf, strconv.FormatInt(int64(tm), 10)...)
		case Wildcard:
			buf = append(buf, '*')
		default:
			// A Slice or Filter writes itself in brackets so drop them.
			start := len(buf)
			buf = m.Append(buf, true, false)
			if start+1 < len(buf) {
				copy(buf[start:], buf[start+1:len(buf)-1])
				buf = buf[:len(buf)-2]
			}
		}
	}
	return append(buf, ']')
}

// get returns the values each member selects from data, in member order. The
// root is used by filters that refer to the root ($) of the document.
func (f MixedUnion) get(data, root any) (values []any) {
	for _, m := range f {
		if tm, ok := m.(*Filter); ok {
			ns, _ := tm.evalWithRoot([]any{}, data, root)
			stack, _ := ns.([]any)
			// The filter pushes matches onto the stack last match first.
			for i := len(stack) - 1; 0 <= i; i-- {
				values = append(values, stack[i])
			}
			continue
		}
		values = append(values, Expr{m}.Get(data)...)
	}
	return
}

// getNodes is the gen.Node version of get.
func (f MixedUnion) getNodes(n, root gen.Node) (values []gen.Node) {
	for _, m := range f {
		if tm, ok := m.(*Filter); ok {
			ns, _ := tm.evalWithRoot([]gen.Node{}, n, root)
			stack, _ := ns.([]gen.Node)
			for i := len(stack) - 1; 0 <= i; i-- {
				values = append(values, stack[i])
			}
			continue
		}
		values = append(values, Expr{m}.GetNodes(n)...)
	}
	return
}

func (f MixedUnion) locate(pp Expr, data any, rest Expr, max int) (locs []Expr) {
	for _, m := range f {
		mx := max
		if 0 < max {
			mx = max - len(locs)
		}
		locs = append(locs, m.locate(pp, data, rest, mx)...)
		if 0 < max && max <= len(locs) {
			break
		}
	}
	return
}

// Walk each member of the union in order.
func (f MixedUnion) Walk(rest, path Expr, nodes []any, cb func(path Expr, nodes []any)) {
	for _, m := range f {
		m.Walk(rest, path, nodes, cb)
	}
}
