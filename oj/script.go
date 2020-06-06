// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"strconv"
)

var (
	eq     = &op{prec: 3, code: '=', name: "==", cnt: 2}
	neq    = &op{prec: 3, code: 'n', name: "!=", cnt: 2}
	lt     = &op{prec: 3, code: '<', name: "<", cnt: 2}
	gt     = &op{prec: 3, code: '>', name: ">", cnt: 2}
	lte    = &op{prec: 3, code: 'l', name: "<=", cnt: 2}
	gte    = &op{prec: 3, code: 'g', name: ">=", cnt: 2}
	or     = &op{prec: 4, code: '|', name: "||", cnt: 2}
	and    = &op{prec: 4, code: '&', name: "&&", cnt: 2}
	not    = &op{prec: 0, code: '!', name: "!", cnt: 1}
	add    = &op{prec: 2, code: '+', name: "+", cnt: 2}
	sub    = &op{prec: 2, code: '-', name: "-", cnt: 2}
	mult   = &op{prec: 1, code: '*', name: "*", cnt: 2}
	divide = &op{prec: 1, code: '/', name: "/", cnt: 2}
	get    = &op{prec: 0, code: 'G', name: "get", cnt: 1}

	opMap = map[string]*op{
		eq.name:     eq,
		neq.name:    neq,
		lt.name:     lt,
		gt.name:     gt,
		lte.name:    lte,
		gte.name:    gte,
		or.name:     or,
		and.name:    and,
		not.name:    not,
		add.name:    add,
		sub.name:    sub,
		mult.name:   mult,
		divide.name: divide,
	}
)

type op struct {
	name string
	prec byte
	cnt  byte
	code byte
}

type precBuf struct {
	prec byte
	buf  []byte
}

// String returns the op name.
func (o *op) String() string {
	return o.name
}

// Script represents JSON Path script used in filters as well.
type Script struct {
	template []interface{}
	stack    []interface{}
}

func NewScript(str string) (s *Script, err error) {
	xp := &xparser{buf: []byte(str)}
	if len(xp.buf) == 0 || xp.buf[0] != '(' {
		return nil, fmt.Errorf("a script must start with a '('")
	}
	xp.pos = 1
	eq, err := xp.readEquation()
	if err == nil && xp.pos < len(xp.buf) {
		err = fmt.Errorf("parse error")
	}
	if err != nil {
		err = fmt.Errorf("%s at %d in %s", err, xp.pos, xp.buf)
	}
	return eq.Script(), nil
}

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (s *Script) Append(buf []byte) []byte {
	buf = append(buf, '(')
	if 0 < len(s.template) {
		bstack := make([]interface{}, len(s.template))
		copy(bstack, s.template)

		for i := len(bstack) - 1; 0 <= i; i-- {
			o, _ := bstack[i].(*op)
			if o == nil {
				continue
			}
			var left interface{}
			var right interface{}
			if 1 < len(bstack)-i {
				left = bstack[i+1]
			}
			if 2 < len(bstack)-i {
				right = bstack[i+2]
			}
			bstack[i] = s.appendOp(o, left, right)
			if i+int(o.cnt)+1 <= len(bstack) {
				copy(bstack[i+1:], bstack[i+int(o.cnt)+1:])
			}
		}
		if pb, _ := bstack[0].(*precBuf); pb == nil {
			buf = s.appendValue(buf, bstack[0], 0)
		} else {
			buf = append(buf, pb.buf...)
		}
	}
	buf = append(buf, ')')

	return buf
}

// String representation of the script.
func (s *Script) String() string {
	return string(s.Append([]byte{}))
}

