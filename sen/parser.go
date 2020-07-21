// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"
	"math"
	"unicode/utf8"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

const (
	stackInitSize = 32 // for container stack { or [
	tmpInitSize   = 32 // for tokens and numbers
	mapInitSize   = 8
	readBufSize   = 4096
)

var emptySlice = []interface{}{}

// Parser is a reusable JSON parser. It can be reused for multiple parsings
// which allows buffer reuse for a performance advantage.
type Parser struct {
	tmp        []byte // used for numbers and strings
	runeBytes  []byte
	stack      []interface{}
	starts     []int
	maps       []map[string]interface{}
	cb         func(interface{}) bool
	resultChan chan interface{}
	line       int
	noff       int // Offset of last newline from start of buf. Can be negative when using a reader.
	ri         int // read index for null, false, and true
	mi         int
	num        gen.Number
	rn         rune
	result     interface{}
	mode       string

	// Reuse maps. Previously returned maps will no longer be valid or rather
	// could be modified during parsing.
	Reuse bool

	// OnlyOne returns an error if more than one JSON is in the string or stream.
	OnlyOne bool
}

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func (p *Parser) Parse(buf []byte, args ...interface{}) (interface{}, error) {
	for _, a := range args {
		switch ta := a.(type) {
		case func(interface{}) bool:
			p.cb = ta
			p.OnlyOne = false
		case chan interface{}:
			p.resultChan = ta
			p.OnlyOne = false
			p.Reuse = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if p.cb == nil && p.resultChan == nil {
		p.OnlyOne = true
	}
	if p.stack == nil {
		p.stack = make([]interface{}, 0, stackInitSize)
		p.tmp = make([]byte, 0, tmpInitSize)
		p.starts = make([]int, 0, 16)
		p.maps = make([]map[string]interface{}, 0, 16)
	} else {
		p.stack = p.stack[:0]
		p.tmp = p.tmp[:0]
		p.starts = p.starts[:0]
	}
	p.result = nil
	p.noff = -1
	p.line = 1
	p.mode = valueMap
	p.mi = 0
	var err error
	// Skip BOM if present.
	if 3 < len(buf) && buf[0] == 0xEF {
		if buf[1] == 0xBB && buf[2] == 0xBF {
			err = p.parseBuffer(buf[3:], true)
		} else {
			return nil, fmt.Errorf("expected BOM at 1:3")
		}
	} else {
		err = p.parseBuffer(buf, true)
	}
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
		case func(interface{}) bool:
			p.cb = ta
			p.OnlyOne = false
		case chan interface{}:
			p.resultChan = ta
			p.OnlyOne = false
			p.Reuse = false
		default:
			return nil, fmt.Errorf("a %T is not a valid option type", a)
		}
	}
	if p.cb == nil && p.resultChan == nil {
		p.OnlyOne = true
	}
	if p.stack == nil {
		p.stack = make([]interface{}, 0, stackInitSize)
		p.tmp = make([]byte, 0, tmpInitSize)
		p.starts = make([]int, 0, 16)
		p.maps = make([]map[string]interface{}, 0, 16)
	} else {
		p.stack = p.stack[:0]
		p.tmp = p.tmp[:0]
		p.starts = p.starts[:0]
	}
	p.result = nil
	p.noff = -1
	p.line = 1
	p.mi = 0
	buf := make([]byte, readBufSize)
	eof := false
	var cnt int
	cnt, err = r.Read(buf)
	buf = buf[:cnt]
	p.mode = valueMap
	if err != nil {
		if err != io.EOF {
			return
		}
		eof = true
	}
	var skip int
	// Skip BOM if present.
	if 3 < len(buf) && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		skip = 3
	}
	for {
		if 0 < skip {
			err = p.parseBuffer(buf[skip:], eof)
		} else {
			err = p.parseBuffer(buf, eof)
		}
		if err != nil {
			p.stack = p.stack[:cap(p.stack)]
			for i := len(p.stack) - 1; 0 <= i; i-- {
				p.stack[i] = nil
			}
			p.stack = p.stack[:0]

			return
		}
		skip = 0
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

func (p *Parser) parseBuffer(buf []byte, last bool) (err error) {
	var b byte
	var i int
	var off int
	depth := len(p.starts)
	for off = 0; off < len(buf); off++ {
		b = buf[off]
		switch p.mode[b] {
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
		case tokenStart:
			start := off
			for i, b = range buf[off:] {
				if tokenMap[b] != tokenOk {
					break
				}
			}
			off += i
			if tokenMap[b] == tokenOk { // end of buf reached
				p.tmp = p.tmp[:0]
				p.tmp = append(p.tmp, buf[start:off+1]...)
				p.mode = tokenMap
				continue
			}
			p.addTokenWith(string(buf[start:off]), off)
			off--
		case strOk:
			p.tmp = append(p.tmp, b)
		case colonColon:
			p.mode = valueMap
			continue
		case skipChar: // skip and continue
			continue
		case openObject:
			if 256 < len(p.mode) {
				switch p.mode[256] {
				case 'n':
					if err = p.add(p.num.AsNum(), off); err != nil {
						return
					}
				case 't':
					p.addToken(off)
				}
			}
			p.starts = append(p.starts, -1)
			var m map[string]interface{}
			if p.Reuse {
				if p.mi < len(p.maps) {
					m = p.maps[p.mi]
					for k := range m {
						delete(m, k)
					}
				} else {
					m = make(map[string]interface{}, mapInitSize)
					p.maps = append(p.maps, m)
				}
				p.mi++
			} else {
				m = make(map[string]interface{}, mapInitSize)
			}
			p.stack = append(p.stack, m)
			depth++
			continue
		case closeObject:
			depth--
			if depth < 0 || 0 <= p.starts[depth] {
				return p.newError(off, "unexpected object close")
			}
			if 256 < len(p.mode) {
				switch p.mode[256] {
				case 'n':
					if err = p.add(p.num.AsNum(), off); err != nil {
						return
					}
				case 't':
					p.addToken(off)
				}
			}
			p.starts = p.starts[0:depth]
			n := p.stack[len(p.stack)-1]
			p.stack = p.stack[:len(p.stack)-1]
			if err = p.add(n, off); err != nil {
				return
			}
		case valDigit:
			p.num.Reset()
			p.mode = digitMap
			p.num.I = uint64(b - '0')
			for i, b = range buf[off+1:] {
				if digitMap[b] != numDigit {
					break
				}
				p.num.I = p.num.I*10 + uint64(b-'0')
				if math.MaxInt64 < p.num.I {
					p.num.FillBig()
					break
				}
			}
			if digitMap[b] == numDigit {
				off++
			}
			off += i
		case valQuote:
			start := off + 1
			if len(buf) <= start {
				p.tmp = p.tmp[:0]
				p.mode = stringMap
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
				p.addString(string(buf[start:off]), off)
			} else {
				p.tmp = p.tmp[:0]
				p.tmp = append(p.tmp, buf[start:off+1]...)
				p.mode = stringMap
				continue
			}
		case numSpc:
			if err = p.add(p.num.AsNum(), off); err != nil {
				return
			}
		case strSlash:
			p.mode = escMap
			continue
		case escOk:
			p.tmp = append(p.tmp, escByteMap[b])
			p.mode = stringMap
			continue
		case val0:
			p.mode = zeroMap
			p.num.Reset()
		case valNeg:
			p.mode = negMap
			p.num.Reset()
			p.num.Neg = true
			continue
		case escU:
			p.mode = uMap
			p.rn = 0
			p.ri = 0
			continue
		case openArray:
			if 256 < len(p.mode) {
				switch p.mode[256] {
				case 'n':
					if err = p.add(p.num.AsNum(), off); err != nil {
						return
					}
				case 't':
					p.addToken(off)
				}
			}
			p.starts = append(p.starts, len(p.stack))
			p.stack = append(p.stack, emptySlice)
			depth++
			continue
		case closeArray:
			depth--
			if depth < 0 || p.starts[depth] < 0 {
				return p.newError(off, "unexpected array close")
			}
			// Only modes with a close array are value, token, and numbers
			// which are all over 256 long.
			switch p.mode[256] {
			case 'n':
				// can not fail appending to an array
				_ = p.add(p.num.AsNum(), off)
			case 't':
				p.addToken(off)
			}
			start := p.starts[len(p.starts)-1] + 1
			p.starts = p.starts[:len(p.starts)-1]
			size := len(p.stack) - start
			n := make([]interface{}, size)
			copy(n, p.stack[start:len(p.stack)])
			p.stack = p.stack[0 : start-1]
			if err = p.add(n, off); err != nil {
				return
			}
			p.mode = valueMap
		case numDot:
			if 0 < len(p.num.BigBuf) {
				p.num.BigBuf = append(p.num.BigBuf, b)
				p.mode = dotMap
				continue
			}
			for i, b = range buf[off+1:] {
				if digitMap[b] != numDigit {
					break
				}
				p.num.Frac = p.num.Frac*10 + uint64(b-'0')
				p.num.Div *= 10.0
				if math.MaxInt64 < p.num.Frac {
					p.num.FillBig()
					break
				}
			}
			off += i
			if digitMap[b] == numDigit {
				off++
			}
			p.mode = fracMap
		case numFrac:
			p.num.AddFrac(b)
			p.mode = fracMap
		case fracE:
			if 0 < len(p.num.BigBuf) {
				p.num.BigBuf = append(p.num.BigBuf, b)
			}
			p.mode = expSignMap
			continue
		case tokenOk:
			p.tmp = append(p.tmp, b)
		case tokenSpc:
			p.addToken(off)
		case tokenColon:
			p.addToken(off)
			p.mode = valueMap
		case tokenNlColon:
			p.addToken(off)
			p.line++
			p.noff = off
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case strQuote:
			p.addString(string(p.tmp), off)
		case numZero:
			p.mode = zeroMap
		case numDigit:
			p.num.AddDigit(b)
		case negDigit:
			p.num.AddDigit(b)
			p.mode = digitMap
		case numNewline:
			if err = p.add(p.num.AsNum(), off); err != nil {
				return
			}
			p.line++
			p.noff = off
			p.mode = valueMap
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case expSign:
			p.mode = expZeroMap
			if b == '-' {
				p.num.NegExp = true
			}
			continue
		case expDigit:
			p.num.AddExp(b)
			p.mode = expMap
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
		case valSlash:
			if 256 < len(p.mode) {
				switch p.mode[256] {
				case 'n':
					if err = p.add(p.num.AsNum(), off); err != nil {
						return
					}
				case 't':
					p.addToken(off)
				}
			}
			p.mode = commentStartMap
		case commentStart:
			p.mode = commentMap
		case commentEnd:
			p.mode = valueMap
		case charErr:
			return p.byteError(off, p.mode, b)
		}
		if depth == 0 && 256 < len(p.mode) && p.mode[256] == 'v' {
			if p.cb == nil && p.resultChan == nil {
				p.result = p.stack[0]
			} else {
				if p.cb != nil {
					p.cb(p.stack[0])
				}
				if p.resultChan != nil {
					p.resultChan <- p.stack[0]
				}
			}
			p.stack = p.stack[:0]
			p.mi = 0
			if p.OnlyOne {
				p.mode = spaceMap
			} else {
				p.mode = valueMap
			}
		}
	}
	if last {
		if 0 < len(p.starts) {
			return p.newError(off, "not closed")
		}
		if len(p.mode) == 256 { // valid finishing maps are one byte longer
			return p.newError(off, "incomplete JSON")
		}
		switch p.mode[256] {
		case 'n': // number
			_ = p.add(p.num.AsNum(), off)
			if p.cb == nil && p.resultChan == nil {
				p.result = p.stack[0]
			} else {
				if p.cb != nil {
					p.cb(p.stack[0])
				}
				if p.resultChan != nil {
					p.resultChan <- p.stack[0]
				}
			}
		case 't': // token
			p.addToken(off)
			if p.cb == nil && p.resultChan == nil {
				p.result = p.stack[0]
			} else {
				if p.cb != nil {
					p.cb(p.stack[0])
				}
				if p.resultChan != nil {
					p.resultChan <- p.stack[0]
				}
			}
		}
	}
	return nil
}

// only for non-string
func (p *Parser) add(n interface{}, off int) error {
	p.mode = valueMap
	if 0 < len(p.starts) {
		if p.starts[len(p.starts)-1] == -1 { // object
			if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
				obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
				obj[string(k)] = n
				p.stack = p.stack[0 : len(p.stack)-1]
			} else {
				return p.newError(off, "expected a key")
			}
		} else { // array
			p.stack = append(p.stack, n)
		}
	} else {
		p.stack = append(p.stack, n)
	}
	return nil
}

