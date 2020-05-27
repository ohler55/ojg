// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

const (
	startMode = 's' // new expression
	fragMode  = 'f' // new fragment should be next
	dot2Mode  = 'o' // just read 2 dots
	openMode  = '[' // last read a [
	closeMode = ']' // expect a ]
	childMode = 'c' // reading a child fragment
	numMode   = '#'

	//   0123456789abcdef0123456789abcdef
	tokenMap = "" +
		"................................" + // 0x00
		"...o.oo....o.o.ooooooooooooooooo" + // 0x20
		".oooooooooooooooooooooooooo...oo" + // 0x40
		".oooooooooooooooooooooooooooooo." + // 0x60
		"................................" + // 0x80
		"................................" + // 0xa0
		"................................" + // 0xc0
		"................................" //   0xe0

)

// ParseExpr parses a string into an Expr.
func ParseExprString(s string) (x Expr, err error) {
	return ParseExpr([]byte(s))
}

// ParseExpr parses a []byte into an Expr.
func ParseExpr(buf []byte) (x Expr, err error) {
	x = Expr{}

	var token []byte
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
					return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
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
					return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
				}
				token = append(token, b)
				mode = childMode
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
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
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
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			default:
				if tokenMap[b] == '.' {
					return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
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
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
		case openMode:
			switch b {
			case '*':
				x = append(x, Wildcard('#'))
				mode = closeMode
			case '\'':
				// TBD
			case '"':
				// TBD
			case '-':
				mode = negMode
			case '0':
				num = 0
				mode = zeroMode
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				num = int(b - '0')
				mode = numMode
			case '?':
				// TBD filter
			case '(':
				// TBD script
			default:
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
		case closeMode:
			if b != ']' {
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
			mode = fragMode
		case zeroMode:
			switch b {
			case ']':
				x = append(x, Nth(num))
				mode = fragMode
			case ',':
				// TBD union
			case ':':
				// TBD slice
			default:
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
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
				x = append(x, Nth(num))
				mode = fragMode
			case ',':
				// TBD union
				//  create a union, replace last nth or child
				//  on ] : or , add nth num or token to union
			case ':':
				// TBD slize
			default:
				return nil, fmt.Errorf("parse expression failed at %d in %q", i, buf)
			}
			// mode for reading a string in a bracket
			// mode for slice and union
		}
		// TBD
	}
	switch mode {
	case childMode:
		if 0 < len(token) {
			x = append(x, Child(token))
		}
		// TBD error on modes expecting a completion
	}
	return
}