func (s *Script) Eval(stack []interface{}, data interface{}) []interface{} {
	// Checking the type each iteration adds 2.5% but allows code not to be
	// duplicated and not to call a separate function. Using just one function
	// call for each iteration adds 6.5%.
	var dlen int
	switch td := data.(type) {
	case []interface{}:
		dlen = len(td)
	case Array:
		dlen = len(td)
	default:
		return stack
	}
	var v interface{}
	for vi := 0; vi < dlen; vi++ {
		switch td := data.(type) {
		case []interface{}:
			v = td[vi]
		case Array:
			v = td[vi]
		}
		// Eval script for each member of the list.
		copy(s.stack, s.template)
		// resolve all expr members
		for i, ev := range s.stack {
			// Normalize into nil, bool, int64, float64, and string early so
			// that each comparison doen't have to.
		Normalize:
			switch x := ev.(type) {
			case Expr:
				// The most common pattern is [?(@.child == value)] where
				// the operation and value vary but the @.child is the
				// most widely used. For that reason an optimization is
				// included for that inclusion of a one level child lookup
				// path.
				if m, ok := v.(map[string]interface{}); ok && len(x) == 2 {
					if _, ok = x[0].(At); ok {
						var c Child
						if c, ok = x[1].(Child); ok {
							ev = m[string(c)]
							s.stack[i] = ev
							goto Normalize
						}
					}
				}
				ev = x.First(v)
				s.stack[i] = ev
				goto Normalize
			case int:
				s.stack[i] = int64(x)
			case int8:
				s.stack[i] = int64(x)
			case int16:
				s.stack[i] = int64(x)
			case int32:
				s.stack[i] = int64(x)
			case uint:
				s.stack[i] = int64(x)
			case uint8:
				s.stack[i] = int64(x)
			case uint16:
				s.stack[i] = int64(x)
			case uint32:
				s.stack[i] = int64(x)
			case uint64:
				s.stack[i] = int64(x)
			case float32:
				s.stack[i] = float64(x)
			case Bool:
				s.stack[i] = bool(x)
			case String:
				s.stack[i] = string(x)
			case Int:
				s.stack[i] = int64(x)
			case Float:
				s.stack[i] = float64(x)

			default:
				// Any other type are already simplified or are not
				// handled and will fail later.
			}
		}
		for i := len(s.stack) - 1; 0 <= i; i-- {
			o, _ := s.stack[i].(*op)
			if o == nil {
				// a value, not an op
				continue
			}
			var left interface{}
			var right interface{}
			if 1 < len(s.stack)-i {
				left = s.stack[i+1]
			}
			if 2 < len(s.stack)-i {
				right = s.stack[i+2]
			}
			switch o.code {
			case eq.code:
				s.stack[i] = left == right
			case neq.code:
				s.stack[i] = left != right
			case lt.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					s.stack[i] = ok && tl < tr
				case float64:
					tr, ok := right.(float64)
					s.stack[i] = ok && tl < tr
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl < tr
				default:
					s.stack[i] = false
				}
			case gt.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					s.stack[i] = ok && tl > tr
				case float64:
					tr, ok := right.(float64)
					s.stack[i] = ok && tl > tr
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl > tr
				default:
					s.stack[i] = false
				}
			case lte.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					s.stack[i] = ok && tl <= tr
				case float64:
					tr, ok := right.(float64)
					s.stack[i] = ok && tl <= tr
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl <= tr
				default:
					s.stack[i] = false
				}
			case gte.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					s.stack[i] = ok && tl >= tr
				case float64:
					tr, ok := right.(float64)
					s.stack[i] = ok && tl >= tr
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl >= tr
				default:
					s.stack[i] = false
				}
			case or.code:
				// If one is a boolean true then true.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				s.stack[i] = lb || rb
			case and.code:
				// If both is a boolean true then true else false.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				s.stack[i] = lb && rb
			case not.code:
				lb, _ := left.(bool)
				s.stack[i] = !lb
			case add.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl + tr
					case float64:
						s.stack[i] = float64(tl) + tr
					default:
						s.stack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl + float64(tr)
					case float64:
						s.stack[i] = tl + tr
					default:
						s.stack[i] = nil
					}
				case string:
					if tr, ok := right.(string); ok {
						s.stack[i] = tl + tr
					} else {
						s.stack[i] = nil
					}
				default:
					s.stack[i] = false
				}
			case sub.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl - tr
					case float64:
						s.stack[i] = float64(tl) - tr
					default:
						s.stack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl - float64(tr)
					case float64:
						s.stack[i] = tl - tr
					default:
						s.stack[i] = nil
					}
				default:
					s.stack[i] = false
				}
			case mult.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl * tr
					case float64:
						s.stack[i] = float64(tl) * tr
					default:
						s.stack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl * float64(tr)
					case float64:
						s.stack[i] = tl * tr
					default:
						s.stack[i] = nil
					}
				default:
					s.stack[i] = false
				}
			case divide.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							s.stack[i] = tl / tr
						} else {
							s.stack[i] = nil
						}
					case float64:
						if tr != 0.0 {
							s.stack[i] = float64(tl) / tr
						} else {
							s.stack[i] = nil
						}
					default:
						s.stack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							s.stack[i] = tl / float64(tr)
						} else {
							s.stack[i] = nil
						}
					case float64:
						if tr != 0.0 {
							s.stack[i] = tl / tr
						} else {
							s.stack[i] = nil
						}
					default:
						s.stack[i] = nil
					}
				default:
					s.stack[i] = false
				}
				if i+int(o.cnt)+1 <= len(s.stack) {
					copy(s.stack[i+1:], s.stack[i+int(o.cnt)+1:])
				}
				// TBD slide tail
			}
		}
		if b, _ := s.stack[0].(bool); b {
			stack = append(stack, v)
		}
	}
	for i := range s.stack {
		s.stack[i] = nil
	}
	return stack
}

func (s *Script) appendOp(o *op, left, right interface{}) (pb *precBuf) {
	pb = &precBuf{prec: o.prec}
	switch o.code {
	case not.code:
		pb.buf = append(pb.buf, o.name...)
		pb.buf = s.appendValue(pb.buf, left, o.prec)
	default:
		pb.buf = s.appendValue(pb.buf, left, o.prec)
		pb.buf = append(pb.buf, ' ')
		pb.buf = append(pb.buf, o.name...)
		pb.buf = append(pb.buf, ' ')
		pb.buf = s.appendValue(pb.buf, right, o.prec)
	}
	return
}

func (s *Script) appendValue(buf []byte, v interface{}, prec byte) []byte {
	switch tv := v.(type) {
	case string:
		buf = append(buf, '\'')
		buf = append(buf, tv...)
		buf = append(buf, '\'')
	case int64:
		buf = append(buf, strconv.FormatInt(tv, 10)...)
	case int:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case int8:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case int16:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case int32:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case uint:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case uint8:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case uint16:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case uint32:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case uint64:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
	case float32:
		buf = append(buf, strconv.FormatFloat(float64(tv), 'g', -1, 32)...)
	case float64:
		buf = append(buf, strconv.FormatFloat(tv, 'g', -1, 64)...)
	case bool:
		if tv {
			buf = append(buf, "true"...)
		} else {
			buf = append(buf, "false"...)
		}
	case Expr:
		buf = tv.Append(buf)
	case *precBuf:
		if prec < tv.prec {
			buf = append(buf, '(')
			buf = append(buf, tv.buf...)
			buf = append(buf, ')')
		} else {
			buf = append(buf, tv.buf...)
		}

	case fmt.Stringer:
		buf = append(buf, tv.String()...)
	default:
		buf = append(buf, fmt.Sprintf("%v", v)...)
	}
	return buf
}
