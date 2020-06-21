// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

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
		ta := strings.Split(fmt.Sprintf("%T", x[len(x)-1]), ".")
		return fmt.Errorf("can not set with an expression ending with a %s", ta[len(ta)-1])
	}
	var v interface{}
	var nv gen.Node
	_, isNode := data.(gen.Node)
	nodeValue, ok := value.(gen.Node)
	if isNode && !ok {
		if value != nil {
			if v = alt.Generify(value); v == nil {
				return fmt.Errorf("can not set a %T in a %T", value, data)
			}
			nodeValue, _ = v.(gen.Node)
		}
	}
	var prev interface{}
	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) {
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack[len(stack)-1] = nil
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case nil:
			case map[string]interface{}:
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = value
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						v = map[string]interface{}{}
						tv[string(tf)] = v
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not deduce what element to add at '%s'", x[:fi+1])
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = nodeValue
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						nv = gen.Object{}
						tv[string(tf)] = nv
						stack = append(stack, nv)
					default:
						return fmt.Errorf("can not deduce what element to add at '%s'", x[:fi+1])
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case nil:
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at '%s'", x[:fi+1])
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at '%s'", x[:fi+1])
				}
			default:
				// TBD reflection
			}
		case Wildcard:
			switch tv := prev.(type) {
			case nil:
			case map[string]interface{}:
				var k string
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
						}
					} else {
						for k = range tv {
							tv[k] = value
						}
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
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				var k string
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
						}
					} else {
						for k = range tv {
							tv[k] = nodeValue
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if int(fi) == len(x)-1 { // last one
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
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
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
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
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
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
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
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				default:
					// TBD reflection
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
					case nil:
					case map[string]interface{}:
						if v, has = tv[string(tu)]; has {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					case gen.Object:
						if v, has = tv[string(tu)]; has {
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
					case nil:
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
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
			case nil:
			case []interface{}:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) <= end || step == 0 {
					continue
				}
				if 0 < step {
					for i := start; i <= end; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) <= end || step == 0 {
					continue
				}
				if 0 < step {
					for i := start; i <= end; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
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
		case *Filter:
			stack, _ = tf.Eval(stack, prev).([]interface{})
		case Root:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not set the root")
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not set an empty expression")
			}
			stack = append(stack, prev)
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

// Set a child node value. An error is returned if it is not possible. If the
// path to the child does not exist array and map elements are added.
func (x Expr) SetOne(data, value interface{}) error {
	if len(x) == 0 {
		return fmt.Errorf("can not set with an empty expression")
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter:
		ta := strings.Split(fmt.Sprintf("%T", x[len(x)-1]), ".")
		return fmt.Errorf("can not set with an expression ending with a %s", ta[len(ta)-1])
	}
	var v interface{}
	var nv gen.Node
	_, isNode := data.(gen.Node)
	nodeValue, ok := value.(gen.Node)
	if isNode && !ok {
		if value != nil {
			if v = alt.Generify(value); v == nil {
				return fmt.Errorf("can not set a %T in a %T", value, data)
			}
			nodeValue, _ = v.(gen.Node)
		}
	}
	var prev interface{}
	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) {
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack[len(stack)-1] = nil
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case nil:
			case map[string]interface{}:
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = value
					}
					return nil
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						v = map[string]interface{}{}
						tv[string(tf)] = v
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not deduce what element to add at '%s'", x[:fi+1])
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						delete(tv, string(tf))
					} else {
						tv[string(tf)] = nodeValue
					}
					return nil
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					}
				} else {
					switch x[fi+1].(type) {
					case Child:
						nv = gen.Object{}
						tv[string(tf)] = nv
						stack = append(stack, nv)
					default:
						return fmt.Errorf("can not deduce what element to add at '%s'", x[:fi+1])
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case nil:
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = value
						}
						return nil
					} else {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at '%s'", x[:fi+1])
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						if value == delFlag {
							tv[i] = nil
						} else {
							tv[i] = nodeValue
						}
						return nil
					} else {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				} else {
					return fmt.Errorf("can not follow out of bounds array index at '%s'", x[:fi+1])
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case nil:
			case map[string]interface{}:
				var k string
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
							return nil
						}
					} else {
						for k = range tv {
							tv[k] = value
							return nil
						}
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
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				var k string
				if int(fi) == len(x)-1 { // last one
					if value == delFlag {
						for k = range tv {
							delete(tv, k)
							return nil
						}
					} else {
						for k = range tv {
							tv[k] = nodeValue
							return nil
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if int(fi) == len(x)-1 { // last one
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
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case nil:
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
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
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				default:
					// TBD reflection
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
					case nil:
					case map[string]interface{}:
						if v, has = tv[string(tu)]; has {
							switch v.(type) {
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					case gen.Object:
						if v, has = tv[string(tu)]; has {
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
					case nil:
					case []interface{}:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
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
			case nil:
			case []interface{}:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end
				}
				if start < 0 || end < 0 || len(tv) <= start || step == 0 {
					continue
				}
				if len(tv) <= end {
					end = len(tv) - 1
				}
				end = start + ((end - start) / step * step)
				if 0 < step {
					for i := end; start <= i; i -= step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := end; i <= start; i -= step {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end
				}
				if start < 0 || end < 0 || len(tv) <= start || step == 0 {
					continue
				}
				if len(tv) <= end {
					end = len(tv) - 1
				}
				end = start + ((end - start) / step * step)
				if 0 < step {
					for i := end; start <= i; i -= step {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := end; i <= start; i -= step {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				// TBD try reflection
				continue
			}
		case *Filter:
			stack, _ = tf.Eval(stack, prev).([]interface{})
		case Root:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not set the root")
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not set an empty expression")
			}
			stack = append(stack, prev)
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
