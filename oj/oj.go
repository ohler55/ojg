// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"io"

	"github.com/ohler55/ojg/alt"
)

// Parse JSON into a gen.Node. Arguments are optional and can be a bool
// or func(interface{}) bool.
//
// A bool indicates the NoComment parser attribute should be set to the bool
// value.
//
// A func argument is the callback for the parser if processing multiple
// JSONs. If no callback function is provided the processing is limited to
// only one JSON.
func Parse(b []byte, args ...interface{}) (n interface{}, err error) {
	p := Parser{}
	return p.Parse(b, args...)
}

// ParseString is similar to Parse except it takes a string
// argument to be parsed instead of a []byte.
func ParseString(s string, args ...interface{}) (n interface{}, err error) {
	p := Parser{}
	return p.Parse([]byte(s), args...)
}

// Load a JSON from a io.Reader into a simple type. An error is returned
// if not valid JSON.
func Load(r io.Reader, args ...interface{}) (interface{}, error) {
	p := Parser{}
	return p.ParseReader(r, args...)
}

// Validate a JSON string. An error is returned if not valid JSON.
func Validate(b []byte) error {
	v := Validator{}
	return v.Validate(b)
}

// ValidateString a JSON string. An error is returned if not valid JSON.
func ValidateString(s string) error {
	v := Validator{}
	return v.Validate([]byte(s))
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func ValidateReader(r io.Reader) error {
	v := Validator{}
	return v.ValidateReader(r)
}

// Unmarshal parses the provided JSON and stores the result in the value
// pointed to by vp.
func Unmarshal(data []byte, vp interface{}, recomposer ...*alt.Recomposer) (err error) {
	p := Parser{}
	var v interface{}
	if v, err = p.Parse(data); err == nil {
		if 0 < len(recomposer) {
			_, err = recomposer[0].Recompose(v, vp)
		} else {
			_, err = alt.Recompose(v, vp)
		}
	}
	return
}
