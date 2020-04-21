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

type parser struct {
	r   io.Reader
	buf []byte

	Err error

	objHand   ObjectHandler
	arrayHand ArrayHandler
	nullHand  NullHandler
	boolHand  BoolHandler
	intHand   IntHandler
	floatHand FloatHandler
	errHand   ErrorHandler
	caller    Caller

	key    []byte
	hasKey bool
	tmp    []byte
	ti     int // tmp index
	line   int
	col    int

	numDot bool
	numE   bool

	noComment bool
	onlyOne   bool
}

func (p *parser) prepStr(s string, handler interface{}) {
	p.prep(handler)
	p.buf = []byte(s)
}

func (p *parser) prepReader(r io.Reader, handler interface{}) {
	p.prep(handler)
	p.r = r
	p.buf = make([]byte, 0, readBufSize)
}

func (p *parser) prep(handler interface{}) {
	p.tmp = make([]byte, 0, tmpMinSize)
	p.key = make([]byte, 0, keyMinSize)
	p.line = 1
	p.objHand, _ = handler.(ObjectHandler)
	p.arrayHand, _ = handler.(ArrayHandler)
	p.nullHand, _ = handler.(NullHandler)
	p.boolHand, _ = handler.(BoolHandler)
	p.intHand, _ = handler.(IntHandler)
	p.floatHand, _ = handler.(FloatHandler)
	p.errHand, _ = handler.(ErrorHandler)
	p.caller, _ = handler.(Caller)
}

