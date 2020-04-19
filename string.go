// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import "fmt"

type sparser struct {
	buf []byte

	next  int // next byte to read in the buf
	line  int
	col   int
	depth int
	Err   error
}

func Valid(buf []byte) error {
	p := sparser{buf: buf, line: 1, col: 1}

	return p.parse()
}

func (p *sparser) parse() error {
	depth := 0
	expComma := false
	started := false // true right after { or [
	token := ""
	tcnt := 0
	for _, b := range p.buf {
		if 0 < len(token) {
			if b != token[tcnt] {
				return p.newError("expected %s", token)
			}
			tcnt++
			if len(token) <= tcnt {
				token = ""
			}
			continue
		}
		switch b {
		case 0:
			return nil
		case ' ', '\n':
			continue
		case '{':
			started = true
			expComma = false
			depth++
		case '}':
			if !started && !expComma {
				return p.newError("extra comma before object close")
			}
			started = false
			expComma = true
			depth--
			if depth < 0 {
				return p.newError("extra character after close: '}'")
			}
		case '[':
			started = true
			expComma = false
			depth++
		case ']':
			if !started && !expComma {
				return p.newError("extra comma before array close")
			}
			started = false
			expComma = true
			depth--
			if depth < 0 {
				return p.newError("extra character after close: ']'")
			}
		case ',':
			if expComma {
				expComma = false
			} else {
				return p.newError("did not expect a comma")
			}
		case 'n':
			if expComma {
				return p.newError("expected a comma")
			}
			started = false
			expComma = true
			token = "null"
			tcnt = 0
		case 't':
			if expComma {
				return p.newError("expected a comma")
			}
			started = false
			expComma = true
			token = "true"
		case 'f':
			if expComma {
				_ = p.newError("expected a comma")
			}
			started = false
			expComma = true
			token = "false"
		case '"':
			if expComma {
				return p.newError("expected a comma")
			}
			started = false
			expComma = true
			//p.readString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+':
			if expComma {
				return p.newError("expected a comma")
			}
			started = false
			expComma = true
			//p.readNum()
		default:
			return p.newError("did not expect '%c'", b)
		}
	}
	return nil
}

func (p *sparser) newError(format string, args ...interface{}) error {
	p.Err = &ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.line,
		Column:  p.col,
	}
	return p.Err
}
