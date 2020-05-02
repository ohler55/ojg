// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"io"

	"github.com/ohler55/ojg/gd"
)

// Parse JSON into a gd.Node. Arguments are optional and can be a bool
// or func(gd.Node) bool.
//
// A bool indicates the NoComment parser attribute should be set to the bool value.
//
// A func argument is the callback for the parser if processing multiple
// JSONs. If no callback function is provided the processing is limited to
// only one JSON.
func Parse(b []byte, args ...interface{}) (n gd.Node, err error) {
	p := Parser{}
	return p.Parse(b, args...)
}

func ParseString(s string, args ...interface{}) (n gd.Node, err error) {
	p := Parser{}
	return p.Parse([]byte(s), args...)
}

// Load
func Load(r io.Reader, args ...interface{}) (gd.Node, error) {

	// TBD

	return nil, nil
}

func ParseSimple(s string, args ...interface{}) (interface{}, error) {

	// TBD

	return nil, nil
}

func LoadSimple(r io.Reader, args ...interface{}) (interface{}, error) {

	// TBD

	return nil, nil
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
