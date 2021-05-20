// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj

import (
	"reflect"
	"unsafe"
)

var boolAppendFuncs = [8]appendFunc{
	appendBool,
	appendBoolAsString,
	appendBoolNotEmpty,
	appendBoolNotEmptyAsString,
	iappendBool,
	iappendBoolAsString,
	iappendBoolNotEmpty,
	iappendBoolNotEmptyAsString,
}

func appendBool(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	if *(*bool)(unsafe.Pointer(addr + fi.offset)) {
		buf = append(buf, "true"...)
	} else {
		buf = append(buf, "false"...)
	}
	return buf, nil, true, false
}

func appendBoolAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	if *(*bool)(unsafe.Pointer(addr + fi.offset)) {
		buf = append(buf, `"true"`...)
	} else {
		buf = append(buf, `"false"`...)
	}
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendBoolNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(addr + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, "true"...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func appendBoolNotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if *(*bool)(unsafe.Pointer(addr + fi.offset)) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, `"true"`...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func iappendBool(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	if rv.FieldByIndex(fi.index).Interface().(bool) {
		buf = append(buf, "true"...)
	} else {
		buf = append(buf, "false"...)
	}
	return buf, nil, true, false
}

func iappendBoolAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	if rv.FieldByIndex(fi.index).Interface().(bool) {
		buf = append(buf, `"true"`...)
	} else {
		buf = append(buf, `"false"`...)
	}
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendBoolNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if rv.FieldByIndex(fi.index).Interface().(bool) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, "true"...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}

func iappendBoolNotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	if rv.FieldByIndex(fi.index).Interface().(bool) {
		buf = append(buf, fi.jkey...)
		buf = append(buf, `"true"`...)
		return buf, nil, true, false
	}
	return buf, nil, false, false
}