func (p *parser) parse() error {
	mode := valueMode
	comma := noComma // expected comma mode
	depth := 0
	if p.r != nil {
		fmt.Printf("*** fill buf\n")
		// TBD read first batch
	}
	// Skip BOM if present.
	if 0 < len(p.buf) && p.buf[0] == 0xEF {
		mode = bomMode
		p.ti = 3
		p.tmp = p.tmp[0:0]
	}
	for _, b := range p.buf {
		p.col++
		switch mode {
		case bomMode:
			if 0 < p.ti {
				p.tmp = append(p.tmp, b)
				p.ti--
			} else if p.tmp[0] != 0xEF || p.tmp[1] != 0xBB || p.tmp[2] != 0xBF {
				return p.newError("invalid BOM, not UTF-8")
			} else {
				mode = valueMode
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
			mode = commentMode
		case commentMode:
			switch b {
			case 0:
				return nil
			case '\n':
				p.line++
				p.col = 0
				mode = valueMode
			}
		case valueMode:
			switch b {
			case 0:
				if 0 < depth {
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
				mode = commentStartMode
			case ',':
				if comma == noComma {
					return p.newError("unexpected comma")
				}
				comma = noComma
			case 'n':
				if comma != noComma && comma != closeOrValue {
					return p.newError("expected a comma or close, not 'n'")
				}
				mode = nullMode
				p.ti = 4
				p.tmp = p.tmp[0:0]
				comma = closeOrComma
			case 'f':
				if comma != noComma && comma != closeOrValue {
					return p.newError("expected a comma or close, not 'f'")
				}
				mode = falseMode
				p.ti = 5
				p.tmp = p.tmp[0:0]
				comma = closeOrComma
			case 't':
				if comma != noComma && comma != closeOrValue {
					return p.newError("expected a comma or close, not 't'")
				}
				mode = trueMode
				p.ti = 4
				p.tmp = p.tmp[0:0]
				comma = closeOrComma
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				if comma != noComma && comma != closeOrValue {
					return p.newError("expected a comma or close, not '%c'", b)
				}
				mode = numMode
				p.tmp = p.tmp[0:0]
				p.numDot = false
				p.numE = false
				_, _ = p.numByte(b)
				comma = closeOrComma
			case '"':
				// TBD value or key
			case '[':
				if comma != noComma && comma != closeOrValue {
					return p.newError("expected a comma or close, not '['")
				}
				depth++
				comma = closeOrValue
			case ']':
				if comma != closeOrComma && comma != closeOrValue {
					return p.newError("unexpected close")
				}
				depth--
				// TBD arrayClose handler
				if depth <= 0 {
					if depth < 0 {
						return p.newError("too many array closes")
					}
					// TBD caller handler
				}
				comma = closeOrComma
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
				mode = valueMode
				if p.nullHand != nil {
					if p.hasKey {
						p.nullHand.KeyNull(string(p.key))
						p.key = p.key[0:0]
						p.hasKey = false
					} else {
						p.nullHand.Null()
					}
				}
			}
		case falseMode:
			p.tmp = append(p.tmp, b)
			p.ti--
			if p.ti == 1 {
				if string(p.tmp) != "alse" {
					return p.newError("f%s is not a valid JSON token", p.tmp)
				}
				mode = valueMode
				if p.boolHand != nil {
					if p.hasKey {
						p.boolHand.KeyBool(string(p.key), false)
						p.key = p.key[0:0]
						p.hasKey = false
					} else {
						p.boolHand.Bool(false)
					}
				}
			}
		case trueMode:
			p.tmp = append(p.tmp, b)
			p.ti--
			if p.ti == 1 {
				if string(p.tmp) != "rue" {
					return p.newError("t%s is not a valid JSON token", p.tmp)
				}
				mode = valueMode
				if p.boolHand != nil {
					if p.hasKey {
						p.boolHand.KeyBool(string(p.key), true)
						p.key = p.key[0:0]
						p.hasKey = false
					} else {
						p.boolHand.Bool(true)
					}
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
						if p.hasKey {
							p.floatHand.KeyFloat(string(p.key), f)
							p.key = p.key[0:0]
							p.hasKey = false
						} else {
							p.floatHand.Float(f)
						}
					}
				} else if p.intHand != nil {
					i, err := strconv.ParseInt(string(p.tmp), 10, 64)
					if err != nil {
						return p.wrapError(err)
					}
					if p.hasKey {
						p.intHand.KeyInt(string(p.key), i)
						p.key = p.key[0:0]
						p.hasKey = false
					} else {
						p.intHand.Int(i)
					}
				}
				mode = valueMode
			}
		case strMode:
			// TBD
		}
		if depth == 0 && mode == valueMode && comma == closeOrComma {
			comma = noComma
			if p.onlyOne {
				mode = spaceMode
			} else {
				mode = valueMode
			}
		}
	}
	return nil
}

func (p *parser) numByte(b byte) (done bool, err error) {
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
		if p.numDot {
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
	case 0, ',', ']', '}':
		prev := p.tmp[len(p.tmp)-1]
		if prev < '0' || '9' < prev {
			err = p.newError("invalid number '%s'", p.tmp)
		}
		done = true
	default:
		err = p.newError("invalid number '%s'", p.tmp)
	}
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

/*
func (p *parser) parse() error {
	if p.skipBOM() == nil {
		cnt := 0
		depth := 0
		p.line = 1
		p.col = 1
		expComma := false
		started := false // true right after { or [
		p.unmark()

	Top:
		for {
			if p.skipSpace() != nil {
				break
			}
			b := p.peek()
			if b != 0 && 0 < p.limit && p.limit <= cnt {
				_ = p.newError("extra characters")
				break
			}
			switch b {
			case 0:
				break Top
			case '{':
				started = true
				expComma = false
				depth++
				if p.objHand != nil {
					p.objHand.ObjectStart()
				} else {
					_ = p.read()
				}
			case '}':
				if !started && !expComma {
					_ = p.newError("extra comma before object close")
					break Top
				}
				started = false
				expComma = true
				depth--
				if depth < 0 {
					_ = p.newError("extra character after close: '}'")
					break Top
				}
				if p.objHand != nil {
					p.objHand.ObjectEnd()
				} else {
					_ = p.read()
				}
			case '[':
				started = true
				expComma = false
				depth++
				if p.arrayHand != nil {
					p.arrayHand.ArrayStart()
				} else {
					_ = p.read()
				}
			case ']':
				if !started && !expComma {
					_ = p.newError("extra comma before array close")
					break Top
				}
				started = false
				expComma = true
				depth--
				if depth < 0 {
					_ = p.newError("extra character after close: ']'")
					break Top
				}
				if p.arrayHand != nil {
					p.arrayHand.ArrayEnd()
				} else {
					_ = p.read()
				}
			case ',':
				_ = p.read()
				if expComma {
					expComma = false
				} else {
					_ = p.newError("did not expect a comma")
				}
			case 'n':
				if expComma {
					_ = p.newError("expected a comma")
					break Top
				} else {
					started = false
					expComma = true
					if p.readToken("null") == nil && p.nullHand != nil {
						p.nullHand.Null(nil)
					}
				}
			case 't':
				if expComma {
					_ = p.newError("expected a comma")
				} else {
					started = false
					expComma = true
					if p.readToken("true") == nil && p.boolHand != nil {
						p.boolHand.Bool(nil, true)
					}
				}
			case 'f':
				if expComma {
					_ = p.newError("expected a comma")
				} else {
					started = false
					expComma = true
					if p.readToken("false") == nil && p.boolHand != nil {
						p.boolHand.Bool(nil, false)
					}
				}
			case '"':
				if expComma {
					_ = p.newError("expected a comma")
				} else {
					started = false
					expComma = true
					p.readString()
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+':
				if expComma {
					_ = p.newError("expected a comma")
				} else {
					started = false
					expComma = true
					p.readNum()
				}
			default:
				_ = p.newError("did not expect '%c'", p.peek())
			}
			if p.Err != nil {
				break
			}
			if depth == 0 {
				cnt++
			}
		}
	}
	return p.Err
}

func (p *parser) readNum() {
	p.mark()
	defer p.unmark()
	float := false
Start:
	for {
		switch p.peek() {
		case '+':
			if p.strict && p.markBuf == p.next {
				_ = p.newError("numbers can not start with a '+'")
				return
			}
			_ = p.read()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			_ = p.read()
		case '.', 'e', 'E':
			float = true
			_ = p.read()
		default:
			break Start
		}
	}
	if p.strict && !float && p.buf[p.markBuf] == '0' && 1 < p.next-p.markBuf {
		_ = p.newError("numbers can not start with a '0' if not 0")
		return
	}
	// TBD better pre-check to avoid creating float or int unless there is a handler
	if float {
		if f, err := strconv.ParseFloat(string(p.buf[p.markBuf:p.next]), 64); err == nil {
			if p.floatHand != nil {
				p.floatHand.Float(nil, f)
			}
		} else {
			_ = p.wrapError(err)
		}
	} else if i, err := strconv.ParseInt(string(p.buf[p.markBuf:p.next]), 10, 64); err == nil {
		if p.intHand != nil {
			p.intHand.Int(nil, i)
		}
	} else {
		_ = p.wrapError(err)
	}
}

func (p *parser) readString() {
	p.mark()
	_ = p.read() // skip over the starting "

	// TBD
	p.unmark()
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
*/
