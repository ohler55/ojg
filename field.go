// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"strconv"
	"unsafe"
)

// Field hold information about a struct field.
type Field struct {
	Type      reflect.Type
	Key       string
	Kind      reflect.Kind
	Elem      *Struct
	jkey      []byte
	index     []int
	offset    uintptr
	asString  bool
	omitEmpty bool
	direct    bool
}

func (fi *Field) Value(rv reflect.Value, omitNil bool, embedded bool) (v interface{}, omit bool) {
	var fv reflect.Value
	var ptr uintptr
	if fi.direct && !embedded {
		ptr = rv.UnsafeAddr() + fi.offset
	} else {
		fv = rv.FieldByIndex(fi.index)
		v = fv.Interface()
	}
	switch fi.Kind {
	case reflect.Bool:
		var b bool
		if ptr != 0 {
			b = *(*bool)(unsafe.Pointer(ptr))
			v = b
		} else {
			b = v.(bool)
		}
		if omit = !b && fi.omitEmpty; !omit && fi.asString {
			if b {
				v = "true"
			} else {
				v = "false"
			}
		}

	case reflect.Int:
		var i int
		if ptr != 0 {
			i = *(*int)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(int)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatInt(int64(i), 10)
		}
	case reflect.Int8:
		var i int8
		if ptr != 0 {
			i = *(*int8)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(int8)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatInt(int64(i), 10)
		}
	case reflect.Int16:
		var i int16
		if ptr != 0 {
			i = *(*int16)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(int16)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatInt(int64(i), 10)
		}
	case reflect.Int32:
		var i int32
		if ptr != 0 {
			i = *(*int32)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(int32)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatInt(int64(i), 10)
		}
	case reflect.Int64:
		var i int64
		if ptr != 0 {
			i = *(*int64)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(int64)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatInt(i, 10)
		}
	case reflect.Uint:
		var i uint
		if ptr != 0 {
			i = *(*uint)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(uint)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatUint(uint64(i), 10)
		}
	case reflect.Uint8:
		var i uint8
		if ptr != 0 {
			i = *(*uint8)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(uint8)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatUint(uint64(i), 10)
		}
	case reflect.Uint16:
		var i uint16
		if ptr != 0 {
			i = *(*uint16)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(uint16)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatUint(uint64(i), 10)
		}
	case reflect.Uint32:
		var i uint32
		if ptr != 0 {
			i = *(*uint32)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(uint32)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatUint(uint64(i), 10)
		}
	case reflect.Uint64:
		var i uint64
		if ptr != 0 {
			i = *(*uint64)(unsafe.Pointer(ptr))
			v = i
		} else {
			i = v.(uint64)
		}
		if omit = i == 0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatUint(i, 10)
		}

	case reflect.Float32:
		var f float32
		if ptr != 0 {
			f = *(*float32)(unsafe.Pointer(ptr))
			v = f
		} else {
			f = v.(float32)
		}
		if omit = f == 0.0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatFloat(float64(f), 'g', -1, 32)
		}
	case reflect.Float64:
		var f float64
		if ptr != 0 {
			f = *(*float64)(unsafe.Pointer(ptr))
			v = f
		} else {
			f = v.(float64)
		}
		if omit = f == 0.0 && fi.omitEmpty; !omit && fi.asString {
			v = strconv.FormatFloat(f, 'g', -1, 64)
		}

	case reflect.String:
		s := v.(string)
		omit = fi.omitEmpty && len(s) == 0 && fi.omitEmpty
	case reflect.Slice, reflect.Array, reflect.Map:
		omit = (fi.omitEmpty || omitNil) && fv.Len() == 0
	case reflect.Interface, reflect.Ptr:
		omit = (fi.omitEmpty || omitNil) && v == nil
	}
	return
}

