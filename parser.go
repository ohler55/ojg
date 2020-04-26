// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ohler55/ojg/gd"
)

const (
	tmpMinSize   = 32 // for tokens and numbers
	keyMinSize   = 64 // for object keys
	stackMinSize = 64 // for container stack { or [
	readBufSize  = 4096

	bomMode          = 'b'
	valueMode        = 'v'
	afterMode        = 'a'
	nullMode         = 'n'
	trueMode         = 't'
	falseMode        = 'f'
	numMode          = 'N'
	strMode          = 's'
	escMode          = 'e'
	uMode            = 'u'
	keyMode          = 'k'
	colonMode        = ':'
	spaceMode        = ' '
	commentStartMode = '/'
	commentMode      = 'c'
)

// Parser is the core of validation and parsing. It can be reused for multiple
// validations or parsings which allows buffer reuse for a performance
// advantage.
type Parser struct {
	r   io.Reader
	buf []byte

	Err error

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool

	objHand   ObjectHandler
	arrayHand ArrayHandler
	nullHand  NullHandler
	boolHand  BoolHandler
	intHand   IntHandler
	floatHand FloatHandler
	strHand   StrHandler
	keyHand   KeyHandler
	errHand   ErrorHandler
	caller    Caller

	handler interface{}

	errorFun       func(h interface{}, err error, line, col int64)
	objectStartFun func(h interface{})
	objectEndFun   func(h interface{})
	arrayStartFun  func(h interface{})
	arrayEndFun    func(h interface{})
	nullFun        func(h interface{})
	intFun         func(h interface{}, value int64)
	floatFun       func(h interface{}, value float64)
	boolFun        func(h interface{}, value bool)
	strFun         func(h interface{}, key string)
	keyFun         func(h interface{}, key string)
	callFun        func(h interface{})

	key      []byte
	tmp      []byte
	stack    []byte // { or [
	runeHex  []uint32
	ri       int // read index for null, false, and true
	line     int
	noff     int // Offset of last newline from start of buf. Can be negative when using a reader.
	off      int
	mode     int
	nextMode int
	numDot   bool
	numE     bool

	noComment bool
	onlyOne   bool
}

