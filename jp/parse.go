// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"math"
	"strconv"
)

const (
	//   0123456789abcdef0123456789abcdef
	tokenMap = "" +
		"................................" + // 0x00
		"...o.o..........oooooooooooo...o" + // 0x20
		".oooooooooooooooooooooooooo...oo" + // 0x40
		".oooooooooooooooooooooooooooooo." + // 0x60
		"oooooooooooooooooooooooooooooooo" + // 0x80
		"oooooooooooooooooooooooooooooooo" + // 0xa0
		"oooooooooooooooooooooooooooooooo" + // 0xc0
		"oooooooooooooooooooooooooooooooo" //   0xe0

	// o for an operatio
	// v for a value start character
	//   0123456789abcdef0123456789abcdef
	eqMap = "" +
		"................................" + // 0x00
		".ov.v.ovv.oo.o.ovvvvvvvvvv..ooo." + // 0x20
		"v..............................." + // 0x40
		"......v.......v.....v.......o.o." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0
)

// Performance is less a concern with Expr parsing as it is usually done just
// once if performance is important. Alternatively, an Expr can be built using
// function calls or bare structs. Parsing is more for convenience. Using this
// approach over modes only adds 10% so a reasonable penalty for
// maintainability.
type parser struct {
	buf []byte
	pos int
}

