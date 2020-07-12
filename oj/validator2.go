// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
)

// Validator2 is a reusable JSON validator. It can be reused for multiple
// validations or parsings which allows buffer reuse for a performance
// advantage.
type Validator2 struct {
	// This and the Parser use the same basic code but without the
	// building. It is a copy since adding the conditionals needed to avoid
	// builing results add 15 to 20% overhead. An additional improvement could
	// be made by not tracking line and column but that would make it
	// validation much less useful.
	stack    []byte // { or [
	ri       int    // read index for null, false, and true
	line     int
	noff     int // Offset of last newline from start of buf. Can be negative when using a reader.
	mode     string
	nextMode string

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool
}

func (p *Validator2) Validate(buf []byte) (err error) {
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMap
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomBBMap
		p.ri = 0
	}
	return p.validateBuffer(buf, true)
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func (p *Validator2) ValidateReader(r io.Reader) error {
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMap
	buf := make([]byte, readBufSize)
	eof := false
	cnt, err := r.Read(buf)
	buf = buf[:cnt]
	if err != nil {
		if err != io.EOF {
			return err
		}
		eof = true
	}
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomBBMap
		p.ri = 0
	}
	for {
		if err := p.validateBuffer(buf, eof); err != nil {
			return err
		}
		if eof {
			break
		}
		buf = buf[:cap(buf)]
		cnt, err := r.Read(buf)
		buf = buf[:cnt]
		if err != nil {
			if err != io.EOF {
				return err
			}
			eof = true
		}
	}
	return nil
}

func (p *Validator2) validateBuffer(buf []byte, last bool) error {
	var b byte
	var i int
	var off int
	depth := 0
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
			if off+4 < len(buf) && string(buf[off:off+4]) == "null" {
				off += 3
				p.mode = afterMap
			} else {
				p.mode = nullMap
				p.ri = 0
			}
		case valTrue:
			if off+4 < len(buf) && string(buf[off:off+4]) == "true" {
				off += 3
				p.mode = afterMap
			} else {
				p.mode = trueMap
				p.ri = 0
			}
		case valFalse:
			if off+5 < len(buf) && string(buf[off:off+5]) == "false" {
				off += 4
				p.mode = afterMap
			} else {
				p.mode = falseMap
				p.ri = 0
			}
		case valNeg:
			p.mode = negMap
			continue
		case val0:
			p.mode = zeroMap
			continue
		case valDigit:
			p.mode = digitMap
			continue
		case valQuote:
			for i, b = range buf[off+1:] {
				if stringMap[b] != skipChar {
					break
				}
			}
			off += i
			if b == '"' {
				off++
				p.mode = afterMap
			} else {
				p.mode = stringMap
				p.nextMode = afterMap
				continue
			}
		case openArray:
			p.stack = append(p.stack, '[')
			depth++
			continue
		case closeArray:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if p.stack[depth] != '[' {
				return p.newError(off, "unexpected array close")
			}
			p.stack = p.stack[0:depth]
			p.mode = afterMap
		case openObject:
			p.stack = append(p.stack, '{')
			p.mode = key1Map
			depth++
			continue
		case closeObject:
			if depth == 0 {
				return p.newError(off, "too many closes")
			}
			depth--
			if p.stack[depth] != '{' {
				return p.newError(off, "unexpected object close")
			}
			p.stack = p.stack[0:depth]
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
				p.mode = afterMap
			}
		case falseOk:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError(off, "expected false")
			}
			if 4 <= p.ri {
				p.mode = afterMap
			}
		case trueOk:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError(off, "expected true")
			}
			if 3 <= p.ri {
				p.mode = afterMap
			}
		case afterComma:
			if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
				p.mode = keyMap
			} else {
				p.mode = commaMap
			}
			continue
		case keyQuote:
			for i, b = range buf[off+1:] {
				if stringMap[b] != skipChar {
					break
				}
			}
			off += i
			if b == '"' {
				off++
				p.mode = colonMap
			} else {
				p.mode = stringMap
				p.nextMode = colonMap
			}
			continue
		case colonColon:
			p.mode = valueMap
			continue
		case numSpc:
			p.mode = afterMap
		case numNewline:
			p.line++
			p.noff = off
			p.mode = afterMap
			for i, b = range buf[off+1:] {
				if spaceMap[b] != skipChar {
					break
				}
			}
			off += i
		case numDot:
			p.mode = dotMap
			continue
		case numComma:
			if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
				p.mode = keyMap
			} else {
				p.mode = commaMap
			}
		case numFrac:
			p.mode = fracMap
		case fracE:
			p.mode = expSignMap
			continue
		case expSign:
			p.mode = expZeroMap
			continue
		case expDigit:
			p.mode = expMap
		case strSlash:
			p.mode = escMap
			continue
		case strQuote:
			p.mode = p.nextMode
		case escOk:
			p.mode = stringMap
			continue
		case escU:
			p.mode = uMap
			p.ri = 0
			continue
		case uOk:
			p.ri++
			if p.ri == 4 {
				p.mode = stringMap
			}
			continue
		case commentStart:
			p.mode = commentMap
		case commentEnd:
			p.line++
			p.noff = off
			p.mode = p.nextMode
		case bomBB:
			p.mode = bomBFMap
			continue
		case bomBF:
			p.mode = valueMap
			continue

		case bomErr:
			return p.newError(off, "expected BOM")
		case valErr:
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
		if depth == 0 && p.mode == afterMap {
			if p.OnlyOne {
				p.mode = spaceMap
			} else {
				p.mode = valueMap
			}
		}
	}
	if last && len(p.mode) == 256 { // valid finishing maps are one byte longer
		return p.newError(off, "incomplete JSON")
	}
	return nil
}

func (p *Validator2) newError(off int, format string, args ...interface{}) error {
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  off - p.noff,
	}
}
