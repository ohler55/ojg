// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"

	"github.com/ohler55/ojg/gd"
)

// Parse a string into a gd.Node. Arguments are optional and can be a bool,
// a *ParseOptions, or func(gd.Node) bool.
//
// A bool is used to indicated if the parsing should be limited to one JSON only.
//
func Parse(s string, args ...interface{}) (n gd.Node, err error) {
	p := Parser{buf: []byte(s)} // TBD add handler
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.onlyOne = ta
		case *ParseOptions:
			p.noComment = ta.NoComment
			p.onlyOne = ta.OnlyOne
			// TBD timeformat?
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
	p := Parser{}
	return p.Validate(s, args...)
}

func ValidateReader(r io.Reader, args ...interface{}) (interface{}, error) {

	// TBD

	return nil, nil
}
