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
	Append func(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool)
	Value  func(fi *Field, rv reflect.Value) (v interface{}, fv *reflect.Value, omit bool)
	jkey   []byte
	Index  []int
	offset uintptr
}

func valBool(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valBoolAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		return "true", nil, false
	}
	return "false", nil, false
}

func valBoolNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, !v
}

func valBoolNotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		return "true", nil, false
	}
	return "false", nil, true
}

func valInt(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valIntAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valIntNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valIntNotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatInt(int64(v), 10), nil, false
}

func valInt8(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valInt8AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valInt8NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valInt8NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatInt(int64(v), 10), nil, false
}

func valInt16(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valInt16AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valInt16NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valInt16NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatInt(int64(v), 10), nil, false
}

func valInt32(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valInt32AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valInt32NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valInt32NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatInt(int64(v), 10), nil, false
}

func valInt64(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valInt64AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatInt(*(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10), nil, false
}

func valInt64NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valInt64NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatInt(v, 10), nil, false
}

func valUint(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valUintAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valUintNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valUintNotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatUint(uint64(v), 10), nil, false
}

func valUint8(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valUint8AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valUint8NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valUint8NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatUint(uint64(v), 10), nil, false
}

func valUint16(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valUint16AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valUint16NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valUint16NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatUint(uint64(v), 10), nil, false
}

func valUint32(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valUint32AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nil, false
}

func valUint32NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valUint32NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatUint(uint64(v), 10), nil, false
}

func valUint64(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valUint64AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatUint(*(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10), nil, false
}

func valUint64NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0
}

func valUint64NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nil, true
	}
	return strconv.FormatUint(v, 10), nil, false
}

func valFloat32(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valFloat32AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatFloat(float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32), nil, false
}

func valFloat32NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0.0
}

func valFloat32NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return nil, nil, true
	}
	return strconv.FormatFloat(float64(v), 'g', -1, 32), nil, false
}

func valFloat64(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nil, false
}

func valFloat64AsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return strconv.FormatFloat(*(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64), nil, false
}

func valFloat64NotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nil, v == 0.0
}

func valFloat64NotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return nil, nil, true
	}
	return strconv.FormatFloat(v, 'g', -1, 64), nil, false
}

func valString(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	return rv.FieldByIndex(fi.Index).String(), nil, false
}

func valStringNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	s := rv.FieldByIndex(fi.Index).String()
	if len(s) == 0 {
		return s, nil, true
	}
	return s, nil, false
}

func valJustVal(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.Index)
	return fv.Interface(), &fv, false
}

func valStruct(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.Index)
	return fv.Interface(), &fv, false
}

func valPtrNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.Index)
	v := fv.Interface()
	return v, &fv, v == nil
}

func valSliceNotEmpty(fi *Field, rv reflect.Value) (interface{}, *reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.Index)
	if fv.Len() == 0 {
		return nil, nil, true
	}
	return fv.Interface(), &fv, false
}

func appendBool(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, "true"...)
	} else {
		buf = append(buf, "false"...)
	}
	return buf, nil, true, false
}

func appendBoolAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendBoolNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, "true"...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func appendBoolNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, `"true"`...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func appendInt(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendIntAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendIntNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendIntNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendInt8(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt8AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt8NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt8NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendInt16(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt16AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt16NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt16NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendInt32(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendInt32AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt32NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, int64(v), 10)

	return buf, nil, true, false
}

func appendInt32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendInt64(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)

	return buf, nil, true, false
}

func appendInt64AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendInt(buf, *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendInt64NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendInt(buf, v, 10)

	return buf, nil, true, false
}

func appendInt64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendUint(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUintAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUintNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUintNotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendUint8(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint8AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint8NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint8NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendUint16(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint16AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint16NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint16NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendUint32(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUint32AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint32NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUint32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendUint64(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)

	return buf, nil, true, false
}

func appendUint64AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint64NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, v, 10)

	return buf, nil, true, false
}

func appendUint64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendFloat32(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32)

	return buf, nil, true, false
}

func appendFloat32AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat32NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 32)

	return buf, nil, true, false
}

