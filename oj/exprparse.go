// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
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
	// - could be either
	//   0123456789abcdef0123456789abcdef
	eqMap = "" +
		"................................" + // 0x00
		".ov.v.ovv.oo.-.ovvvvvvvvvv..ooo." + // 0x20
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
type xparser struct {
	buf []byte
	pos int
}

// ParseExpr parses a string into an Expr.
func ParseExprString(s string) (x Expr, err error) {
	return ParseExpr([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func ParseExpr(buf []byte) (x Expr, err error) {
	xp := &xparser{buf: buf}
	x, err = xp.readExpr()
	if err == nil && xp.pos < len(buf) {
		err = fmt.Errorf("parse error")
	}
	if err != nil {
		err = fmt.Errorf("%s at %d in %s", err, xp.pos+1, buf)
	}
	return
}

func (xp *xparser) readExpr() (x Expr, err error) {
	x = Expr{}
	var f Frag
	first := true
	lastDescent := false
	for {
		if f, err = xp.nextFrag(first, lastDescent); err != nil || f == nil {
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

func (xp *xparser) nextFrag(first, lastDescent bool) (f Frag, err error) {
	if xp.pos < len(xp.buf) {
		b := xp.buf[xp.pos]
		xp.pos++
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
			f, err = xp.afterDot()
		case '[':
			f, err = xp.afterBracket()
		case ']':
			// done
		default:
			xp.pos--
			if tokenMap[b] == 'o' {
				if first {
					f, err = xp.afterDot()
				} else if lastDescent {
					f, err = xp.afterDotDot()
				}
			}
		}
		// Any other character is the end of the Expr, figure out later if
		// that is an error.
	}
	return
}

func (xp *xparser) afterDot() (Frag, error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	var token []byte
	b := xp.buf[xp.pos]
	xp.pos++
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
	for xp.pos < len(xp.buf) {
		b := xp.buf[xp.pos]
		xp.pos++
		if tokenMap[b] == '.' {
			xp.pos--
			break
		}
		token = append(token, b)
	}
	return Child(token), nil
}

func (xp *xparser) afterDotDot() (Frag, error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	var token []byte
	b := xp.buf[xp.pos]
	xp.pos++
	if tokenMap[b] == '.' {
		return nil, fmt.Errorf("an expression fragment can not start with a '%c'", b)
	}
	token = append(token, b)
	for xp.pos < len(xp.buf) {
		b := xp.buf[xp.pos]
		xp.pos++
		if tokenMap[b] == '.' {
			xp.pos--
			break
		}
		token = append(token, b)
	}
	return Child(token), nil
}

func (xp *xparser) afterBracket() (Frag, error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	b := xp.skipSpace()
	switch b {
	case '*':
		// expect ]
		b := xp.skipSpace()
		if b != ']' {
			return nil, fmt.Errorf("not terminated")
		}
		return Wildcard('#'), nil
	case '\'', '"':
		var err error
		var s string
		if s, err = xp.readStr(b); err != nil {
			return nil, err
		}
		b = xp.skipSpace()
		switch b {
		case ']':
			return Child(s), nil
		case ',':
			return xp.readUnion(s, b)
		default:
			return nil, fmt.Errorf("invalid bracket fragment")
		}
	case ':':
		return xp.readSlice(0)
	case '?':
		return xp.readFilter()
	case '(':
		return nil, fmt.Errorf("scripts not implemented yet")
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var err error
		var i int
		if i, b, err = xp.readInt(b); err != nil {
			return nil, err
		}
	Next:
		switch b {
		case ' ':
			b = xp.skipSpace()
			goto Next
		case ']':
			return Nth(i), nil
		case ',':
			return xp.readUnion(i, b)
		case ':':
			return xp.readSlice(i)
		default:
			return nil, fmt.Errorf("invalid bracket fragment")
		}
	default:

	}
	return nil, nil
}

func (xp *xparser) readInt(b byte) (int, byte, error) {
	if b == '0' {
		if xp.pos < len(xp.buf) {
			b = xp.buf[xp.pos]
			xp.pos++
		}
		return 0, b, nil
	}
	neg := b == '-'
	if neg {
		if len(xp.buf) <= xp.pos {
			return 0, 0, fmt.Errorf("expected a number")
		}
		b = xp.buf[xp.pos]
		xp.pos++
	}
	var i int
	for {
		if b < '0' || '9' < b {
			break
		}
		i = i*10 + int(b-'0')
		if len(xp.buf) <= xp.pos {
			break
		}
		b = xp.buf[xp.pos]
		xp.pos++
	}
	if neg {
		i = -i
	}
	return i, b, nil
}

func (xp *xparser) readNum(b byte) (interface{}, error) {
	var num []byte

	num = append(num, b)
	// Read digits first
	for xp.pos < len(xp.buf) {
		b = xp.buf[xp.pos]
		if b < '0' || '9' < b {
			break
		}
		num = append(num, b)
		xp.pos++
	}
	switch b {
	case '.':
		num = append(num, b)
		xp.pos++
		for xp.pos < len(xp.buf) {
			b = xp.buf[xp.pos]
			if b < '0' || '9' < b {
				break
			}
			num = append(num, b)
			xp.pos++
		}
		if b == 'e' || b == 'E' {
			xp.pos++
			num = append(num, b)
			if len(xp.buf) <= xp.pos {
				return 0, fmt.Errorf("expected a number")
			}
			b = xp.buf[xp.pos]
		} else {
			f, err := strconv.ParseFloat(string(num), 64)
			return f, err
		}
	case 'e', 'E':
		xp.pos++
		if len(xp.buf) <= xp.pos {
			return 0, fmt.Errorf("expected a number")
		}
		num = append(num, b)
		b = xp.buf[xp.pos]
	default:
		i, err := strconv.ParseInt(string(num), 10, 64)
		return int(i), err
	}
	if b == '+' || b == '-' {
		num = append(num, b)
		xp.pos++
		if len(xp.buf) <= xp.pos {
			return 0, fmt.Errorf("expected a number")
		}
	}
	for xp.pos < len(xp.buf) {
		b = xp.buf[xp.pos]
		if b < '0' || '9' < b {
			break
		}
		num = append(num, b)
		xp.pos++
	}
	f, err := strconv.ParseFloat(string(num), 64)

	return f, err
}

func (xp *xparser) readSlice(i int) (Frag, error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	f := Slice{i}
	b := xp.buf[xp.pos]
	if b == ']' {
		f = append(f, -1)
		return f, nil
	}
	b = xp.skipSpace()
	var err error
	// read the end
	if b == ':' {
		f = append(f, -1)
		if len(xp.buf) <= xp.pos {
			return nil, fmt.Errorf("not terminated")
		}
		b = xp.buf[xp.pos]
		xp.pos++
		if i, b, err = xp.readInt(b); err != nil {
			return nil, err
		}
		f = append(f, i)
	} else if i, b, err = xp.readInt(b); err == nil {
		f = append(f, i)
		if b == ':' {
			if len(xp.buf) <= xp.pos {
				return nil, fmt.Errorf("not terminated")
			}
			b = xp.buf[xp.pos]
			xp.pos++
			if i, b, err = xp.readInt(b); err != nil {
				return nil, err
			}
			f = append(f, i)
		}
	}
	if b != ']' {
		return nil, fmt.Errorf("invalid slice syntax")
	}
	return f, nil
}

func (xp *xparser) readUnion(v interface{}, b byte) (Frag, error) {
	if len(xp.buf) <= xp.pos {
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
		b = xp.skipSpace()
		switch b {
		case '\'', '"':
			var s string
			if s, err = xp.readStr(b); err != nil {
				return nil, err
			}
			b = xp.skipSpace()
			f = append(f, s)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var i int
			if i, b, err = xp.readInt(b); err != nil {
				return nil, err
			}
			f = append(f, i)
			if b == ' ' {
				b = xp.skipSpace()
			}
		default:
			return nil, fmt.Errorf("invalid union syntax")
		}
	}
}

func (xp *xparser) readStr(term byte) (string, error) {
	start := xp.pos
	esc := false
	for xp.pos < len(xp.buf) {
		b := xp.buf[xp.pos]
		xp.pos++
		if b == term && !esc {
			break
		}
		if b == '\\' {
			esc = !esc
		} else {
			esc = false
		}
	}
	return string(xp.buf[start : xp.pos-1]), nil
}

func (xp *xparser) readFilter() (*Filter, error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	b := xp.buf[xp.pos]
	xp.pos++
	if b != '(' {
		return nil, fmt.Errorf("expected a '(' in filter")
	}
	eq, err := xp.readEquation()
	if len(xp.buf) <= xp.pos || xp.buf[xp.pos] != ']' {
		return nil, fmt.Errorf("not terminated")
	}
	xp.pos++
	if err == nil {
		return eq.Filter(), nil
	}
	return nil, err
}

func (xp *xparser) readEquation() (eq *Equation, err error) {
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}
	eq = &Equation{}

	b := xp.nextNonSpace()
	if b == '!' {
		eq.o = not
		xp.pos++
		eq.left, err = xp.readEqValue()
		b := xp.nextNonSpace()
		if b != ')' {
			return nil, fmt.Errorf("not terminated")
		}
		xp.pos++
		return
	}
	if eq.left, err = xp.readEqValue(); err != nil {
		return
	}
	if eq.o, err = xp.readEqOp(); err != nil {
		return
	}
	if eq.right, err = xp.readEqValue(); err != nil {
		return
	}
	for {
		b = xp.nextNonSpace()
		if b == ')' {
			xp.pos++
			return
		}
		var o *op
		if o, err = xp.readEqOp(); err != nil {
			return
		}
		if eq.o.prec <= o.prec {
			eq = &Equation{left: eq, o: o}
			if eq.right, err = xp.readEqValue(); err != nil {
				return
			}
		} else {
			eq.right = &Equation{left: eq.right, o: o}
			if eq.right.right, err = xp.readEqValue(); err != nil {
				return
			}
		}
	}
}

func (xp *xparser) readEqValue() (eq *Equation, err error) {
	b := xp.nextNonSpace()
	switch b {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var v interface{}
		xp.pos++
		if v, err = xp.readNum(b); err != nil {
			return
		}
		eq = &Equation{result: v}
	case '\'', '"':
		xp.pos++
		var s string
		if s, err = xp.readStr(b); err != nil {
			return
		}
		eq = &Equation{result: s}
	case 'n':
		if err = xp.readEqToken([]byte("null")); err != nil {
			return
		}
		eq = &Equation{result: nil}
	case 't':
		if err = xp.readEqToken([]byte("true")); err != nil {
			return
		}
		eq = &Equation{result: true}

	case 'f':
		if err = xp.readEqToken([]byte("false")); err != nil {
			return
		}
		eq = &Equation{result: false}
	case '@', '$':
		var x Expr
		x, err = xp.readExpr()
		eq = &Equation{result: x}
	case '(':
		xp.pos++
		eq, err = xp.readEquation()
	default:
		err = fmt.Errorf("expected a value")
	}
	return
}

func (xp *xparser) readEqToken(token []byte) (err error) {
	for _, t := range token {
		if len(xp.buf) <= xp.pos || xp.buf[xp.pos] != t {
			return fmt.Errorf("expected %s", token)
		}
		xp.pos++
	}
	return nil
}

func (xp *xparser) readEqOp() (o *op, err error) {
	var token []byte
	b := xp.nextNonSpace()
	for {
		if eqMap[b] != 'o' {
			break
		}
		token = append(token, b)
		if b == '-' && 1 < len(token) {
			err = fmt.Errorf("%q is not a valid operation", token)
			return
		}
		xp.pos++
		if len(xp.buf) <= xp.pos {
			err = fmt.Errorf("equation not terminated")
		}
		b = xp.buf[xp.pos]
	}
	o = opMap[string(token)]
	if o == nil {
		err = fmt.Errorf("%q is not a valid operation", token)
	}
	return
}

func (xp *xparser) skipSpace() (b byte) {
	for xp.pos < len(xp.buf) {
		b = xp.buf[xp.pos]
		xp.pos++
		if b != ' ' {
			break
		}
	}
	return
}

func (xp *xparser) nextNonSpace() (b byte) {
	for xp.pos < len(xp.buf) {
		b = xp.buf[xp.pos]
		if b != ' ' {
			break
		}
		xp.pos++
	}
	return
}
