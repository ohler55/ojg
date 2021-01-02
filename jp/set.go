// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

type delFlagType struct{}

var delFlag = &delFlagType{}

// Del removes matching nodes.
func (x Expr) Del(n interface{}) error {
	return x.Set(n, delFlag)
}

// Del removes at most one node.
func (x Expr) DelOne(n interface{}) error {
	return x.SetOne(n, delFlag)
}

// Set all matching child node values. An error is returned if it is not
// possible. If the path to the child does not exist array and map elements
// are added.
func (x Expr) Set(data, value interface{}) error {
	fun := "set"
	if value == delFlag {
		fun = "delete"
	}
	if len(x) == 0 {
		return fmt.Errorf("can not %s with an empty expression", fun)
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter:
		ta := strings.Split(fmt.Sprintf("%T", x[len(x)-1]), ".")
		return fmt.Errorf("can not %s with an expression ending with a %s", fun, ta[len(ta)-1])
	}
	var v interface{}
	var nv gen.Node
	_, isNode := data.(gen.Node)
	nodeValue, ok := value.(gen.Node)
	if isNode && !ok {
		if value != nil {
			if v = alt.Generify(value); v == nil {
				return fmt.Errorf("can not %s a %T in a %T", fun, value, data)
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
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
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
				if int(fi) == len(x)-1 { // last one
					if value != delFlag {
						x.reflectSetChild(tv, string(tf), value)
					}
				} else if v, has = x.reflectGetChild(tv, string(tf)); has {
					switch v.(type) {
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				}
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
						case bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
							nil, gen.Bool, gen.Int, gen.Float, gen.String:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							default:
								return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
							}
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
				var has bool
				if int(fi) == len(x)-1 { // last one
					if value != delFlag {
						x.reflectSetNth(tv, i, value)
					}
				} else if v, has = x.reflectGetNth(tv, i); has {
					switch v.(type) {
					case bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
						nil, gen.Bool, gen.Int, gen.Float, gen.String:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
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
						}
					} else {
						for k = range tv {
							tv[k] = value
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
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
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
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
				if int(fi) != len(x)-1 {
					for _, v := range x.reflectGetWild(tv) {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						}
					}
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
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
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
						var has bool
						if v, has = x.reflectGetChild(tv, string(tu)); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
						}
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
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
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
						var has bool
						if v, has = x.reflectGetNth(tv, i); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
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
				if int(fi) != len(x)-1 {
					for _, v := range x.reflectGetSlice(tv, start, end, step) {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case *Filter:
			stack, _ = tf.Eval(stack, prev).([]interface{})
		case Root:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not %s the root", fun)
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not %s an empty expression", fun)
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
	fun := "set"
	if value == delFlag {
		fun = "delete"
	}
	if len(x) == 0 {
		return fmt.Errorf("can not %s with an empty expression", fun)
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter:
		ta := strings.Split(fmt.Sprintf("%T", x[len(x)-1]), ".")
		return fmt.Errorf("can not %s with an expression ending with a %s", fun, ta[len(ta)-1])
	}
	var v interface{}
	var nv gen.Node
	_, isNode := data.(gen.Node)
	nodeValue, ok := value.(gen.Node)
	if isNode && !ok {
		if value != nil {
			if v = alt.Generify(value); v == nil {
				return fmt.Errorf("can not %s a %T in a %T", fun, value, data)
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
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
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
				if int(fi) == len(x)-1 { // last one
					if value != delFlag {
						x.reflectSetChild(tv, string(tf), value)
					}
					return nil
				}
				if v, has = x.reflectGetChild(tv, string(tf)); has {
					switch v.(type) {
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
				}
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
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							default:
								return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
							}
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
			default:
				if int(fi) == len(x)-1 { // last one
					if value != delFlag {
						x.reflectSetNth(tv, i, value)
					}
					return nil
				}
				var has bool
				if v, has = x.reflectGetNth(tv, i); has {
					switch v.(type) {
					case bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
						nil, gen.Bool, gen.Int, gen.Float, gen.String:
						return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
					case map[string]interface{}, []interface{}, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						switch reflect.TypeOf(v).Kind() {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							return fmt.Errorf("can not follow a %T at '%s'", v, x[:fi+1])
						}
					}
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
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
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
					for i := len(tv) - 1; 0 <= i; i-- {
						v = tv[i]
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
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
					for i := len(tv) - 1; 0 <= i; i-- {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				if int(fi) != len(x)-1 {
					for _, v := range x.reflectGetWild(tv) {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(fragIndex)
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				switch tv := prev.(type) {
				case nil:
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for i := len(tv) - 1; 0 <= i; i-- {
						v = tv[i]
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, fragIndex(di|descentFlag))
					for i := len(tv) - 1; 0 <= i; i-- {
						v = tv[i]
						switch v.(type) {
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fragIndex(fi|descentChildFlag))
						}
					}
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
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
						}
					case gen.Object:
						if v, has = tv[string(tu)]; has {
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						var has bool
						if v, has = x.reflectGetChild(tv, string(tu)); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
						}
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
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
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
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					default:
						var has bool
						if v, has = x.reflectGetNth(tv, i); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]interface{}, []interface{}, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								switch reflect.TypeOf(v).Kind() {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
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
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := end; i <= start; i -= step {
						v = tv[i]
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
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
				if int(fi) != len(x)-1 {
					for _, v := range x.reflectGetSlice(tv, start, end, step) {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32, int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]interface{}, []interface{}, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							switch reflect.TypeOf(v).Kind() {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case *Filter:
			stack, _ = tf.Eval(stack, prev).([]interface{})
		case Root:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not %s the root", fun)
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				return fmt.Errorf("can not %s an empty expression", fun)
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

func (x Expr) reflectSetChild(data interface{}, key string, v interface{}) {
	rd := reflect.ValueOf(data)
	rt := rd.Type()
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rd = rd.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return
	}
	rv := rd.FieldByNameFunc(func(k string) bool { return strings.EqualFold(k, key) })
	vv := reflect.ValueOf(v)
	vt := vv.Type()
	if rv.CanSet() && vt.AssignableTo(rv.Type()) {
		rv.Set(vv)
	}
}

func (x Expr) reflectSetNth(data interface{}, i int, v interface{}) {
	rd := reflect.ValueOf(data)
	rt := rd.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		size := rd.Len()
		if i < 0 {
			i = size + i
		}
		if 0 <= i && i < size {
			rv := rd.Index(i)
			vv := reflect.ValueOf(v)
			vt := vv.Type()
			if rv.CanSet() && vt.AssignableTo(rv.Type()) {
				rv.Set(vv)
			}
		}
	}
}
