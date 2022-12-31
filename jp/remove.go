// Copyright (c) 2022, Peter Ohler, All rights reserved.

package jp

import (
	"fmt"
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
	return x.remove(n, false)
}

// MustRemoveOne removes matching nodes and panics on error expression error
// but silently makes no changes if there is no match for the
// expression. Removed slice elements are removed and the remaining elements
// are moveed to fill in the removed element. The slice is shortened.
func (x Expr) MustRemoveOne(n any) any {
	return x.remove(n, true)
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
	result = x.remove(n, false)

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
	result = x.remove(n, true)

	return
}

func (x Expr) remove(data any, one bool) any {
	if len(x) == 0 {
		panic("can not remove with an empty expression")
	}
	last := x[len(x)-1]
	for _, f := range x {
		if _, ok := f.(Descent); ok {
			ta := strings.Split(fmt.Sprintf("%T", f), ".")
			panic(fmt.Sprintf("can not remove with an expression where the second to last fragment is a %s",
				ta[len(ta)-1]))
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

done:
	for 1 < len(stack) {
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
				if int(fi) == len(wx)-1 { // last one
					if nv, changed := removeLast(last, tv[string(tf)], one); changed {
						tv[string(tf)] = nv
						if one && changed {
							break done
						}
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
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
				if int(fi) == len(wx)-1 { // last one
					if nv, changed := removeNodeLast(last, tv[string(tf)], one); changed {
						tv[string(tf)] = nv
						if one && changed {
							break done
						}
					}
				} else if v, has = tv[string(tf)]; has {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			default:
				if v, has = wx.reflectGetChild(tv, string(tf)); has {
					if int(fi) == len(wx)-1 { // last one
						if nv, changed := removeLast(last, v, one); changed {
							wx.reflectSetChild(tv, string(tf), nv)
							if one && changed {
								break done
							}
						}
					} else {
						switch v.(type) {
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
					if int(fi) == len(wx)-1 { // last one
						if nv, changed := removeLast(last, tv[i], one); changed {
							tv[i] = nv
							if one && changed {
								break done
							}
						}
					} else {
						v = tv[i]
						switch v.(type) {
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
					if int(fi) == len(wx)-1 { // last one
						if nv, changed := removeNodeLast(last, tv[i], one); changed {
							tv[i] = nv
							if one && changed {
								break done
							}
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
				if v, has = wx.reflectGetNth(tv, i); has {
					if int(fi) == len(wx)-1 { // last one
						if nv, changed := removeLast(last, v, one); changed {
							wx.reflectSetNth(tv, i, nv)
							if one && changed {
								break done
							}
						}
					} else {
						switch v.(type) {
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
				if int(fi) == len(wx)-1 { // last one
					for k = range tv {
						if nv, changed := removeLast(last, tv[k], one); changed {
							tv[k] = nv
							if one && changed {
								break done
							}
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
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
				if int(fi) == len(wx)-1 { // last one
					for i := range tv {
						if nv, changed := removeLast(last, tv[i], one); changed {
							tv[i] = nv
							if one && changed {
								break done
							}
						}
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
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
				if int(fi) == len(wx)-1 { // last one
					for k = range tv {
						if nv, changed := removeNodeLast(last, tv[k], one); changed {
							tv[k] = nv
							if one && changed {
								break done
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
				if int(fi) == len(wx)-1 { // last one
					for i := range tv {
						if nv, changed := removeNodeLast(last, tv[i], one); changed {
							tv[i] = nv
							if one && changed {
								break done
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
				if int(fi) == len(wx)-1 { // last one
					rv := reflect.ValueOf(tv)
					switch rv.Kind() {
					case reflect.Slice:
						cnt := rv.Len()
						for i := 0; i < cnt; i++ {
							iv := rv.Index(i)
							if nv, changed := removeLast(last, iv.Interface(), one); changed {
								iv.Set(reflect.ValueOf(nv))
								if one && changed {
									break done
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
							if nv, changed := removeLast(last, ev.Interface(), one); changed {
								rv.SetMapIndex(k, reflect.ValueOf(nv))
								if one && changed {
									break done
								}
							}
						}
					}
				} else {
					for _, v := range wx.reflectGetWild(tv) {
						switch v.(type) {
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
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					var has bool
					switch tv := prev.(type) {
					case map[string]any:
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeLast(last, tv[tu], one); changed {
								tv[tu] = nv
								if one && changed {
									break done
								}
							}
						} else if v, has = tv[tu]; has {
							switch v.(type) {
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
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeNodeLast(last, tv[tu], one); changed {
								tv[tu] = nv
								if one && changed {
									break done
								}
							}
						} else if v, has = tv[tu]; has {
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						var has bool
						if v, has = wx.reflectGetChild(tv, tu); has {
							if int(fi) == len(wx)-1 { // last one
								if nv, changed := removeLast(last, v, one); changed {
									wx.reflectSetChild(tv, tu, nv)
									if one && changed {
										break done
									}
								}
							} else {
								switch v.(type) {
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
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case []any:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
							if int(fi) == len(wx)-1 { // last one
								if nv, changed := removeLast(last, v, one); changed {
									tv[i] = nv
									if one && changed {
										break done
									}
								}
							} else {
								switch v.(type) {
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
							if int(fi) == len(wx)-1 { // last one
								if nv, changed := removeNodeLast(last, tv[i], one); changed {
									tv[i] = nv
									if one && changed {
										break done
									}
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
						if int(fi) == len(wx)-1 { // last one
							rv := reflect.ValueOf(tv)
							if rv.Kind() == reflect.Slice {
								cnt := rv.Len()
								if i < 0 {
									i = cnt + i
								}
								if 0 <= i && i < cnt {
									iv := rv.Index(i)
									if nv, changed := removeLast(last, iv.Interface(), one); changed {
										iv.Set(reflect.ValueOf(nv))
										if one && changed {
											break done
										}
									}
								}
							}
						} else {
							var has bool
							if v, has = wx.reflectGetNth(tv, i); has {
								switch v.(type) {
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
				if len(tv) <= end {
					end = len(tv) - 1
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) <= end || step == 0 {
					continue
				}
				if 0 < step {
					for i := start; i <= end; i += step {
						v = tv[i]
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeLast(last, v, one); changed {
								tv[i] = nv
								if one && changed {
									break done
								}
							}
						} else {
							switch v.(type) {
							case map[string]any, []any, gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeLast(last, v, one); changed {
								tv[i] = nv
								if one && changed {
									break done
								}
							}
						} else {
							switch v.(type) {
							case map[string]any, []any, gen.Object, gen.Array:
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
				if len(tv) <= end {
					end = len(tv) - 1
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) <= end || step == 0 {
					continue
				}
				if 0 < step {
					for i := start; i <= end; i += step {
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeNodeLast(last, tv[i], one); changed {
								tv[i] = nv
								if one && changed {
									break done
								}
							}
						} else {
							v = tv[i]
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := start; end <= i; i += step {
						if int(fi) == len(wx)-1 { // last one
							if nv, changed := removeNodeLast(last, tv[i], one); changed {
								tv[i] = nv
								if one && changed {
									break done
								}
							}
						} else {
							v = tv[i]
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			default:
				if int(fi) == len(wx)-1 {
					rv := reflect.ValueOf(tv)
					if rv.Kind() == reflect.Slice {
						cnt := rv.Len()
						if start < 0 {
							start = cnt + start
						}
						if end < 0 {
							end = cnt + end
						}
						if cnt <= end {
							end = cnt - 1
						}
						if start < 0 || end < 0 || cnt <= start || cnt <= end || step == 0 {
							continue
						}
						if 0 < step {
							for i := start; i <= end; i += step {
								iv := rv.Index(i)
								if nv, changed := removeLast(last, iv.Interface(), one); changed {
									iv.Set(reflect.ValueOf(nv))
									if one && changed {
										break done
									}
								}
							}
						} else {
							for i := start; end <= i; i += step {
								iv := rv.Index(i)
								if nv, changed := removeLast(last, iv.Interface(), one); changed {
									iv.Set(reflect.ValueOf(nv))
									if one && changed {
										break done
									}
								}
							}
						}
					}
				} else {
					for _, v := range wx.reflectGetSlice(tv, start, end, step) {
						switch v.(type) {
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

			// TBD prev is a slice
			// iterate and match each, remember i

			fmt.Printf("*** filter - %v\n", prev)

			stack, _ = tf.Eval(stack, prev).([]any)
		case Root:
			if int(fi) == len(wx)-1 { // last one
				if nv, changed := removeLast(last, data, one); changed {
					wrap[0] = nv
					if one && changed {
						break done
					}
				}
			} else {
				stack = append(stack, data)
			}
		case At, Bracket:
			if int(fi) == len(wx)-1 { // last one
				if nv, changed := removeLast(last, data, one); changed {
					wrap[0] = nv
					if one && changed {
						break done
					}
				}
			}
			stack = append(stack, prev)
		}
		if int(fi) < len(wx)-1 {
			if _, ok := stack[len(stack)-1].(fragIndex); !ok {
				fi++
				f = wx[fi]
				stack = append(stack, fi)
			}
		}
	}
	return wrap[0]
}

func removeLast(f Frag, value any, one bool) (out any, changed bool) {
	out = value
	switch tf := f.(type) {
	case Child:
		key := string(tf)
		switch tv := value.(type) {
		case map[string]any:
			if _, changed = tv[key]; changed {
				delete(tv, key)
			}
		case gen.Object:
			if _, changed = tv[key]; changed {
				delete(tv, string(tf))
			}
		default:
			if rt := reflect.TypeOf(value); rt != nil {
				// Can't remove a field from a struct so only a map can be modified.
				if rt.Kind() == reflect.Map {
					rv := reflect.ValueOf(value)
					rk := reflect.ValueOf(key)
					if rv.MapIndex(rk).IsValid() {
						rv.SetMapIndex(rk, reflect.Value{})
						changed = true
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
				out = append(tv[:i], tv[i+1:]...)
				changed = true
			}
		case gen.Array:
			if i < 0 {
				i = len(tv) + i
			}
			if 0 <= i && i < len(tv) {
				out = append(tv[:i], tv[i+1:]...)
				changed = true
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
							out = nv.Interface()
							changed = true
						}
					}
				}
			}
		}
	case Wildcard:
		switch tv := value.(type) {
		case []any:
			if 0 < len(tv) {
				changed = true
				if one {
					out = tv[1:]
				} else {
					out = []any{}
				}
			}
		case map[string]any:
			if 0 < len(tv) {
				changed = true
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				if one {
					delete(tv, keys[0])
				} else {
					for _, k := range keys {
						delete(tv, k)
					}
				}
			}
		case gen.Array:
			if 0 < len(tv) {
				changed = true
				if one {
					out = tv[1:]
				} else {
					out = gen.Array{}
				}
			}
		case gen.Object:
			if 0 < len(tv) {
				changed = true
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				if one {
					delete(tv, keys[0])
				} else {
					for _, k := range keys {
						delete(tv, k)
					}
				}
			}
		default:
			// TBD reflect
		}
	case Union:
		// TBD
	case Slice:
		// TBD
	case *Filter:
		// TBD find indices then remove those until max
	}
	return
}

func removeNodeLast(f Frag, value gen.Node, one bool) (out gen.Node, changed bool) {
	out = value
	switch tf := f.(type) {
	case Child:
		if tv, ok := value.(gen.Object); ok {
			if _, changed = tv[string(tf)]; changed {
				delete(tv, string(tf))
			}
		}
	case Nth:
		i := int(tf)
		if tv, ok := value.(gen.Array); ok {
			if i < 0 {
				i = len(tv) + i
			}
			if 0 <= i && i < len(tv) {
				out = append(tv[:i], tv[i+1:]...)
				changed = true
			}
		}
	case Wildcard:
		switch tv := value.(type) {
		case gen.Array:
			if 0 < len(tv) {
				changed = true
				if one {
					out = tv[1:]
				} else {
					out = gen.Array{}
				}
			}
		case gen.Object:
			if 0 < len(tv) {
				changed = true
				keys := make([]string, 0, len(tv))
				for k := range tv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				if one {
					delete(tv, keys[0])
				} else {
					for _, k := range keys {
						delete(tv, k)
					}
				}
			}
		}
	case Union:
		// TBD
	case Slice:
		// TBD
	case *Filter:
		// TBD find indices then remove those until max
	}
	return
}
