// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"math"
	"math/big"
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
	numMode          = 'N'
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
	vstack      []gd.Node
	arrayStarts []int
	r           io.Reader
	cb          func(gd.Node) bool
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
	err = p.parse(b)

	return
}

// This is a huge function only because there was a significant performance
// improvement by reducing function calls. The code is predominantly switch
// statements with the first layer being the various parsing modes and the
// second level deciding what to do with a byte read while in that mode.
func (p *Parser) parse(buf []byte) error {
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
	if cap(p.vstack) < 64 {
		p.vstack = make([]gd.Node, 0, 64)
	}
	if cap(p.arrayStarts) < 64 {
		p.arrayStarts = make([]int, 0, 16)
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	if p.r != nil {
		fmt.Printf("*** fill buf\n")
		// TBD read first batch
	}
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
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
				p.stack = append(p.stack, '[')
				if p.simple {
					// TBD
				} else {
					p.arrayStarts = append(p.arrayStarts, len(p.vstack))
					p.vstack = append(p.vstack, emptyGdArray)
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = keyMode
				if p.simple {
					// TBD
				} else {
					n := gd.Object{}
					p.vstack = append(p.vstack, n)
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
					// TBD
				} else {
					p.add(nil)
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
					// TBD
				} else {
					p.add(gd.Bool(false))
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
					// TBD
				} else {
					p.add(gd.Bool(true))
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
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
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
						// TBD
					} else {
						p.vstack = append(p.vstack, keyStr(p.tmp))
					}
				} else {
					if p.simple {
						// TBD
					} else {
						p.add(gd.String(string(p.tmp)))
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
				// TBD
			} else {
				p.cb(p.vstack[0])
			}
			if p.OnlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	switch p.mode {
	case afterMode, valueMode:
		if p.simple {
			// TBD
		} else {
			p.cb(p.vstack[0])
		}
	case zeroMode, digitMode, fracMode, expMode:
		if err := p.appendNum(); err != nil {
			return err
		}
		if p.simple {
			// TBD
		} else {
			p.cb(p.vstack[0])
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

func (p *Parser) wrapError(err error) error {
	return &ParseError{
		Message: err.Error(),
		Line:    p.line,
		Column:  p.off - p.noff,
	}
}

func (p *Parser) add(n gd.Node) {
	if 2 <= len(p.vstack) {
		if k, ok := p.vstack[len(p.vstack)-1].(keyStr); ok {
			obj, _ := p.vstack[len(p.vstack)-2].(gd.Object)
			obj[string(k)] = n
			p.vstack = p.vstack[0 : len(p.vstack)-1]

			return
		}
	}
	p.vstack = append(p.vstack, n)
}

type number struct {
	i      uint64
	frac   uint64
	div    uint64
	exp    uint64
	neg    bool
	negExp bool
	bigBuf []byte
}

func (n *number) reset() {
	n.i = 0
	n.frac = 0
	n.div = 1
	n.neg = false
	n.negExp = false
	if 0 < len(n.bigBuf) {
		n.bigBuf = n.bigBuf[:0]
	}
}

const bigLimit = math.MaxInt64 / 10

func (n *number) addDigit(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.i <= bigLimit {
		n.i = n.i*10 + uint64(b-'0')
		if math.MaxInt64 < n.i {
			// fill bigBuf
		}
	} else { // big
		// fill bigBuf
		// TBD
	}
}

func (n *number) addFrac(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.frac <= bigLimit {
		n.frac = n.frac*10 + uint64(b-'0')
		if math.MaxInt64 < n.frac {
			// fill bigBuf
		}
	} else { // big
		// fill bigBuf
		// TBD
	}
}

func (n *number) addExp(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.exp <= 102 {
		n.exp = n.exp*10 + uint64(b-'0')
		if 1022 < n.exp {
			// fill bigBuf
		}
	} else { // big
		// fill bigBuf
		// TBD
	}
}

func (n *number) asInt() int64 {
	i := int64(n.i)
	if n.neg {
		i = -i
	}
	return i
}

func (n *number) asFloat() (float64, error) {
	f := float64(n.i)
	if 0 < n.frac {
		f += float64(n.frac) / float64(n.div)
	}
	if n.neg {
		f = -f
	}
	if 0 < n.exp {
		x := int(n.exp)
		if n.negExp {
			x = -x
		}
		f *= math.Pow10(int(x))
	}
	return f, nil
}

func (n *number) asBig() (f *big.Float, err error) {
	f, _, err = big.ParseFloat(string(n.bigBuf), 10, 0, big.ToNearestAway)
	return
}

func (p *Parser) appendNum() error {
	if 0 < len(p.num.bigBuf) {
		if p.simple {
			// TBD
		} else {
			// TBD
			//p.add(gd.Big(p.num.asBig()))
		}
	} else if p.num.frac == 0 && p.num.exp == 0 {
		if p.simple {
			// TBD
		} else {
			p.add(gd.Int(p.num.asInt()))
		}
	} else if f, err := p.num.asFloat(); err == nil {
		if p.simple {
			// TBD
		} else {
			p.add(gd.Float(f))
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
	if p.simple {
		// TBD
	} else {
		start := p.arrayStarts[len(p.arrayStarts)-1] + 1
		p.arrayStarts = p.arrayStarts[:len(p.arrayStarts)-1]
		size := len(p.vstack) - start
		n := gd.Array(make([]gd.Node, size))
		copy(n, p.vstack[start:len(p.vstack)])
		p.vstack = p.vstack[0 : start-1]
		p.add(n)
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
	n := p.vstack[len(p.vstack)-1]
	p.vstack = p.vstack[:len(p.vstack)-1]
	p.add(n)

	return nil
}

func (p *Parser) printStack(label string) {
	fmt.Printf("*** stack at %s - %v\n", label, p.arrayStarts)
	for _, v := range p.vstack {
		fmt.Printf("  %T %v\n", v, v)
	}
}
