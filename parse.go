// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"

	"github.com/ohler55/ojg/gd"
)

// should make callback when end of top
type Handler interface {
	ObjectStart()
	ObjectEnd()
	ArrayStart()
	ArrayEnd()
	Value(key string, p *Parser) // p.Value() if interested
	Error(err error, line, col int64)
}

type Parser struct {
	r       io.Reader
	buf     []byte
	strict  bool
	handler Handler
	line    int64
	col     int64
}

func (p *Parser) Value() interface{} {
	// TBD
	return nil
}

func (p *Parser) parse() (interface{}, error) {
	// TBD
	return nil, nil
}

// TBD options
//  strict - don't allow comments
// callback

// Parse a string into a gd.Node. Arguments are optional and can be a boolean
// to indicated if the parsing should be strict or tolerant. If tolerant then
// C style // comments are allowed. That is the default. If the input includes
// multiple JSON documents then a callback function of the form func(gd.Node)
// bool can be provided. The bool return if true will abort processing.
func Parse(s string, args ...interface{}) (gd.Node, error) {
	p := Parser{buf: []byte(s)} // TBD add handler
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.strict = ta
		case func(gd.Node) bool:
			// TBD set in handler
		default:
			return nil, fmt.Errorf("%T is not a valid argument type", a)
		}
	}
	// TBD

	r, err := p.parse()
	n, _ := r.(gd.Node)

	return n, err
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
