// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

import (
	"encoding/json"
	"math"
	"strconv"
)

// BigLimit is the limit before a number is converted into a Big
// instance. (9223372036854775807 / 10 = 922337203685477580)
const BigLimit = math.MaxInt64 / 10

// Number is used internally by parsers.
type Number struct {
	I          uint64
	Frac       uint64
	Div        uint64
	Exp        uint64
	Neg        bool
	NegExp     bool
	BigBuf     []byte
	ForceFloat bool
}

// Reset the number.
func (n *Number) Reset() {
	n.I = 0
	n.Frac = 0
	n.Div = 1
	n.Exp = 0
	n.Neg = false
	n.NegExp = false
	if 0 < len(n.BigBuf) {
		n.BigBuf = n.BigBuf[:0]
	}
}

// AddDigit to a number.
func (n *Number) AddDigit(b byte) {
	if 0 < len(n.BigBuf) {
		n.BigBuf = append(n.BigBuf, b)
	} else if n.I <= BigLimit {
		n.I = n.I*10 + uint64(b-'0')
		if math.MaxInt64 < n.I {
			n.FillBig()
		}
	} else {
		n.FillBig()
		n.BigBuf = append(n.BigBuf, b)
	}
}

// AddFrac adds a fractional digit.
func (n *Number) AddFrac(b byte) {
	if 0 < len(n.BigBuf) {
		n.BigBuf = append(n.BigBuf, b)
	} else if n.Frac <= BigLimit {
		n.Frac = n.Frac*10 + uint64(b-'0')
		n.Div *= 10.0
		if math.MaxInt64 < n.Frac {
			n.FillBig()
		}
	} else { // big
		n.FillBig()
		n.BigBuf = append(n.BigBuf, b)
	}
}

// AddExp adds an exponent digit.
func (n *Number) AddExp(b byte) {
	if 0 < len(n.BigBuf) {
		n.BigBuf = append(n.BigBuf, b)
	} else if n.Exp <= 102 {
		n.Exp = n.Exp*10 + uint64(b-'0')
		if 1022 < n.Exp {
			n.FillBig()
		}
	} else { // big
		n.FillBig()
		n.BigBuf = append(n.BigBuf, b)
	}
}

// FillBig fills the internal buffer with a big number.
func (n *Number) FillBig() {
	if n.Neg {
		n.BigBuf = append(n.BigBuf, '-')
	}
	n.BigBuf = append(n.BigBuf, strconv.FormatUint(n.I, 10)...)
	if 0 < n.Frac {
		n.BigBuf = append(n.BigBuf, '.')
		if 1000000000000000000 <= n.Frac { // nearest multiple of 10 below max int64
			n.BigBuf = append(n.BigBuf, strconv.FormatUint(n.Frac, 10)...)
		} else {
			s := strconv.FormatUint(n.Frac+n.Div, 10)
			n.BigBuf = append(n.BigBuf, s[1:]...)
		}
	}
	if 0 < n.Exp {
		n.BigBuf = append(n.BigBuf, 'e')
		if n.NegExp {
			n.BigBuf = append(n.BigBuf, '-')
		}
		n.BigBuf = append(n.BigBuf, strconv.FormatUint(n.Exp, 10)...)
	}
}

// AsNum returns the number as best fit.
func (n *Number) AsNum() (num interface{}) {
	if 0 < len(n.BigBuf) {
		num = json.Number(n.BigBuf)
	} else if n.Div == 1 && n.Exp == 0 {
		i := int64(n.I)
		if n.Neg {
			i = -i
		}
		if n.ForceFloat {
			num = float64(i)
		} else {
			num = i
		}
	} else {
		f := float64(n.I)
		if 0 < n.Frac {
			f += float64(n.Frac) / float64(n.Div)
		}
		if n.Neg {
			f = -f
		}
		if 0 < n.Exp {
			x := int(n.Exp)
			if n.NegExp {
				x = -x
			}
			f *= math.Pow10(int(x))
		}
		num = f
	}
	return
}

// AsNode returns the number as best fit.
func (n *Number) AsNode() (num Node) {
	if 0 < len(n.BigBuf) {
		num = Big(n.BigBuf)
	} else if n.Frac == 0 && n.Exp == 0 {
		i := int64(n.I)
		if n.Neg {
			i = -i
		}
		num = Int(i)
	} else {
		f := float64(n.I)
		if 0 < n.Frac {
			f += float64(n.Frac) / float64(n.Div)
		}
		if n.Neg {
			f = -f
		}
		if 0 < n.Exp {
			x := int(n.Exp)
			if n.NegExp {
				x = -x
			}
			f *= math.Pow10(int(x))
		}
		num = Float(f)
	}
	return
}