// ParseExprString parses a string into an Expr.
func ParseString(s string) (x Expr, err error) {
	return Parse([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func Parse(buf []byte) (x Expr, err error) {
	p := &parser{buf: buf}
	x, err = p.readExpr()
	if err == nil && p.pos < len(buf) {
		err = fmt.Errorf("parse error")
	}
	if err != nil {
		err = fmt.Errorf("%s at %d in %s", err, p.pos+1, buf)
	}
	return
}

func (p *parser) readExpr() (x Expr, err error) {
	x = Expr{}
	var f Frag
	first := true
	lastDescent := false
	for {
		if f, err = p.nextFrag(first, lastDescent); err != nil || f == nil {
			return
		}
		first = false
		if _, ok := f.(Descent); ok {
			lastDescent = true
		} else {
			lastDescent = false
		}
		x = append(x, f)
	}
}

func (p *parser) nextFrag(first, lastDescent bool) (f Frag, err error) {
	if p.pos < len(p.buf) {
		b := p.buf[p.pos]
		p.pos++
		switch b {
		case '$':
			if first {
				f = Root('$')
			}
		case '@':
			if first {
				f = At('@')
			}
		case '.':
			f, err = p.afterDot()
		case '*':
			return Wildcard('*'), nil
		case '[':
			f, err = p.afterBracket()
		case ']':
			// done
		default:
			p.pos--
			if tokenMap[b] == 'o' {
				if first {
					f, err = p.afterDot()
				} else if lastDescent {
					f, err = p.afterDotDot()
				}
			}
		}
		// Any other character is the end of the Expr, figure out later if
		// that is an error.
	}
	return
}

func (p *parser) afterDot() (Frag, error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	var token []byte
	b := p.buf[p.pos]
	p.pos++
	switch b {
	case '*':
		return Wildcard('*'), nil
	case '.':
		return Descent('.'), nil
	default:
		if tokenMap[b] == '.' {
			return nil, fmt.Errorf("an expression fragment can not start with a '%c'", b)
		}
		token = append(token, b)
	}
	for p.pos < len(p.buf) {
		b := p.buf[p.pos]
		p.pos++
		if tokenMap[b] == '.' {
			p.pos--
			break
		}
		token = append(token, b)
	}
	return Child(token), nil
}

func (p *parser) afterDotDot() (Frag, error) {
	var token []byte
	b := p.buf[p.pos]
	p.pos++
	token = append(token, b)
	for p.pos < len(p.buf) {
		b := p.buf[p.pos]
		p.pos++
		if tokenMap[b] == '.' {
			p.pos--
			break
		}
		token = append(token, b)
	}
	return Child(token), nil
}

func (p *parser) afterBracket() (Frag, error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	b := p.skipSpace()
	switch b {
	case '*':
		// expect ]
		b := p.skipSpace()
		if b != ']' {
			return nil, fmt.Errorf("not terminated")
		}
		return Wildcard('#'), nil
	case '\'', '"':
		s := p.readStr(b)
		b = p.skipSpace()
		switch b {
		case ']':
			return Child(s), nil
		case ',':
			return p.readUnion(s, b)
		default:
			return nil, fmt.Errorf("invalid bracket fragment")
		}
	case ':':
		return p.readSlice(0)
	case '?':
		return p.readFilter()
	case '(':
		return nil, fmt.Errorf("scripts not implemented yet")
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var err error
		var i int
		if i, b, err = p.readInt(b); err != nil {
			return nil, err
		}
	Next:
		switch b {
		case ' ':
			b = p.skipSpace()
			goto Next
		case ']':
			return Nth(i), nil
		case ',':
			return p.readUnion(int64(i), b)
		case ':':
			return p.readSlice(i)
		default:
			return nil, fmt.Errorf("invalid bracket fragment")
		}
	default:
		p.pos--
		return nil, fmt.Errorf("parse error")
	}
}

func (p *parser) readInt(b byte) (int, byte, error) {
	// Allow numbers to begin with a zero.
	/*
		if b == '0' {
			if p.pos < len(p.buf) {
				b = p.buf[p.pos]
				p.pos++
			}
			return 0, b, nil
		}
	*/
	neg := b == '-'
	if neg {
		if len(p.buf) <= p.pos {
			return 0, 0, fmt.Errorf("expected a number")
		}
		b = p.buf[p.pos]
		p.pos++
	}
	start := p.pos
	var i int
	for {
		if b < '0' || '9' < b {
			break
		}
		i = i*10 + int(b-'0')
		if len(p.buf) <= p.pos {
			break
		}
		b = p.buf[p.pos]
		p.pos++
	}
	if p.pos == start {
		return 0, 0, fmt.Errorf("expected a number")
	}
	if neg {
		i = -i
	}
	return i, b, nil
}

func (p *parser) readNum(b byte) (interface{}, error) {
	var num []byte

	num = append(num, b)
	// Read digits first
	for p.pos < len(p.buf) {
		b = p.buf[p.pos]
		if b < '0' || '9' < b {
			break
		}
		num = append(num, b)
		p.pos++
	}
	switch b {
	case '.':
		num = append(num, b)
		p.pos++
		for p.pos < len(p.buf) {
			b = p.buf[p.pos]
			if b < '0' || '9' < b {
				break
			}
			num = append(num, b)
			p.pos++
		}
		if b == 'e' || b == 'E' {
			p.pos++
			num = append(num, b)
			if len(p.buf) <= p.pos {
				return 0, fmt.Errorf("expected a number")
			}
			b = p.buf[p.pos]
		} else {
			f, err := strconv.ParseFloat(string(num), 64)
			return f, err
		}
	case 'e', 'E':
		p.pos++
		if len(p.buf) <= p.pos {
			return 0, fmt.Errorf("expected a number")
		}
		num = append(num, b)
		b = p.buf[p.pos]
	default:
		i, err := strconv.ParseInt(string(num), 10, 64)
		return int(i), err
	}
	if b == '+' || b == '-' {
		num = append(num, b)
		p.pos++
		if len(p.buf) <= p.pos {
			return 0, fmt.Errorf("expected a number")
		}
	}
	for p.pos < len(p.buf) {
		b = p.buf[p.pos]
		if b < '0' || '9' < b {
			break
		}
		num = append(num, b)
		p.pos++
	}
	f, err := strconv.ParseFloat(string(num), 64)

	return f, err
}

func (p *parser) readSlice(i int) (Frag, error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	f := Slice{i}
	b := p.buf[p.pos]
	if b == ']' {
		f = append(f, math.MaxInt64)
		p.pos++
		return f, nil
	}
	b = p.skipSpace()
	var err error
	// read the end
	if b == ':' {
		f = append(f, math.MaxInt64)
		if len(p.buf) <= p.pos {
			return nil, fmt.Errorf("not terminated")
		}
		b = p.buf[p.pos]
		p.pos++
		if b != ']' {
			if i, b, err = p.readInt(b); err != nil {
				return nil, err
			}
			f = append(f, i)
		}
	} else if i, b, err = p.readInt(b); err == nil {
		f = append(f, i)
		if b == ':' {
			if len(p.buf) <= p.pos {
				return nil, fmt.Errorf("not terminated")
			}
			b = p.buf[p.pos]
			p.pos++
			if b != ']' {
				if i, b, err = p.readInt(b); err != nil {
					return nil, err
				}
				f = append(f, i)
			}
		}
	}
	if b != ']' {
		return nil, fmt.Errorf("invalid slice syntax")
	}
	return f, nil
}

func (p *parser) readUnion(v interface{}, b byte) (Frag, error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	f := Union{v}
	var err error
	for {
		switch b {
		case ',':
			// next union member
		case ']':
			return f, nil
		default:
			return nil, fmt.Errorf("invalid union syntax")
		}
		b = p.skipSpace()
		switch b {
		case '\'', '"':
			var s string
			s = p.readStr(b)
			b = p.skipSpace()
			f = append(f, s)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var i int
			if i, b, err = p.readInt(b); err != nil {
				return nil, err
			}
			f = append(f, int64(i))
			if b == ' ' {
				b = p.skipSpace()
			}
		default:
			return nil, fmt.Errorf("invalid union syntax")
		}
	}
}

func (p *parser) readStr(term byte) string {
	start := p.pos
	esc := false
	for p.pos < len(p.buf) {
		b := p.buf[p.pos]
		p.pos++
		if b == term && !esc {
			break
		}
		if b == '\\' {
			esc = !esc
		} else {
			esc = false
		}
	}
	return string(p.buf[start : p.pos-1])
}

func (p *parser) readFilter() (*Filter, error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	b := p.buf[p.pos]
	p.pos++
	if b != '(' {
		return nil, fmt.Errorf("expected a '(' in filter")
	}
	eq, err := p.readEquation()
	if err != nil {
		return nil, err
	}
	if len(p.buf) <= p.pos || p.buf[p.pos] != ']' {
		return nil, fmt.Errorf("not terminated")
	}
	p.pos++

	return eq.Filter(), nil
}

func (p *parser) readEquation() (eq *Equation, err error) {
	if len(p.buf) <= p.pos {
		return nil, fmt.Errorf("not terminated")
	}
	eq = &Equation{}

	b := p.nextNonSpace()
	if b == '!' {
		eq.o = not
		p.pos++
		if eq.left, err = p.readEqValue(); err != nil {
			return
		}
		b := p.nextNonSpace()
		if b != ')' {
			return nil, fmt.Errorf("not terminated")
		}
		p.pos++
		return
	}
	if eq.left, err = p.readEqValue(); err != nil {
		return
	}
	if eq.o, err = p.readEqOp(); err != nil {
		return
	}
	if eq.right, err = p.readEqValue(); err != nil {
		return
	}
	for {
		b = p.nextNonSpace()
		if b == ')' {
			p.pos++
			return
		}
		var o *op
		if o, err = p.readEqOp(); err != nil {
			return
		}
		if eq.o.prec <= o.prec {
			eq = &Equation{left: eq, o: o}
			if eq.right, err = p.readEqValue(); err != nil {
				return
			}
		} else {
			eq.right = &Equation{left: eq.right, o: o}
			if eq.right.right, err = p.readEqValue(); err != nil {
				return
			}
		}
	}
}

func (p *parser) readEqValue() (eq *Equation, err error) {
	b := p.nextNonSpace()
	switch b {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var v interface{}
		p.pos++
		if v, err = p.readNum(b); err != nil {
			return
		}
		eq = &Equation{result: v}
	case '\'', '"':
		p.pos++
		var s string
		s = p.readStr(b)
		eq = &Equation{result: s}
	case 'n':
		if err = p.readEqToken([]byte("null")); err != nil {
			return
		}
		eq = &Equation{result: nil}
	case 't':
		if err = p.readEqToken([]byte("true")); err != nil {
			return
		}
		eq = &Equation{result: true}

	case 'f':
		if err = p.readEqToken([]byte("false")); err != nil {
			return
		}
		eq = &Equation{result: false}
	case '@', '$':
		var x Expr
		x, err = p.readExpr()
		eq = &Equation{result: x}
	case '(':
		p.pos++
		eq, err = p.readEquation()
	default:
		err = fmt.Errorf("expected a value")
	}
	return
}

func (p *parser) readEqToken(token []byte) (err error) {
	for _, t := range token {
		if len(p.buf) <= p.pos || p.buf[p.pos] != t {
			return fmt.Errorf("expected %s", token)
		}
		p.pos++
	}
	return nil
}

func (p *parser) readEqOp() (o *op, err error) {
	var token []byte
	b := p.nextNonSpace()
	for {
		if eqMap[b] != 'o' {
			break
		}
		token = append(token, b)
		if b == '-' && 1 < len(token) {
			err = fmt.Errorf("'%s' is not a valid operation", token)
			return
		}
		p.pos++
		if len(p.buf) <= p.pos {
			return nil, fmt.Errorf("equation not terminated")
		}
		b = p.buf[p.pos]
	}
	o = opMap[string(token)]
	if o == nil {
		err = fmt.Errorf("'%s' is not a valid operation", token)
	}
	return
}

func (p *parser) skipSpace() (b byte) {
	for p.pos < len(p.buf) {
		b = p.buf[p.pos]
		p.pos++
		if b != ' ' {
			break
		}
	}
	return
}

func (p *parser) nextNonSpace() (b byte) {
	for p.pos < len(p.buf) {
		b = p.buf[p.pos]
		if b != ' ' {
			break
		}
		p.pos++
	}
	return
}