// Validate a JSON string. An error is returned if not valid JSON.
func (p *Parser) Validate(s string) error {
	p.prepStr(s)
	err := p.parse()

	return err
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func (p *Parser) ValidateReader(r io.Reader) error {
	p.prepReader(r)
	return p.parse()
}

func boo(hand interface{}, value bool) {
}

type Boo struct {
}

func (b *Boo) Bool(value bool) {
}

// Parse a JSON string. An error is returned if not valid JSON.
func (p *Parser) Parse(s string) (gd.Node, error) {
	p.boolFun = boo
	p.boolHand = &Boo{}
	p.prepStr(s)
	err := p.parse()

	// TBD

	return nil, err
}

func (p *Parser) prepStr(s string) {
	p.prep()
	p.buf = []byte(s)
}

func (p *Parser) prepReader(r io.Reader) {
	p.prep()
	p.r = r
	p.buf = make([]byte, 0, readBufSize)
}

func (p *Parser) prep() {
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
}

func (p *Parser) SetHandler(h interface{}) {
	p.objHand, _ = h.(ObjectHandler)
	p.arrayHand, _ = h.(ArrayHandler)
	p.nullHand, _ = h.(NullHandler)
	p.boolHand, _ = h.(BoolHandler)
	p.intHand, _ = h.(IntHandler)
	p.floatHand, _ = h.(FloatHandler)
	p.strHand, _ = h.(StrHandler)
	p.keyHand, _ = h.(KeyHandler)
	p.errHand, _ = h.(ErrorHandler)
	p.caller, _ = h.(Caller)
}

// This is a huge function only because there was a significant performance
// improvement by reducing function calls. The code is predominantly switch
// statements with the first layer being the various parsing modes and the
// second level deciding what to do with a byte read while in that mode.
func (p *Parser) parse() error {
	p.mode = valueMode
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
		case valueMode:
			switch b {
			case ' ', '\t', '\r':
				// ignore and continue
			case '\n':
				p.line++
				p.noff = p.off
			case 'n':
				p.mode = nullMode
				p.ri = 0
			case 'f':
				p.mode = falseMode
				p.ri = 0
			case 't':
				p.mode = trueMode
				p.ri = 0
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
				p.mode = numMode
				p.tmp = p.tmp[0:0]
				p.numDot = false
				p.numE = false
				p.tmp = append(p.tmp, b)
			case '"':
				p.mode = strMode
				p.nextMode = afterMode
			case '[':
				p.stack = append(p.stack, '[')
				// TBD arrayOpen handler
			case ']':
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '[' {
					return p.newError("expected an array close")
				}
				p.stack = p.stack[0:depth]
				p.mode = afterMode
				// TBD arrayClose handler
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = keyMode
				// TBD objHand.ObjectStart
			case '}':
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '{' {
					return p.newError("expected an object close")
				}
				p.stack = p.stack[0:depth]
				p.mode = afterMode
				// TBD objClose handler
			case '/':
				if p.noComment {
					return p.newError("comments not allowed")
				}
				p.nextMode = p.mode
				p.mode = commentStartMode
			default:
				return p.newError("unexpected character '%c'", b)
			}
		case afterMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = p.off
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
			case ']':
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
			case '}':
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '{' {
					return p.newError("expected an object close")
				}
				p.stack = p.stack[0:depth]
				// TBD p.objHand.ObjectEnd
			default:
				return p.newError("expected a comma or close, not '%c'", b)
			}
		case keyMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = p.off
			case '"':
				p.mode = strMode
				p.nextMode = colonMode
			case '}':
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '{' {
					return p.newError("expected an object close")
				}
				p.stack = p.stack[0:depth]
				// TBD p.obj.Hand.ObjectClose()
			default:
				return p.newError("expected a string start or object close, not '%c'", b)
			}
		case colonMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = p.off
			case ':':
				p.mode = valueMode
			default:
				return p.newError("expected a colon, not '%c'", b)
			}
		case nullMode:
			p.ri++
			if "null"[p.ri] != b {
				return p.newError("expected null")
			}
			if 3 <= p.ri {
				p.mode = afterMode
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
				p.mode = afterMode
				if false {
					if p.boolHand != nil {
						p.boolHand.Bool(false)
					}
				} else {
					if p.boolFun != nil {
						p.boolFun(p.handler, false)
					}
				}
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError("expected false")
			}
			if 3 <= p.ri {
				p.mode = afterMode
				if false {
					if p.boolHand != nil {
						p.boolHand.Bool(true)
					}
				} else {

					if p.boolFun != nil {
						p.boolFun(p.handler, true)
					}
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
				p.mode = afterMode
			case '\n':
				done = true
				p.line++
				p.noff = p.off
				p.mode = afterMode
			case ',':
				done = true
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
			case ']':
				done = true
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '[' {
					return p.newError("expected an array close")
				}
				p.stack = p.stack[0:depth]
				p.mode = afterMode
				// TBD arrayClose handler
			case '}':
				done = true
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '{' {
					return p.newError("expected an object close")
				}
				p.stack = p.stack[0:depth]
				p.mode = afterMode
				// TBD objClose handler
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
			}
		case strMode:
			if b < 0x20 {
				return p.newError("invalid JSON character 0x%02x", b)
			}
			switch b {
			case '\\':
				p.mode = escMode
			case '"':
				p.mode = p.nextMode
				// TBD if strHand then p.strHand.Str(string(p.tmp))
			default:
				// TBD if strHand then append to p.tmp
			}
		case escMode:
			p.mode = strMode
			switch b {
			case 'n':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\n')
				}
			case '"':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '"')
				}
			case '\\':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\\')
				}
			case '/':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '/')
				}
			case 'b':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\b')
				}
			case 'f':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\f')
				}
			case 'r':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\r')
				}
			case 't':
				if p.strHand != nil {
					p.tmp = append(p.tmp, '\t')
				}
			case 'u':
				p.mode = uMode
				if 0 < cap(p.runeHex) {
					p.runeHex = p.runeHex[0:0]
				} else {
					p.runeHex = make([]uint32, 0, 4)
				}
			default:
				return p.newError("invalid JSON escape character '\\%c'", b)
			}
		case uMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.runeHex = append(p.runeHex, uint32(b-'0'))
			case 'a', 'b', 'c', 'd', 'e', 'f':
				p.runeHex = append(p.runeHex, uint32(b-'a'+10))
			case 'A', 'B', 'C', 'D', 'E', 'F':
				p.runeHex = append(p.runeHex, uint32(b-'A'+10))
			default:
				return p.newError("invalid JSON unicode character '%c'", b)
			}
			if len(p.runeHex) == 4 {
				if p.strHand != nil {
					// TBD build rune then append to p.tmp as bytes
				}
				p.mode = strMode
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
				p.mode = p.nextMode
			}
		case bomMode:
			p.ri++
			if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
				return p.newError("expected BOM")
			}
			if 2 <= p.ri {
				p.mode = valueMode
			}
		}
		if len(p.stack) == 0 && p.mode == afterMode {
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
