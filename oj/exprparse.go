// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

const (
	//   0123456789abcdef0123456789abcdef
	tokenMap = "" +
		"................................" + // 0x00
		"...o.oo....o.o.ooooooooooooooooo" + // 0x20
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
	// Using a xparser function adds 50% overhead so pass the xparser as an
	// arg instead.
	fun      func(*xparser, byte) error
	x        Expr
	token    []byte
	slice    []int
	num      int
	depth    int
	union    []interface{}
	eqs      []*Equation
	opName   []byte
	isFilter bool
}

// ParseExpr parses a string into an Expr.
func ParseExprString(s string) (x Expr, err error) {
	return ParseExpr([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func ParseExpr(buf []byte) (Expr, error) {
	xp := &xparser{x: Expr{}}
	err := xp.parse(buf)
	return xp.x, err
}

func (xp *xparser) parse(buf []byte) (err error) {
	xp.fun = startFun
	for i, b := range buf {
		if err = xp.fun(xp, b); err != nil {
			return fmt.Errorf("%s at %d in %q", err, i, buf)
			break
		}
	}
	return xp.fun(xp, 0)
}

func startFun(xp *xparser, b byte) (err error) {
	switch b {
	case '$':
		xp.x = append(xp.x, Root('$'))
		xp.fun = fragFun
	case '@':
		xp.x = append(xp.x, At('@'))
		xp.fun = fragFun
	case '[':
		xp.x = append(xp.x, Bracket(' '))
		xp.fun = openFun
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("an expression can not start with a '%c'", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = childFun
	}
	return nil
}

func fragFun(xp *xparser, b byte) (err error) {
	switch b {
	case 0:
	case '.':
		xp.fun = dotFun
	case '[':
		xp.fun = openFun
	default:
		err = fmt.Errorf("expected a '.' or a '['")
	}
	return
}

func childFun(xp *xparser, b byte) (err error) {
	switch b {
	case 0:
		xp.x = append(xp.x, Child(xp.token))
	case '.':
		xp.x = append(xp.x, Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = dotFun
	case '[':
		xp.x = append(xp.x, Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = openFun
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' character can not be in a non-bracketed child", b)
		}
		xp.token = append(xp.token, b)
	}
	return
}

func openFun(xp *xparser, b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '*':
		xp.x = append(xp.x, Wildcard('#'))
		xp.fun = closeFun
	case '\'':
		xp.token = xp.token[:0]
		xp.fun = quoteFun
	case '"':
		xp.token = xp.token[:0]
		xp.fun = quote2Fun
	case '-':
		xp.fun = negFun
	case '0':
		xp.num = 0
		xp.fun = zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = numFun
	case ':':
		xp.slice = xp.slice[:0]
		xp.slice = append(xp.slice, 0)
		xp.fun = colonFun
	case '?':
		xp.fun = filterFun
		xp.isFilter = true
	case '(':
		xp.depth = 1
		xp.fun = scriptFun
		xp.isFilter = false
	default:
		err = fmt.Errorf("a '%c' can not follow a '['", b)
	}
	return
}

func closeFun(xp *xparser, b byte) (err error) {
	switch b {
	case ']':
		xp.fun = fragFun
	case ' ':
		// keep going
	default:
		err = fmt.Errorf("expected a ']'")
	}
	return
}

func closeCommaFun(xp *xparser, b byte) (err error) {
	// used after a close quote only
	switch b {
	case ' ':
		// keep going
	case ']':
		if 0 < len(xp.union) {
			c, _ := xp.x[len(xp.x)-1].(Child)
			xp.union = append(xp.union, string(c))
			u := make(Union, len(xp.union))
			copy(u, xp.union)
			xp.x[len(xp.x)-1] = u
			xp.union = xp.union[:0]
		}
		xp.fun = fragFun
	case ',':
		c, _ := xp.x[len(xp.x)-1].(Child)
		xp.union = append(xp.union, string(c))
		xp.x = xp.x[:len(xp.x)-1]
		xp.fun = unionFun
	default:
		err = fmt.Errorf("expected a ']'")
	}
	return
}

func quoteFun(xp *xparser, b byte) (err error) {
	switch b {
	case '\\':
		xp.fun = escFun
	case '\'':
		xp.x = append(xp.x, Child(xp.token))
		xp.fun = closeCommaFun
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func quote2Fun(xp *xparser, b byte) (err error) {
	switch b {
	case '\\':
		xp.fun = esc2Fun
	case '"':
		xp.x = append(xp.x, Child(xp.token))
		xp.fun = closeCommaFun
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func escFun(xp *xparser, b byte) (err error) {
	if b != '\'' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = quoteFun
	return
}

func esc2Fun(xp *xparser, b byte) (err error) {
	if b != '"' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = quote2Fun
	return
}

func colonFun(xp *xparser, b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '-':
		xp.fun = negFun
	case '0':
		xp.num = 0
		xp.fun = zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = numFun
	case ']':
		if 0 < len(xp.slice) {
			xp.slice = append(xp.slice, -1)
			ia := make([]int, len(xp.slice))
			copy(ia, xp.slice)
			xp.slice = xp.slice[:0]
			xp.x = append(xp.x, Slice(ia))
		} else {
			xp.x = append(xp.x, Nth(xp.num))
		}
		xp.fun = fragFun
	default:
		err = fmt.Errorf("invalid slice format")
	}
	return
}

func negFun(xp *xparser, b byte) (err error) {
	switch b {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = -int(b - '0')
		xp.fun = numFun
	default:
		err = fmt.Errorf("parse expression failed")
	}
	return
}

func zeroFun(xp *xparser, b byte) (err error) {
	switch b {
	case ' ':
		xp.fun = numDoneFun
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, 0)
		xp.fun = unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, 0)
		xp.fun = colonFun
	default:
		err = fmt.Errorf("unexpected character")
	}
	return
}

func numFun(xp *xparser, b byte) (err error) {
	fmt.Printf("*** numFun %c depth: %d - %d\n", b, xp.depth, len(xp.eqs))
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if 0 <= xp.num {
			xp.num = xp.num*10 + int(b-'0')
		} else {
			xp.num = xp.num*10 - int(b-'0')
		}
	case ' ':
		if 0 < len(xp.eqs) {
			e := xp.eqs[len(xp.eqs)-1]
			if e.o == nil {
				e.result = xp.num
				xp.fun = opFun
			} else {
				e.right = &Equation{result: xp.num}
				xp.fun = eqCloseFun
			}
		} else {
			xp.fun = numDoneFun
		}
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, xp.num)
		xp.fun = unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, xp.num)
		xp.fun = colonFun
	case ')':
		xp.depth--
		if 0 < len(xp.eqs) {
			fmt.Printf("*** len eqs: %d\n", len(xp.eqs))
			e := xp.eqs[len(xp.eqs)-1]
			if e.o == nil {
				e.result = xp.num
			} else {
				e.right = &Equation{result: xp.num}
			}
			if xp.depth <= 0 {
				if xp.isFilter {
					xp.x = append(xp.x, e.Filter())
				} else {
					xp.x = append(xp.x, &ScriptFrag{Script: e.Script()})
				}
				xp.fun = closeScriptFun
			} else {
				xp.fun = opFun
			}
		} else {
			err = fmt.Errorf("invalid syntax")
		}
	default:
		if 0 < len(xp.eqs) && eqMap[b] == 'o' {
			fmt.Printf("*** start op with %c\n", b)
			e := xp.eqs[len(xp.eqs)-1]
			if e.o == nil {
				e.result = xp.num
			} else {
				e.right = &Equation{result: xp.num}
			}
			xp.token = xp.token[:0]
			xp.token = append(xp.token, b)
			xp.fun = opFun
		} else {
			err = fmt.Errorf("invalid number")
		}
	}
	return
}

func numDoneFun(xp *xparser, b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, xp.num)
		xp.fun = unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, xp.num)
		xp.fun = colonFun
	default:
		err = fmt.Errorf("invalid number")
	}
	return
}

