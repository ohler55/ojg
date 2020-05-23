// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

const (
	equal  = "=="
	length = "length"
)

var (
	eq     = &op{prec: 3, code: '=', name: "=="}
	neq    = &op{prec: 3, code: 'n', name: "!="}
	lt     = &op{prec: 3, code: '<', name: "<"}
	gt     = &op{prec: 3, code: '>', name: ">"}
	lte    = &op{prec: 3, code: 'l', name: "<="}
	gte    = &op{prec: 3, code: 'g', name: ">="}
	or     = &op{prec: 4, code: '|', name: "||"}
	and    = &op{prec: 4, code: '&', name: "&&"}
	not    = &op{prec: 0, code: '!', name: "!"}
	add    = &op{prec: 2, code: '+', name: "+"}
	sub    = &op{prec: 2, code: '-', name: "-"}
	mult   = &op{prec: 1, code: '*', name: "*"}
	divide = &op{prec: 1, code: '/', name: "/"}
	// TBD add more as desired
)

type op struct {
	name string
	prec byte
	code byte
}

// String returns the op name.
func (o *op) String() string {
	return o.name
}

// Script represents JSON Path script used in filters as well.
type Script []interface{}

type Filter Script

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (s Script) Append(buf []byte) []byte {

	// TBD

	return buf
}

func (f Filter) Eval(stack []interface{}, data interface{}) []interface{} {
	estack := make([]interface{}, len(f))
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
		// Eval filter for each member of the list.
		copy(estack, f)
		// resolve all expr members
		for i, ev := range estack {
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
							estack[i] = ev
							goto Normalize
						}
					}
				}
				ev = x.First(v)
				estack[i] = ev
				goto Normalize
			case int:
				estack[i] = int64(x)
			case int8:
				estack[i] = int64(x)
			case int16:
				estack[i] = int64(x)
			case int32:
				estack[i] = int64(x)
			case uint:
				estack[i] = int64(x)
			case uint8:
				estack[i] = int64(x)
			case uint16:
				estack[i] = int64(x)
			case uint32:
				estack[i] = int64(x)
			case uint64:
				estack[i] = int64(x)
			case float32:
				estack[i] = float64(x)
			case Bool:
				estack[i] = bool(x)
			case String:
				estack[i] = string(x)
			case Int:
				estack[i] = int64(x)
			case Float:
				estack[i] = float64(x)

			default:
				// Any other type are already simplified or are not
				// handled and will fail later.
			}
		}
		for i := len(estack) - 1; 0 <= i; i-- {
			o, _ := estack[i].(*op)
			if o == nil {
				// a value, not an op
				continue
			}
			var left interface{}
			var right interface{}
			var acnt int
			if 1 < len(estack)-i {
				left = estack[i+1]
				acnt++
			}
			if 2 < len(estack)-i {
				right = estack[i+2]
				acnt++
			}
			switch o.code {
			case eq.code:
				estack[i] = left == right
			case neq.code:
				estack[i] = left != right
			case lt.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					estack[i] = ok && tl < tr
				case float64:
					tr, ok := right.(float64)
					estack[i] = ok && tl < tr
				case string:
					tr, ok := right.(string)
					estack[i] = ok && tl < tr
				default:
					estack[i] = false
				}
			case gt.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					estack[i] = ok && tl > tr
				case float64:
					tr, ok := right.(float64)
					estack[i] = ok && tl > tr
				case string:
					tr, ok := right.(string)
					estack[i] = ok && tl > tr
				default:
					estack[i] = false
				}
			case lte.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					estack[i] = ok && tl <= tr
				case float64:
					tr, ok := right.(float64)
					estack[i] = ok && tl <= tr
				case string:
					tr, ok := right.(string)
					estack[i] = ok && tl <= tr
				default:
					estack[i] = false
				}
			case gte.code:
				switch tl := left.(type) {
				case int64:
					tr, ok := right.(int64)
					estack[i] = ok && tl >= tr
				case float64:
					tr, ok := right.(float64)
					estack[i] = ok && tl >= tr
				case string:
					tr, ok := right.(string)
					estack[i] = ok && tl >= tr
				default:
					estack[i] = false
				}
			case or.code:
				// If one is a boolean true then true.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				estack[i] = lb || rb
			case and.code:
				// If both is a boolean true then true else false.
				lb, _ := left.(bool)
				rb, _ := right.(bool)
				estack[i] = lb && rb
			case not.code:
				lb, _ := left.(bool)
				estack[i] = !lb
			case add.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl + tr
					case float64:
						estack[i] = float64(tl) + tr
					default:
						estack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl + float64(tr)
					case float64:
						estack[i] = tl + tr
					default:
						estack[i] = nil
					}
				case string:
					if tr, ok := right.(string); ok {
						estack[i] = tl + tr
					} else {
						estack[i] = nil
					}
				default:
					estack[i] = false
				}
			case sub.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl - tr
					case float64:
						estack[i] = float64(tl) - tr
					default:
						estack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl - float64(tr)
					case float64:
						estack[i] = tl - tr
					default:
						estack[i] = nil
					}
				default:
					estack[i] = false
				}
			case mult.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl * tr
					case float64:
						estack[i] = float64(tl) * tr
					default:
						estack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						estack[i] = tl * float64(tr)
					case float64:
						estack[i] = tl * tr
					default:
						estack[i] = nil
					}
				default:
					estack[i] = false
				}
			case divide.code:
				switch tl := left.(type) {
				case int64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							estack[i] = tl / tr
						} else {
							estack[i] = nil
						}
					case float64:
						if tr != 0.0 {
							estack[i] = float64(tl) / tr
						} else {
							estack[i] = nil
						}
					default:
						estack[i] = nil
					}
				case float64:
					switch tr := right.(type) {
					case int64:
						if tr != 0 {
							estack[i] = tl / float64(tr)
						} else {
							estack[i] = nil
						}
					case float64:
						if tr != 0.0 {
							estack[i] = tl / tr
						} else {
							estack[i] = nil
						}
					default:
						estack[i] = nil
					}
				default:
					estack[i] = false
				}
			}
		}
		if b, _ := estack[0].(bool); b {
			stack = append(stack, v)
		}
	}
	for i, _ := range estack {
		estack[i] = nil
	}
	return stack
}

// TBD remove
func (s Script) Foo() Script {
	s = append(s, lt)
	s = append(s, A().C("a"))
	s = append(s, int64(52))
	return s
}
