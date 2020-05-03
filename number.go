// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"math"
	"strconv"

	"github.com/ohler55/ojg/gd"
)

// 9223372036854775807 / 10 = 922337203685477580
const bigLimit = math.MaxInt64 / 10

type number struct {
	i      uint64
	frac   uint64
	div    uint64
	exp    uint64
	neg    bool
	negExp bool
	bigBuf []byte
}

func (n *number) reset() {
	n.i = 0
	n.frac = 0
	n.div = 1
	n.neg = false
	n.negExp = false
	if 0 < len(n.bigBuf) {
		n.bigBuf = n.bigBuf[:0]
	}
}

func (n *number) addDigit(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.i <= bigLimit {
		n.i = n.i*10 + uint64(b-'0')
		if math.MaxInt64 < n.i {
			n.fillBig()
		}
	} else {
		n.fillBig()
		n.bigBuf = append(n.bigBuf, b)
	}
}

func (n *number) addFrac(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.frac <= bigLimit {
		n.frac = n.frac*10 + uint64(b-'0')
		n.div *= 10.0
		if math.MaxInt64 < n.frac {
			n.fillBig()
		}
	} else { // big
		n.fillBig()
		n.bigBuf = append(n.bigBuf, b)
	}
}

func (n *number) addExp(b byte) {
	if 0 < len(n.bigBuf) {
		n.bigBuf = append(n.bigBuf, b)
	} else if n.exp <= 102 {
		n.exp = n.exp*10 + uint64(b-'0')
		if 1022 < n.exp {
			n.fillBig()
		}
	} else { // big
		n.fillBig()
		n.bigBuf = append(n.bigBuf, b)
	}
}

func (n *number) fillBig() {
	if n.neg {
		n.bigBuf = append(n.bigBuf, '-')
	}
	n.bigBuf = append(n.bigBuf, strconv.FormatUint(n.i, 10)...)
	if 0 < n.frac {
		n.bigBuf = append(n.bigBuf, '.')
		if 1000000000000000000 <= n.frac { // nearest multiple of 10 below max int64
			n.bigBuf = append(n.bigBuf, strconv.FormatUint(n.frac, 10)...)
		} else {
			s := strconv.FormatUint(n.frac+n.div, 10)
			n.bigBuf = append(n.bigBuf, s[1:]...)
		}
	}
	if 0 < n.exp {
		n.bigBuf = append(n.bigBuf, 'e')
		if n.negExp {
			n.bigBuf = append(n.bigBuf, '-')
		}
		n.bigBuf = append(n.bigBuf, strconv.FormatUint(n.exp, 10)...)
	}
}

func (n *number) asInt() int64 {
	i := int64(n.i)
	if n.neg {
		i = -i
	}
	return i
}

func (n *number) asFloat() (float64, error) {
	f := float64(n.i)
	if 0 < n.frac {
		f += float64(n.frac) / float64(n.div)
	}
	if n.neg {
		f = -f
	}
	if 0 < n.exp {
		x := int(n.exp)
		if n.negExp {
			x = -x
		}
		f *= math.Pow10(int(x))
	}
	return f, nil
}

func (n *number) asBig() gd.Big {
	return gd.Big(n.bigBuf)
}