func (xp *xparser) closeNumBracket() {
	if 0 < len(xp.slice) {
		xp.slice = append(xp.slice, xp.num)
		ia := make([]int, len(xp.slice))
		copy(ia, xp.slice)
		xp.slice = xp.slice[:0]
		xp.x = append(xp.x, Slice(ia))
	} else if 0 < len(xp.union) {
		xp.union = append(xp.union, xp.num)
		u := make(Union, len(xp.union))
		copy(u, xp.union)
		xp.x = append(xp.x, u)
		xp.union = xp.union[:0]
	} else {
		xp.x = append(xp.x, Nth(xp.num))
	}
	xp.fun = fragFun
}

func unionFun(xp *xparser, b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '\'':
		xp.token = xp.token[:0]
		xp.fun = quoteFun
	case '"':
		xp.token = xp.token[:0]
		xp.fun = quote2Fun
	case '-':
		xp.fun = negFun
	case '0':
		xp.num = 0
		xp.fun = zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = numFun
	case ']':
		u := make(Union, len(xp.union))
		copy(u, xp.union)
		xp.x = append(xp.x, u)
		xp.union = xp.union[:0]
	}
	return
}

func dotFun(xp *xparser, b byte) (err error) {
	switch b {
	case '.':
		xp.x = append(xp.x, Descent('.'))
		xp.fun = dot2Fun
	case '*':
		xp.x = append(xp.x, Wildcard('*'))
		xp.fun = fragFun
	case '[':
		err = fmt.Errorf("unexpected '[' after a '.'")
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' character can not be in a non-bracketed child", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = childFun
	}
	return
}

