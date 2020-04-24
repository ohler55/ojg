// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"strconv"
)

const (
	tmpMinSize   = 32 // for tokens and numbers
	keyMinSize   = 64 // for object keys
	stackMinSize = 64 // for container stack { or [
	readBufSize  = 4096

	bomMode = iota
	valueMode
	nullMode
	trueMode
	falseMode
	numMode
	strMode
	spaceMode
	commentStartMode
	commentMode
)

const (
	// comma modes are expected modes
	noComma      = iota // expect a value
	closeOrComma        // close, ] or } or a comma
	closeOrValue        // close, ] or } or a value
	closeOrKey
	colonOnly
)

type Parser struct {
	r   io.Reader
	buf []byte

	Err error

	objHand   ObjectHandler
	arrayHand ArrayHandler
	nullHand  NullHandler
	boolHand  BoolHandler
	intHand   IntHandler
	floatHand FloatHandler
	keyHand   KeyHandler
	errHand   ErrorHandler
	caller    Caller

	key     []byte
	tmp     []byte
	stack   []byte // { or [
	ri      int    // read index for null, false, and true
	line    int
	noff    int // Offset of last newline from start of buf. Can be negative when using a reader.
	off     int
	mode    int
	comma   int
	modeFun func(*Parser, byte) error
	numDot  bool
	numE    bool

	noComment bool
	onlyOne   bool
}

func (p *Parser) Validate(s string, args ...interface{}) error {
	p.prepStr(s, nil)
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.onlyOne = ta
		case *ParseOptions:
			p.noComment = ta.NoComment
			p.onlyOne = ta.OnlyOne
		default:
			return fmt.Errorf("%T is not a valid argument type", a)
		}
	}
	err := p.parse()

	return err
}

func (p *Parser) prepStr(s string, handler interface{}) {
	p.prep(handler)
	p.buf = []byte(s)
}

func (p *Parser) prepReader(r io.Reader, handler interface{}) {
	p.prep(handler)
	p.r = r
	p.buf = make([]byte, 0, readBufSize)
}

func (p *Parser) prep(handler interface{}) {
	if cap(p.tmp) < tmpMinSize {
		p.tmp = make([]byte, 0, tmpMinSize)
	} else {
		p.tmp = p.tmp[0:0]
	}
	if cap(p.key) < keyMinSize {
		p.key = make([]byte, 0, keyMinSize)
	} else {
		p.key = p.key[0:0]
	}
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[0:0]
	}
	p.noff = -1
	p.line = 1
	p.Err = nil
	p.objHand, _ = handler.(ObjectHandler)
	p.arrayHand, _ = handler.(ArrayHandler)
	p.nullHand, _ = handler.(NullHandler)
	p.boolHand, _ = handler.(BoolHandler)
	p.intHand, _ = handler.(IntHandler)
	p.floatHand, _ = handler.(FloatHandler)
	p.keyHand, _ = handler.(KeyHandler)
	p.errHand, _ = handler.(ErrorHandler)
	p.caller, _ = handler.(Caller)
}

