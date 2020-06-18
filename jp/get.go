// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"github.com/ohler55/ojg/gen"
)

const (
	fragIndexMask    = 0x0000ffff
	descentFlag      = 0x00010000
	descentChildFlag = 0x00020000
)

type fragIndex int

// The easy way to implement the Get is to have each fragment handle the
// getting using recursion. The overhead of a go function call is rather high
// though so instead a psuedo call stack is implemented here that grows and
// shrinks as the getting takes place. The fragment index if placed on the
// stack as well mostly for a small degree of simplicity in what a few people
// might find a complex approach to the solution. Its at least twice as fast
// as the recursive function call approach and in some cases such as the
// recursive descent more than an order of magnitude faster.

// Get the elements of the data identified by the path.
func (x Expr) Get(data interface{}) (results []interface{}) {
	if len(x) == 0 {
		return
	}
	var v interface{}
	var prev interface{}
	var has bool

	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		has = false
		switch tf := f.(type) {
		case Child:
			switch tv := prev.(type) {
			case map[string]interface{}:
				v, has = tv[string(tf)]
			case gen.Object:
				v, has = tv[string(tf)]
			default:
				v, has = x.reflectGetChild(tv, string(tf))
			}
			if has {
				if int(fi) == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					v = tv[i]
					has = true
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					v = tv[i]
					has = true
				}
			default:
				v, has = x.reflectGetNth(tv, i)
			}
			if has {
				if int(fi) == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case []interface{}:
				if int(fi) == len(x)-1 { // last one
					results = append(results, tv...)
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				for _, v := range x.reflectGetWild(tv) {
					if int(fi) == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						results = append(results, tv...)
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				}
			} else {
				if int(fi) == len(x)-1 { // last one
					if top {
						results = append(results, prev)
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if int(fi) == len(x)-1 { // last one
				results = append(results, data)
			} else {
				stack = append(stack, data)
			}
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				results = append(results, prev)
			} else {
				stack = append(stack, prev)
			}
		case Union:
			for _, u := range tf {
				has = false
				switch tu := u.(type) {
				case string:
					switch tv := prev.(type) {
					case map[string]interface{}:
						v, has = tv[string(tu)]
					case gen.Object:
						v, has = tv[string(tu)]
					default:
						v, has = x.reflectGetChild(tv, string(tu))
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							has = true
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							has = true
						}
					default:
						v, has = x.reflectGetNth(tv, i)
					}
				}
				if has {
					if int(fi) == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Slice:
			start := 0
			end := -1
			step := 1
			if 0 < len(tf) {
				start = tf[0]
			}
			if 1 < len(tf) {
				end = tf[1]
			}
			if 2 < len(tf) {
				step = tf[2]
			}
			switch tv := prev.(type) {
			case []interface{}:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end + 1
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) < end || step == 0 {
					continue
				}
				var v interface{}
				if 0 < step {
					for i := start; i < end; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			case gen.Array:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end + 1
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) < end || step == 0 {
					continue
				}
				var v interface{}
				if 0 < step {
					for i := start; i < end; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			default:
				for _, v := range x.reflectGetSlice(tv, start, end, step) {
					if int(fi) == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
			if int(fi) == len(x)-1 { // last one
				results = append(results, stack[len(stack)-1])
			}
		}
		if int(fi) < len(x)-1 {
			if _, ok := stack[len(stack)-1].(fragIndex); !ok {
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		}
	}
	// Free up anything still on the stack.
	stack = stack[0:cap(stack)]
	for i := len(stack) - 1; 0 <= i; i-- {
		stack[i] = nil
	}
	return
}

// First element of the data identified by the path.
func (x Expr) First(data interface{}) interface{} {
	if len(x) == 0 {
		return nil
	}
	var v interface{}
	var prev interface{}
	var has bool

	stack := make([]interface{}, 0, 64)
	defer func() {
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	}()
	stack = append(stack, data)
	f := x[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		has = false
		switch tf := f.(type) {
		case Child:
			switch tv := prev.(type) {
			case map[string]interface{}:
				v, has = tv[string(tf)]
			case gen.Object:
				v, has = tv[string(tf)]
			default:
				v, has = x.reflectGetChild(tv, string(tf))
			}
			if has {
				if int(fi) == len(x)-1 { // last one
					return v
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					v = tv[i]
					has = true
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					v = tv[i]
					has = true
				}
			default:
				v, has = x.reflectGetNth(tv, i)
			}
			if has {
				if int(fi) == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case map[string]interface{}, []interface{}, gen.Object, gen.Array:
					stack = append(stack, v)
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case []interface{}:
				if int(fi) == len(x)-1 { // last one
					if 0 < len(tv) {
						return tv[0]
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if int(fi) == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				if v, has = x.reflectGetWildOne(tv); has {
					if int(fi) == len(x)-1 { // last one
						return v
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if int(fi) == len(x)-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				}
			} else {
				if int(fi) == len(x)-1 { // last one
					if top {
						return prev
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if int(fi) == len(x)-1 { // last one
				return data
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				return prev
			}
			stack = append(stack, prev)
		case Union:
			for _, u := range tf {
				has = false
				switch tu := u.(type) {
				case string:
					switch tv := prev.(type) {
					case map[string]interface{}:
						v, has = tv[string(tu)]
					case gen.Object:
						v, has = tv[string(tu)]
					default:
						v, has = x.reflectGetChild(tv, string(tu))
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							has = true
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							has = true
						}
					default:
						v, has = x.reflectGetNth(tv, i)
					}
				}
				if has {
					if int(fi) == len(x)-1 { // last one
						return v
					}
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Slice:
			start := 0
			if 0 < len(tf) {
				start = tf[0]
			}
			switch tv := prev.(type) {
			case []interface{}:
				if start < 0 {
					start = len(tv) + start
				}
				if start < 0 || len(tv) <= start {
					continue
				}
				v := tv[start]
				if int(fi) == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case map[string]interface{}, []interface{}, gen.Object, gen.Array:
					stack = append(stack, v)
				}
			case gen.Array:
				if start < 0 {
					start = len(tv) + start
				}
				if start < 0 || len(tv) <= start {
					continue
				}
				v := tv[start]
				if int(fi) == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case gen.Object, gen.Array:
					stack = append(stack, v)
				}
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
			if int(fi) == len(x)-1 { // last one
				return stack[len(stack)-1]
			}
		}
		if int(fi) < len(x)-1 {
			if _, ok := stack[len(stack)-1].(fragIndex); !ok {
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		}
	}
	return nil
}

func (x Expr) reflectGetChild(data interface{}, key string) (v interface{}, has bool) {

	// TBD

	return
}

func (x Expr) reflectGetNth(data interface{}, i int) (v interface{}, has bool) {

	// TBD

	return
}

func (x Expr) reflectGetWild(data interface{}) (va []interface{}) {

	// TBD

	return
}

func (x Expr) reflectGetWildOne(data interface{}) (v interface{}, has bool) {

	// TBD

	return
}

func (x Expr) reflectGetSlice(data interface{}, start, end, step int) (va []interface{}) {

	// TBD

	return
}
