// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

const (
	startMode  = 's' // new expression
	fragMode   = 'f' // new fragment should be next
	dot2Mode   = 'o' // just read 2 dots
	openMode   = '[' // last read a [
	closeMode  = ']' // expect a ]
	childMode  = 'c' // reading a child fragment
	numMode    = '#'
	quoteMode  = 'q'
	quote2Mode = 'Q'
	esc2Mode   = 'E'

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

// ParseExpr parses a string into an Expr.
func ParseExprString(s string) (x Expr, err error) {
	return ParseExpr([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func ParseExpr(buf []byte) (x Expr, err error) {
	x = Expr{}

	var token []byte
	var slice []int
	var num int

	mode := startMode
	for i, b := range buf {
		switch mode {
		case startMode:
			switch b {
			case '$':
				x = append(x, Root('$'))
				mode = fragMode
			case '@':
				x = append(x, At('@'))
				mode = fragMode
			case '[':
				x = append(x, Bracket(' '))
				mode = openMode
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("an expression can not start with a '%c'at %d in %q", b, i, buf)
				}
				token = token[:0]
				token = append(token, b)
				mode = childMode
			}
		case childMode:
			switch b {
			case '.':
				x = append(x, Child(token))
				mode = dotMode
			case '[':
				x = append(x, Child(token))
				mode = openMode
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("a '%c' character can not be in a non-bracketed child at %d in %q", b, i, buf)
				}
				token = append(token, b)
			}
		case dotMode:
			switch b {
			case '.':
				x = append(x, Descent('.'))
				mode = dot2Mode
			case '*':
				x = append(x, Wildcard('*'))
				mode = fragMode
			case '[':
				return nil, fmt.Errorf("unexpected '[' after a '.' at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("a '%c' character can not be in a non-bracketed child at %d in %q", b, i, buf)
				}
				token = token[:0]
				token = append(token, b)
				mode = childMode
			}
		case dot2Mode:
			switch b {
			case '*':
				x = append(x, Wildcard('*'))
				mode = fragMode
			case '[':
				return nil, fmt.Errorf("a '[' can not follow '..' at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("a '%c' can not follow a '..' at %d in %q", b, i, buf)
				}
				token = token[:0]
				token = append(token, b)
				mode = childMode
			}
		case fragMode:
			switch b {
			case '.':
				mode = dotMode
			case '[':
				mode = openMode
			default:
				return nil, fmt.Errorf("expected a '.' or a '[' at %d in %q", i, buf)
			}
		case openMode:
			switch b {
			case '*':
				x = append(x, Wildcard('#'))
				mode = closeMode
			case '\'':
				token = token[:0]
				mode = quoteMode
			case '"':
				token = token[:0]
				mode = quote2Mode
			case '-':
				mode = negMode
			case '0':
				num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				num = int(b - '0')
				mode = numMode
			case ':':
				slice = slice[:0]
				slice = append(slice, 0)
				mode = colonMode
			case '?':
				// TBD filter
			case '(':
				// TBD script
			default:
				return nil, fmt.Errorf("a '%c' can not follow a '[' at %d in %q", b, i, buf)
			}
		case closeMode:
			if b != ']' {
				return nil, fmt.Errorf("expected a ']' at %d in %q", i, buf)
			}
			mode = fragMode
		case zeroMode:
			switch b {
			case ']':
				if 0 < len(slice) {
					slice = append(slice, num)
					ia := make([]int, len(slice))
					copy(ia, slice)
					slice = slice[:0]
					x = append(x, Slice(ia))
				} else {
					x = append(x, Nth(num))
				}
				mode = fragMode
			case ',':
				// TBD union, error is a slice exists
			case ':':
				if 2 < len(slice) {
					return nil, fmt.Errorf("too many numbers in the slice at %d in %q", i, buf)
				}
				slice = append(slice, 0)
				mode = colonMode
			default:
				return nil, fmt.Errorf("unexpected character at %d in %q", i, buf)
			}
		case colonMode:
			switch b {
			case '-':
				mode = negMode
			case '0':
				num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				num = int(b - '0')
				mode = numMode
			case ']':
				if 0 < len(slice) {
					slice = append(slice, -1)
					ia := make([]int, len(slice))
					copy(ia, slice)
					slice = slice[:0]
					x = append(x, Slice(ia))
				} else {
					x = append(x, Nth(num))
				}
				mode = fragMode
			default:
				return nil, fmt.Errorf("invalid slice format at %d in %q", i, buf)
			}
		case negMode:
			switch b {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				num = -int(b - '0')
				mode = numMode
			default:
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
		case numMode:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if 0 <= num {
					num = num*10 + int(b-'0')
				} else {
					num = num*10 - int(b-'0')
				}
				mode = numMode
			case ']':
				if 0 < len(slice) {
					slice = append(slice, num)
					ia := make([]int, len(slice))
					copy(ia, slice)
					slice = slice[:0]
					x = append(x, Slice(ia))
				} else {
					x = append(x, Nth(num))
				}
				mode = fragMode
			case ',':
				// TBD union
				//  create a union, replace last nth or child
				//  on ] : or , add nth num or token to union
			case ':':
				if 2 < len(slice) {
					return nil, fmt.Errorf("too many numbers in the slice at %d in %q", i, buf)
				}
				slice = append(slice, num)
				mode = colonMode
			default:
				return nil, fmt.Errorf("invalid number at %d in %q", i, buf)
			}
		case quoteMode:
			switch b {
			case '\\':
				mode = escMode
			case '\'':
				x = append(x, Child(token))
				// TBD close of comma mode
				mode = closeMode
			default:
				token = append(token, b)
			}
		case escMode:
			if b != '\'' {
				token = append(token, '\\')
			}
			token = append(token, b)
			mode = quoteMode
		case quote2Mode:
			switch b {
			case '\\':
				mode = esc2Mode
			case '"':
				x = append(x, Child(token))
				// TBD close of comma mode
				mode = closeMode
			default:
				token = append(token, b)
			}
		case esc2Mode:
			if b != '"' {
				token = append(token, '\\')
			}
			token = append(token, b)
			mode = quoteMode

			// mode for union
		}
	}
	switch mode {
	case childMode:
		if 0 < len(token) {
			x = append(x, Child(token))
		}
	case fragMode:
		// normal termination
	default:
		return nil, fmt.Errorf("path not terminated for %q", buf)

		// TBD error on modes expecting a completion
	}
	return
}
