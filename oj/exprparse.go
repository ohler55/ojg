// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

const (
	startMode        = 's' // new expression
	fragMode         = 'f' // new fragment should be next
	dot2Mode         = 'o' // just read 2 dots
	openMode         = '[' // last read a [
	closeMode        = ']' // expect a ]
	childMode        = 'c' // reading a child fragment
	numMode          = '#'
	numDoneMode      = 'D'
	quoteMode        = 'q'
	quote2Mode       = 'Q'
	esc2Mode         = 'E'
	filterMode       = '?' // scan filter until matching )
	filterQuoteMode  = '\''
	filterQuote2Mode = '"'
	filterEscMode    = '1'
	filterEsc2Mode   = '2'
	unionMode        = 'u'
	closeCommaMode   = ',' // expect a ] or comma

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
)

// Performance is less a concern with Expr parsing as it is usually done just
// once if performance is important. Alternatively, an Expr can be built using
// function calls or bare structs. Parsing is more for convenience.
type xparser struct {
	// TBD try with xparser as an arg.
	fun   func(b byte) error
	x     Expr
	token []byte
	slice []int
	num   int
	depth int
	union []interface{}
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
	xp.fun = xp.startFun
	for i, b := range buf {
		if err = xp.fun(b); err != nil {
			return fmt.Errorf("%s at %d in %q", err, i, buf)
			break
		}
	}
	return xp.fun(0)
}

func (xp *xparser) startFun(b byte) (err error) {
	switch b {
	case '$':
		xp.x = append(xp.x, Root('$'))
		xp.fun = xp.fragFun
	case '@':
		xp.x = append(xp.x, At('@'))
		xp.fun = xp.fragFun
	case '[':
		xp.x = append(xp.x, Bracket(' '))
		xp.fun = xp.openFun
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("an expression can not start with a '%c'", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = xp.childFun
	}
	return nil
}

func (xp *xparser) fragFun(b byte) (err error) {
	switch b {
	case 0:
	case '.':
		xp.fun = xp.dotFun
	case '[':
		xp.fun = xp.openFun
	default:
		err = fmt.Errorf("expected a '.' or a '['")
	}
	return
}

func (xp *xparser) childFun(b byte) (err error) {
	switch b {
	case 0:
		xp.x = append(xp.x, Child(xp.token))
	case '.':
		xp.x = append(xp.x, Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = xp.dotFun
	case '[':
		xp.x = append(xp.x, Child(xp.token))
		xp.token = xp.token[:0]
		xp.fun = xp.openFun
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' character can not be in a non-bracketed child", b)
		}
		xp.token = append(xp.token, b)
	}
	return
}

func (xp *xparser) openFun(b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '*':
		xp.x = append(xp.x, Wildcard('#'))
		xp.fun = xp.closeFun
	case '\'':
		xp.token = xp.token[:0]
		xp.fun = xp.quoteFun
	case '"':
		xp.token = xp.token[:0]
		xp.fun = xp.quote2Fun
	case '-':
		xp.fun = xp.negFun
	case '0':
		xp.num = 0
		xp.fun = xp.zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = xp.numFun
	case ':':
		xp.slice = xp.slice[:0]
		xp.slice = append(xp.slice, 0)
		xp.fun = xp.colonFun
	case '?':
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.depth = 0
		xp.fun = xp.filterFun
	case '(':
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.depth = 1
		xp.fun = xp.filterFun
	default:
		err = fmt.Errorf("a '%c' can not follow a '['", b)
	}
	return
}

func (xp *xparser) closeFun(b byte) (err error) {
	switch b {
	case ']':
		xp.fun = xp.fragFun
	case ' ':
		// keep going
	default:
		err = fmt.Errorf("expected a ']'")
	}
	return
}

func (xp *xparser) closeCommaFun(b byte) (err error) {
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
		xp.fun = xp.fragFun
	case ',':
		c, _ := xp.x[len(xp.x)-1].(Child)
		xp.union = append(xp.union, string(c))
		xp.x = xp.x[:len(xp.x)-1]
		xp.fun = xp.unionFun
	default:
		err = fmt.Errorf("expected a ']'")
	}
	return
}