func (fi *Field) Append(buf []byte, rv reflect.Value, omitNil bool, embedded bool) ([]byte, interface{}, bool, bool) {
	var v interface{}
	var fv reflect.Value
	var ptr uintptr
	if fi.direct && !embedded {
		ptr = rv.UnsafeAddr() + fi.offset
	} else {
		fv = rv.FieldByIndex(fi.index)
		v = fv.Interface()
	}
	var skip bool
	switch fi.Kind {
	case reflect.Bool:
		var b bool
		if ptr != 0 {
			b = *(*bool)(unsafe.Pointer(ptr))
		} else {
			b = v.(bool)
		}
		if b {
			buf = append(buf, fi.jkey...)
			if fi.asString {
				buf = append(buf, `"true"`...)
			} else {
				buf = append(buf, "true"...)
			}
		} else if fi.omitEmpty {
			buf = append(buf, fi.jkey...)
			if fi.asString {
				buf = append(buf, `"false"`...)
			} else {
				buf = append(buf, "false"...)
			}
		} else {
			return buf, nil, false, false
		}
	case reflect.Int:
		var i int
		if ptr != 0 {
			i = *(*int)(unsafe.Pointer(ptr))
		} else {
			i = v.(int)
		}
		if buf, skip = fi.fillInt(buf, int64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Int8:
		var i int8
		if ptr != 0 {
			i = *(*int8)(unsafe.Pointer(ptr))
		} else {
			i = v.(int8)
		}
		if buf, skip = fi.fillInt(buf, int64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Int16:
		var i int16
		if ptr != 0 {
			i = *(*int16)(unsafe.Pointer(ptr))
		} else {
			i = v.(int16)
		}
		if buf, skip = fi.fillInt(buf, int64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Int32:
		var i int32
		if ptr != 0 {
			i = *(*int32)(unsafe.Pointer(ptr))
		} else {
			i = v.(int32)
		}
		if buf, skip = fi.fillInt(buf, int64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Int64:
		var i int64
		if ptr != 0 {
			i = *(*int64)(unsafe.Pointer(ptr))
		} else {
			i = v.(int64)
		}
		if buf, skip = fi.fillInt(buf, i); skip {
			return buf, nil, false, false
		}
	case reflect.Uint:
		var i uint
		if ptr != 0 {
			i = *(*uint)(unsafe.Pointer(ptr))
		} else {
			i = v.(uint)
		}
		if buf, skip = fi.fillUint(buf, uint64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Uint8:
		var i uint8
		if ptr != 0 {
			i = *(*uint8)(unsafe.Pointer(ptr))
		} else {
			i = v.(uint8)
		}
		if buf, skip = fi.fillUint(buf, uint64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Uint16:
		var i uint16
		if ptr != 0 {
			i = *(*uint16)(unsafe.Pointer(ptr))
		} else {
			i = v.(uint16)
		}
		if buf, skip = fi.fillUint(buf, uint64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Uint32:
		var i uint32
		if ptr != 0 {
			i = *(*uint32)(unsafe.Pointer(ptr))
		} else {
			i = v.(uint32)
		}
		if buf, skip = fi.fillUint(buf, uint64(i)); skip {
			return buf, nil, false, false
		}
	case reflect.Uint64:
		var i uint64
		if ptr != 0 {
			i = *(*uint64)(unsafe.Pointer(ptr))
		} else {
			i = v.(uint64)
		}
		if buf, skip = fi.fillUint(buf, i); skip {
			return buf, nil, false, false
		}

	case reflect.Float32:
		var f float32
		if ptr != 0 {
			f = *(*float32)(unsafe.Pointer(ptr))
		} else {
			f = v.(float32)
		}
		if f == 0.0 && fi.omitEmpty {
			return buf, nil, false, false
		}
		buf = append(buf, fi.jkey...)
		if fi.asString {
			buf = append(buf, '"')
			buf = strconv.AppendFloat(buf, float64(f), 'g', -1, 32)
			buf = append(buf, '"')
		} else {
			buf = strconv.AppendFloat(buf, float64(f), 'g', -1, 32)
		}
	case reflect.Float64:
		var f float64
		if ptr != 0 {
			f = *(*float64)(unsafe.Pointer(ptr))
		} else {
			f = v.(float64)
		}
		if f == 0.0 && fi.omitEmpty {
			return buf, nil, false, false
		}
		buf = append(buf, fi.jkey...)
		if fi.asString {
			buf = append(buf, '"')
			buf = strconv.AppendFloat(buf, f, 'g', -1, 64)
			buf = append(buf, '"')
		} else {
			buf = strconv.AppendFloat(buf, f, 'g', -1, 64)
		}

	case reflect.String:
		s := v.(string)
		if len(s) == 0 && fi.omitEmpty {
			return buf, nil, false, false
		}
		buf = append(buf, fi.jkey...)
		buf = AppendJSONString(buf, s, false) // TBD html safe flag needed

	case reflect.Slice, reflect.Array, reflect.Map:
		if fi.omitEmpty && fv.Len() == 0 {
			return buf, nil, false, false
		}
		buf = append(buf, fi.jkey...)
		return buf, v, false, true

	case reflect.Interface, reflect.Ptr:
		if fi.omitEmpty && v == nil {
			return buf, nil, false, false
		}
		buf = append(buf, fi.jkey...)
		return buf, v, false, true

	default:
		buf = append(buf, fi.jkey...)
		return buf, v, false, true
	}
	buf = append(buf, ',')

	return buf, nil, true, false
}

func (fi *Field) fillInt(buf []byte, i int64) ([]byte, bool) {
	if i == 0 && fi.omitEmpty {
		return buf, true
	}
	buf = append(buf, fi.jkey...)
	if fi.asString {
		buf = append(buf, '"')
		buf = strconv.AppendInt(buf, i, 10)
		buf = append(buf, '"')
	} else {
		buf = strconv.AppendInt(buf, i, 10)
	}
	return buf, false
}

func (fi *Field) fillUint(buf []byte, i uint64) ([]byte, bool) {
	if i == 0 && fi.omitEmpty {
		return buf, true
	}
	buf = append(buf, fi.jkey...)
	if fi.asString {
		buf = append(buf, '"')
		buf = strconv.AppendUint(buf, i, 10)
		buf = append(buf, '"')
	} else {
		buf = strconv.AppendUint(buf, i, 10)
	}
	return buf, false
}

func (fi *Field) setOmitEmpty() {
	switch fi.Kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		fi.omitEmpty = true
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Interface, reflect.Ptr:
		fi.omitEmpty = true
	}
}

func (fi *Field) setup() {
	switch fi.Kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		fi.direct = true
	case reflect.Struct:
		fi.Elem = getTypeStruct(fi.Type)
	case reflect.Ptr, reflect.Slice, reflect.Array:
		et := fi.Type.Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.Elem = getTypeStruct(et)
		}
	}
	fi.jkey = AppendJSONString(fi.jkey, fi.Key, false)
	fi.jkey = append(fi.jkey, ':')
}
