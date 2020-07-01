// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"strconv"

	"github.com/ohler55/ojg/gen"
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
	//rx     = &op{prec: 0, code: '~', name: "=~", cnt: 2}

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

// Script represents JSON Path script used in filters as well.
type Script struct {
	template []interface{}
	stack    []interface{}
}

// NewScript parses the string argument and returns a script or an error.
func NewScript(str string) (s *Script, err error) {
	p := &parser{buf: []byte(str)}
	if len(p.buf) == 0 || p.buf[0] != '(' {
		return nil, fmt.Errorf("a script must start with a '('")
	}
	p.pos = 1
	eq, err := p.readEquation()
	if err != nil {
		return nil, fmt.Errorf("%s at %d in %s", err, p.pos, p.buf)
	}
	return eq.Script(), nil
}

// Append a string representation of the fragment to the buffer and then
// return the expanded buffer.
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
		if pb, _ := bstack[0].(*precBuf); pb != nil {
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

// Match returns true if the script returns true when evaluated against the
// data argument.
func (s *Script) Match(data interface{}) bool {
	stack := []interface{}{}
	if node, ok := data.(gen.Node); ok {
		stack, _ = s.Eval(stack, gen.Array{node}).([]interface{})
	} else {
		stack, _ = s.Eval(stack, []interface{}{data}).([]interface{})
	}
	return 0 < len(stack)
}

// Eval is primarily used by the Expr parser but is public for testing.
func (s *Script) Eval(stack interface{}, data interface{}) interface{} {
	// Checking the type each iteration adds 2.5% but allows code not to be
	// duplicated and not to call a separate function. Using just one more
	// function call for each iteration adds 6.5%.
	var dlen int
	switch td := data.(type) {
	case []interface{}:
		dlen = len(td)
	case gen.Array:
		dlen = len(td)
	case map[string]interface{}:
		dlen = len(td)
		da := make([]interface{}, 0, dlen)
		for _, v := range td {
			da = append(da, v)
		}
		data = da
	case gen.Object:
		dlen = len(td)
		da := make(gen.Array, 0, dlen)
		for _, v := range td {
			da = append(da, v)
		}
		data = da
	default:
		return stack
	}
	var v interface{}
	for vi := 0; vi < dlen; vi++ {
		switch td := data.(type) {
		case []interface{}:
			v = td[vi]
		case gen.Array:
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
			case gen.Bool:
				s.stack[i] = bool(x)
			case gen.String:
				s.stack[i] = string(x)
			case gen.Int:
				s.stack[i] = int64(x)
			case gen.Float:
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
				if left == right {
					s.stack[i] = true
				} else {
					s.stack[i] = false
					switch tl := left.(type) {
					case int64:
						if tr, ok := right.(float64); ok {
							s.stack[i] = ok && float64(tl) == tr
						}
					case float64:
						tr, ok := right.(int64)
						s.stack[i] = ok && tl == float64(tr)
					}
				}
			case neq.code:
				if left == right {
					s.stack[i] = false
				} else {
					s.stack[i] = true
					switch tl := left.(type) {
					case int64:
						if tr, ok := right.(float64); ok {
							s.stack[i] = ok && float64(tl) != tr
						}
					case float64:
						tr, ok := right.(int64)
						s.stack[i] = ok && tl != float64(tr)
					}
				}
			case lt.code:
				s.stack[i] = false
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl < tr
					case float64:
						s.stack[i] = float64(tl) < tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl < float64(tr)
					case float64:
						s.stack[i] = tl < tr
					}
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl < tr
				}
			case gt.code:
				s.stack[i] = false
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl > tr
					case float64:
						s.stack[i] = float64(tl) > tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl > float64(tr)
					case float64:
						s.stack[i] = tl > tr
					}
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl > tr
				}
			case lte.code:
				s.stack[i] = false
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl <= tr
					case float64:
						s.stack[i] = float64(tl) <= tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl <= float64(tr)
					case float64:
						s.stack[i] = tl <= tr
					}
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl <= tr
				}
			case gte.code:
				s.stack[i] = false
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl >= tr
					case float64:
						s.stack[i] = float64(tl) >= tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl >= float64(tr)
					case float64:
						s.stack[i] = tl >= tr
					}
				case string:
					tr, ok := right.(string)
					s.stack[i] = ok && tl >= tr
				}
			case or.code:
				// If one is a boolean true then true.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				s.stack[i] = lb || rb
			case and.code:
				// If both are a boolean true then true else false.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				s.stack[i] = lb && rb
			case not.code:
				lb, _ := left.(bool)
				s.stack[i] = !lb
			case add.code:
				s.stack[i] = nil
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl + tr
					case float64:
						s.stack[i] = float64(tl) + tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl + float64(tr)
					case float64:
						s.stack[i] = tl + tr
					}
				case string:
					if tr, ok := right.(string); ok {
						s.stack[i] = tl + tr
					}
				}
			case sub.code:
				s.stack[i] = nil
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl - tr
					case float64:
						s.stack[i] = float64(tl) - tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl - float64(tr)
					case float64:
						s.stack[i] = tl - tr
					}
				}
			case mult.code:
				s.stack[i] = nil
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl * tr
					case float64:
						s.stack[i] = float64(tl) * tr
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						s.stack[i] = tl * float64(tr)
					case float64:
						s.stack[i] = tl * tr
					}
				}
			case divide.code:
				s.stack[i] = nil
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							s.stack[i] = tl / tr
						}
					case float64:
						if tr != 0.0 {
							s.stack[i] = float64(tl) / tr

						}
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							s.stack[i] = tl / float64(tr)
						}
					case float64:
						if tr != 0.0 {
							s.stack[i] = tl / tr
						}
					}
				}
			}
			if i+int(o.cnt)+1 <= len(s.stack) {
				copy(s.stack[i+1:], s.stack[i+int(o.cnt)+1:])
			}
		}
		if b, _ := s.stack[0].(bool); b {
			switch tstack := stack.(type) {
			case []interface{}:
				stack = append(tstack, v)
			case []gen.Node:
				if n, ok := v.(gen.Node); ok {
					stack = append(tstack, n)
				}
			}
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
	case nil:
		buf = append(buf, "null"...)
	case string:
		buf = append(buf, '\'')
		buf = append(buf, tv...)
		buf = append(buf, '\'')
	case int64:
		buf = append(buf, strconv.FormatInt(tv, 10)...)
	case int:
		buf = append(buf, strconv.FormatInt(int64(tv), 10)...)
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
	}
	return buf
}
