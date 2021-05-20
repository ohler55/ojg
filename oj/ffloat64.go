// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj

import (
	"reflect"
	"strconv"
	"unsafe"
)

var float64AppendFuncs = [8]appendFunc{
	appendFloat64,
	appendFloat64AsString,
	appendFloat64NotEmpty,
	appendFloat64NotEmptyAsString,
	iappendFloat64,
	iappendFloat64AsString,
	iappendFloat64NotEmpty,
	iappendFloat64NotEmptyAsString,
}

func appendFloat64(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(addr + fi.offset)), 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64AsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, *(*float64)(unsafe.Pointer(addr + fi.offset)), 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func appendFloat64NotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float64)(unsafe.Pointer(addr + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)

	return buf, nil, true, false
}

func appendFloat64NotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := *(*float64)(unsafe.Pointer(addr + fi.offset))
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendFloat64(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, rv.FieldByIndex(fi.index).Interface().(float64), 'g', -1, 64)

	return buf, nil, true, false
}

func iappendFloat64AsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, rv.FieldByIndex(fi.index).Interface().(float64), 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}

func iappendFloat64NotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(float64)
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)

	return buf, nil, true, false
}

func iappendFloat64NotEmptyAsString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface().(float64)
	if v == 0.0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
	buf = append(buf, '"')

	return buf, nil, true, false
}
