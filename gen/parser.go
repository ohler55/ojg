// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"fmt"
	"io"
	"unicode/utf8"
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
	key1Mode         = 'K'
	keyMode          = 'k'
	colonMode        = ':'
	commaMode        = ','
	spaceMode        = ' '
	commentStartMode = '/'
	commentMode      = 'c'

	//   0123456789abcdef0123456789abcdef
	strMap = "" +
		"................................" + // 0x00
		"oo.ooooooooooooooooooooooooooooo" + // 0x20
		"oooooooooooooooooooooooooooo.ooo" + // 0x40
		"ooooooooooooooooooooooooooooooo." + // 0x60
		"oooooooooooooooooooooooooooooooo" + // 0x80
		"oooooooooooooooooooooooooooooooo" + // 0xa0
		"oooooooooooooooooooooooooooooooo" + // 0xc0
		"oooooooooooooooooooooooooooooooo" //   0xe0

	//   0123456789abcdef0123456789abcdef
	charTypeMap = "" +
		".........ss....................." + // 0x00
		"s...............dddddddddd......" + // 0x20
		"................................" + // 0x40
		"................................" + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
)

// Parser a JSON parser. It can be reused for multiple parsings which allows
// buffer reuse for a performance advantage.
type Parser struct {
	tmp       []byte // used for numbers and strings
	stack     []byte // { or [
	runeBytes []byte
	nstack    []Node
	starts    []int
	cb        func(Node) bool
	ri        int // read index for null, false, and true
	line      int
	noff      int // Offset of last newline from start of buf. Can be negative when using a reader.
	num       Number
	rn        rune
	mode      byte
	nextMode  byte
	onlyOne   bool

	// NoComments returns an error if a comment is encountered.
	NoComment bool
}

func (p *Parser) Parse(buf []byte, args ...interface{}) (node Node, err error) {
	var callback func(Node) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(Node) bool:
			callback = ta
			p.onlyOne = false
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n Node) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]byte, 0, stackMinSize)
		p.nstack = make([]Node, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[0:0]
		p.nstack = p.nstack[:0]
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
func (p *Parser) ParseReader(r io.Reader, args ...interface{}) (node Node, err error) {
	var callback func(Node) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(Node) bool:
			callback = ta
			p.onlyOne = false
		}
	}
	if callback == nil {
		p.onlyOne = true
		callback = func(n Node) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	if cap(p.tmp) < tmpMinSize { // indicates not initialized
		p.tmp = make([]byte, 0, tmpMinSize)
		p.stack = make([]byte, 0, stackMinSize)
		p.nstack = make([]Node, 0, 64)
		p.starts = make([]int, 0, 16)
	} else {
		p.tmp = p.tmp[0:0]
		p.stack = p.stack[0:0]
		p.nstack = p.nstack[:0]
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
			break
		}
		if eof {
			break
		}
		buf = buf[:cap(buf)]
		cnt, err = r.Read(buf)
		buf = buf[:cnt]
		if err != nil {
			if err != io.EOF {
				break
			}
			eof = true
		}
	}
	for i := len(p.stack) - 1; 0 <= i; i-- {
		p.stack = nil
	}
	p.stack = p.stack[0:0]
	return
}

