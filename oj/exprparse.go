// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

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

	// Using a xparser function adds 50% overhead so pass the xparser as an
	// arg instead.
	fun      func(*xparser, byte) error
	xa       []Expr
	token    []byte
	slice    []int
	num      int
	depth    int
	union    []interface{}
	eqs      []*Equation
	opName   []byte
	isFilter bool
	script   bool
}

// TBD remove after implemented and tested
func (xp *xparser) where(fun string) {
	if false {
		var b byte
		if xp.pos < len(xp.buf) {
			b = xp.buf[xp.pos]
		}
		fmt.Printf("*** %s - %d '%c'\n", fun, xp.pos, b)
	}
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
		err = fmt.Errorf("%s at %d in %s", err, xp.pos, buf)
	}
	return
}

func (xp *xparser) readExpr() (x Expr, err error) {
	xp.where("readExpr")
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
	return
}

func (xp *xparser) nextFrag(first, lastDescent bool) (f Frag, err error) {
	xp.where("nextFrag")
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
		default:
			if first || lastDescent {
				xp.pos--
				f, err = xp.afterDot()
			}
		}
		// Any other character is the end of the Expr, figure out later if
		// that is an error.
	}
	return
}

func (xp *xparser) afterDot() (Frag, error) {
	xp.where("afterDot")
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

func (xp *xparser) afterBracket() (Frag, error) {
	xp.where("afterBracket")
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

func (xp *xparser) readSlice(i int) (Frag, error) {
	xp.where("readSlice")
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
	xp.where("readUnion")
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
	return f, nil
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
	if err == nil {
		return eq.Filter(), nil
	}
	return nil, err
}

func (xp *xparser) readEquation() (*Equation, error) {
	xp.where("readEquation")
	if len(xp.buf) <= xp.pos {
		return nil, fmt.Errorf("not terminated")
	}

	// TBD read value, then op then value

	return nil, nil
}
