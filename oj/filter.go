// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Filter is a script used as a filter.
type Filter struct {
	Script
}

func NewFilter(str string) (f *Filter, err error) {
	f = &Filter{}
	err = f.Parse([]byte(str))
	return
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
