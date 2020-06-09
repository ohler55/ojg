// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import "fmt"

// Filter is a script used as a filter.
type Filter struct {
	Script
}

func NewFilter(str string) (f *Filter, err error) {
	p := &parser{buf: []byte(str)}
	if len(p.buf) == 0 || p.buf[0] != '(' {
		return nil, fmt.Errorf("a filter must start with a '('")
	}
	p.pos = 1
	eq, err := p.readEquation()
	if err == nil && p.pos < len(p.buf) {
		err = fmt.Errorf("parse error")
	}
	if err != nil {
		return nil, fmt.Errorf("%s at %d in %s", err, p.pos, p.buf)
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