func (xp *xparser) quoteFun(b byte) (err error) {
	switch b {
	case '\\':
		xp.fun = xp.escFun
	case '\'':
		xp.x = append(xp.x, Child(xp.token))
		xp.fun = xp.closeCommaFun
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func (xp *xparser) quote2Fun(b byte) (err error) {
	switch b {
	case '\\':
		xp.fun = xp.esc2Fun
	case '"':
		xp.x = append(xp.x, Child(xp.token))
		xp.fun = xp.closeCommaFun
	default:
		xp.token = append(xp.token, b)
	}
	return
}

func (xp *xparser) escFun(b byte) (err error) {
	if b != '\'' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = xp.quoteFun
	return
}

func (xp *xparser) esc2Fun(b byte) (err error) {
	if b != '"' {
		xp.token = append(xp.token, '\\')
	}
	xp.token = append(xp.token, b)
	xp.fun = xp.quote2Fun
	return
}

func (xp *xparser) filterFun(b byte) (err error) {
	// TBD
	return
}

func (xp *xparser) colonFun(b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '-':
		xp.fun = xp.negFun
	case '0':
		xp.num = 0
		xp.fun = xp.zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = xp.numFun
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
		xp.fun = xp.fragFun
	default:
		err = fmt.Errorf("invalid slice format")
	}
	return
}

func (xp *xparser) negFun(b byte) (err error) {
	switch b {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = -int(b - '0')
		xp.fun = xp.numFun
	default:
		err = fmt.Errorf("parse expression failed")
	}
	return
}

func (xp *xparser) zeroFun(b byte) (err error) {
	switch b {
	case ' ':
		xp.fun = xp.numDoneFun
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, 0)
		xp.fun = xp.unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, 0)
		xp.fun = xp.colonFun
	default:
		err = fmt.Errorf("unexpected character")
	}
	return
}

func (xp *xparser) numFun(b byte) (err error) {
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if 0 <= xp.num {
			xp.num = xp.num*10 + int(b-'0')
		} else {
			xp.num = xp.num*10 - int(b-'0')
		}
	case ' ':
		xp.fun = xp.numDoneFun
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, xp.num)
		xp.fun = xp.unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, xp.num)
		xp.fun = xp.colonFun
	default:
		err = fmt.Errorf("invalid number")
	}
	return
}

func (xp *xparser) numDoneFun(b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case ']':
		xp.closeNumBracket()
	case ',':
		xp.union = append(xp.union, xp.num)
		xp.fun = xp.unionFun
	case ':':
		if 2 < len(xp.slice) {
			err = fmt.Errorf("too many numbers in the slice")
		}
		xp.slice = append(xp.slice, xp.num)
		xp.fun = xp.colonFun
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
	xp.fun = xp.fragFun
}

func (xp *xparser) unionFun(b byte) (err error) {
	switch b {
	case ' ':
		// keep going
	case '\'':
		xp.token = xp.token[:0]
		xp.fun = xp.quoteFun
	case '"':
		xp.token = xp.token[:0]
		xp.fun = xp.quote2Fun
	case '-':
		xp.fun = xp.negFun
	case '0':
		xp.num = 0
		xp.fun = xp.zeroFun
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		xp.num = int(b - '0')
		xp.fun = xp.numFun
	case ']':
		u := make(Union, len(xp.union))
		copy(u, xp.union)
		xp.x = append(xp.x, u)
		xp.union = xp.union[:0]
	}
	return
}

func (xp *xparser) dotFun(b byte) (err error) {
	switch b {
	case '.':
		xp.x = append(xp.x, Descent('.'))
		xp.fun = xp.dot2Fun
	case '*':
		xp.x = append(xp.x, Wildcard('*'))
		xp.fun = xp.fragFun
	case '[':
		err = fmt.Errorf("unexpected '[' after a '.'")
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' character can not be in a non-bracketed child", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = xp.childFun
	}
	return
}

func (xp *xparser) dot2Fun(b byte) (err error) {
	switch b {
	case '*':
		xp.x = append(xp.x, Wildcard('*'))
		xp.fun = xp.fragFun
	case '[':
		err = fmt.Errorf("a '[' can not follow '..'")
	default:
		if tokenMap[b] == '.' {
			err = fmt.Errorf("a '%c' can not follow a '..'", b)
		}
		xp.token = xp.token[:0]
		xp.token = append(xp.token, b)
		xp.fun = xp.childFun
	}
	return
}

