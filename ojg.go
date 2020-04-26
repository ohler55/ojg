// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"io"

	"github.com/ohler55/ojg/gd"
)

// Parse a string into a gd.Node. Arguments are optional and can be a bool,
// a *ParseOptions, or func(gd.Node) bool.
//
// A bool is used to indicated if the parsing should be limited to one JSON only.
//
func Parse(s string, args ...interface{}) (n gd.Node, err error) {
	p := Parser{}
	return p.Parse(s)
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
func Validate(s string, args ...interface{}) error {
	p := Parser{}
	return p.Validate(s)
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func ValidateReader(r io.Reader, args ...interface{}) error {
	p := Parser{}
	return p.ValidateReader(r)
}