func (p *Parser) addToken(off int) {
	s := string(p.tmp)
	p.mode = valueMap
	if 0 < len(p.starts) {
		if p.starts[len(p.starts)-1] == -1 { // object
			if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
				obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
				switch s {
				case "null":
					obj[string(k)] = nil
				case "true":
					obj[string(k)] = true
				case "false":
					obj[string(k)] = false
				default:
					obj[string(k)] = s
				}
				p.stack = p.stack[0 : len(p.stack)-1]
			} else {
				p.stack = append(p.stack, gen.Key(s))
				p.mode = colonMap
			}
			return
		}
	}
	// Array or just a value
	switch s {
	case "null":
		p.stack = append(p.stack, nil)
	case "true":
		p.stack = append(p.stack, true)
	case "false":
		p.stack = append(p.stack, false)
	default:
		p.stack = append(p.stack, s)
	}
}

func (p *Parser) addTokenWith(s string, off int) {
	p.mode = valueMap
	if 0 < len(p.starts) {
		if p.starts[len(p.starts)-1] == -1 { // object
			if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
				obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
				switch s {
				case "null":
					obj[string(k)] = nil
				case "true":
					obj[string(k)] = true
				case "false":
					obj[string(k)] = false
				default:
					obj[string(k)] = s
				}
				p.stack = p.stack[0 : len(p.stack)-1]
			} else {
				p.stack = append(p.stack, gen.Key(s))
				p.mode = colonMap
			}
			return
		}
	}
	// Array or just a value
	switch s {
	case "null":
		p.stack = append(p.stack, nil)
	case "true":
		p.stack = append(p.stack, true)
	case "false":
		p.stack = append(p.stack, false)
	default:
		p.stack = append(p.stack, s)
	}
}

