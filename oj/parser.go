// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/ohler55/ojg/gen"
)

const (
	stackInitSize = 32 // for container stack { or [
	tmpInitSize   = 32 // for tokens and numbers
	mapInitSize   = 8
	readBufSize   = 4096
)

// Parser is a reusable JSON parser. It can be reused for multiple
// validations or parsings which allows buffer reuse for a performance
// advantage.
type Parser struct {
	tmp       []byte // used for numbers and strings
	runeBytes []byte
	stack     []interface{}
	starts    []int
	cb        func(interface{}) bool
	ri        int // read index for null, false, and true
	line      int
	noff      int // Offset of last newline from start of buf. Can be negative when using a reader.
	num       gen.Number
	rn        rune
	result    interface{}
	mode      string
	nextMode  string

	onlyOne bool

	// NoComments returns an error if a comment is encountered.
	NoComment bool
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func (p *Parser) Parse(buf []byte, args ...interface{}) (interface{}, error) {
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(interface{}) bool:
			p.cb = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if p.cb == nil {
		p.onlyOne = true
	}
	if p.stack == nil {
		p.stack = make([]interface{}, 0, stackInitSize)
		p.tmp = make([]byte, 0, tmpInitSize)
		p.starts = make([]int, 0, 16)
	} else {
		p.stack = p.stack[:0]
		p.tmp = p.tmp[:0]
		p.starts = p.starts[:0]
	}
	p.result = nil
	p.noff = -1
	p.line = 1
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomEFMap
		p.ri = 0
	} else {
		p.mode = valueMap
	}
	err := p.parseBuffer(buf, true)
	p.stack = p.stack[:cap(p.stack)]
	for i := len(p.stack) - 1; 0 <= i; i-- {
		p.stack[i] = nil
	}
	p.stack = p.stack[:0]

	return p.result, err
}

