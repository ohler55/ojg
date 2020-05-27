// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

type delFlagType struct{}

var delFlag = &delFlagType{}

// Del removes matching nodes.
func (x Expr) Del(n interface{}) {
	_ = x.Set(n, delFlag)
}

// Del removes at most one node.
func (x Expr) DelOne(n interface{}) {
	_ = x.SetOne(n, delFlag)
}

// Set all matching child node values. An error is returned if it is not
// possible. If the path to the child does not exist array and map elements
// are added.
func (x Expr) Set(data, value interface{}) error {
	if len(x) == 0 {
		return fmt.Errorf("can not set with an empty expression")
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter:
		return fmt.Errorf("can not set with an expression ending with a %T", x[len(x)-1])
	}
	_, isNode := value.(Node)
	nodeValue, ok := value.(Node)
	if isNode && !ok {
		return fmt.Errorf("can not set a %T as an oj.Node in a %T", value, data)
	}

	var v interface{}
	var prev interface{}
	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) {
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
				if fi == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = value
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						tv[string(tf)] = map[string]interface{}{}
					default:
						return fmt.Errorf("can not deduce what element to add at %s", x[:fi+1])
					}
				}
			case Object:
				if fi == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = nodeValue
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case Object, Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						tv[string(tf)] = Object{}
					default:
						return fmt.Errorf("can not deduce what element to add at %s", x[:fi+1])
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
				if 0 <= i && i < len(tv) {
					if fi == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at %s", x[:fi+1])
				}
			case Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if fi == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at %s", x[:fi+1])
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				var k string
				if fi == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
						}
					} else {
						for k, v = range tv {
							tv[k] = value
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case []interface{}:
				if fi == len(x)-1 { // last one
					for i := range tv {
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Object:
				var k string
				if fi == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
						}
					} else {
						for k, v = range tv {
							tv[k] = nodeValue
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
				if fi == len(x)-1 { // last one
					for i := range tv {
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
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
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				}
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
							switch v.(type) {
							case map[string]interface{}, []interface{}, Object, Array:
								stack = append(stack, v)
							}
						}
					case Object:
						if v, has = tv[string(tu)]; has {
							switch v.(type) {
							case map[string]interface{}, []interface{}, Object, Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					case Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
		case Root:
			if fi == len(x)-1 { // last one
				return fmt.Errorf("can not set the root")
			}
			stack = append(stack, data)
		case At, Bracket:
			if fi == len(x)-1 { // last one
				return fmt.Errorf("can not set an empty expression")
			}
			stack = append(stack, prev)
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
	return nil
}

// Set a child node value. An error is returned if it is not possible. If the
// path to the child does not exist array and map elements are added.
func (x Expr) SetOne(data, value interface{}) error {
	if len(x) == 0 {
		return fmt.Errorf("can not set with an empty expression")
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter:
		return fmt.Errorf("can not set with an expression ending with a %T", x[len(x)-1])
	}
	_, isNode := value.(Node)
	nodeValue, ok := value.(Node)
	if isNode && !ok {
		return fmt.Errorf("can not set a %T as an oj.Node in a %T", value, data)
	}

	var v interface{}
	var prev interface{}
	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) {
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
				if fi == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = value
					}
					return nil
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						tv[string(tf)] = map[string]interface{}{}
					default:
						return fmt.Errorf("can not deduce what element to add at %s", x[:fi+1])
					}
				}
			case Object:
				if fi == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = nodeValue
					}
					return nil
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case Object, Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						tv[string(tf)] = Object{}
					default:
						return fmt.Errorf("can not deduce what element to add at %s", x[:fi+1])
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
				if 0 <= i && i < len(tv) {
					if fi == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
						return nil
					} else {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at %s", x[:fi+1])
				}
			case Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if fi == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
						return nil
					} else {
						v = tv[i]
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at %s", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at %s", x[:fi+1])
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				var k string
				if fi == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
							return nil
						}
					} else {
						for k, v = range tv {
							tv[k] = value
							return nil
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case []interface{}:
				if fi == len(x)-1 { // last one
					for i := range tv {
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
						return nil
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Object:
				var k string
				if fi == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
							return nil
						}
					} else {
						for k, v = range tv {
							tv[k] = nodeValue
							return nil
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
				if fi == len(x)-1 { // last one
					for i := range tv {
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
						return nil
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
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
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				}
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
							switch v.(type) {
							case map[string]interface{}, []interface{}, Object, Array:
								stack = append(stack, v)
							}
						}
					case Object:
						if v, has = tv[string(tu)]; has {
							switch v.(type) {
							case map[string]interface{}, []interface{}, Object, Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					case Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v interface{}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
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
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end < i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case *Filter:
			stack = tf.Eval(stack, prev)
		case Root:
			if fi == len(x)-1 { // last one
				return fmt.Errorf("can not set the root")
			}
			stack = append(stack, data)
		case At, Bracket:
			if fi == len(x)-1 { // last one
				return fmt.Errorf("can not set an empty expression")
			}
			stack = append(stack, prev)
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
	return nil
}