func (p *Parser) addString(s string, off int) {
	p.mode = valueMap
	if 0 < len(p.starts) {
		if p.starts[len(p.starts)-1] == -1 { // object
			if k, ok := p.stack[len(p.stack)-1].(gen.Key); ok {
				obj, _ := p.stack[len(p.stack)-2].(map[string]interface{})
				obj[string(k)] = s
				p.stack = p.stack[0 : len(p.stack)-1]
			} else {
				p.stack = append(p.stack, gen.Key(s))
				p.mode = colonMap
			}
			return
		}
	}
	// Array or just a value
	p.stack = append(p.stack, s)
}

func (p *Parser) newError(off int, format string, args ...interface{}) error {
	return &oj.ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  off - p.noff,
	}
}

func (p *Parser) byteError(off int, mode string, b byte) error {
	err := &oj.ParseError{
		Line:   p.line,
		Column: off - p.noff,
	}
	switch mode {
	case colonMap:
		err.Message = fmt.Sprintf("expected a colon, not '%c'", b)
	case negMap, zeroMap, digitMap, dotMap, fracMap, expSignMap, expZeroMap, expMap:
		err.Message = "invalid number"
	case stringMap:
		err.Message = fmt.Sprintf("invalid JSON character 0x%02x", b)
	case escMap:
		err.Message = fmt.Sprintf("invalid JSON escape character '\\%c'", b)
	case uMap:
		err.Message = fmt.Sprintf("invalid JSON unicode character '%c'", b)
	case spaceMap:
		err.Message = fmt.Sprintf("extra characters after close, '%c'", b)
	default:
		err.Message = fmt.Sprintf("unexpected character '%c'", b)
	}
	return err
}