func (p *Parser) parseBuffer(buf []byte, last bool) error {
	var b byte
	var i int
	var off int
	for off = 0; off < len(buf); off++ {
		b = buf[off]
		switch p.mode {
		case valueMode:
			switch b {
			case ' ', '\t', '\r':
				// ignore and continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case 'n':
				if off+4 < len(buf) && string(buf[off:off+4]) == "null" {
					off += 3
					p.mode = afterMode
					p.nadd(nil)
				} else {
					p.mode = nullMode
					p.ri = 0
				}
			case 'f':
				if off+5 < len(buf) && string(buf[off:off+5]) == "false" {
					off += 4
					p.mode = afterMode
					p.nadd(False)
				} else {
					p.mode = falseMode
					p.ri = 0
				}
			case 't':
				if off+4 < len(buf) && string(buf[off:off+4]) == "true" {
					off += 3
					p.mode = afterMode
					p.nadd(True)
				} else {
					p.mode = trueMode
					p.ri = 0
				}
			case '-':
				p.mode = negMode
				p.num.Reset()
				p.num.Neg = true
			case '0':
				p.mode = zeroMode
				p.num.Reset()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.Reset()
				p.num.I = uint64(b - '0')
			case '"':
				start := off + 1
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.nadd(String(buf[start:off]))
					p.mode = afterMode
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					p.nextMode = afterMode
				}
			case '[':
				p.stack = append(p.stack, '[')
				p.starts = append(p.starts, len(p.nstack))
				p.nstack = append(p.nstack, EmptyArray)
			case ']':
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = key1Mode
				n := Object{}
				p.nstack = append(p.nstack, n)
			case '}':
				if err := p.objectEnd(off); err != nil {
					return err
				}
			case '/':
				if p.NoComment {
					return p.newError(off, "comments not allowed")
				}
				p.nextMode = p.mode
				p.mode = commentStartMode
			default:
				return p.newError(off, "unexpected character '%c'", b)
			}
		case commaMode:
			switch b {
			case ' ', '\t', '\r':
				// ignore and continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case 'n':
				if off+4 < len(buf) && string(buf[off:off+4]) == "null" {
					off += 3
					p.mode = afterMode
					p.nadd(nil)
				} else {
					p.mode = nullMode
					p.ri = 0
				}
			case 'f':
				if off+5 < len(buf) && string(buf[off:off+5]) == "false" {
					off += 4
					p.mode = afterMode
					p.nadd(False)
				} else {
					p.mode = falseMode
					p.ri = 0
				}
			case 't':
				if off+4 < len(buf) && string(buf[off:off+4]) == "true" {
					off += 3
					p.mode = afterMode
					p.nadd(True)
				} else {
					p.mode = trueMode
					p.ri = 0
				}
			case '-':
				p.mode = negMode
				p.num.Reset()
				p.num.Neg = true
			case '0':
				p.mode = zeroMode
				p.num.Reset()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.Reset()
				p.num.I = uint64(b - '0')
			case '"':
				start := off + 1
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.nadd(String(buf[start:off]))
					p.mode = afterMode
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					p.nextMode = afterMode
				}
			case '[':
				p.stack = append(p.stack, '[')
				p.starts = append(p.starts, len(p.nstack))
				p.nstack = append(p.nstack, EmptyArray)
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = key1Mode
				n := Object{}
				p.nstack = append(p.nstack, n)
			case '/':
				if p.NoComment {
					return p.newError(off, "comments not allowed")
				}
				p.nextMode = p.mode
				p.mode = commentStartMode
			default:
				return p.newError(off, "unexpected character '%c'", b)
			}
		case afterMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '}':
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "expected a comma or close, not '%c'", b)
			}
		case key1Mode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '"':
				start := off + 1
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.nstack = append(p.nstack, Key(buf[start:off]))
					p.mode = colonMode
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					p.nextMode = colonMode
				}
			case '}':
				// If in key mode } is always okay
				_ = p.objectEnd(off)
			default:
				return p.newError(off, "expected a string start or object close, not '%c'", b)
			}
		case keyMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case '"':
				start := off + 1
				for i, b = range buf[start:] {
					if strMap[b] != 'o' {
						break
					}
				}
				off += i
				if b == '"' {
					off++
					p.nstack = append(p.nstack, Key(buf[start:off]))
					p.mode = colonMode
				} else {
					p.tmp = p.tmp[:0]
					p.tmp = append(p.tmp, buf[start:off+1]...)
					p.mode = strMode
					p.nextMode = colonMode
				}
			default:
				return p.newError(off, "expected a string start, not '%c'", b)
			}
		case colonMode:
			switch b {
			case ' ', '\t', '\r':
				continue
			case '\n':
				p.line++
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ':':
				p.mode = valueMode
			default:
				return p.newError(off, "expected a colon, not '%c'", b)
			}
		case nullMode:
			p.ri++
			if "null"[p.ri] != b {
				return p.newError(off, "expected null")
			}
			if 3 <= p.ri {
				p.mode = afterMode
				p.nadd(nil)
			}
		case falseMode:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError(off, "expected false")
			}
			if 4 <= p.ri {
				p.mode = afterMode
				p.nadd(Bool(false))
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError(off, "expected true")
			}
			if 3 <= p.ri {
				p.mode = afterMode
				p.nadd(Bool(true))
			}
		case negMode:
			switch b {
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
				p.num.AddDigit(b)
			default:
				return p.newError(off, "invalid number")
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
				p.noff = off
				p.mode = afterMode
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case digitMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddDigit(b)
			case '.':
				p.mode = dotMode
				if 0 < len(p.num.BigBuf) {
					p.num.BigBuf = append(p.num.BigBuf, b)
				}
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.mode = afterMode
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case dotMode:
			if '0' <= b && b <= '9' {
				p.mode = fracMode
				p.num.AddFrac(b)
			} else {
				return p.newError(off, "invalid number")
			}
		case fracMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddFrac(b)
			case 'e', 'E':
				p.mode = expSignMode
				if 0 < len(p.num.BigBuf) {
					p.num.BigBuf = append(p.num.BigBuf, b)
				}
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.mode = afterMode
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case expSignMode:
			switch b {
			case '-':
				p.mode = expZeroMode
				p.num.NegExp = true
			case '+':
				p.mode = expZeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = expMode
				p.num.AddExp(b)
			default:
				return p.newError(off, "invalid number")
			}
		case expZeroMode:
			if '0' <= b && b <= '9' {
				p.mode = expMode
				p.num.AddExp(b)
			} else {
				return p.newError(off, "invalid number")
			}
		case expMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.num.AddExp(b)
			case ' ', '\t', '\r':
				p.mode = afterMode
				p.appendNum()
			case '\n':
				p.line++
				p.noff = off
				p.mode = afterMode
				p.appendNum()
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
				p.appendNum()
			case ']':
				p.appendNum()
				if err := p.arrayEnd(off); err != nil {
					return err
				}
			case '}':
				p.appendNum()
				if err := p.objectEnd(off); err != nil {
					return err
				}
			default:
				return p.newError(off, "invalid number")
			}
		case strMode:
			if b < 0x20 {
				return p.newError(off, "invalid JSON character 0x%02x", b)
			}
			switch b {
			case '\\':
				p.mode = escMode
			case '"':
				p.mode = p.nextMode
				if p.mode == colonMode {
					p.nstack = append(p.nstack, Key(p.tmp))
				} else {
					p.nadd(String(string(p.tmp)))
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
				return p.newError(off, "invalid JSON escape character '\\%c'", b)
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
				return p.newError(off, "invalid JSON unicode character '%c'", b)
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
				p.noff = off
				for i, b = range buf[off+1:] {
					if charTypeMap[b] != 's' {
						break
					}
				}
				off += i
			default:
				return p.newError(off, "extra characters after close, '%c'", b)
			}
		case commentStartMode:
			if b != '/' {
				return p.newError(off, "unexpected character '%c'", b)
			}
			p.mode = commentMode
		case commentMode:
			if b == '\n' {
				p.line++
				p.noff = off
				p.mode = p.nextMode
			}
		case bomMode:
			if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
				return p.newError(off, "expected BOM")
			}
			p.ri++
			if 3 <= p.ri {
				p.mode = valueMode
			}
		}
		if len(p.stack) == 0 && p.mode == afterMode {
			p.cb(p.nstack[0])
			p.nstack = p.nstack[:0]
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
			/*
				if 0 < len(p.nstack) {
					p.cb(p.nstack[0])
				}
			*/
		case zeroMode, digitMode, fracMode, expMode:
			p.appendNum()
			if 0 < len(p.nstack) {
				p.cb(p.nstack[0])
			}
		case spaceMode:
			// just reading white space
		default:
			//fmt.Printf("*** final mode: %c\n", p.mode)
			return p.newError(off, "incomplete JSON")
		}
	}
	return nil
}

