// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/ohler55/ojg/gen"
)

const (
	tmpMinSize   = 32 // for tokens and numbers
	stackMinSize = 32 // for container stack { or [
	readBufSize  = 4096

	bomMode          = 'b'
	valueMode        = 'v'
	afterMode        = 'a'
	nullMode         = 'n'
	trueMode         = 't'
	falseMode        = 'f'
	negMode          = '-'
	zeroMode         = '0'
	digitMode        = 'd'
	dotMode          = '.'
	fracMode         = 'F'
	expSignMode      = '+'
	expZeroMode      = 'X'
	expMode          = 'x'
	strMode          = 's'
	escMode          = 'e'
	uMode            = 'u'
	keyMode          = 'k'
	colonMode        = ':'
	spaceMode        = ' '
	commentStartMode = '/'
	commentMode      = 'c'
)

// Parser a JSON parser. It can be reused for multiple parsings which allows
// buffer reuse for a performance advantage.
type Parser struct {
	tmp       []byte // used for numbers and strings
	runeBytes []byte
	stack     []interface{}
	starts    []int
	cb        func(interface{}) bool
	ri        int // read index for null, false, and true
	line      int
	noff      int // Offset of last newline from start of buf. Can be negative when using a reader.
	off       int
	num       number
	rn        rune
	mode      byte
	nextMode  byte
	onlyOne   bool

	// NoComments returns an error if a comment is encountered.
	NoComment bool
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func (p *Parser) Parse(buf []byte, args ...interface{}) (data interface{}, err error) {
	var callback func(interface{}) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(interface{}) bool:
			callback = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n interface{}) bool {
			data = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]interface{}, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[:0]
		p.starts = p.starts[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	err = p.parseBuffer(buf, true)
	for i := len(p.stack) - 1; 0 <= i; i-- {
		p.stack = nil
	}
	p.stack = p.stack[:0]
	return
}

// ParseReader a JSON io.Reader. An error is returned if not valid JSON.
func (p *Parser) ParseReader(r io.Reader, args ...interface{}) (node interface{}, err error) {
	var callback func(interface{}) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(interface{}) bool:
			callback = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n interface{}) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]interface{}, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[:0]
		p.starts = p.starts[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	buf := make([]byte, readBufSize)
	eof := false
	var cnt int
	cnt, err = r.Read(buf)
	buf = buf[:cnt]
	if err != nil {
		if err != io.EOF {
			return
		}
		eof = true
	}
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	for {
		if err = p.parseBuffer(buf, eof); err != nil {
			return
		}
		if eof {
			break
		}
		buf = buf[:cap(buf)]
		cnt, err = r.Read(buf)
		buf = buf[:cnt]
		if err != nil {
			if err != io.EOF {
				return
			}
			eof = true
		}
	}
	return
}

