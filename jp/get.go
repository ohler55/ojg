// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import "github.com/ohler55/ojg/gen"

const (
	fragIndexMask    = 0x0000ffff
	descentFlag      = 0x00010000
	descentChildFlag = 0x00020000
)

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

	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(int); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]interface{}:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 <= i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 <= i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				// TBD try reflection
				continue
			}
		case Descent:
			di, _ := stack[len(stack)-1].(int)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
					if top {
						results = append(results, prev)
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == len(x)-1 { // last one
				results = append(results, data)
			} else {
				stack = append(stack, data)
			}
		case At, Bracket:
			if fi == len(x)-1 { // last one
				results = append(results, prev)
			} else {
				stack = append(stack, prev)
			}
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					var has bool
					switch tv := prev.(type) {
					case map[string]interface{}:
						if v, has = tv[string(tu)]; has {
							if fi == len(x)-1 { // last one
								results = append(results, v)
							} else {
								switch v.(type) {
								case map[string]interface{}, []interface{}, gen.Object, gen.Array:
									stack = append(stack, v)
								}
							}
						}
					case gen.Object:
						if v, has = tv[string(tu)]; has {
							if fi == len(x)-1 { // last one
								results = append(results, v)
							} else {
								switch v.(type) {
								case map[string]interface{}, []interface{}, gen.Object, gen.Array:
									stack = append(stack, v)
								}
							}
						}
					default:
						// TBD try reflection
						continue
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						if fi == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						if fi == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						// TBD reflection
						continue
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
						if fi == len(x)-1 { // last one
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
						if fi == len(x)-1 { // last one
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
						if fi == len(x)-1 { // last one
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
						if fi == len(x)-1 { // last one
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
				// TBD try reflection
				continue
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
		}
		if fi < len(x)-1 {
			if _, ok := stack[len(stack)-1].(int); !ok {
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

	stack := make([]interface{}, 0, 64)
	defer func() {
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	}()
	stack = append(stack, data)
	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(int); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]interface{}:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						return v
					}
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			case gen.Object:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						return v
					}
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			default:
				// TBD try reflection
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 <= i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case map[string]interface{}, []interface{}, gen.Object, gen.Array:
					stack = append(stack, v)
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 <= i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
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
			}
		case Descent:
			di, _ := stack[len(stack)-1].(int)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
					if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
					if top {
						return prev
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == len(x)-1 { // last one
				return data
			}
			stack = append(stack, data)
		case At, Bracket:
			if fi == len(x)-1 { // last one
				return prev
			}
			stack = append(stack, prev)
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					var has bool
					switch tv := prev.(type) {
					case map[string]interface{}:
						if v, has = tv[string(tu)]; has {
							if fi == len(x)-1 { // last one
								return v
							}
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					case gen.Object:
						if v, has = tv[string(tu)]; has {
							if fi == len(x)-1 { // last one
								return v
							}
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						// TBD try reflection
						continue
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						if fi == len(x)-1 { // last one
							return v
						}
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						if fi == len(x)-1 { // last one
							return v
						}
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					default:
						// TBD reflection
						continue
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
				if fi == len(x)-1 { // last one
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
				if fi == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case gen.Object, gen.Array:
					stack = append(stack, v)
				}
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
		}
		if fi < len(x)-1 {
			if _, ok := stack[len(stack)-1].(int); !ok {
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		}
	}
	return nil
}
