// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
)

type parser struct {
	r      io.Reader
	buf    []byte
	strict bool

	// TBD use separate handles for each callback type
	objHand   ObjectHandler
	arrayHand ArrayHandler
	nullHand  NullHandler
	boolHand  BoolHandler
	intHand   IntHandler
	floatHand FloatHandler
	errHand   ErrorHandler
	caller    Caller

	next  int // next byte to read in the buf
	mark  int // mark set for start of some value
	line  int
	col   int
	limit int
	Err   error
}

func (p *parser) parse() error {
	if p.skipBOM() == nil {
		cnt := 0
		depth := 0
		p.line = 1
		p.col = 1
	Top:
		for {
			if p.skipSpace() != nil {
				break
			}
			b := p.peek()
			// TBD maybe move to end of switch
			if b != 0 && 0 < p.limit && p.limit <= cnt {
				_ = p.newError("extra characters")
				break Top
			}
			switch b {
			case 0:
				break Top
			case '{':
				depth++
				if p.objHand != nil {
					p.objHand.ObjectStart()
				} else {
					_ = p.read()
				}
			case '}':
				depth--
				if depth <= 0 {
					if depth < 0 {
						_ = p.newError("extra character after close: '}'")
						break Top
					}
					cnt++
				}
				if p.objHand != nil {
					p.objHand.ObjectEnd()
				} else {
					_ = p.read()
				}
			case '[':
				depth++
				if p.arrayHand != nil {
					p.arrayHand.ArrayStart()
				} else {
					_ = p.read()
				}
			case ']':
				depth--
				if depth < 0 {
					if depth < 0 {
						_ = p.newError("extra character after close: '}'")
						break Top
					}
					cnt++
				}
				if p.arrayHand != nil {
					p.arrayHand.ArrayEnd()
				} else {
					_ = p.read()
				}
			case ',':
				_ = p.read()
				// TBD when XStart set needComma to false
				//  comma needs are based on handler stack
				//  maybe a comma stack to keep track for each

			case 'n':
				if p.readToken("null") != nil {
					break Top
				}
				if p.nullHand != nil {
					p.nullHand.Null(nil)
				}
				if depth == 0 {
					cnt++
				}
			case 't':
				if p.readToken("true") != nil {
					break Top
				}
				if p.boolHand != nil {
					p.boolHand.Bool(nil, true)
				}
				if depth == 0 {
					cnt++
				}
			case 'f':
				if p.readToken("false") != nil {
					break Top
				}
				if p.boolHand != nil {
					p.boolHand.Bool(nil, false)
				}
				if depth == 0 {
					cnt++
				}
			case '"':
				// p.readString
				if depth == 0 {
					cnt++
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+':
				// p.readNum
				if depth == 0 {
					cnt++
				}
			default:
				_ = p.newError("did not expect '%c'", p.peek())
				break Top
			}
			// TBD if caller then call
		}
	}
	return p.Err
}

// Read a byte. Return 0 on error or EOF. Error will be placed in parser.err.
func (p *parser) read() (b byte) {
	b = p.peek()
	switch b {
	case 0:
		p.next--
	case '\n':
		p.line++
		p.col = 1
	default:
		p.col++
	}
	p.next++

	return
}

func (p *parser) peek() (b byte) {
	if len(p.buf) <= p.next {
		if p.r != nil {
			fmt.Printf("*** should read from p.r\n")
			// TBD read some more after make space for the read (slide to mark or start)
			//  update next and mark as needed
		}
		if len(p.buf) <= p.next {
			return // EOF so p.err not set
		}
	}
	b = p.buf[p.next]

	return
}

func (p *parser) newError(format string, args ...interface{}) error {
	p.Err = &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.col,
	}
	return p.Err
}

func (p *parser) wrapError(err error) error {
	p.Err = &ParseError{
		Message: err.Error(),
		Line:    p.line,
		Column:  p.col,
	}
	return p.Err
}

func (p *parser) skipBOM() (err error) {
	// Only a UTF-8 BOM is allowed for JSON.
	if b := p.peek(); b == 0xEF {
		p.next++
		for _, bx := range []byte{0xBB, 0xBF} {
			if err == nil {
				if b = p.read(); b != bx {
					err = p.newError("BOM invalid")
				}
			}
		}
	}
	return
}

func (p *parser) skipSpace() error {
Top:
	for {
		switch p.peek() {
		case 0:
			break Top
		case '/':
			if p.strict {
				_ = p.newError("did not expect '/'")
				break Top
			} else {
				p.next++
				if p.read() != '/' {
					_ = p.newError("did not expect '/'")
					break Top
				}
			Comment:
				for {
					switch p.read() {
					case 0:
						break Top
					case '\n', '\r':
						break Comment
					}
				}
			}
		case ' ', '\n', '\r', '\t':
			// space, continue
			_ = p.read()
		default:
			break Top
		}
	}
	return p.Err
}

func (p *parser) readToken(token string) error {
	for _, b := range []byte(token) {
		if p.read() != b {
			return p.newError("expected '%s'", token)
		}
	}
	return nil
}
