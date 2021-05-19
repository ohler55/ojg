// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"unsafe"
)

const (
	strMask   = byte(0x01)
	omitMask  = byte(0x02)
	embedMask = byte(0x04)
)

var boolValFuncs = [8]func(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool){
	valBool,
	valBoolAsString,
	valBoolNotEmpty,
	valBoolNotEmptyAsString,
	valBool, // index based
	valBoolAsString,
	valBoolNotEmpty,
	valBoolNotEmptyAsString,
}

func valBool(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	return *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valBoolAsString(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		return "true", nilValue, false
	}
	return "false", nilValue, false
}

func valBoolNotEmpty(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	v := *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, !v
}

func valBoolNotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	if *(*bool)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)) {
		return "true", nilValue, false
	}
	return "false", nilValue, true
}