/*
func (xp *xparser) parsex(buf []byte) (err error) {
	mode := startMode
	for i, b := range buf {
		switch mode {
		case startMode:
			switch b {
			case '$':
				xp.x = append(xp.x, Root('$'))
				mode = fragMode
			case '@':
				xp.x = append(xp.x, At('@'))
				mode = fragMode
			case '[':
				xp.x = append(xp.x, Bracket(' '))
				mode = openMode
			default:
				if tokenMap[b] == '.' {
					return fmt.Errorf("an expression can not start with a '%c'at %d in %q", b, i, buf)
				}
				xp.token = xp.token[:0]
				xp.token = append(xp.token, b)
				mode = childMode
			}
		case childMode:
			switch b {
			case '.':
				xp.x = append(xp.x, Child(xp.token))
				mode = dotMode
			case '[':
				xp.x = append(xp.x, Child(xp.token))
				mode = openMode
			default:
				if tokenMap[b] == '.' {
					return fmt.Errorf("a '%c' character can not be in a non-bracketed child at %d in %q", b, i, buf)
				}
				xp.token = append(xp.token, b)
			}
		case dotMode:
			switch b {
			case '.':
				xp.x = append(xp.x, Descent('.'))
				mode = dot2Mode
			case '*':
				xp.x = append(xp.x, Wildcard('*'))
				mode = fragMode
			case '[':
				return fmt.Errorf("unexpected '[' after a '.' at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return fmt.Errorf("a '%c' character can not be in a non-bracketed child at %d in %q", b, i, buf)
				}
				xp.token = xp.token[:0]
				xp.token = append(xp.token, b)
				mode = childMode
			}
		case dot2Mode:
			switch b {
			case '*':
				xp.x = append(xp.x, Wildcard('*'))
				mode = fragMode
			case '[':
				return fmt.Errorf("a '[' can not follow '..' at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return fmt.Errorf("a '%c' can not follow a '..' at %d in %q", b, i, buf)
				}
				xp.token = xp.token[:0]
				xp.token = append(xp.token, b)
				mode = childMode
			}
		case fragMode:
			switch b {
			case '.':
				mode = dotMode
			case '[':
				mode = openMode
			default:
				return fmt.Errorf("expected a '.' or a '[' at %d in %q", i, buf)
			}
		case openMode:
			switch b {
			case ' ':
				// keep going
			case '*':
				xp.x = append(xp.x, Wildcard('#'))
				mode = closeMode
			case '\'':
				xp.token = xp.token[:0]
				mode = quoteMode
			case '"':
				xp.token = xp.token[:0]
				mode = quote2Mode
			case '-':
				mode = negMode
			case '0':
				xp.num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				xp.num = int(b - '0')
				mode = numMode
			case ':':
				xp.slice = xp.slice[:0]
				xp.slice = append(xp.slice, 0)
				mode = colonMode
			case '?':
				xp.token = xp.token[:0]
				xp.token = append(xp.token, b)
				xp.depth = 0
				mode = filterMode
			case '(':
				xp.token = xp.token[:0]
				xp.token = append(xp.token, b)
				xp.depth = 1
				mode = filterMode
			default:
				return fmt.Errorf("a '%c' can not follow a '[' at %d in %q", b, i, buf)
			}
		case closeMode:
			switch b {
			case ']':
				mode = fragMode
			case ' ':
				// keep going
			default:
				return fmt.Errorf("expected a ']' at %d in %q", i, buf)
			}
		case zeroMode:
			switch b {
			case ' ':
				mode = numDoneMode
			case ']':
				if 0 < len(xp.slice) {
					xp.slice = append(xp.slice, 0)
					ia := make([]int, len(xp.slice))
					copy(ia, xp.slice)
					xp.slice = xp.slice[:0]
					xp.x = append(xp.x, Slice(ia))
				} else if 0 < len(xp.union) {
					xp.union = append(xp.union, 0)
					u := make(Union, len(xp.union))
					copy(u, xp.union)
					xp.x = append(xp.x, u)
					xp.union = xp.union[:0]
				} else {
					xp.x = append(xp.x, Nth(0))
				}
				mode = fragMode
			case ',':
				xp.union = append(xp.union, 0)
				mode = unionMode
			case ':':
				if 2 < len(xp.slice) {
					return fmt.Errorf("too many numbers in the slice at %d in %q", i, buf)
				}
				xp.slice = append(xp.slice, 0)
				mode = colonMode
			default:
				return fmt.Errorf("unexpected character at %d in %q", i, buf)
			}
		case colonMode:
			switch b {
			case ' ':
				// keep going
			case '-':
				mode = negMode
			case '0':
				xp.num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				xp.num = int(b - '0')
				mode = numMode
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
				mode = fragMode
			default:
				return fmt.Errorf("invalid slice format at %d in %q", i, buf)
			}
		case negMode:
			switch b {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				xp.num = -int(b - '0')
				mode = numMode
			default:
				return fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
		case numMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if 0 <= xp.num {
					xp.num = xp.num*10 + int(b-'0')
				} else {
					xp.num = xp.num*10 - int(b-'0')
				}
			case ' ':
				mode = numDoneMode
			case ']':
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
				mode = fragMode
			case ',':
				xp.union = append(xp.union, xp.num)
				mode = unionMode
			case ':':
				if 2 < len(xp.slice) {
					return fmt.Errorf("too many numbers in the slice at %d in %q", i, buf)
				}
				xp.slice = append(xp.slice, xp.num)
				mode = colonMode
			default:
				return fmt.Errorf("invalid number at %d in %q", i, buf)
			}
		case numDoneMode:
			switch b {
			case ' ':
				// keep going
			case ']':
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
				mode = fragMode
			case ',':
				xp.union = append(xp.union, xp.num)
				mode = unionMode
			case ':':
				if 2 < len(xp.slice) {
					return fmt.Errorf("too many numbers in the slice at %d in %q", i, buf)
				}
				xp.slice = append(xp.slice, xp.num)
				mode = colonMode
			default:
				return fmt.Errorf("invalid number at %d in %q", i, buf)
			}
		case quoteMode:
			switch b {
			case '\\':
				mode = escMode
			case '\'':
				xp.x = append(xp.x, Child(xp.token))
				mode = closeCommaMode
			default:
				xp.token = append(xp.token, b)
			}
		case escMode:
			if b != '\'' {
				xp.token = append(xp.token, '\\')
			}
			xp.token = append(xp.token, b)
			mode = quoteMode
		case quote2Mode:
			switch b {
			case '\\':
				mode = esc2Mode
			case '"':
				xp.x = append(xp.x, Child(xp.token))
				mode = closeCommaMode
			default:
				xp.token = append(xp.token, b)
			}
		case esc2Mode:
			if b != '"' {
				xp.token = append(xp.token, '\\')
			}
			xp.token = append(xp.token, b)
			mode = quoteMode
		case unionMode:
			switch b {
			case ' ':
				// keep going
			case '\'':
				xp.token = xp.token[:0]
				mode = quoteMode
			case '"':
				xp.token = xp.token[:0]
				mode = quote2Mode
			case '-':
				mode = negMode
			case '0':
				xp.num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				xp.num = int(b - '0')
				mode = numMode
			case ']':
				u := make(Union, len(xp.union))
				copy(u, xp.union)
				xp.x = append(xp.x, u)
				xp.union = xp.union[:0]
			}
		case closeCommaMode:
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
				mode = fragMode
			case ',':
				c, _ := xp.x[len(xp.x)-1].(Child)
				xp.union = append(xp.union, string(c))
				xp.x = xp.x[:len(xp.x)-1]
				mode = unionMode
			default:
				return fmt.Errorf("expected a ']' at %d in %q", i, buf)
			}
		case filterMode:
			switch b {
			case ' ':
				// keep going
			case '(':
				xp.depth++
				xp.token = append(xp.token, b)
			case ')':
				xp.depth--
				xp.token = append(xp.token, b)
				if xp.depth <= 0 {
					if xp.token[0] == '?' {
						f := &Filter{}
						if err := f.Parse(xp.token[1:]); err != nil {
							return fmt.Errorf("%s at %d in %q", err, i, buf)
						}
						xp.x = append(xp.x, f)
					} else {
						sf := &ScriptFrag{}
						if err := sf.Parse(xp.token); err != nil {
							return fmt.Errorf("%s at %d in %q", err, i, buf)
						}
						xp.x = append(xp.x, sf)
					}
					mode = closeMode
				}
			case '\'':
				mode = filterQuoteMode
				xp.token = append(xp.token, b)
			case '"':
				mode = filterQuote2Mode
				xp.token = append(xp.token, b)
			default:
				xp.token = append(xp.token, b)
			}
		case filterQuoteMode:
			switch b {
			case '\\':
				mode = filterEscMode
			case '\'':
				mode = filterMode
			}
			xp.token = append(xp.token, b)
		case filterEscMode:
			xp.token = append(xp.token, b)
			mode = quoteMode
		case filterQuote2Mode:
			switch b {
			case '\\':
				mode = filterEsc2Mode
			case '"':
				mode = filterMode
			}
			xp.token = append(xp.token, b)
		case filterEsc2Mode:
			xp.token = append(xp.token, b)
			mode = quote2Mode

		}
	}
	switch mode {
	case childMode:
		if 0 < len(xp.token) {
			xp.x = append(xp.x, Child(xp.token))
		}
	case fragMode:
		// normal termination
	default:
		return fmt.Errorf("path not terminated for %q", buf)

		// TBD error on modes expecting a completion
	}
	return
}
*/