// ParseReader a JSON io.Reader. An error is returned if not valid JSON.
func (p *Parser) ParseReader(r io.Reader, args ...interface{}) (data interface{}, err error) {
	for _, a := range args {
		switch ta := a.(type) {
		case bool:
			p.NoComment = ta
		case func(interface{}) bool:
			p.cb = ta
			p.onlyOne = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if p.cb == nil {
		p.onlyOne = true
	}
	if p.stack == nil {
		p.stack = make([]interface{}, 0, stackInitSize)
		p.tmp = make([]byte, 0, tmpInitSize)
		p.starts = make([]int, 0, 16)
	} else {
		p.stack = p.stack[:0]
		p.tmp = p.tmp[:0]
		p.starts = p.starts[:0]
	}
	p.result = nil
	p.noff = -1
	p.line = 1
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
		p.mode = bomEFMap
		p.ri = 0
	} else {
		p.mode = valueMap
	}
	for {
		if err = p.parseBuffer(buf, eof); err != nil {
			p.stack = p.stack[:cap(p.stack)]
			for i := len(p.stack) - 1; 0 <= i; i-- {
				p.stack[i] = nil
			}
			p.stack = p.stack[:0]

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
	data = p.result

	return
}

// TBD return node and error, node can be nil if callback else return last value
func (p *Parser) parseBuffer(buf []byte, last bool) error {
	var b byte
	var i int
	var off int
	depth := len(p.starts)
	for off = 0; off < len(buf); off++ {
		b = buf[off]
		switch p.mode[b] {
		case skipChar: // skip and continue
			continue
		case skipNewline:
			p.line++
			p.noff = off
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
			continue
		case valNull:
			if off+4 <= len(buf) && string(buf[off:off+4]) == "null" {
				off += 3
				p.mode = afterMap
				p.add(nil)
			} else {
				p.mode = nullMap
				p.ri = 0
			}
		case valTrue:
			if off+4 <= len(buf) && string(buf[off:off+4]) == "true" {
				off += 3
				p.mode = afterMap
				p.add(true)
			} else {
				p.mode = trueMap
				p.ri = 0
			}
		case valFalse:
			if off+5 <= len(buf) && string(buf[off:off+5]) == "false" {
				off += 4
				p.mode = afterMap
				p.add(false)
			} else {
				p.mode = falseMap
				p.ri = 0
			}
		case valNeg:
			p.mode = negMap
			p.num.Reset()
			p.num.Neg = true
			continue
		case val0:
			p.mode = zeroMap
			p.num.Reset()
		case valDigit:
			p.mode = digitMap
			p.num.Reset()
			p.num.I = uint64(b - '0')
		case valQuote:
			start := off + 1
			if len(buf) <= start {
				p.tmp = p.tmp[:0]
				p.mode = stringMap
				p.nextMode = afterMap
				continue
			}
			for i, b = range buf[off+1:] {
				if stringMap[b] != strOk {
					break
				}
			}
			off += i
			if b == '"' {
				off++
				p.add(string(buf[start:off]))
				p.mode = afterMap
			} else {
				p.tmp = p.tmp[:0]
				p.tmp = append(p.tmp, buf[start:off+1]...)
				p.mode = stringMap
				p.nextMode = afterMap
				continue
			}
		case openArray:
			p.starts = append(p.starts, len(p.stack))
			p.stack = append(p.stack, emptySlice)
			depth++
			continue
		case closeArray:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if p.starts[depth] < 0 {
				return p.newError(off, "unexpected array close")
			}
			p.mode = afterMap
			start := p.starts[len(p.starts)-1] + 1
			p.starts = p.starts[:len(p.starts)-1]
			size := len(p.stack) - start
			n := make([]interface{}, size)
			copy(n, p.stack[start:len(p.stack)])
			p.stack = p.stack[0 : start-1]
			p.add(n)
		case openObject:
			p.starts = append(p.starts, -1)
			p.mode = key1Map
			p.stack = append(p.stack, make(map[string]interface{}, mapInitSize))
			depth++
			continue
		case closeObject:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if 0 <= p.starts[depth] {
				return p.newError(off, "unexpected object close")
			}
			p.starts = p.starts[0:depth]
			n := p.stack[len(p.stack)-1]
			p.stack = p.stack[:len(p.stack)-1]
			p.add(n)
			p.mode = afterMap
		case valSlash:
			if p.NoComment {
				return p.newError(off, "comments not allowed")
			}
			p.nextMode = p.mode
			p.mode = commentStartMap
			continue
		case nullOk:
			p.ri++
			if "null"[p.ri] != b {
				return p.newError(off, "expected null")
			}
			if 3 <= p.ri {
				p.add(nil)
				p.mode = afterMap
			}
		case falseOk:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError(off, "expected false")
			}
			if 4 <= p.ri {
				p.add(false)
				p.mode = afterMap
			}
		case trueOk:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError(off, "expected true")
			}
			if 3 <= p.ri {
				p.add(true)
				p.mode = afterMap
			}
		case afterComma:
			if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
				p.mode = keyMap
			} else {
				p.mode = commaMap
			}
			continue
		case keyQuote:
			start := off + 1
			if len(buf) <= start {
				p.tmp = p.tmp[:0]
				p.mode = stringMap
				p.nextMode = colonMap
				continue
			}
			for i, b = range buf[off+1:] {
				if stringMap[b] != strOk {
					break
				}
			}
			off += i
			if b == '"' {
				off++
				p.stack = append(p.stack, gen.Key(buf[start:off]))
				p.mode = colonMap
			} else {
				p.tmp = p.tmp[:0]
				p.tmp = append(p.tmp, buf[start:off+1]...)
				p.mode = stringMap
				p.nextMode = colonMap
			}
			continue
		case colonColon:
			p.mode = valueMap
			continue
		case numSpc:
			p.add(p.num.AsNum())
			p.mode = afterMap
		case numNewline:
			p.add(p.num.AsNum())
			p.line++
			p.noff = off
			p.mode = afterMap
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case numZero:
			p.mode = zeroMap
		case negDigit:
			p.num.AddDigit(b)
			p.mode = digitMap
		case numDigit:
			p.num.AddDigit(b)
		case numDot:
			if 0 < len(p.num.BigBuf) {
				p.num.BigBuf = append(p.num.BigBuf, b)
			}
			p.mode = dotMap
			continue
		case numComma:
			p.add(p.num.AsNum())
			if 0 < len(p.starts) && p.starts[len(p.starts)-1] == -1 {
				p.mode = keyMap
			} else {
				p.mode = commaMap
			}
		case numFrac:
			p.num.AddFrac(b)
			p.mode = fracMap
		case fracE:
			if 0 < len(p.num.BigBuf) {
				p.num.BigBuf = append(p.num.BigBuf, b)
			}
			p.mode = expSignMap
			continue
		case expSign:
			p.mode = expZeroMap
			if b == '-' {
				p.num.NegExp = true
			}
			continue
		case expDigit:
			p.num.AddExp(b)
			p.mode = expMap
		case numCloseArray:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if p.starts[depth] < 0 {
				return p.newError(off, "unexpected array close")
			}
			p.add(p.num.AsNum())
			p.mode = afterMap
			start := p.starts[len(p.starts)-1] + 1
			p.starts = p.starts[:len(p.starts)-1]
			size := len(p.stack) - start
			n := make([]interface{}, size)
			copy(n, p.stack[start:len(p.stack)])
			p.stack = p.stack[0 : start-1]
			p.add(n)
		case numCloseObject:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if 0 <= p.starts[depth] {
				return p.newError(off, "unexpected object close")
			}
			p.add(p.num.AsNum())
			p.starts = p.starts[0:depth]
			n := p.stack[len(p.stack)-1]
			p.stack = p.stack[:len(p.stack)-1]
			p.add(n)
			p.mode = afterMap
		case strOk:
			p.tmp = append(p.tmp, b)
		case strSlash:
			p.mode = escMap
			continue
		case strQuote:
			p.mode = p.nextMode
			if p.mode[0] == colonErr {
				p.stack = append(p.stack, gen.Key(p.tmp))
			} else {
				p.add(string(p.tmp))
			}
		case escOk:
			switch b {
			case '/':
				p.tmp = append(p.tmp, '/')
			case 'b':
				p.tmp = append(p.tmp, '\b')
			case 'f':
				p.tmp = append(p.tmp, '\f')
			case 'r':
				p.tmp = append(p.tmp, '\r')
			}
			p.mode = stringMap
			continue
		case escN:
			p.tmp = append(p.tmp, '\n')
			p.mode = stringMap
			continue
		case escT:
			p.tmp = append(p.tmp, '\t')
			p.mode = stringMap
			continue
		case escQ:
			p.tmp = append(p.tmp, '"')
			p.mode = stringMap
			continue
		case escBackSlash:
			p.tmp = append(p.tmp, '\\')
			p.mode = stringMap
			continue
		case escU:
			p.mode = uMap
			p.rn = 0
			p.ri = 0
			continue
		case uOk:
			p.ri++
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.rn = p.rn<<4 | rune(b-'0')
			case 'a', 'b', 'c', 'd', 'e', 'f':
				p.rn = p.rn<<4 | rune(b-'a'+10)
			case 'A', 'B', 'C', 'D', 'E', 'F':
				p.rn = p.rn<<4 | rune(b-'A'+10)
			}
			if p.ri == 4 {
				if len(p.runeBytes) < 6 {
					p.runeBytes = make([]byte, 6)
				}
				n := utf8.EncodeRune(p.runeBytes, p.rn)
				p.tmp = append(p.tmp, p.runeBytes[:n]...)
				p.mode = stringMap
			}
			continue
		case commentStart:
			p.mode = commentMap
		case commentEnd:
			p.line++
			p.noff = off
			p.mode = p.nextMode
		case bomEF:
			p.mode = bomBBMap
			continue
		case bomBB:
			p.mode = bomBFMap
			continue
		case bomBF:
			p.mode = valueMap
			continue

		case bomErr:
			return p.newError(off, "expected BOM")
		case valErr, commentErr:
			return p.newError(off, "unexpected character '%c'", b)
		case nullErr:
			return p.newError(off, "expected null")
		case trueErr:
			return p.newError(off, "expected true")
		case falseErr:
			return p.newError(off, "expected false")
		case afterErr:
			return p.newError(off, "expected a comma or close, not '%c'", b)
		case key1Err:
			return p.newError(off, "expected a string start or object close, not '%c'", b)
		case keyErr:
			return p.newError(off, "expected a string start, not '%c'", b)
		case colonErr:
			return p.newError(off, "expected a colon, not '%c'", b)
		case numErr:
			return p.newError(off, "invalid number")
		case strLowErr:
			return p.newError(off, "invalid JSON character 0x%02x", b)
		case strErr:
			return p.newError(off, "invalid JSON unicode character '%c'", b)
		case escErr:
			return p.newError(off, "invalid JSON escape character '\\%c'", b)
		case spcErr:
			return p.newError(off, "extra characters after close, '%c'", b)
		}
		if depth == 0 && 256 < len(p.mode) && p.mode[256] == 'a' {
			if p.cb != nil {
				p.cb(p.stack[0])
				p.stack = p.stack[:0]
			}
			if p.onlyOne {
				p.mode = spaceMap
			} else {
				p.mode = valueMap
			}
		}
	}
	if last {
		if len(p.mode) == 256 { // valid finishing maps are one byte longer
			return p.newError(off, "incomplete JSON")
		}
		switch p.mode[256] {
		case 'a':
			/*
				// never gets here
				if p.cb == nil {
					p.result = p.stack[0]
				} else {
					p.cb(p.stack[0])
				}
			*/
		case 'n':
			p.add(p.num.AsNum())
			if 0 < len(p.stack) {
				if p.cb == nil {
					p.result = p.stack[0]
				} else {
					p.cb(p.stack[0])
				}
			}
		case 's': // reading space
			if 0 < len(p.stack) {
				if p.cb == nil {
					p.result = p.stack[0]
				} else {
					p.cb(p.stack[0])
				}
			}
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

func (p *Parser) add(n interface{}) {
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
