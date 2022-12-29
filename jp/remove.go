// Copyright (c) 2022, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/gen"
)

// MustRemove removes matching nodes and panics on error expression error but
// silently makes no changes if there is no match for the expression. Removed
// slice elements are removed and the remaining elements are moveed to fill in
// the removed element. The slice is shortened.
func (x Expr) MustRemove(n any) any {
	return x.remove(n, math.MaxInt)
}

// MustRemoveOne removes matching nodes and panics on error expression error
// but silently makes no changes if there is no match for the
// expression. Removed slice elements are removed and the remaining elements
// are moveed to fill in the removed element. The slice is shortened.
func (x Expr) MustRemoveOne(n any) any {
	return x.remove(n, 1)
}

// Remove removes matching nodes. An error is returned for an expression error
// but silently makes no changes if there is no match for the
// expression. Removed slice elements are removed and the remaining elements
// are moveed to fill in the removed element. The slice is shortened.
func (x Expr) Remove(n any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	result = x.remove(n, math.MaxInt)

	return
}

// RemoveOne removes at most one node. An error is returned for an expression
// error but silently makes no changes if there is no match for the
// expression. Removed slice elements are removed and the remaining elements
// are moveed to fill in the removed element. The slice is shortened.
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
	last := x[len(x)-1]
	for i, f := range x {
		switch f.(type) {
		case Descent:
			ta := strings.Split(fmt.Sprintf("%T", f), ".")
			panic(fmt.Sprintf("can not remove with an expression where the second to last fragment is a %s", ta[len(ta)-1]))
		case Union, Slice:
			if i == len(x)-1 {
				ta := strings.Split(fmt.Sprintf("%T", f), ".")
				panic(fmt.Sprintf("can not remove with an expression ending with a %s", ta[len(ta)-1]))
			}
		}
	}
	wx := make(Expr, len(x))
	copy(wx[1:], x)
	wx[0] = Nth(0)
	wrap := []any{data}
	var (
		v    any
		prev any
	)
	stack := make([]any, 0, 64)
	stack = append(stack, wrap)

	f := wx[0]
	fi := fragIndex(0) // frag index
	stack = append(stack, fi)

	for 1 < len(stack) && 0 < max {
		prev = stack[len(stack)-2]
		if ii, up := prev.(fragIndex); up {
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = wx[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack[len(stack)-1] = nil
		stack = stack[:len(stack)-1]

		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]any:
				if int(fi) == len(x)-1 { // last one
					if nv, mx := removeLast(last, tv[string(tf)], max); max != mx {
						tv[string(tf)] = nv
						max = mx
					}
				} else if v, has = tv[string(tf)]; has {
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
						case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
							stack = append(stack, v)
						}
					}
				}
			case gen.Object:
				if int(fi) == len(x)-1 { // last one
					if nv, mx := removeNodeLast(last, tv[string(tf)], max); max != mx {
						tv[string(tf)] = nv
						max = mx
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			default:
				if v, has = x.reflectGetChild(tv, string(tf)); has {
					if int(fi) == len(x)-1 { // last one
						if nv, mx := removeLast(last, v, max); max != mx {
							x.reflectSetChild(tv, string(tf), nv)
							max = mx
						}
					} else {
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
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
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
						if nv, mx := removeLast(last, tv[i], max); max != mx {
							tv[i] = nv
							max = mx
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
							nil, gen.Bool, gen.Int, gen.Float, gen.String:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
						}
					}
				}
			case gen.Array:
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					if int(fi) == len(x)-1 { // last one
						if nv, mx := removeNodeLast(last, tv[i], max); max != mx {
							tv[i] = nv
							max = mx
						}
					} else {
						v = tv[i]
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			default:
				var has bool
				if v, has = x.reflectGetNth(tv, i); has {
					if int(fi) == len(x)-1 { // last one
						if nv, mx := removeLast(last, v, max); max != mx {
							x.reflectSetNth(tv, i, nv)
							max = mx
						}
					} else {
						switch v.(type) {
						case bool, string, float64, float32,
							int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
							nil, gen.Bool, gen.Int, gen.Float, gen.String:
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						default:
							kind := reflect.Invalid
							if rt := reflect.TypeOf(v); rt != nil {
								kind = rt.Kind()
							}
							switch kind {
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
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
						if nv, mx := removeLast(last, tv[k], max); max != mx {
							tv[k] = nv
							max = mx
							if max <= 0 {
								break
							}
						}
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
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
						}
					}
				}
			case []any:
				if int(fi) == len(x)-1 { // last one
					for i := range tv {
						if nv, mx := removeLast(last, tv[i], max); max != mx {
							tv[i] = nv
							max = mx
							if max <= 0 {
								break
							}
						}
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
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
						}
					}
				}
			case gen.Object:
				var k string
				if int(fi) == len(x)-1 { // last one
					for k = range tv {
						if nv, mx := removeNodeLast(last, tv[k], max); max != mx {
							tv[k] = nv
							max = mx
							if max <= 0 {
								break
							}
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
						if nv, mx := removeNodeLast(last, tv[i], max); max != mx {
							tv[i] = nv
							max = mx
							if max <= 0 {
								break
							}
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
				if int(fi) == len(x)-1 { // last one
					rv := reflect.ValueOf(tv)
					switch rv.Kind() {
					case reflect.Slice:
						cnt := rv.Len()
						for i := 0; i < cnt; i++ {
							iv := rv.Index(i)
							if nv, mx := removeLast(last, iv.Interface(), max); max != mx {
								iv.Set(reflect.ValueOf(nv))
								max = mx
								if max <= 0 {
									break
								}
							}
						}
					case reflect.Map:
						keys := rv.MapKeys()
						sort.Slice(keys, func(i, j int) bool {
							return strings.Compare(keys[i].String(), keys[j].String()) < 0
						})
						for _, k := range keys {
							ev := rv.MapIndex(k)
							if nv, mx := removeLast(last, ev.Interface(), max); max != mx {
								rv.SetMapIndex(k, reflect.ValueOf(nv))
								max = mx
							}
							if max <= 0 {
								break
							}
						}
					}
				} else {
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
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case Union:

			// TBD handle

			for _, u := range tf {
				switch tu := u.(type) {
				case string:

					// TBD if last...

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
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
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
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
									stack = append(stack, v)
								}
							}
						}
					}
				case int64:

					// TBD if last...

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
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
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
								case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
									stack = append(stack, v)
								}
							}
						}
					}
				}
			}
		case Slice:

			// TBD handle

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

						// TBD handle last

						v = tv[i]
						switch v.(type) {
						case map[string]any, []any, gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				} else {
					for i := start; end <= i; i += step {

						// TBD handle last

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
							case reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Array, reflect.Map:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case *Filter:
			// TBD handle
			stack, _ = tf.Eval(stack, prev).([]any)
		case Root:
			if int(fi) == len(x)-1 { // last one
				if nv, mx := removeLast(last, data, max); max != mx {
					wrap[0] = nv
					max = mx
				}
			} else {
				stack = append(stack, data)
			}
		case At, Bracket:
			if int(fi) == len(x)-1 { // last one
				if nv, mx := removeLast(last, data, max); max != mx {
					wrap[0] = nv
					max = mx
				}
			}
			stack = append(stack, prev)
		}
		if int(fi) < len(x)-1 {
			if _, ok := stack[len(stack)-1].(fragIndex); !ok {
				fi++
				f = wx[fi]
				stack = append(stack, fi)
			}
		}
	}
	return wrap[0]
}

