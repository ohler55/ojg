// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
)

const stackMinSize = 32 // for container stack { or [

// Validator is a reusable JSON validator. It can be reused for multiple
// validations or parsings which allows buffer reuse for a performance
// advantage.
type Validator struct {
	// This and the Parser use the same basic code but without the
	// building. It is a copy since adding the conditionals needed to avoid
	// builing results add 15 to 20% overhead. An additional improvement could
	// be made by not tracking line and column but that would make it
	// validation much less useful.
	stack    []byte // { or [
	ri       int    // read index for null, false, and true
	line     int
	noff     int // Offset of last newline from start of buf. Can be negative when using a reader.
	off      int
	mode     byte
	nextMode byte

	// NoComments returns an error if a comment is encountered.
	NoComment bool

	// OnlyOne returns an error if more than one JSON is in the string or
	// stream.
	OnlyOne bool
}

func (p *Validator) Validate(buf []byte) (err error) {
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
	// Skip BOM if present.
	if 0 < len(buf) && buf[0] == 0xEF {
		p.mode = bomMode
		p.ri = 0
	}
	return p.validateBuffer(buf, true)
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func (p *Validator) ValidateReader(r io.Reader) error {
	if cap(p.stack) < stackMinSize {
		p.stack = make([]byte, 0, stackMinSize)
	} else {
		p.stack = p.stack[:0]
	}
	p.noff = -1
	p.line = 1
	p.mode = valueMode
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
		p.mode = bomMode
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

func (p *Validator) validateBuffer(buf []byte, last bool) error {
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
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
			case '"':
				p.mode = strMode
				p.nextMode = afterMode
			case '[':
				p.stack = append(p.stack, '[')
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = key1Mode
			case '}':
				if err := p.objEnd(); err != nil {
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
		case commaMode:
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
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
			case '"':
				p.mode = strMode
				p.nextMode = afterMode
			case '[':
				p.stack = append(p.stack, '[')
			case '{':
				p.stack = append(p.stack, '{')
				p.mode = key1Mode
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
				// keep going
			case '\n':
				p.line++
				p.noff = p.off
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objEnd(); err != nil {
					return err
				}
			default:
				return p.newError("expected a comma or close, not '%c'", b)
			}
		case key1Mode:
			switch b {
			case ' ', '\t', '\r':
				// keep going
			case '\n':
				p.line++
				p.noff = p.off
			case '"':
				p.mode = strMode
				p.nextMode = colonMode
			case '}':
				_ = p.objEnd()
			default:
				return p.newError("expected a string start or object close, not '%c'", b)
			}
		case keyMode:
			switch b {
			case ' ', '\t', '\r':
				// keep going
			case '\n':
				p.line++
				p.noff = p.off
			case '"':
				p.mode = strMode
				p.nextMode = colonMode
			default:
				return p.newError("expected a string start, not '%c'", b)
			}
		case colonMode:
			switch b {
			case ' ', '\t', '\r':
				// keep going
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
			}
		case falseMode:
			p.ri++
			if "false"[p.ri] != b {
				return p.newError("expected false")
			}
			if 4 <= p.ri {
				p.mode = afterMode
			}
		case trueMode:
			p.ri++
			if "true"[p.ri] != b {
				return p.newError("expected true")
			}
			if 3 <= p.ri {
				p.mode = afterMode
			}
		case negMode:
			switch b {
			case '0':
				p.mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = digitMode
			default:
				return p.newError("invalid number")
			}
		case zeroMode:
			switch b {
			case '.':
				p.mode = dotMode
			case ' ', '\t', '\r':
				p.mode = afterMode
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case digitMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// no change in mode
			case '.':
				p.mode = dotMode
			case ' ', '\t', '\r':
				p.mode = afterMode
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case dotMode:
			if '0' <= b && b <= '9' {
				p.mode = fracMode
			} else {
				return p.newError("invalid number")
			}
		case fracMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// no change in mode
			case 'e', 'E':
				p.mode = expSignMode
			case ' ', '\t', '\r':
				p.mode = afterMode
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objEnd(); err != nil {
					return err
				}
			default:
				return p.newError("invalid number")
			}
		case expSignMode:
			switch b {
			case '-', '+':
				p.mode = expZeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.mode = expMode
			default:
				return p.newError("invalid number")
			}
		case expZeroMode:
			if '0' <= b && b <= '9' {
				p.mode = expMode
			} else {
				return p.newError("invalid number")
			}
		case expMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// okay
			case ' ', '\t', '\r':
				p.mode = afterMode
			case '\n':
				p.line++
				p.noff = p.off
				p.mode = afterMode
			case ',':
				if 0 < len(p.stack) && p.stack[len(p.stack)-1] == '{' {
					p.mode = keyMode
				} else {
					p.mode = commaMode
				}
			case ']':
				if err := p.arrayEnd(); err != nil {
					return err
				}
			case '}':
				if err := p.objEnd(); err != nil {
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
			}
		case escMode:
			switch b {
			case 'n', '"', '\\', '/', 'b', 'f', 'r', 't':
				p.mode = strMode
			case 'u':
				p.mode = uMode
				p.ri = 0
			default:
				return p.newError("invalid JSON escape character '\\%c'", b)
			}
		case uMode:
			p.ri++
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case 'a', 'b', 'c', 'd', 'e', 'f':
			case 'A', 'B', 'C', 'D', 'E', 'F':
			default:
				return p.newError("invalid JSON unicode character '%c'", b)
			}
			if p.ri == 4 {
				p.mode = strMode
			}
		case spaceMode:
			switch b {
			case ' ', '\t', '\r':
				// keep going
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
		if len(p.stack) == 0 && p.mode == afterMode {
			if p.OnlyOne {
				p.mode = spaceMode
			} else {
				p.mode = valueMode
			}
		}
	}
	if last {
		switch p.mode {
		case afterMode, zeroMode, digitMode, fracMode, expMode, valueMode:
			// okay
		default:
			return p.newError("incomplete JSON")
		}
	}
	return nil
}

func (p *Validator) newError(format string, args ...interface{}) error {
	return &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.off - p.noff,
	}
}

func (p *Validator) arrayEnd() error {
	depth := len(p.stack)
	if depth == 0 {
		return p.newError("too many closes")
	}
	depth--
	if p.stack[depth] != '[' {
		return p.newError("unexpected array close")
	}
	p.stack = p.stack[0:depth]
	p.mode = afterMode
	return nil
}

func (p *Validator) objEnd() error {
	depth := len(p.stack)
	if depth == 0 {
		return p.newError("too many closes")
	}
	depth--
	if p.stack[depth] != '{' {
		return p.newError("unexpected object close")
	}
	p.stack = p.stack[0:depth]
	p.mode = afterMode
	return nil
}