func (p *Parser) newError(off int, format string, args ...interface{}) error {
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  off - p.noff,
	}
}

func (p *Parser) nadd(n Node) {
	if 2 <= len(p.nstack) {
		if k, ok := p.nstack[len(p.nstack)-1].(Key); ok {
			obj, _ := p.nstack[len(p.nstack)-2].(Object)
			obj[string(k)] = n
			p.nstack = p.nstack[0 : len(p.nstack)-1]

			return
		}
	}
	p.nstack = append(p.nstack, n)
}

func (p *Parser) appendNum() {
	if 0 < len(p.num.BigBuf) {
		p.nadd(Big(p.num.AsBig()))
	} else if p.num.Frac == 0 && p.num.Exp == 0 {
		p.nadd(Int(p.num.AsInt()))
	} else {
		p.nadd(Float(p.num.AsFloat()))
	}
}

func (p *Parser) arrayEnd(off int) error {
	depth := len(p.stack)
	if depth == 0 {
		return p.newError(off, "too many closes")
	}
	depth--
	if p.stack[depth] != '[' {
		return p.newError(off, "unexpected array close")
	}
	p.stack = p.stack[0:depth]
	p.mode = afterMode
	start := p.starts[len(p.starts)-1] + 1
	p.starts = p.starts[:len(p.starts)-1]
	size := len(p.nstack) - start
	n := Array(make([]Node, size))
	copy(n, p.nstack[start:len(p.nstack)])
	p.nstack = p.nstack[0 : start-1]
	p.nadd(n)

	return nil
}

func (p *Parser) objectEnd(off int) error {
	depth := len(p.stack)
	if depth == 0 {
		return p.newError(off, "too many closes")
	}
	depth--
	if p.stack[depth] != '{' {
		return p.newError(off, "unexpected object close")
	}
	p.stack = p.stack[0:depth]
	p.mode = afterMode
	n := p.nstack[len(p.nstack)-1]
	p.nstack = p.nstack[:len(p.nstack)-1]
	p.nadd(n)

	return nil
}
