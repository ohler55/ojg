// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

// Filter is a script used as a filter.
type Filter struct {
	Script
}

func NewFilter(str string) (f *Filter, err error) {
	xp := &xparser{buf: []byte(str)}
	if len(xp.buf) == 0 || xp.buf[0] != '(' {
		return nil, fmt.Errorf("a filter must start with a '('")
	}
	xp.pos = 1
	eq, err := xp.readEquation()
	if err == nil && xp.pos < len(xp.buf) {
		err = fmt.Errorf("parse error")
	}
	if err != nil {
		return nil, fmt.Errorf("%s at %d in %s", err, xp.pos, xp.buf)
	}
	return eq.Filter(), nil
}

// String representation of the filter.
func (f *Filter) String() string {
	return string(f.Append([]byte{}, true, false))
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Filter) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, "[?"...)
	buf = f.Script.Append(buf)
	buf = append(buf, ']')

	return buf
}
