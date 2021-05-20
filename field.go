// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"strconv"
	"unsafe"
)

// Field hold information about a struct field.
type Field struct {
	Type   reflect.Type
	Key    string
	Kind   reflect.Kind
	Elem   *Struct
	Append func(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool)
	jkey   []byte
	Index  []int
	offset uintptr
}

// KeyLen returns the length of the key plus syntax. For example a JSON key of
// _key_ would become "key": with a KeyLen of 6.
func (f *Field) KeyLen() int {
	return len(f.jkey)
}

func appendBool(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, "true"...)
	} else {
		buf = append(buf, "false"...)
	}
	return buf, nil, true, false
}

func appendBoolAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, `"true"`...)
	} else {
		buf = append(buf, `"false"`...)
	}
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendBoolNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, "true"...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func appendBoolNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, `"true"`...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func appendInt(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendIntAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendIntNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendIntNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt8(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt8AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt8NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt8NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt16(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt16AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt16NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt16NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt32(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt32AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt32NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt64(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)

	return buf, nil, true, false
}

func appendInt64AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt64NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, v, 10)

	return buf, nil, true, false
}

func appendInt64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, v, 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUintAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUintNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUintNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint8(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint8AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint8NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint8NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint16(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint16AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint16NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint16NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint32(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint32AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint32NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint64(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)

	return buf, nil, true, false
}

func appendUint64AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint64NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, v, 10)

	return buf, nil, true, false
}

func appendUint64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, v, 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat32(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32)

	return buf, nil, true, false
}

func appendFloat32AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat32NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 32)

	return buf, nil, true, false
}

func appendFloat32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 32)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat64(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64AsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat64NotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).String()
	buf = append(buf, fi.jkey...)
	buf = AppendJSONString(buf, v, safe)

	return buf, nil, true, false
}

func appendStringNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	s := rv.FieldByIndex(fi.Index).String()
	if len(s) == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = AppendJSONString(buf, s, safe)

	return buf, nil, true, false
}

func appendSENString(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).String()
	buf = append(buf, fi.jkey...)
	buf = AppendSENString(buf, v, safe)

	return buf, nil, true, false
}

func appendSENStringNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	s := rv.FieldByIndex(fi.Index).String()
	if len(s) == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = AppendSENString(buf, s, safe)

	return buf, nil, true, false
}

