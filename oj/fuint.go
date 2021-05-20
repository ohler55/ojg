// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj

import (
	"reflect"
	"strconv"
	"unsafe"
)

var uintAppendFuncs = [8]appendFunc{
	appendUint,
	appendUintAsString,
	appendUintNotEmpty,
	appendUintNotEmptyAsString,
	iappendUint,
	iappendUintAsString,
	iappendUintNotEmpty,
	iappendUintNotEmptyAsString,
}

func appendUint(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(addr + fi.offset))), 10)

	return buf, nil, true, false
}

func appendUintAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(*(*uint)(unsafe.Pointer(addr + fi.offset))), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendUintNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint)(unsafe.Pointer(addr + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func appendUintNotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*uint)(unsafe.Pointer(addr + fi.offset))
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendUint(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(rv.FieldByIndex(fi.index).Interface().(uint)), 10)

	return buf, nil, true, false
}

func iappendUintAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(rv.FieldByIndex(fi.index).Interface().(uint)), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendUintNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(uint)
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendUint(buf, uint64(v), 10)

	return buf, nil, true, false
}

func iappendUintNotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(uint)
	if v == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendUint(buf, uint64(v), 10)
	buf = append(buf, '"')

	return buf, nil, true, false
}
