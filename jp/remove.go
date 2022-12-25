// Copyright (c) 2022, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
)

// MustRemove removes matching nodes and panics on error. Removed slice
// elements are removed and the remaining elements are moveed to fill in the
// removed element. The slice is shortened.
func (x Expr) MustRemove(n any) any {
	return x.remove(n, math.MaxInt)
}

// Remove removes matching nodes. Removed slice elements are removed and the
// remaining elements are moveed to fill in the removed element. The slice is
// shortened.
func (x Expr) Remove(n any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	result = x.remove(n, math.MaxInt)

	return
}

// RemoveOne removes at most one node. Removed slice elements are removed and
// the remaining elements are moveed to fill in the removed element. The slice
// is shortened.
func (x Expr) RemoveOne(n any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	result = x.remove(n, 1)

	return
}

func (x Expr) remove(data any, max int) any {
	if len(x) == 0 {
		panic("can not remove with an empty expression")
	}
	switch x[len(x)-1].(type) {
	case Descent, Union, Slice, *Filter: // TBD filter is okay
		ta := strings.Split(fmt.Sprintf("%T", x[len(x)-1]), ".")
		panic(fmt.Sprintf("can not remove with an expression ending with a %s", ta[len(ta)-1]))
	}
	// TBD keep a parents stack
	// push on map or slice
	// pop when last on stack is removed and a map or slice

	var v any
	var prev any
	stack := make([]any, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) && 0 < max {
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			// TBD maybe pop from parents
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		// TBD maybe pop from parents
		stack[len(stack)-1] = nil
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]any:
				if int(fi) == len(x)-1 { // last one
					delete(tv, string(tf))
					max--
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32,
						int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
					case map[string]any, []any, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						kind := reflect.Invalid
						if rt := reflect.TypeOf(v); rt != nil {
							kind = rt.Kind()
						}
						switch kind {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
						}
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					delete(tv, string(tf))
					max--
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
					}
				}
			default:
				if int(fi) == len(x)-1 { // last one
					// TBD nothing to do?
				} else if v, has = x.reflectGetChild(tv, string(tf)); has {
					switch v.(type) {
					case nil, gen.Bool, gen.Int, gen.Float, gen.String,
						bool, string, float64, float32,
						int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
					case map[string]any, []any, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						kind := reflect.Invalid
						if rt := reflect.TypeOf(v); rt != nil {
							kind = rt.Kind()
						}
						switch kind {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
						}
					}
				}
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []any:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						tv[i] = nil // TBD remove from list
						max--
					} else {
						v = tv[i]
						switch v.(type) {
						case bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
							nil, gen.Bool, gen.Int, gen.Float, gen.String:
							panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							default:
								panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
							}
						}
					}
				} else {
					panic(fmt.Sprintf("can not follow out of bounds array index at '%s'", x[:fi+1]))
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						tv[i] = nil // TBD remove
						max--
					} else {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
						}
					}
				} else {
					panic(fmt.Sprintf("can not follow out of bounds array index at '%s'", x[:fi+1]))
				}
			default:
				var has bool
				if int(fi) == len(x)-1 { // last one
					// TBD
				} else if v, has = x.reflectGetNth(tv, i); has {
					switch v.(type) {
					case bool, string, float64, float32,
						int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
						nil, gen.Bool, gen.Int, gen.Float, gen.String:
						panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
					case map[string]any, []any, gen.Object, gen.Array:
						stack = append(stack, v)
					default:
						kind := reflect.Invalid
						if rt := reflect.TypeOf(v); rt != nil {
							kind = rt.Kind()
						}
						switch kind {
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
							stack = append(stack, v)
						default:
							panic(fmt.Sprintf("can not follow a %T at '%s'", v, x[:fi+1]))
						}
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]any:
				var k string
				if int(fi) == len(x)-1 { // last one
					for k = range tv {
						delete(tv, k)
						max--
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			case []any:
				if int(fi) == len(x)-1 { // last one
					for i := range tv {
						tv[i] = nil // TBD remove
						max--
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			case gen.Object:
				var k string
				if int(fi) == len(x)-1 { // last one
					for k = range tv {
						delete(tv, k)
						max--
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
						tv[i] = nil // TBD remove
						max--
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
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
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
				case map[string]any:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fi|descentChildFlag)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case []any:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fi|descentChildFlag)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fi|descentChildFlag)
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					for _, v = range tv {
						switch v.(type) {
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
							stack = append(stack, fi|descentChildFlag)
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
					case map[string]any:
						if v, has = tv[tu]; has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32,
								int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								kind := reflect.Invalid
								if rt := reflect.TypeOf(v); rt != nil {
									kind = rt.Kind()
								}
								switch kind {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
						}
					case gen.Object:
						if v, has = tv[tu]; has {
							switch v.(type) {
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						var has bool
						if v, has = x.reflectGetChild(tv, tu); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32,
								int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								kind := reflect.Invalid
								if rt := reflect.TypeOf(v); rt != nil {
									kind = rt.Kind()
								}
								switch kind {
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
									stack = append(stack, v)
								}
							}
						}
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []any:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32,
								int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								kind := reflect.Invalid
								if rt := reflect.TypeOf(v); rt != nil {
									kind = rt.Kind()
								}
								switch kind {
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
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					default:
						var has bool
						if v, has = x.reflectGetNth(tv, i); has {
							switch v.(type) {
							case nil, gen.Bool, gen.Int, gen.Float, gen.String,
								bool, string, float64, float32,
								int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							default:
								kind := reflect.Invalid
								if rt := reflect.TypeOf(v); rt != nil {
									kind = rt.Kind()
								}
								switch kind {
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
			case []any:
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
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]any, []any, gen.Object, gen.Array:
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
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
						switch v.(type) {
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				if int(fi) != len(x)-1 {
					for _, v := range x.reflectGetSlice(tv, start, end, step) {
						switch v.(type) {
						case nil, gen.Bool, gen.Int, gen.Float, gen.String,
							bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case *Filter:
			// TBD if last one then set or remove
			stack, _ = tf.Eval(stack, prev).([]any)
		case Root:
			if int(fi) == len(x)-1 { // last one
				panic("can not remove the root")
			}
			stack = append(stack, data)
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				panic("can not fun an empty expression")
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

	// TBD
	// use common remove but with counter for 1 or maxint or -1 which needs an extra check

	return data
}