func (p *Parser) parseBuffer(buf []byte, last bool) error {
	var b byte
	for p.off, b = range buf {
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
			case '-':
				p.mode = negMode
				p.num.reset()
				p.num.neg = true
			case '0':
				p.mode = zeroMode
				p.num.reset()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.reset()
				p.num.i = uint64(b - '0')
			case '"':
				p.mode = strMode
				p.nextMode = afterMode
				p.tmp = p.tmp[0:0]
			case '[':
				p.starts = append(p.starts, len(p.stack))
				p.stack = append(p.stack, emptySlice)
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '{':
				p.starts = append(p.starts, -1)
				p.mode = keyMode
				n := map[string]interface{}{}
				p.stack = append(p.stack, n)
			case '}':
				if err := p.objectEnd(); err != nil {
					return err
				}
			case '/':
				if p.NoComment {
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
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objectEnd(); err != nil {
					return err
				}
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
				p.tmp = p.tmp[0:0]
			case '}':
				// If in key mode } is always okay
				_ = p.objectEnd()
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
				p.iadd(nil)
			}
		case falseMode:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError("expected false")
			}
			if 4 <= p.ri {
				p.mode = afterMode
				p.iadd(false)
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError("expected true")
			}
			if 3 <= p.ri {
				p.mode = afterMode
				p.iadd(true)
			}
		case negMode:
			switch b {
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.addDigit(b)
			default:
				return p.newError("invalid number")
			}
		case zeroMode:
			switch b {
			case '.':
				p.mode = dotMode
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
				p.appendNum()
			case ',':
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case digitMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.addDigit(b)
			case '.':
				p.mode = dotMode
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
				p.appendNum()
			case ',':
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case dotMode:
			if '0' <= b && b <= '9' {
				p.mode = fracMode
				p.num.addFrac(b)
			} else {
				return p.newError("invalid number")
			}
		case fracMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.addFrac(b)
			case 'e', 'E':
				p.mode = expSignMode
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
				p.appendNum()
			case ',':
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case expSignMode:
			switch b {
			case '-':
				p.mode = expZeroMode
				p.num.negExp = true
			case '+':
				p.mode = expZeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = expMode
				p.num.addExp(b)
			default:
				return p.newError("invalid number")
			}
		case expZeroMode:
			if '0' <= b && b <= '9' {
				p.mode = expMode
				p.num.addExp(b)
			} else {
				return p.newError("invalid number")
			}
		case expMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.addExp(b)
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
				p.appendNum()
			case ',':
				if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
					p.mode = keyMode
				} else {
					p.mode = valueMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
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
				if p.mode == colonMode {
					p.stack = append(p.stack, gen.Key(p.tmp))
				} else {
					p.iadd(string(p.tmp))
				}
			default:
				p.tmp = append(p.tmp, b)
			}
		case escMode:
			p.mode = strMode
			switch b {
			case 'n':
				p.tmp = append(p.tmp, '\n')
			case '"':
				p.tmp = append(p.tmp, '"')
			case '\\':
				p.tmp = append(p.tmp, '\\')
			case '/':
				p.tmp = append(p.tmp, '/')
			case 'b':
				p.tmp = append(p.tmp, '\b')
			case 'f':
				p.tmp = append(p.tmp, '\f')
			case 'r':
				p.tmp = append(p.tmp, '\r')
			case 't':
				p.tmp = append(p.tmp, '\t')
			case 'u':
				p.mode = uMode
				p.rn = 0
				p.ri = 0
			default:
				return p.newError("invalid JSON escape character '\\%c'", b)
			}
		case uMode:
			p.ri++
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.rn = p.rn<<4 | rune(b-'0')
			case 'a', 'b', 'c', 'd', 'e', 'f':
				p.rn = p.rn<<4 | rune(b-'a'+10)
			case 'A', 'B', 'C', 'D', 'E', 'F':
				p.rn = p.rn<<4 | rune(b-'A'+10)
			default:
				return p.newError("invalid JSON unicode character '%c'", b)
			}
			if p.ri == 4 {
				if len(p.runeBytes) < 6 {
					p.runeBytes = make([]byte, 6)
				}
				n := utf8.EncodeRune(p.runeBytes, p.rn)
				p.tmp = append(p.tmp, p.runeBytes[:n]...)
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
			if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
				return p.newError("expected BOM")
			}
			p.ri++
			if 3 <= p.ri {
				p.mode = valueMode
			}
		}
		if len(p.starts) == 0 && (p.mode == afterMode || p.mode == keyMode) {
			p.cb(p.stack[0])
			p.stack = p.stack[:0]
			if p.onlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	if last {
		switch p.mode {
		case afterMode, valueMode:
			if 0 < len(p.stack) {
				p.cb(p.stack[0])
			}
		case zeroMode, digitMode, fracMode, expMode:
			p.appendNum()
			if 0 < len(p.stack) {
				p.cb(p.stack[0])
			}
		case spaceMode:
			// just reading white space
		default:
			//fmt.Printf("*** final mode: %c\n", p.mode)
			return p.newError("incomplete JSON")
		}
	}
	return nil
}

func (p *Parser) newError(format string, args ...interface{}) error {
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.off - p.noff,
	}
}

func (p *Parser) iadd(n interface{}) {
	if 2 <= len(p.stack) {
		if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
			obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
			obj[string(k)] = n
			p.stack = p.stack[0 : len(p.stack)-1]

			return
		}
	}
	p.stack = append(p.stack, n)
}

func (p *Parser) appendNum() {
	if 0 < len(p.num.bigBuf) {
		p.iadd(string(p.num.asBig()))
	} else if p.num.frac == 0 && p.num.exp == 0 {
		p.iadd(p.num.asInt())
	} else {
		p.iadd(p.num.asFloat())
	}
}

func (p *Parser) arrayEnd() error {
	depth := len(p.starts)
	if depth == 0 {
		return p.newError("too many closes")
	}
	depth--
	if p.starts[depth] < 0 {
		return p.newError("unexpected array close")
	}
	p.mode = afterMode
	start := p.starts[len(p.starts)-1] + 1
	p.starts = p.starts[:len(p.starts)-1]
	size := len(p.stack) - start
	n := make([]interface{}, size)
	copy(n, p.stack[start:len(p.stack)])
	p.stack = p.stack[0 : start-1]
	p.iadd(n)

	return nil
}

func (p *Parser) objectEnd() error {
	depth := len(p.starts)
	if depth == 0 {
		return p.newError("too many closes")
	}
	depth--
	if 0 <= p.starts[depth] {
		return p.newError("unexpected object close")
	}
	p.starts = p.starts[0:depth]
	p.mode = afterMode
	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	p.iadd(n)

	return nil
}