func appendFloat32NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendFloat64(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64AsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat64NotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64NotEmptyAsString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
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

func appendString(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).String()
	buf = append(buf, fi.jkey...)
	buf = AppendJSONString(buf, v, safe)

	return buf, nil, true, false
}

func appendStringNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	var fv reflect.Value
	fv = rv.FieldByIndex(fi.Index)
	s := fv.String()
	if len(s) == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = AppendJSONString(buf, s, safe)

	return buf, nil, true, false
}

func appendJustKey(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).Interface()
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendPtrNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.Index).Interface()
	if v == nil {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendSliceNotEmpty(fi *Field, buf []byte, rv reflect.Value, safe bool) ([]byte, interface{}, bool, bool) {
	fv := rv.FieldByIndex(fi.Index)
	if fv.Len() == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, fv.Interface(), false, true
}

func newField(f reflect.StructField, key string, omitEmpty, asString, anon bool) *Field {
	fi := Field{
		Type:   f.Type,
		Key:    key,
		Kind:   f.Type.Kind(),
		Index:  f.Index,
		offset: f.Offset,
	}
	switch fi.Kind {
	case reflect.Bool:
		if asString {
			if omitEmpty {
				fi.Append = appendBoolNotEmptyAsString
				fi.Value = valBoolNotEmptyAsString
			} else {
				fi.Append = appendBoolAsString
				fi.Value = valBoolAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendBoolNotEmpty
				fi.Value = valBoolNotEmpty
			} else {
				fi.Append = appendBool
				fi.Value = valBool
			}
		}
	case reflect.Int:
		if asString {
			if omitEmpty {
				fi.Append = appendIntNotEmptyAsString
				fi.Value = valIntNotEmptyAsString
			} else {
				fi.Append = appendIntAsString
				fi.Value = valIntAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendIntNotEmpty
				fi.Value = valIntNotEmpty
			} else {
				fi.Append = appendInt
				fi.Value = valInt
			}
		}
	case reflect.Int8:
		if asString {
			if omitEmpty {
				fi.Append = appendInt8NotEmptyAsString
				fi.Value = valInt8NotEmptyAsString
			} else {
				fi.Append = appendInt8AsString
				fi.Value = valInt8AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt8NotEmpty
				fi.Value = valInt8NotEmpty
			} else {
				fi.Append = appendInt8
				fi.Value = valInt8
			}
		}
	case reflect.Int16:
		if asString {
			if omitEmpty {
				fi.Append = appendInt16NotEmptyAsString
				fi.Value = valInt16NotEmptyAsString
			} else {
				fi.Append = appendInt16AsString
				fi.Value = valInt16AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt16NotEmpty
				fi.Value = valInt16NotEmpty
			} else {
				fi.Append = appendInt16
				fi.Value = valInt16
			}
		}
	case reflect.Int32:
		if asString {
			if omitEmpty {
				fi.Append = appendInt32NotEmptyAsString
				fi.Value = valInt32NotEmptyAsString
			} else {
				fi.Append = appendInt32AsString
				fi.Value = valInt32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt32NotEmpty
				fi.Value = valInt32NotEmpty
			} else {
				fi.Append = appendInt32
				fi.Value = valInt32
			}
		}
	case reflect.Int64:
		if asString {
			if omitEmpty {
				fi.Append = appendInt64NotEmptyAsString
				fi.Value = valInt64NotEmptyAsString
			} else {
				fi.Append = appendInt64AsString
				fi.Value = valInt64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendInt64NotEmpty
				fi.Value = valInt64NotEmpty
			} else {
				fi.Append = appendInt64
				fi.Value = valInt64
			}
		}
	case reflect.Uint:
		if asString {
			if omitEmpty {
				fi.Append = appendUintNotEmptyAsString
				fi.Value = valUintNotEmptyAsString
			} else {
				fi.Append = appendUintAsString
				fi.Value = valUintAsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUintNotEmpty
				fi.Value = valUintNotEmpty
			} else {
				fi.Append = appendUint
				fi.Value = valUint
			}
		}
	case reflect.Uint8:
		if asString {
			if omitEmpty {
				fi.Append = appendUint8NotEmptyAsString
				fi.Value = valUint8NotEmptyAsString
			} else {
				fi.Append = appendUint8AsString
				fi.Value = valUint8AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint8NotEmpty
				fi.Value = valUint8NotEmpty
			} else {
				fi.Append = appendUint8
				fi.Value = valUint8
			}
		}
	case reflect.Uint16:
		if asString {
			if omitEmpty {
				fi.Append = appendUint16NotEmptyAsString
				fi.Value = valUint16NotEmptyAsString
			} else {
				fi.Append = appendUint16AsString
				fi.Value = valUint16AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint16NotEmpty
				fi.Value = valUint16NotEmpty
			} else {
				fi.Append = appendUint16
				fi.Value = valUint16
			}
		}
	case reflect.Uint32:
		if asString {
			if omitEmpty {
				fi.Append = appendUint32NotEmptyAsString
				fi.Value = valUint32NotEmptyAsString
			} else {
				fi.Append = appendUint32AsString
				fi.Value = valUint32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint32NotEmpty
				fi.Value = valUint32NotEmpty
			} else {
				fi.Append = appendUint32
				fi.Value = valUint32
			}
		}
	case reflect.Uint64:
		if asString {
			if omitEmpty {
				fi.Append = appendUint64NotEmptyAsString
				fi.Value = valUint64NotEmptyAsString
			} else {
				fi.Append = appendUint64AsString
				fi.Value = valUint64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendUint64NotEmpty
				fi.Value = valUint64NotEmpty
			} else {
				fi.Append = appendUint64
				fi.Value = valUint64
			}
		}
	case reflect.Float32:
		if asString {
			if omitEmpty {
				fi.Append = appendFloat32NotEmptyAsString
				fi.Value = valFloat32NotEmptyAsString
			} else {
				fi.Append = appendFloat32AsString
				fi.Value = valFloat32AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendFloat32NotEmpty
				fi.Value = valFloat32NotEmpty
			} else {
				fi.Append = appendFloat32
				fi.Value = valFloat32
			}
		}
	case reflect.Float64:
		if asString {
			if omitEmpty {
				fi.Append = appendFloat64NotEmptyAsString
				fi.Value = valFloat64NotEmptyAsString
			} else {
				fi.Append = appendFloat64AsString
				fi.Value = valFloat64AsString
			}
		} else {
			if omitEmpty {
				fi.Append = appendFloat64NotEmpty
				fi.Value = valFloat64NotEmpty
			} else {
				fi.Append = appendFloat64
				fi.Value = valFloat64
			}
		}
	case reflect.String:
		if omitEmpty {
			fi.Append = appendStringNotEmpty
			fi.Value = valStringNotEmpty
		} else {
			fi.Append = appendString
			fi.Value = valString
		}
	case reflect.Struct:
		fi.Elem = getTypeStruct(fi.Type)
		fi.Append = appendJustKey
		fi.Value = valStruct
		// TBD put back in
		// fi.Value = valJustVal
	case reflect.Ptr:
		et := fi.Type.Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.Elem = getTypeStruct(et)
		}
		if omitEmpty {
			fi.Append = appendPtrNotEmpty
			fi.Value = valPtrNotEmpty
		} else {
			fi.Append = appendJustKey
			fi.Value = valJustVal
		}
	case reflect.Interface:
		if omitEmpty {
			fi.Append = appendPtrNotEmpty
			fi.Value = valPtrNotEmpty
		} else {
			fi.Append = appendJustKey
			fi.Value = valJustVal
		}
	case reflect.Slice, reflect.Array:
		et := fi.Type.Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.Elem = getTypeStruct(et)
		}
		if omitEmpty {
			fi.Append = appendSliceNotEmpty
			fi.Value = valSliceNotEmpty
		} else {
			fi.Append = appendJustKey
			fi.Value = valJustVal
		}
	}
	fi.jkey = AppendJSONString(fi.jkey, fi.Key, false)
	fi.jkey = append(fi.jkey, ':')

	return &fi
}
