// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"strconv"
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

	expect func(*parser) bool

	next     int // next byte to read in the buf
	markBuf  int // mark set for start of some value
	markLine int
	markCol  int
	line     int
	col      int
	limit    int
	depth    int
	cnt      int
	Err      error
}

const (
	bomMode = iota
	tokenMode
)

func (p *parser) mark() {
	p.markBuf = p.next
	p.markLine = p.line
	p.markCol = p.col
}

func (p *parser) unmark() {
	p.markBuf = -1
	p.markLine = -1
	p.markCol = -1
}

func (p *parser) parsex() error {
	if p.skipBOM() == nil {
		p.depth = 0
		p.line = 1
		p.col = 1
		p.unmark()
		p.expect = value
	}
	for {
		if p.skipSpace() != nil {
			break
		}
		if p.depth == 0 && 0 < p.cnt && 0 < p.limit && p.limit <= p.cnt && p.peek() != 0 {
			_ = p.newError("extra characters")
			break
		}
		if !p.expect(p) {
			break
		}
		if p.Err != nil {
			break
		}
	}
	return p.Err
}

func value(p *parser) bool {
	//fmt.Printf("*** value() - %c\n", p.peek())
	switch p.peek() {
	case 0:
		return false
	case '{':
		p.depth++
		if p.objHand != nil {
			p.objHand.ObjectStart()
		} else {
			_ = p.read()
		}
		// TBD keep track of container {} or []
		p.expect = valueOrClose
	case '[':
		p.depth++
		if p.arrayHand != nil {
			p.arrayHand.ArrayStart()
		} else {
			_ = p.read()
		}
		// TBD keep track of container {} or []
		p.expect = valueOrClose
	case 'n':
		if p.readToken("null") == nil && p.nullHand != nil {
			p.nullHand.Null(nil)
		}
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case 't':
		if p.readToken("true") == nil && p.boolHand != nil {
			p.boolHand.Bool(nil, true)
		}
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case 'f':
		if p.readToken("false") == nil && p.boolHand != nil {
			p.boolHand.Bool(nil, false)
		}
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case '"':
		p.readString()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+':
		p.readNum()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	default:
		_ = p.newError("did not expect '%c'", p.peek())
	}
	return true
}

func valueOrClose(p *parser) bool {
	//fmt.Printf("*** valueOrClose() - %c\n", p.peek())
	switch p.peek() {
	case '}':
		// TBD verify container type
		p.depth--
		_ = p.read()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case ']':
		// TBD verify container type
		p.depth--
		_ = p.read()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	default:
		return value(p)
	}
	return true
}

func commaOrClose(p *parser) bool {
	//fmt.Printf("*** commaOrClose() - %c\n", p.peek())
	switch p.peek() {
	case '}':
		// TBD verify container type
		p.depth--
		_ = p.read()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case ']':
		// TBD verify container type
		p.depth--
		_ = p.read()
		if 0 < p.depth {
			p.expect = commaOrClose
		} else {
			p.expect = value
			p.cnt++
		}
	case ',':
		_ = p.read()
		p.expect = value
	default:
		_ = p.newError("expected a comma or close, not '%c'", p.peek())
	}
	return true
}

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

func (p *parser) newError(format string, args ...interface{}) error {
	line := p.line
	col := p.col
	if 0 <= p.markLine {
		line = p.markLine
		col = p.markCol
	}
	p.Err = &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
	return p.Err
}

func (p *parser) wrapError(err error) error {
	line := p.line
	col := p.col
	if 0 <= p.markLine {
		line = p.markLine
		col = p.markCol
	}
	p.Err = &ParseError{
		Message: err.Error(),
		Line:    line,
		Column:  col,
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