func removeLast(f Frag, value any, max int) (any, int) {
	switch tf := f.(type) {
	case Child:
		key := string(tf)
		switch tv := value.(type) {
		case map[string]any:
			if _, has := tv[key]; has {
				delete(tv, key)
				max--
			}
		case gen.Object:
			if _, has := tv[key]; has {
				delete(tv, string(tf))
				max--
			}
		default:
			if rt := reflect.TypeOf(value); rt != nil {
				// Can't remove a field from a struct so only a map can be modified.
				if rt.Kind() == reflect.Map {
					rv := reflect.ValueOf(value)
					rk := reflect.ValueOf(key)
					if rv.MapIndex(rk).IsValid() {
						rv.SetMapIndex(rk, reflect.Value{})
						max--
					}
				}
			}
		}
	case Nth:
		i := int(tf)
		switch tv := value.(type) {
		case []any:
			if i < 0 {
				i = len(tv) + i
			}
			if 0 <= i && i < len(tv) {
				value = append(tv[:i], tv[i+1:]...)
				max--
			}
		case gen.Array:
			if i < 0 {
				i = len(tv) + i
			}
			if 0 <= i && i < len(tv) {
				value = append(tv[:i], tv[i+1:]...)
				max--
			}
		default:
			if rt := reflect.TypeOf(value); rt != nil {
				if rt.Kind() == reflect.Slice {
					rv := reflect.ValueOf(value)
					cnt := rv.Len()
					if 0 < cnt {
						if i < 0 {
							i = cnt + i
						}
						if 0 <= i && i < cnt {
							nv := reflect.MakeSlice(rt, cnt-1, cnt-1)
							for j := 0; j < i; j++ {
								nv.Index(j).Set(rv.Index(j))
							}
							for j := i + 1; j < cnt; j++ {
								nv.Index(j - 1).Set(rv.Index(j))
							}
							value = nv.Interface()
							max--
						}
					}
				}
			}
		}
	case Wildcard:
		switch tv := value.(type) {
		case []any:
			if len(tv) <= max {
				max -= len(tv)
				value = []any{}
			} else {
				for ; 0 < max; max-- {
					tv = tv[1:]
				}
				value = tv
			}
		case map[string]any:
			if len(tv) <= max {
				max -= len(tv)
				value = map[string]any{}
			} else {
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					delete(tv, k)
					max--
					if max <= 0 {
						break
					}
				}
			}
		case gen.Array:
			if len(tv) <= max {
				max -= len(tv)
				value = gen.Array{}
			} else {
				for ; 0 < max; max-- {
					tv = tv[1:]
				}
				value = tv
			}
		case gen.Object:
			if len(tv) <= max {
				max -= len(tv)
				value = gen.Object{}
			} else {
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					delete(tv, k)
					max--
					if max <= 0 {
						break
					}
				}
			}
		default:
			// TBD reflect
		}
	case *Filter:
		// TBD find indices then remove those until max
	}
	return value, max
}

func removeNodeLast(f Frag, value gen.Node, max int) (gen.Node, int) {
	switch tf := f.(type) {
	case Child:
		if tv, ok := value.(gen.Object); ok {
			if _, has := tv[string(tf)]; has {
				delete(tv, string(tf))
				max--
			}
		}
	case Nth:
		i := int(tf)
		if tv, ok := value.(gen.Array); ok {
			if i < 0 {
				i = len(tv) + i
			}
			if 0 <= i && i < len(tv) {
				value = append(tv[:i], tv[i+1:]...)
				max--
			}
		}
	case Wildcard:
		switch tv := value.(type) {
		case gen.Array:
			if len(tv) <= max {
				max -= len(tv)
				value = gen.Array{}
			} else {
				for ; 0 < max; max-- {
					tv = tv[1:]
				}
				value = tv
			}
		case gen.Object:
			if len(tv) <= max {
				max -= len(tv)
				value = gen.Object{}
			} else {
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					delete(tv, k)
					max--
					if max <= 0 {
						break
					}
				}
			}
		}
	case *Filter:
		// TBD find indices then remove those until max
	}
	return value, max
}