func (p *Parser) parse() error {
	p.mode = valueMode
	p.comma = noComma // expected comma mode
	if p.r != nil {
		fmt.Printf("*** fill buf\n")
		// TBD read first batch
	}
	// Skip BOM if present.
	if 0 < len(p.buf) && p.buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	var b byte
	for p.off, b = range p.buf {
		switch p.mode {
		case bomMode:
			p.ri++
			if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
				return p.newError("expected BOM")
			}
			if 2 <= p.ri {
				p.mode = valueMode
			}
		case spaceMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = p.off
			default:
				return p.newError("extra characters after close, '%c'", b)
			}
		case commentStartMode:
			if b != '/' {
				return p.newError("unexpected character '%c'", b)
			}
			p.mode = commentMode
		case commentMode:
			if b == '\n' {
				p.line++
				p.noff = p.off
				p.mode = valueMode
			}
		case valueMode:
			switch b {
			case ' ', '\t', '\r':
				// ignore and continue
			case '\n':
				p.line++
				p.noff = p.off
			case ',':
				/*
					if p.comma == noComma {
						return p.newError("unexpected comma")
					}
				*/
				p.comma = noComma
			case 'n':
				/*
					if p.comma != noComma && p.comma != closeOrValue {
						return p.newError("expected a comma or close, not 'n'")
					}
				*/
				p.mode = nullMode
				p.ri = 0
				p.comma = closeOrComma
			case 'f':
				/*
					if p.comma != noComma && p.comma != closeOrValue {
						return p.newError("expected a comma or close, not 'f'")
					}
				*/
				p.mode = falseMode
				p.ri = 0
				p.comma = closeOrComma
			case 't':
				/*
					if p.comma != noComma && p.comma != closeOrValue {
						return p.newError("expected a comma or close, not 't'")
					}
				*/
				p.mode = trueMode
				p.ri = 0
				p.comma = closeOrComma
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				/*
					if p.comma != noComma && p.comma != closeOrValue {
						return p.newError("expected a comma or close, not '%c'", b)
					}
				*/
				p.mode = numMode
				p.tmp = p.tmp[0:0]
				p.numDot = false
				p.numE = false
				p.tmp = append(p.tmp, b)
				p.comma = closeOrComma
			case '"':
				// TBD value or key
			case '[':
				/*
					if p.comma != noComma && p.comma != closeOrValue {
						return p.newError("expected a comma or close, not '['")
					}
				*/
				p.stack = append(p.stack, '[')
				// TBD arrayOpen handler
				p.comma = closeOrValue
			case ']':
				/*
					if p.comma != closeOrComma && p.comma != closeOrValue {
						return p.newError("unexpected close")
					}
				*/
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '[' {
					return p.newError("expected an array close")
				}
				p.stack = p.stack[0:depth]
				// TBD arrayClose handler
				p.comma = closeOrComma
			case '{':
				// TBD need a comma mode for key and :
			case '}':
				// TBD
			case '/':
				if p.noComment {
					return p.newError("comments not allowed")
				}
				p.mode = commentStartMode
			default:
				return p.newError("unexpected character '%c'", b)
			}
		case nullMode:
			p.ri++
			if "null"[p.ri] != b {
				return p.newError("expected null")
			}
			if 3 <= p.ri {
				p.mode = valueMode
				if p.nullHand != nil {
					p.nullHand.Null()
				}
			}
		case falseMode:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError("expected false")
			}
			if 4 <= p.ri {
				p.mode = valueMode
				if p.boolHand != nil {
					p.boolHand.Bool(false)
				}
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError("expected false")
			}
			if 3 <= p.ri {
				p.mode = valueMode
				if p.boolHand != nil {
					p.boolHand.Bool(true)
				}
			}
		case numMode:
			done := false
			p.tmp = append(p.tmp, b)
			switch b {
			case '0':
				// ok as first if no other after
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if len(p.tmp) == 2 && p.tmp[0] == '0' {
					return p.newError("invalid number '%s'", p.tmp)
				}
			case ' ', '\t', '\r':
				done = true
			case '\n':
				done = true
				p.line++
				p.noff = p.off
			case ',':
				done = true
				/*
					if p.comma == noComma {
						return p.newError("unexpected comma")
					}
				*/
				p.comma = noComma
			case ']':
				done = true
				/*
					if p.comma != closeOrComma && p.comma != closeOrValue {
						return p.newError("unexpected close")
					}
				*/
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '[' {
					return p.newError("expected an array close")
				}
				p.stack = p.stack[0:depth]
				// TBD arrayClose handler
				p.comma = closeOrComma
			case '-':
				if 1 < len(p.tmp) {
					prev := p.tmp[len(p.tmp)-2]
					if prev != 'e' && prev != 'E' {
						return p.newError("invalid number '%s'", p.tmp)
					}
				}
			case '.':
				if p.numDot || p.numE {
					return p.newError("invalid number '%s'", p.tmp)
				}
				p.numDot = true
			case 'e', 'E':
				if p.numE {
					return p.newError("invalid number '%s'", p.tmp)
				}
				p.numE = true
			case '+':
				if len(p.tmp) == 1 {
					return p.newError("invalid number '%s'", p.tmp)
				} else {
					prev := p.tmp[len(p.tmp)-2]
					if prev != 'e' && prev != 'E' {
						return p.newError("invalid number '%s'", p.tmp)
					}
				}
			default:
				return p.newError("invalid number '%s'", p.tmp)
			}
			if done {
				if p.numDot || p.numE {
					if p.floatHand != nil {
						f, err := strconv.ParseFloat(string(p.tmp), 64)
						if err != nil {
							return p.wrapError(err)
						}
						p.floatHand.Float(f)
					}
				} else if p.intHand != nil {
					i, err := strconv.ParseInt(string(p.tmp), 10, 64)
					if err != nil {
						return p.wrapError(err)
					}
					p.intHand.Int(i)
				}
				p.mode = valueMode
			}
		case strMode:
			// TBD
		}
		if len(p.stack) == 0 && p.mode == valueMode && p.comma == closeOrComma {
			p.comma = noComma
			if p.onlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	return nil
}

func (p *Parser) newError(format string, args ...interface{}) error {
	p.Err = &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.off - p.noff,
	}
	return p.Err
}

func (p *Parser) wrapError(err error) error {
	p.Err = &ParseError{
		Message: err.Error(),
		Line:    p.line,
		Column:  p.off - p.noff,
	}
	return p.Err
}
