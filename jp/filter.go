// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"

	"github.com/ohler55/ojg"
)

// Filter is a script used as a filter.
type Filter struct {
	Script
}

// NewFilter creates a new Filter.
func NewFilter(str string) (f *Filter, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	f = MustNewFilter(str)
	return
}

// MustNewFilter creates a new Filter and panics on error.
func MustNewFilter(str string) (f *Filter) {
	p := &parser{buf: []byte(str)}
	if len(p.buf) <= 5 ||
		p.buf[0] != '[' || p.buf[1] != '?' || p.buf[2] != '(' ||
		p.buf[len(p.buf)-2] != ')' || p.buf[len(p.buf)-1] != ']' {
		panic(fmt.Errorf("a filter must start with a '[?(' and end with ')]'"))
	}
	p.buf = p.buf[3 : len(p.buf)-1]
	eq := p.readEquation()

	return eq.Filter()
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
