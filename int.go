// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"strconv"
	"unsafe"
)

var intValFuncs = [8]func(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool){
	valInt,
	valIntAsString,
	valIntNotEmpty,
	valIntNotEmptyAsString,
	valInt, // index based
	valIntAsString,
	valIntNotEmpty,
	valIntNotEmptyAsString,
}

func valInt(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	return *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valIntAsString(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valIntNotEmpty(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valIntNotEmptyAsString(fi *Field, rv reflect.Value) (interface{}, reflect.Value, bool) {
	v := *(*int)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatInt(int64(v), 10), nilValue, false
}
