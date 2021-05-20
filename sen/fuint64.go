// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"reflect"
	"strconv"
	"unsafe"
)

var uint64AppendFuncs = [8]appendFunc{
	appendUint64,
	appendUint64AsString,
	appendUint64NotEmpty,
	appendUint64NotEmptyAsString,
	iappendUint64,
	iappendUint64AsString,
	iappendUint64NotEmpty,
	iappendUint64NotEmptyAsString,
}

func appendUint64(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(addr + fi.offset)), 10)

	return buf, nil, true, false
}

func appendUint64AsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, *(*uint64)(unsafe.Pointer(addr + fi.offset)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUint64NotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint64)(unsafe.Pointer(addr + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, v, 10)

	return buf, nil, true, false
}

func appendUint64NotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint64)(unsafe.Pointer(addr + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, v, 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendUint64(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, rv.FieldByIndex(fi.index).Interface().(uint64), 10)

	return buf, nil, true, false
}

func iappendUint64AsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, rv.FieldByIndex(fi.index).Interface().(uint64), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendUint64NotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(uint64)
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, v, 10)

	return buf, nil, true, false
}

func iappendUint64NotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(uint64)
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, v, 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}