func appendJustKey(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).Interface()
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendPtrNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).Interface()
	if v == nil {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendSliceNotEmpty(fi *Field, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	fv := rv.FieldByIndex(fi.Index)
	if fv.Len() == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, fv.Interface(), false, true
}

func newField(f reflect.StructField, key string, omitEmpty, asString, pretty, sen, embedded bool) *Field {
	fi := Field{
		Type:   f.Type,
		Key:    key,
		Kind:   f.Type.Kind(),
		Index:  f.Index,
		offset: f.Offset,
	}
	if embedded {
		// TBD might have to move the embedded check into each
		fi.Append = appendJustKey
		if sen {
			fi.jkey = AppendSENString(fi.jkey, fi.Key, false)
		} else {
			fi.jkey = AppendJSONString(fi.jkey, fi.Key, false)
		}
		fi.jkey = append(fi.jkey, ':')
		if pretty {
			fi.jkey = append(fi.jkey, ' ')
		}
		return &fi
	}
	// TBD calc fx in struct.go
	switch fi.Kind {
	case reflect.Bool:
		if asString {
			if omitEmpty {
				fi.Append = appendBoolNotEmptyAsString
			} else {
				fi.Append = appendBoolAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendBoolNotEmpty
			} else {
				fi.Append = appendBool
			}
		}

	case reflect.Int:
		if asString {
			if omitEmpty {
				fi.Append = appendIntNotEmptyAsString
			} else {
				fi.Append = appendIntAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendIntNotEmpty
			} else {
				fi.Append = appendInt
			}
		}
	case reflect.Int8:
		if asString {
			if omitEmpty {
				fi.Append = appendInt8NotEmptyAsString
			} else {
				fi.Append = appendInt8AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt8NotEmpty
			} else {
				fi.Append = appendInt8
			}
		}
	case reflect.Int16:
		if asString {
			if omitEmpty {
				fi.Append = appendInt16NotEmptyAsString
			} else {
				fi.Append = appendInt16AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt16NotEmpty
			} else {
				fi.Append = appendInt16
			}
		}
	case reflect.Int32:
		if asString {
			if omitEmpty {
				fi.Append = appendInt32NotEmptyAsString
			} else {
				fi.Append = appendInt32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt32NotEmpty
			} else {
				fi.Append = appendInt32
			}
		}
	case reflect.Int64:
		if asString {
			if omitEmpty {
				fi.Append = appendInt64NotEmptyAsString
			} else {
				fi.Append = appendInt64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt64NotEmpty
			} else {
				fi.Append = appendInt64
			}
		}
	case reflect.Uint:
		if asString {
			if omitEmpty {
				fi.Append = appendUintNotEmptyAsString
			} else {
				fi.Append = appendUintAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUintNotEmpty
			} else {
				fi.Append = appendUint
			}
		}
	case reflect.Uint8:
		if asString {
			if omitEmpty {
				fi.Append = appendUint8NotEmptyAsString
			} else {
				fi.Append = appendUint8AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint8NotEmpty
			} else {
				fi.Append = appendUint8
			}
		}
	case reflect.Uint16:
		if asString {
			if omitEmpty {
				fi.Append = appendUint16NotEmptyAsString
			} else {
				fi.Append = appendUint16AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint16NotEmpty
			} else {
				fi.Append = appendUint16
			}
		}
	case reflect.Uint32:
		if asString {
			if omitEmpty {
				fi.Append = appendUint32NotEmptyAsString
			} else {
				fi.Append = appendUint32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint32NotEmpty
			} else {
				fi.Append = appendUint32
			}
		}
	case reflect.Uint64:
		if asString {
			if omitEmpty {
				fi.Append = appendUint64NotEmptyAsString
			} else {
				fi.Append = appendUint64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint64NotEmpty
			} else {
				fi.Append = appendUint64
			}
		}
	case reflect.Float32:
		if asString {
			if omitEmpty {
				fi.Append = appendFloat32NotEmptyAsString
			} else {
				fi.Append = appendFloat32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendFloat32NotEmpty
			} else {
				fi.Append = appendFloat32
			}
		}
	case reflect.Float64:
		if asString {
			if omitEmpty {
				fi.Append = appendFloat64NotEmptyAsString
			} else {
				fi.Append = appendFloat64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendFloat64NotEmpty
			} else {
				fi.Append = appendFloat64
			}
		}
	case reflect.String:
		if omitEmpty {
			if sen {
				fi.Append = appendSENStringNotEmpty
			} else {
				fi.Append = appendStringNotEmpty
			}
		} else {
			if sen {
				fi.Append = appendSENString
			} else {
				fi.Append = appendString
			}
		}
	case reflect.Struct:
		fi.Elem = getTypeStruct(fi.Type, true)
		fi.Append = appendJustKey
	case reflect.Ptr:
		et := fi.Type.Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.Elem = getTypeStruct(et, false)
		}
		if omitEmpty {
			fi.Append = appendPtrNotEmpty
		} else {
			fi.Append = appendJustKey
		}
	case reflect.Interface:
		if omitEmpty {
			fi.Append = appendPtrNotEmpty
		} else {
			fi.Append = appendJustKey
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		et := fi.Type.Elem()
		embedded := true
		if et.Kind() == reflect.Ptr {
			embedded = false
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.Elem = getTypeStruct(et, embedded)
		}
		if omitEmpty {
			fi.Append = appendSliceNotEmpty
		} else {
			fi.Append = appendJustKey
		}
	}
	if sen {
		fi.jkey = AppendSENString(fi.jkey, fi.Key, false)
	} else {
		fi.jkey = AppendJSONString(fi.jkey, fi.Key, false)
	}
	fi.jkey = append(fi.jkey, ':')
	if pretty {
		fi.jkey = append(fi.jkey, ' ')
	}
	return &fi
}
