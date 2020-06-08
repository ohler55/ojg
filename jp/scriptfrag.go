// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

// ScriptFrag is a script used as a script fragment.
type ScriptFrag struct {
	Script *Script
}

// String representation of the scriptFrag.
func (f *ScriptFrag) String() string {
	return string(f.Append([]byte{}, true, false))
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f ScriptFrag) Append(buf []byte, _, _ bool) []byte {
	buf = append(buf, "["...)
	buf = f.Script.Append(buf)
	buf = append(buf, ']')

	return buf
}
