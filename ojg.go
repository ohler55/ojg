// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"

	"github.com/ohler55/ojg/gd"
)

// Parse a string into a gd.Node. Arguments are optional and can be a bool,
// int, or func(gd.Node) bool.
//
// A bool is used to indicated if the parsing should be strict or tolerant. If tolerant then
// C style // comments are allowed. That is the default.
//
// An int indicates the parsing should be limited to that number of top level
// JSON elements and return and error if more than that limit is encountered.
//
// If the input includes multiple JSON documents then a callback function of
// the form func(gd.Node) bool can be provided. The bool return if true will
// abort processing.
func Parse(s string, args ...interface{}) (n gd.Node, err error) {
	p := parser{buf: []byte(s)} // TBD add handler
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.strict = ta
		case int:
			p.limit = ta
		case string:
			// TBD timeformat
		case func(gd.Node) bool:
			// TBD set in handler
		default:
			return nil, fmt.Errorf("%T is not a valid argument type", a)
		}
	}
	// TBD callback to set n

	err = p.parse()

	return
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

func Validate(s string, args ...interface{}) error {
	p := parser{buf: []byte(s)}
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.strict = ta
		case int:
			p.limit = ta
		default:
			return fmt.Errorf("%T is not a valid argument type", a)
		}
	}
	err := p.parse()

	return err
}

func ValidateReader(r io.Reader, args ...interface{}) (interface{}, error) {

	// TBD

	return nil, nil
}
