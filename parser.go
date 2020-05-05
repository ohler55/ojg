// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/ohler55/ojg/gd"
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

var emptyGdArray = gd.Array([]gd.Node{})

// Parser a JSON parser. It can be reused for multiple parsings which allows
// buffer reuse for a performance advantage.
type Parser struct {
	tmp         []byte // used for numbers and strings
	stack       []byte // { or [
	runeBytes   []byte
	nstack      []gd.Node
	istack      []interface{}
	arrayStarts []int
	cb          func(gd.Node) bool
	icb         func(interface{}) bool
	ri          int // read index for null, false, and true
	line        int
	noff        int // Offset of last newline from start of buf. Can be negative when using a reader.
	off         int
	num         number
	rn          rune
	mode        byte
	nextMode    byte
	simple      bool

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool
}

// Parse a JSON string. An error is returned if not valid JSON.
func (p *Parser) Parse(b []byte, args ...interface{}) (node gd.Node, err error) {
	var callback func(gd.Node) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(gd.Node) bool:
			callback = ta
			p.OnlyOne = false
		}
	}
	if callback == nil {
		callback = func(n gd.Node) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	p.simple = false
	err = p.parse(b, nil)

	return
}

// ParseSimple a JSON string in to simple types. An error is returned if not valid JSON.
func (p *Parser) ParseSimple(b []byte, args ...interface{}) (data interface{}, err error) {
	var callback func(interface{}) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(interface{}) bool:
			callback = ta
			p.OnlyOne = false
		}
	}
	if callback == nil {
		callback = func(n interface{}) bool {
			data = n
			return false // tells the parser to stop
		}
	}
	p.icb = callback
	p.simple = true
	err = p.parse(b, nil)

	return
}

// ParseReader a JSON io.Reader. An error is returned if not valid JSON.
func (p *Parser) ParseReader(r io.Reader, args ...interface{}) (node gd.Node, err error) {
	var callback func(gd.Node) bool

	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(gd.Node) bool:
			callback = ta
			p.OnlyOne = false
		}
	}
	if callback == nil {
		callback = func(n gd.Node) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	p.simple = false
	err = p.parse(nil, r)

	return
}

