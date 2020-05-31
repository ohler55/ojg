// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

// TBD remove after implemented and tested
const debug = false

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
}

// ParseExpr parses a string into an Expr.
func ParseExprString(s string) (x Expr, err error) {
	return ParseExpr([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func ParseExpr(buf []byte) (Expr, error) {
	xp := &xparser{}
	xp.xa = append(xp.xa, Expr{})
	err := xp.parse(buf)
	return xp.xa[0], err
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
	if debug {
		fmt.Printf("*** startFun %c\n", b)
	}
	switch b {
	case '$':
		xp.exprAppend(Root('$'))
		xp.fun = fragFun
	case '@':
		xp.exprAppend(At('@'))
		xp.fun = fragFun
	case '[':
		xp.exprAppend(Bracket(' '))
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
	if debug {
		fmt.Printf("*** fragFun %c\n", b)
	}
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
	if debug {
		fmt.Printf("*** childFun %c\n", b)
	}
	switch b {
	case 0:
		xp.exprAppend(Child(xp.token))
	case '.':
		xp.exprAppend(Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = dotFun
	case '[':
		xp.exprAppend(Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = openFun
	default:
		if tokenMap[b] == '.' {
			if 1 < len(xp.xa) { // processing an Expr in a Filter
				if eqMap[b] == 'o' || b == ' ' { // an operation char or space
					xp.exprAppend(Child(xp.token))
					xp.token = xp.token[:0]
					e := xp.eqs[len(xp.eqs)-1]
					x := xp.popExpr()
					if e.o == nil {
						e.left = &Equation{o: get, left: &Equation{result: x}}
					} else {
						e.right = &Equation{o: get, left: &Equation{result: x}}
						// TBD close if needed
						//xp.fun = eqCloseFun
					}
					if b != ' ' {
						xp.token = append(xp.token, b)
					}
					xp.fun = opFun
					return
				}
			}
			err = fmt.Errorf("a '%c' character can not be in a non-bracketed child", b)
		} else {
			xp.token = append(xp.token, b)
		}
	}
	return
}

func openFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** openFun %c\n", b)
	}
	switch b {
	case ' ':
		// keep going
	case '*':
		xp.exprAppend(Wildcard('#'))
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
	if debug {
		fmt.Printf("*** closeFun %c\n", b)
	}
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
	if debug {
		fmt.Printf("*** closeCommaFun %c\n", b)
	}
	// used after a close quote only
	switch b {
	case ' ':
		// keep going
	case ']':
		if 0 < len(xp.union) {
			x := xp.currentExpr()
			c, _ := x[len(x)-1].(Child)
			xp.union = append(xp.union, string(c))
			u := make(Union, len(xp.union))
			copy(u, xp.union)
			x[len(x)-1] = u
			xp.union = xp.union[:0]
		}
		xp.fun = fragFun
	case ',':
		x := xp.currentExpr()
		c, _ := x[len(x)-1].(Child)
		xp.union = append(xp.union, string(c))
		xp.xa[len(xp.xa)-1] = x[:len(x)-1]
		xp.fun = unionFun
	default:
		err = fmt.Errorf("expected a ']'")
	}
	return
}

func quoteFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** quoteFun %c\n", b)
	}
	switch b {
	case '\\':
		xp.fun = escFun
	case '\'':
		if 0 < len(xp.eqs) {
			_ = xp.setEqValue(string(xp.token))
			xp.fun = opFun
		} else {
			xp.exprAppend(Child(xp.token))
			xp.fun = closeCommaFun
		}
		xp.token = xp.token[:0]
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func quote2Fun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** quote2Fun %c\n", b)
	}
	switch b {
	case '\\':
		xp.fun = esc2Fun
	case '"':
		if 0 < len(xp.eqs) {
			_ = xp.setEqValue(string(xp.token))
			xp.fun = opFun
		} else {
			xp.exprAppend(Child(xp.token))
			xp.fun = closeCommaFun
		}
		xp.token = xp.token[:0]
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func escFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** escFun %c\n", b)
	}
	if b != '\'' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = quoteFun
	return
}

func esc2Fun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** esc2Fun %c\n", b)
	}
	if b != '"' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = quote2Fun
	return
}

func colonFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** colonFun %c\n", b)
	}
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
			xp.exprAppend(Slice(ia))
		} else {
			xp.exprAppend(Nth(xp.num))
		}
		xp.fun = fragFun
	default:
		err = fmt.Errorf("invalid slice format")
	}
	return
}

func negFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** negFun %c\n", b)
	}
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
	if debug {
		fmt.Printf("*** zeroFun %c\n", b)
	}
	switch b {
	case ' ':
		xp.numSpace()
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
	case ')':
		err = xp.numCloseParen()
	default:
		err = xp.numDefault(b)
	}
	return
}

func numFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** numFun %c\n", b)
	}
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if 0 <= xp.num {
			xp.num = xp.num*10 + int(b-'0')
		} else {
			xp.num = xp.num*10 - int(b-'0')
		}
	case ' ':
		xp.numSpace()
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
		err = xp.numCloseParen()
	default:
		err = xp.numDefault(b)
	}
	return
}

func numDoneFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** numDoneFun %c\n", b)
	}
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

func unionFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** unionFun %c\n", b)
	}
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
		xp.exprAppend(u)
		xp.union = xp.union[:0]
	}
	return
}

func dotFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** dotFun %c\n", b)
	}
	switch b {
	case '.':
		xp.exprAppend(Descent('.'))
		xp.fun = dot2Fun
	case '*':
		xp.exprAppend(Wildcard('*'))
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
	if debug {
		fmt.Printf("*** dot2Fun %c\n", b)
	}
	switch b {
	case '*':
		xp.exprAppend(Wildcard('*'))
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
	if debug {
		fmt.Printf("*** filterFun %c\n", b)
	}
	if b != '(' {
		err = fmt.Errorf("a filter must begin with '?('")
	}
	xp.fun = scriptFun
	xp.depth = 1
	return
}

func scriptFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** scriptFun %c\n", b)
	}
	// starts after the ( which was already read
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
	if debug {
		fmt.Printf("*** opFun %c\n", b)
	}
	switch b {
	case ' ':
		// keep going
	case ')':
		xp.depth--
		e := xp.eqs[len(xp.eqs)-1]
		xp.eqs = xp.eqs[:len(xp.eqs)-1]
		if xp.depth <= 0 {
			if xp.isFilter {
				xp.exprAppend(e.Filter())
			} else {
				xp.exprAppend(&ScriptFrag{Script: e.Script()})
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
			e := xp.eqs[len(xp.eqs)-1]
			if e.o = opMap[string(xp.token)]; e.o == nil {
				err = fmt.Errorf("invalid operation, %q", xp.token)
			}
			if e.left == nil {
				e.left = &Equation{result: e.result}
				e.result = nil
			}
			xp.token = xp.token[:0]
			xp.startValue(b)
		default:
			err = fmt.Errorf("invalid operation or value")
		}
	}
	return
}

func eqCloseFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** eqCloseFun %c\n", b)
	}
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
	if debug {
		fmt.Printf("*** closeScriptFun %c\n", b)
	}
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

func falseFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** falseFun %c\n", b)
	}
	switch b {
	case 'a', 'l', 's', 'e':
		xp.token = append(xp.token, b)
		if len(xp.token) == 5 {
			if "false" == string(xp.token) {
				_ = xp.setEqValue(false)
				xp.fun = opFun
				return
			}
			err = fmt.Errorf("expected 'false', not '%s'", xp.token)
		}
	default:
		err = fmt.Errorf("espected 'false'")
	}
	return
}

func trueFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** trueFun %c\n", b)
	}
	switch b {
	case 'r', 'u', 'e':
		xp.token = append(xp.token, b)
		if len(xp.token) == 4 {
			if "true" == string(xp.token) {
				_ = xp.setEqValue(true)
				xp.fun = opFun
				return
			}
			err = fmt.Errorf("expected 'true', not '%s'", xp.token)
		}
	default:
		err = fmt.Errorf("espected 'true'")
	}
	return
}

func nullFun(xp *xparser, b byte) (err error) {
	if debug {
		fmt.Printf("*** nullFun %c\n", b)
	}
	switch b {
	case 'u', 'l':
		xp.token = append(xp.token, b)
		if len(xp.token) == 4 {
			if "null" == string(xp.token) {
				_ = xp.setEqValue(nil)
				xp.fun = opFun
				return
			}
			err = fmt.Errorf("expected 'null', not '%s'", xp.token)
		}
	default:
		err = fmt.Errorf("espected 'null'")
	}
	return
}

func (xp *xparser) currentExpr() Expr {
	return xp.xa[len(xp.xa)-1]
}

func (xp *xparser) popExpr() (x Expr) {
	x = xp.xa[len(xp.xa)-1]
	xp.xa[len(xp.xa)-1] = nil
	xp.xa = xp.xa[:len(xp.xa)-1]
	return
}

func (xp *xparser) startValue(b byte) {
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
		xp.token = xp.token[:0]
		xp.fun = quoteFun
	case '"':
		xp.token = xp.token[:0]
		xp.fun = quote2Fun
	case 'n': // null
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = nullFun
	case 't': // true
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = trueFun
	case 'f': // false
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = falseFun
	case '@':
		xp.xa = append(xp.xa, Expr{})
		xp.exprAppend(At('@'))
		xp.fun = fragFun
	case '$':
		xp.xa = append(xp.xa, Expr{})
		xp.exprAppend(Root('$'))
		xp.fun = fragFun
	case '(':
		// TBD new equation
	}
	return
}

func (xp *xparser) closeNumBracket() {
	if 0 < len(xp.slice) {
		xp.slice = append(xp.slice, xp.num)
		ia := make([]int, len(xp.slice))
		copy(ia, xp.slice)
		xp.slice = xp.slice[:0]
		xp.exprAppend(Slice(ia))
	} else if 0 < len(xp.union) {
		xp.union = append(xp.union, xp.num)
		u := make(Union, len(xp.union))
		copy(u, xp.union)
		xp.exprAppend(u)
		xp.union = xp.union[:0]
	} else {
		xp.exprAppend(Nth(xp.num))
	}
	xp.fun = fragFun
}

func (xp *xparser) numSpace() {
	if 0 < len(xp.eqs) {
		e := xp.setEqValue(xp.num)
		if e.o == nil {
			xp.fun = opFun
		} else {
			xp.fun = eqCloseFun
		}
	} else {
		xp.fun = numDoneFun
	}
}

func (xp *xparser) numCloseParen() (err error) {
	xp.depth--
	if 0 < len(xp.eqs) {
		e := xp.setEqValue(xp.num)
		if xp.depth <= 0 {
			if xp.isFilter {
				xp.exprAppend(e.Filter())
			} else {
				xp.exprAppend(&ScriptFrag{Script: e.Script()})
			}
			xp.fun = closeScriptFun
		} else {
			xp.fun = opFun
		}
	} else {
		err = fmt.Errorf("invalid syntax")
	}
	return
}

func (xp *xparser) numDefault(b byte) (err error) {
	if 0 < len(xp.eqs) && eqMap[b] == 'o' {
		_ = xp.setEqValue(xp.num)
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = opFun
	} else {
		err = fmt.Errorf("invalid number")
	}
	return
}

func (xp *xparser) setEqValue(v interface{}) (e *Equation) {
	e = xp.eqs[len(xp.eqs)-1]
	if e.o == nil {
		e.result = v
	} else {
		e.right = &Equation{result: v}
	}
	return
}

func (xp *xparser) exprAppend(f Frag) {
	xp.xa[len(xp.xa)-1] = append(xp.xa[len(xp.xa)-1], f)
}
