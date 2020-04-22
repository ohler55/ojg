// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"strconv"
)

const (
	tmpMinSize  = 32 // for tokens and numbers
	keyMinSize  = 64 // for object keys
	readBufSize = 4096

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

	key  []byte
	tmp  []byte
	ti   int // tmp index
	line int
	col  int

	depth int
	mode  int
	comma int

	numDot bool
	numE   bool

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
	p.col = 0
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
	p.depth = 0
	if p.r != nil {
		fmt.Printf("*** fill buf\n")
		// TBD read first batch
	}
	// Skip BOM if present.
	if 0 < len(p.buf) && p.buf[0] == 0xEF {
		p.mode = bomMode
		p.ti = 3
		p.tmp = p.tmp[0:0]
	}
	for _, b := range p.buf {
		p.col++
		switch p.mode {
		case bomMode:
			if 0 < p.ti {
				p.tmp = append(p.tmp, b)
				p.ti--
			} else if p.tmp[0] != 0xEF || p.tmp[1] != 0xBB || p.tmp[2] != 0xBF {
				return p.newError("invalid BOM, not UTF-8")
			} else {
				p.mode = valueMode
			}
		case spaceMode:
			switch b {
			case 0:
				return nil
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.col = 0
			default:
				return p.newError("extra characters after close, '%c'", b)
			}
		case commentStartMode:
			if b != '/' {
				return p.newError("unexpected character '%c'", b)
			}
			p.mode = commentMode
		case commentMode:
			switch b {
			case 0:
				return nil
			case '\n':
				p.line++
				p.col = 0
				p.mode = valueMode
			}
		case valueMode:
			switch b {
			case 0:
				if 0 < p.depth {
					return p.newError("element not closed")
				}
				return nil
			case ' ', '\t', '\r':
				// ignore and continue
			case '\n':
				p.line++
				p.col = 0
			case '/':
				if p.noComment {
					return p.newError("comments not allowed")
				}
				p.mode = commentStartMode
			case ',':
				if p.comma == noComma {
					return p.newError("unexpected comma")
				}
				p.comma = noComma
			case 'n':
				if p.comma != noComma && p.comma != closeOrValue {
					return p.newError("expected a comma or close, not 'n'")
				}
				p.mode = nullMode
				p.ti = 4
				p.tmp = p.tmp[0:0]
				p.comma = closeOrComma
			case 'f':
				if p.comma != noComma && p.comma != closeOrValue {
					return p.newError("expected a comma or close, not 'f'")
				}
				p.mode = falseMode
				p.ti = 5
				p.tmp = p.tmp[0:0]
				p.comma = closeOrComma
			case 't':
				if p.comma != noComma && p.comma != closeOrValue {
					return p.newError("expected a comma or close, not 't'")
				}
				p.mode = trueMode
				p.ti = 4
				p.tmp = p.tmp[0:0]
				p.comma = closeOrComma
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				if p.comma != noComma && p.comma != closeOrValue {
					return p.newError("expected a comma or close, not '%c'", b)
				}
				p.mode = numMode
				p.tmp = p.tmp[0:0]
				p.numDot = false
				p.numE = false
				_, _ = p.numByte(b)
				p.comma = closeOrComma
			case '"':
				// TBD value or key
			case '[':
				if p.comma != noComma && p.comma != closeOrValue {
					return p.newError("expected a comma or close, not '['")
				}
				p.depth++
				p.comma = closeOrValue
			case ']':
				if p.comma != closeOrComma && p.comma != closeOrValue {
					return p.newError("unexpected close")
				}
				p.depth--
				// TBD arrayClose handler
				if p.depth <= 0 {
					if p.depth < 0 {
						return p.newError("too many array closes")
					}
					// TBD caller handler
				}
				p.comma = closeOrComma
			case '{':
				// TBD need a comma mode for key and :
			case '}':
				// TBD
			default:
				return p.newError("unexpected character '%c'", b)
			}
		case nullMode:
			p.tmp = append(p.tmp, b)
			p.ti--
			if p.ti == 1 {
				if string(p.tmp) != "ull" {
					return p.newError("n%s is not a valid JSON token", p.tmp)
				}
				p.mode = valueMode
				if p.nullHand != nil {
					p.nullHand.Null()
				}
			}
		case falseMode:
			p.tmp = append(p.tmp, b)
			p.ti--
			if p.ti == 1 {
				if string(p.tmp) != "alse" {
					return p.newError("f%s is not a valid JSON token", p.tmp)
				}
				p.mode = valueMode
				if p.boolHand != nil {
					p.boolHand.Bool(false)
				}
			}
		case trueMode:
			p.tmp = append(p.tmp, b)
			p.ti--
			if p.ti == 1 {
				if string(p.tmp) != "rue" {
					return p.newError("t%s is not a valid JSON token", p.tmp)
				}
				p.mode = valueMode
				if p.boolHand != nil {
					p.boolHand.Bool(true)
				}
			}
		case numMode:
			done, err := p.numByte(b)
			if err != nil {
				return err
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
				if err := p.valueByte(b); err != nil {
					return err
				}
				// TBD have to redo b
			}
		case strMode:
			// TBD
		}
		if p.depth == 0 && p.mode == valueMode && p.comma == closeOrComma {
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

func (p *Parser) numByte(b byte) (done bool, err error) {
	if len(p.tmp) == 0 {
		p.numDot = false
		p.numE = false
	}
	p.tmp = append(p.tmp, b)
	switch b {
	case '-':
		if len(p.tmp) == 1 {
			// okay
		} else {
			prev := p.tmp[len(p.tmp)-2]
			if prev != 'e' && prev != 'E' {
				err = p.newError("invalid number '%s'", p.tmp)
			}
		}
	case '+':
		if len(p.tmp) == 1 {
			err = p.newError("invalid number '%s'", p.tmp)
		} else {
			prev := p.tmp[len(p.tmp)-2]
			if prev != 'e' && prev != 'E' {
				err = p.newError("invalid number '%s'", p.tmp)
			}
		}
	case '.':
		if p.numDot || p.numE {
			err = p.newError("invalid number '%s'", p.tmp)
		}
		p.numDot = true
	case 'e', 'E':
		if p.numE {
			err = p.newError("invalid number '%s'", p.tmp)
		}
		p.numE = true
	case '0':
		// ok as first if no other after
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if len(p.tmp) == 2 && p.tmp[0] == '0' {
			err = p.newError("invalid number '%s'", p.tmp)
		}
	case 0, ',', ']', '}', ' ', '\n':
		if len(p.tmp) < 2 {
			err = p.newError("invalid number '%s'", p.tmp)
		} else {
			prev := p.tmp[len(p.tmp)-2]
			if prev < '0' || '9' < prev {
				err = p.newError("invalid number '%s'", p.tmp)
			}
		}
		done = true
	default:
		err = p.newError("invalid number '%s'", p.tmp)
	}
	// TBD include other byte checkd for 0, ], }, ...
	// try inline

	return
}

func (p *Parser) newError(format string, args ...interface{}) error {
	p.Err = &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.col,
	}
	return p.Err
}

func (p *Parser) wrapError(err error) error {
	p.Err = &ParseError{
		Message: err.Error(),
		Line:    p.line,
		Column:  p.col,
	}
	return p.Err
}

func (p *Parser) valueByte(b byte) error {
	switch b {
	case 0:
		if 0 < p.depth {
			return p.newError("element not closed")
		}
		return nil
	case ' ', '\t', '\r':
		// ignore and continue
	case '\n':
		p.line++
		p.col = 0
	case '/':
		if p.noComment {
			return p.newError("comments not allowed")
		}
		p.mode = commentStartMode
	case ',':
		if p.comma == noComma {
			return p.newError("unexpected comma")
		}
		p.comma = noComma
	case 'n':
		if p.comma != noComma && p.comma != closeOrValue {
			return p.newError("expected a comma or close, not 'n'")
		}
		p.mode = nullMode
		p.ti = 4
		p.tmp = p.tmp[0:0]
		p.comma = closeOrComma
	case 'f':
		if p.comma != noComma && p.comma != closeOrValue {
			return p.newError("expected a comma or close, not 'f'")
		}
		p.mode = falseMode
		p.ti = 5
		p.tmp = p.tmp[0:0]
		p.comma = closeOrComma
	case 't':
		if p.comma != noComma && p.comma != closeOrValue {
			return p.newError("expected a comma or close, not 't'")
		}
		p.mode = trueMode
		p.ti = 4
		p.tmp = p.tmp[0:0]
		p.comma = closeOrComma
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		if p.comma != noComma && p.comma != closeOrValue {
			return p.newError("expected a comma or close, not '%c'", b)
		}
		p.mode = numMode
		p.tmp = p.tmp[0:0]
		p.numDot = false
		p.numE = false
		_, _ = p.numByte(b)
		p.comma = closeOrComma
	case '"':
		// TBD value or key
	case '[':
		if p.comma != noComma && p.comma != closeOrValue {
			return p.newError("expected a comma or close, not '['")
		}
		p.depth++
		p.comma = closeOrValue
	case ']':
		if p.comma != closeOrComma && p.comma != closeOrValue {
			return p.newError("unexpected close")
		}
		p.depth--
		// TBD arrayClose handler
		if p.depth <= 0 {
			if p.depth < 0 {
				return p.newError("too many array closes")
			}
			// TBD caller handler
		}
		p.comma = closeOrComma
	case '{':
		// TBD need a comma mode for key and :
	case '}':
		// TBD
	default:
		return p.newError("unexpected character '%c'", b)
	}
	return nil
}