// This is a huge function only because there was a significant performance
// improvement by reducing function calls. The code is predominantly switch
// statements with the first layer being the various parsing modes and the
// second level deciding what to do with a byte read while in that mode.
func (p *Parser) parse(buf []byte, r io.Reader) error {
	if cap(p.tmp) < tmpMinSize {
		p.tmp = make([]byte, 0, tmpMinSize)
	} else {
		p.tmp = p.tmp[0:0]
	}
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[0:0]
	}
	if p.simple {
		if cap(p.istack) < 64 {
			p.istack = make([]interface{}, 0, 64)
		}
	} else {
		if cap(p.nstack) < 64 {
			p.nstack = make([]gd.Node, 0, 64)
		}
	}
	if cap(p.arrayStarts) < 64 {
		p.arrayStarts = make([]int, 0, 16)
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	if r != nil {
		if cap(buf) < readBufSize {
			buf = make([]byte, readBufSize)
		}
		buf = buf[:cap(buf)]
		if cnt, err := r.Read(buf); err == nil {
			buf = buf[:cnt]
		} else if err != io.EOF {
			return err
		}
	}
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	var b byte
	for {
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
					p.stack = append(p.stack, '[')
					if p.simple {
						p.arrayStarts = append(p.arrayStarts, len(p.istack))
						p.istack = append(p.istack, emptyGdArray)
					} else {
						p.arrayStarts = append(p.arrayStarts, len(p.nstack))
						p.nstack = append(p.nstack, emptyGdArray)
					}
				case ']':
					if err := p.arrayEnd(); err != nil {
						return err
					}
				case '{':
					p.stack = append(p.stack, '{')
					p.mode = keyMode
					if p.simple {
						n := map[string]interface{}{}
						p.istack = append(p.istack, n)
					} else {
						n := gd.Object{}
						p.nstack = append(p.nstack, n)
					}
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
					if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
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
					if err := p.objectEnd(); err != nil {
						return err
					}
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
					if p.simple {
						p.iadd(nil)
					} else {
						p.nadd(nil)
					}
				}
			case falseMode:
				p.ri++
				if "false"[p.ri] != b {
					return p.newError("expected false")
				}
				if 4 <= p.ri {
					p.mode = afterMode
					if p.simple {
						p.iadd(false)
					} else {
						p.nadd(gd.Bool(false))
					}
				}
			case trueMode:
				p.ri++
				if "true"[p.ri] != b {
					return p.newError("expected false")
				}
				if 3 <= p.ri {
					p.mode = afterMode
					if p.simple {
						p.iadd(true)
					} else {
						p.nadd(gd.Bool(true))
					}
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
					if err := p.appendNum(); err != nil {
						return err
					}
				case '\n':
					p.line++
					p.noff = p.off
					p.mode = afterMode
					if err := p.appendNum(); err != nil {
						return err
					}
				case ',':
					if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
						p.mode = keyMode
					} else {
						p.mode = valueMode
					}
					if err := p.appendNum(); err != nil {
						return err
					}
				case ']':
					if err := p.appendNum(); err != nil {
						return err
					}
					if err := p.arrayEnd(); err != nil {
						return err
					}
				case '}':
					if err := p.appendNum(); err != nil {
						return err
					}
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
					if err := p.appendNum(); err != nil {
						return err
					}
				case '\n':
					p.line++
					p.noff = p.off
					p.mode = afterMode
					if err := p.appendNum(); err != nil {
						return err
					}
				case ',':
					if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
						p.mode = keyMode
					} else {
						p.mode = valueMode
					}
					if err := p.appendNum(); err != nil {
						return err
					}
				case ']':
					if err := p.appendNum(); err != nil {
						return err
					}
					if err := p.arrayEnd(); err != nil {
						return err
					}
				case '}':
					if err := p.appendNum(); err != nil {
						return err
					}
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
					if err := p.appendNum(); err != nil {
						return err
					}
				case '\n':
					p.line++
					p.noff = p.off
					p.mode = afterMode
					if err := p.appendNum(); err != nil {
						return err
					}
				case ',':
					if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
						p.mode = keyMode
					} else {
						p.mode = valueMode
					}
					if err := p.appendNum(); err != nil {
						return err
					}
				case ']':
					if err := p.appendNum(); err != nil {
						return err
					}
					if err := p.arrayEnd(); err != nil {
						return err
					}
				case '}':
					if err := p.appendNum(); err != nil {
						return err
					}
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
					if err := p.appendNum(); err != nil {
						return err
					}
				case '\n':
					p.line++
					p.noff = p.off
					p.mode = afterMode
					if err := p.appendNum(); err != nil {
						return err
					}
				case ',':
					if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
						p.mode = keyMode
					} else {
						p.mode = valueMode
					}
					if err := p.appendNum(); err != nil {
						return err
					}
				case ']':
					if err := p.appendNum(); err != nil {
						return err
					}
					if err := p.arrayEnd(); err != nil {
						return err
					}
				case '}':
					if err := p.appendNum(); err != nil {
						return err
					}
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
						if p.simple {
							p.istack = append(p.istack, nKey(p.tmp))
						} else {
							p.nstack = append(p.nstack, nKey(p.tmp))
						}
					} else {
						if p.simple {
							p.iadd(string(p.tmp))
						} else {
							p.nadd(gd.String(string(p.tmp)))
						}
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
				p.ri++
				if []byte{0xEF, 0xBB, 0xBF}[p.ri] != b {
					return p.newError("expected BOM")
				}
				if 2 <= p.ri {
					p.mode = valueMode
				}
			}
			if len(p.stack) == 0 && (p.mode == afterMode || p.mode == keyMode) {
				if p.simple {
					p.icb(p.istack[0])
				} else {
					p.cb(p.nstack[0])
				}
				if p.OnlyOne {
					p.mode = spaceMode
				} else {
					p.mode = valueMode
				}
			}
		}
		if r != nil {
			buf = buf[:cap(buf)]
			if cnt, err := r.Read(buf); err == nil {
				buf = buf[:cnt]
			} else if err == io.EOF && cnt == 0 {
				break
			} else {
				return err
			}
		} else {
			break
		}
	}
	switch p.mode {
	case afterMode, valueMode:
		if p.simple {
			p.icb(p.istack[0])
		} else {
			p.cb(p.nstack[0])
		}
	case zeroMode, digitMode, fracMode, expMode:
		if err := p.appendNum(); err != nil {
			return err
		}
		if p.simple {
			p.icb(p.istack[0])
		} else {
			p.cb(p.nstack[0])
		}
	default:
		//fmt.Printf("*** final mode: %c\n", p.mode)
		return p.newError("incomplete JSON")
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

func (p *Parser) nadd(n gd.Node) {
	if 2 <= len(p.nstack) {
		if k, ok := p.nstack[len(p.nstack)-1].(nKey); ok {
			obj, _ := p.nstack[len(p.nstack)-2].(gd.Object)
			obj[string(k)] = n
			p.nstack = p.nstack[0 : len(p.nstack)-1]

			return
		}
	}
	p.nstack = append(p.nstack, n)
}

func (p *Parser) iadd(n interface{}) {
	if 2 <= len(p.istack) {
		if k, ok := p.istack[len(p.istack)-1].(nKey); ok {
			obj, _ := p.istack[len(p.istack)-2].(map[string]interface{})
			obj[string(k)] = n
			p.istack = p.istack[0 : len(p.istack)-1]

			return
		}
	}
	p.istack = append(p.istack, n)
}

func (p *Parser) appendNum() error {
	if 0 < len(p.num.bigBuf) {
		if p.simple {
			p.iadd(string(p.num.asBig()))
		} else {
			p.nadd(gd.Big(p.num.asBig()))
		}
	} else if p.num.frac == 0 && p.num.exp == 0 {
		if p.simple {
			p.iadd(p.num.asInt())
		} else {
			p.nadd(gd.Int(p.num.asInt()))
		}
	} else if f, err := p.num.asFloat(); err == nil {
		if p.simple {
			p.iadd(f)
		} else {
			p.nadd(gd.Float(f))
		}
	} else {
		return err
	}
	return nil
}

func (p *Parser) arrayEnd() error {
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
	start := p.arrayStarts[len(p.arrayStarts)-1] + 1
	p.arrayStarts = p.arrayStarts[:len(p.arrayStarts)-1]
	if p.simple {
		size := len(p.istack) - start
		n := make([]interface{}, size)
		copy(n, p.istack[start:len(p.istack)])
		p.istack = p.istack[0 : start-1]
		p.iadd(n)
	} else {
		size := len(p.nstack) - start
		n := gd.Array(make([]gd.Node, size))
		copy(n, p.nstack[start:len(p.nstack)])
		p.nstack = p.nstack[0 : start-1]
		p.nadd(n)
	}
	return nil
}

func (p *Parser) objectEnd() error {
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
	if p.simple {
		n := p.istack[len(p.istack)-1]
		p.istack = p.istack[:len(p.istack)-1]
		p.iadd(n)
	} else {
		n := p.nstack[len(p.nstack)-1]
		p.nstack = p.nstack[:len(p.nstack)-1]
		p.nadd(n)
	}
	return nil
}
