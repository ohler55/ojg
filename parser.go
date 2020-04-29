// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/ohler55/ojg/gd"
)

const (
	tmpMinSize   = 32 // for tokens and numbers
	keyMinSize   = 32 // for object keys
	stackMinSize = 32 // for container stack { or [
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

	validateEval = 'v'
	nodeEval     = 'n'
	simpleEval   = 's'
)

// Parser is the core of validation and parsing. It can be reused for multiple
// validations or parsings which allows buffer reuse for a performance
// advantage.
type Parser struct {
	buf         []byte
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
	mode        int
	nextMode    int
	eval        int
	rn          rune
	numDot      bool
	numE        bool

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool
}

// Validate a JSON string. An error is returned if not valid JSON.
func (p *Parser) Validate(s string) error {
	p.buf = []byte(s)
	p.eval = validateEval

	return p.parse()
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func (p *Parser) ValidateReader(r io.Reader) error {
	p.r = r
	p.buf = make([]byte, 0, readBufSize)
	p.eval = validateEval
	return p.parse()
}

// Parse a JSON string. An error is returned if not valid JSON.
func (p *Parser) Parse(s string, args ...interface{}) (node gd.Node, err error) {
	p.buf = []byte(s)

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
	if cap(p.vstack) < 64 {
		p.vstack = make([]gd.Node, 0, 64)
	}
	if cap(p.arrayStarts) < 64 {
		p.arrayStarts = make([]int, 0, 16)
	}
	if callback == nil {
		callback = func(n gd.Node) bool {
			node = n
			return false // tells the parser to stop
		}
	}
	p.cb = callback
	p.eval = nodeEval
	err = p.parse()

	return
}

// This is a huge function only because there was a significant performance
// improvement by reducing function calls. The code is predominantly switch
// statements with the first layer being the various parsing modes and the
// second level deciding what to do with a byte read while in that mode.
func (p *Parser) parse() error {
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
	p.noff = -1
	p.line = 1
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
				p.tmp = p.tmp[0:0]
			case '[':
				p.stack = append(p.stack, '[')
				switch p.eval {
				case nodeEval:
					p.arrayStarts = append(p.arrayStarts, len(p.vstack))
				case simpleEval:
					// TBD
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
				p.mode = afterMode
				switch p.eval {
				case nodeEval:
					start := p.arrayStarts[len(p.arrayStarts)-1]
					size := len(p.vstack) - start
					n := gd.Array(make([]gd.Node, size))
					copy(n, p.vstack[start:len(p.vstack)])
					p.vstack = p.vstack[0:start]
					p.add(n)
				case simpleEval:
					// TBD
				}
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = keyMode
				switch p.eval {
				case nodeEval:
					top := len(p.vstack) == 0
					n := gd.Object{}
					p.add(n)
					if !top {
						p.vstack = append(p.vstack, n)
					}
				case simpleEval:
					// TBD
				}
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
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '[' {
					return p.newError("expected an array close")
				}
				p.stack = p.stack[0:depth]
				switch p.eval {
				case nodeEval:
					start := p.arrayStarts[len(p.arrayStarts)-1]
					size := len(p.vstack) - start
					n := gd.Array(make([]gd.Node, size))
					copy(n, p.vstack[start:len(p.vstack)])
					p.vstack = p.vstack[0:start]
					p.add(n)
				case simpleEval:
					// TBD
				}
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
				depth := len(p.stack)
				if depth == 0 {
					return p.newError("too many closes")
				}
				depth--
				if p.stack[depth] != '{' {
					return p.newError("expected an object close")
				}
				p.stack = p.stack[0:depth]
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
				switch p.eval {
				case nodeEval:
					p.add(nil)
				case simpleEval:
					// TBD
				}
			}
		case falseMode:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError("expected false")
			}
			if 4 <= p.ri {
				p.mode = afterMode
				switch p.eval {
				case nodeEval:
					p.add(gd.Bool(false))
				case simpleEval:
					// TBD
				}
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError("expected false")
			}
			if 3 <= p.ri {
				p.mode = afterMode
				switch p.eval {
				case nodeEval:
					p.add(gd.Bool(true))
				case simpleEval:
					// TBD
				}
			}
		case numMode:
			done := false
			switch b {
			case '0':
				// ok as first if no other after
				p.tmp = append(p.tmp, b)
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.tmp = append(p.tmp, b)
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
				switch p.eval {
				case nodeEval:
					start := p.arrayStarts[len(p.arrayStarts)-1]
					size := len(p.vstack) - start
					n := gd.Array(make([]gd.Node, size))
					copy(n, p.vstack[start:len(p.vstack)])
					p.vstack = p.vstack[0:start]
					p.add(n)
				case simpleEval:
					// TBD
				}
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
			case '-':
				p.tmp = append(p.tmp, b)
				if 1 < len(p.tmp) {
					prev := p.tmp[len(p.tmp)-2]
					if prev != 'e' && prev != 'E' {
						return p.newError("invalid number '%s'", p.tmp)
					}
				}
			case '.':
				p.tmp = append(p.tmp, b)
				if p.numDot || p.numE {
					return p.newError("invalid number '%s'", p.tmp)
				}
				p.numDot = true
			case 'e', 'E':
				p.tmp = append(p.tmp, b)
				if p.numE {
					return p.newError("invalid number '%s'", p.tmp)
				}
				p.numE = true
			case '+':
				p.tmp = append(p.tmp, b)
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
			if done && p.eval != validateEval {
				if p.numDot || p.numE {
					f, err := strconv.ParseFloat(string(p.tmp), 64)
					if err != nil {
						return p.wrapError(err)
					}
					switch p.eval {
					case nodeEval:
						p.add(gd.Float(f))
					case simpleEval:
						// TBD
					}
				} else {
					i, err := strconv.ParseInt(string(p.tmp), 10, 64)
					if err != nil {
						return p.wrapError(err)
					}
					switch p.eval {
					case nodeEval:
						p.add(gd.Int(i))
					case simpleEval:
						// TBD
					}
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
				if p.eval != validateEval {
					if p.mode == colonMode {
						switch p.eval {
						case nodeEval:
							p.vstack = append(p.vstack, keyStr(p.tmp))
						case simpleEval:
							// TBD
						}
					} else {
						switch p.eval {
						case nodeEval:
							p.add(gd.String(string(p.tmp)))
						case simpleEval:
							// TBD
						}
					}
				}
			default:
				if p.eval != validateEval {
					p.tmp = append(p.tmp, b)
				}
			}
		case escMode:
			p.mode = strMode
			switch b {
			case 'n':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\n')
				}
			case '"':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '"')
				}
			case '\\':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\\')
				}
			case '/':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '/')
				}
			case 'b':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\b')
				}
			case 'f':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\f')
				}
			case 'r':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\r')
				}
			case 't':
				if p.eval != validateEval {
					p.tmp = append(p.tmp, '\t')
				}
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
				if p.eval != validateEval {
					if len(p.runeBytes) < 6 {
						fmt.Println("*** allocating rune")
						p.runeBytes = make([]byte, 6)
					}
					n := utf8.EncodeRune(p.runeBytes, p.rn)
					p.tmp = append(p.tmp, p.runeBytes[:n]...)
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
		if len(p.stack) == 0 && (p.mode == afterMode || p.mode == keyMode) {
			switch p.eval {
			case nodeEval:
				p.cb(p.vstack[0])
			case simpleEval:
				// TBD
			}
			if p.OnlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	if p.mode == numMode && p.eval != validateEval {
		if p.numDot || p.numE {
			f, err := strconv.ParseFloat(string(p.tmp), 64)
			if err != nil {
				return p.wrapError(err)
			}
			switch p.eval {
			case nodeEval:
				p.add(gd.Float(f))
			case simpleEval:
				// TBD
			}
		} else {
			i, err := strconv.ParseInt(string(p.tmp), 10, 64)
			if err != nil {
				return p.wrapError(err)
			}
			switch p.eval {
			case nodeEval:
				p.add(gd.Int(i))
			case simpleEval:
				// TBD
			}
		}
		switch p.eval {
		case nodeEval:
			p.cb(p.vstack[0])
		case simpleEval:
			// TBD
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
