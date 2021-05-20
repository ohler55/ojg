// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj

import (
	"reflect"

	"github.com/ohler55/ojg"
)

const (
	strMask   = byte(0x01)
	omitMask  = byte(0x02)
	embedMask = byte(0x04)
)

type appendFunc func(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool)

// Field hold information about a struct field.
type finfo struct {
	rt      reflect.Type
	key     string
	kind    reflect.Kind
	elem    *sinfo
	Append  appendFunc
	iAppend appendFunc
	jkey    []byte
	index   []int
	offset  uintptr
}

// KeyLen returns the length of the key plus syntax. For example a JSON key of
// _key_ would become "key": with a KeyLen of 6.
func (f *finfo) KeyLen() int {
	return len(f.jkey)
}

func appendString(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).String()
	buf = append(buf, fi.jkey...)
	buf = ojg.AppendJSONString(buf, v, safe)

	return buf, nil, true, false
}

func appendStringNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	s := rv.FieldByIndex(fi.index).String()
	if len(s) == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	buf = ojg.AppendJSONString(buf, s, safe)

	return buf, nil, true, false
}

func appendJustKey(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface()
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendPtrNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	v := rv.FieldByIndex(fi.index).Interface()
	if v == nil {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, v, false, true
}

func appendSliceNotEmpty(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, interface{}, bool, bool) {
	fv := rv.FieldByIndex(fi.index)
	if fv.Len() == 0 {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	return buf, fv.Interface(), false, true
}

func newFinfo(f reflect.StructField, key string, omitEmpty, asString, pretty, embedded bool) *finfo {
	fi := finfo{
		rt:     f.Type,
		key:    key,
		kind:   f.Type.Kind(),
		index:  f.Index,
		offset: f.Offset,
	}
	var fx byte
	if omitEmpty {
		fx |= omitMask
	}
	if asString {
		fx |= strMask
	}
	if embedded {
		fx |= embedMask
	}
	switch fi.kind {
	case reflect.Bool:
		fi.Append = boolAppendFuncs[fx]
		fi.iAppend = boolAppendFuncs[fx|embedMask]

	case reflect.Int:
		fi.Append = intAppendFuncs[fx]
		fi.iAppend = intAppendFuncs[fx|embedMask]
	case reflect.Int8:
		fi.Append = int8AppendFuncs[fx]
		fi.iAppend = int8AppendFuncs[fx|embedMask]
	case reflect.Int16:
		fi.Append = int16AppendFuncs[fx]
		fi.iAppend = int16AppendFuncs[fx|embedMask]
	case reflect.Int32:
		fi.Append = int32AppendFuncs[fx]
		fi.iAppend = int32AppendFuncs[fx|embedMask]
	case reflect.Int64:
		fi.Append = int64AppendFuncs[fx]
		fi.iAppend = int64AppendFuncs[fx|embedMask]

	case reflect.Uint:
		fi.Append = uintAppendFuncs[fx]
		fi.iAppend = uintAppendFuncs[fx|embedMask]
	case reflect.Uint8:
		fi.Append = uint8AppendFuncs[fx]
		fi.iAppend = uint8AppendFuncs[fx|embedMask]
	case reflect.Uint16:
		fi.Append = uint16AppendFuncs[fx]
		fi.iAppend = uint16AppendFuncs[fx|embedMask]
	case reflect.Uint32:
		fi.Append = uint32AppendFuncs[fx]
		fi.iAppend = uint32AppendFuncs[fx|embedMask]
	case reflect.Uint64:
		fi.Append = uint64AppendFuncs[fx]
		fi.iAppend = uint64AppendFuncs[fx|embedMask]

	case reflect.Float32:
		fi.Append = float32AppendFuncs[fx]
		fi.iAppend = float32AppendFuncs[fx|embedMask]
	case reflect.Float64:
		fi.Append = float64AppendFuncs[fx]
		fi.iAppend = float64AppendFuncs[fx|embedMask]

	case reflect.String:
		if omitEmpty {
			fi.Append = appendStringNotEmpty
		} else {
			fi.Append = appendString
		}
	case reflect.Struct:
		fi.elem = getTypeStruct(fi.rt, true)
		fi.Append = appendJustKey
	case reflect.Ptr:
		et := fi.rt.Elem()
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.elem = getTypeStruct(et, false)
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
		et := fi.rt.Elem()
		embedded := true
		if et.Kind() == reflect.Ptr {
			embedded = false
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct {
			fi.elem = getTypeStruct(et, embedded)
		}
		if omitEmpty {
			fi.Append = appendSliceNotEmpty
		} else {
			fi.Append = appendJustKey
		}
	}
	fi.jkey = ojg.AppendJSONString(fi.jkey, fi.key, false)
	fi.jkey = append(fi.jkey, ':')
	if pretty {
		fi.jkey = append(fi.jkey, ' ')
	}
	return &fi
}