func dot2Fun(xp *xparser, b byte) (err error) {
	switch b {
	case '*':
		xp.x = append(xp.x, Wildcard('*'))
		xp.fun = fragFun
	case '[':
		err = fmt.Errorf("a '[' can not follow '..'")
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' can not follow a '..'", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = childFun
	}
	return
}

func filterFun(xp *xparser, b byte) (err error) {
	if b != '(' {
		err = fmt.Errorf("a filter must begin with '?('")
	}
	xp.fun = scriptFun
	xp.depth = 1
	return
}

func scriptFun(xp *xparser, b byte) (err error) {
	// starts after the ( which was already read
	fmt.Printf("*** scriptFun %c\n", b)

	switch b {
	case ' ':
		// Skip spaces waiting for value start then create equation.
	case '-':
		xp.eqs = append(xp.eqs, &Equation{})
		xp.fun = negFun
		// TBD
	default:
		if eqMap[b] == 'v' {
			xp.eqs = append(xp.eqs, &Equation{})
			xp.startValue(b)
		} else {
			err = fmt.Errorf("invalid equation value")
		}
	}
	return
}

func opFun(xp *xparser, b byte) (err error) {
	fmt.Printf("*** opFun %c len: %d\n", b, len(xp.eqs))
	switch b {
	case ' ':
		// keep going
	case ')':
		xp.depth--
		e := xp.eqs[len(xp.eqs)-1]
		xp.eqs = xp.eqs[:len(xp.eqs)-1]
		if xp.depth <= 0 {
			if xp.isFilter {
				xp.x = append(xp.x, e.Filter())
			} else {
				xp.x = append(xp.x, &ScriptFrag{Script: e.Script()})
			}
			xp.fun = closeScriptFun
		} else {
			if e.o == nil {
				// no change in fun
			} else {
				// TBD add to parent right
			}
		}
		// TBD close equation, set in parent, pop bsaed on prec
	default:
		switch eqMap[b] {
		case 'o':
			xp.token = append(xp.token, b)
		case 'v', '-':
			fmt.Printf("*** opFun value %c len: %d\n", b, len(xp.eqs))
			e := xp.eqs[len(xp.eqs)-1]
			if e.o = opMap[string(xp.token)]; e.o == nil {
				err = fmt.Errorf("invalid operation, %q", xp.token)
			}
			if e.left == nil {
				e.left = &Equation{result: e.result}
				e.result = nil
			}
			// TBD copy result to left if left is nil
			xp.token = xp.token[:0]
			xp.startValue(b)
		default:
			err = fmt.Errorf("invalid operation or value")
		}
	}
	return
}

func eqCloseFun(xp *xparser, b byte) (err error) {
	fmt.Printf("*** eqCloseFun %c\n", b)
	switch b {
	case ' ':
		// keep going
	case ')':
		// TBD close equation, set in parent, pop
	default:
		// if op byte then compare precidence
		//  if xx then create equation before and set current as left
		//  else current right becomes left (or result) or new equation
		// TBD
		err = fmt.Errorf("????")
	}
	// TBD
	return
}

func closeScriptFun(xp *xparser, b byte) (err error) {
	fmt.Printf("*** closeScriptFun %c\n", b)
	switch b {
	case ' ':
		// keep going
	case ']':
		xp.fun = fragFun
	default:
		err = fmt.Errorf("espected at ']'")
	}
	return
}

func (xp *xparser) startValue(b byte) {
	fmt.Printf("*** start value %c - %d\n", b, len(xp.eqs))
	switch b {
	case '-':
		xp.fun = negFun
	case '0':
		xp.num = 0
		xp.fun = zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = numFun
	case '\'':
	case '"':
	case 'n': // null
	case 't': // true
	case 'f': // false
	case '@':
	case '$':
	case '(':
	}
	// TBD
	return
}
