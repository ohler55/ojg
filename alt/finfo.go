// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt

import (
	"reflect"
	"strconv"
	"unsafe"
)

const (
	strMask   = byte(0x01)
	omitMask  = byte(0x02)
	embedMask = byte(0x04)
)

var nilValue reflect.Value

type valFunc func(fi *finfo, rv reflect.Value, addr uintptr) (v interface{}, fv reflect.Value, omit bool)

type finfo struct {
	rt     reflect.Type
	key    string
	value  valFunc
	ivalue valFunc
	index  []int
	offset uintptr
}

func valInt8(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valInt8AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valInt8NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valInt8NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatInt(int64(v), 10), nilValue, false
}

func valInt16(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valInt16AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valInt16NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valInt16NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatInt(int64(v), 10), nilValue, false
}

func valInt32(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valInt32AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatInt(int64(*(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valInt32NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valInt32NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatInt(int64(v), 10), nilValue, false
}

func valInt64(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valInt64AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatInt(*(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10), nilValue, false
}

func valInt64NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valInt64NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*int64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatInt(v, 10), nilValue, false
}

func valUint(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valUintAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valUintNotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valUintNotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatUint(uint64(v), 10), nilValue, false
}

func valUint8(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valUint8AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valUint8NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valUint8NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint8)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatUint(uint64(v), 10), nilValue, false
}

func valUint16(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valUint16AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valUint16NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valUint16NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint16)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatUint(uint64(v), 10), nilValue, false
}

func valUint32(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valUint32AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatUint(uint64(*(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 10), nilValue, false
}

func valUint32NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valUint32NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatUint(uint64(v), 10), nilValue, false
}

func valUint64(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valUint64AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatUint(*(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 10), nilValue, false
}

func valUint64NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0
}

func valUint64NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*uint64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0 {
		return nil, nilValue, true
	}
	return strconv.FormatUint(v, 10), nilValue, false
}

func valFloat32(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valFloat32AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatFloat(float64(*(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))), 'g', -1, 32), nilValue, false
}

func valFloat32NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0.0
}

func valFloat32NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*float32)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return nil, nilValue, true
	}
	return strconv.FormatFloat(float64(v), 'g', -1, 32), nilValue, false
}

func valFloat64(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), nilValue, false
}

func valFloat64AsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return strconv.FormatFloat(*(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset)), 'g', -1, 64), nilValue, false
}

func valFloat64NotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	return v, nilValue, v == 0.0
}

func valFloat64NotEmptyAsString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	v := *(*float64)(unsafe.Pointer(rv.UnsafeAddr() + fi.offset))
	if v == 0.0 {
		return nil, nilValue, true
	}
	return strconv.FormatFloat(v, 'g', -1, 64), nilValue, false
}

func valString(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	return rv.FieldByIndex(fi.index).String(), nilValue, false
}

func valStringNotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	s := rv.FieldByIndex(fi.index).String()
	if len(s) == 0 {
		return s, nilValue, true
	}
	return s, nilValue, false
}

func valJustVal(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.index)
	return fv.Interface(), fv, false
}

func valPtrNotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.index)
	v := fv.Interface()
	return v, fv, v == nil
}

func valSliceNotEmpty(fi *finfo, rv reflect.Value, addr uintptr) (interface{}, reflect.Value, bool) {
	fv := rv.FieldByIndex(fi.index)
	if fv.Len() == 0 {
		return nil, nilValue, true
	}
	return fv.Interface(), fv, false
}

func newFinfo(f reflect.StructField, key string, fx byte) *finfo {
	fi := finfo{
		rt:     f.Type,
		key:    key,
		index:  f.Index,
		value:  valJustVal, // replace as necessary later
		ivalue: valJustVal, // replace as necessary later
		offset: f.Offset,
	}
	// TBD remove once converted
	omitEmpty := (fx & omitMask) != 0
	asString := (fx & strMask) != 0
	embedded := (fx & embedMask) != 0
	if embedded {
		fi.value = valJustVal
		return &fi
	}

	switch f.Type.Kind() {
	case reflect.Bool:
		fi.value = boolValFuncs[fx]
		fi.ivalue = boolValFuncs[fx|embedMask]

	case reflect.Int:
		fi.value = intValFuncs[fx]
		fi.ivalue = intValFuncs[fx|embedMask]

	case reflect.Int8:
		if asString {
			if omitEmpty {
				fi.value = valInt8NotEmptyAsString
			} else {
				fi.value = valInt8AsString
			}
		} else {
			if omitEmpty {
				fi.value = valInt8NotEmpty
			} else {
				fi.value = valInt8
			}
		}
	case reflect.Int16:
		if asString {
			if omitEmpty {
				fi.value = valInt16NotEmptyAsString
			} else {
				fi.value = valInt16AsString
			}
		} else {
			if omitEmpty {
				fi.value = valInt16NotEmpty
			} else {
				fi.value = valInt16
			}
		}
	case reflect.Int32:
		if asString {
			if omitEmpty {
				fi.value = valInt32NotEmptyAsString
			} else {
				fi.value = valInt32AsString
			}
		} else {
			if omitEmpty {
				fi.value = valInt32NotEmpty
			} else {
				fi.value = valInt32
			}
		}
	case reflect.Int64:
		if asString {
			if omitEmpty {
				fi.value = valInt64NotEmptyAsString
			} else {
				fi.value = valInt64AsString
			}
		} else {
			if omitEmpty {
				fi.value = valInt64NotEmpty
			} else {
				fi.value = valInt64
			}
		}
	case reflect.Uint:
		if asString {
			if omitEmpty {
				fi.value = valUintNotEmptyAsString
			} else {
				fi.value = valUintAsString
			}
		} else {
			if omitEmpty {
				fi.value = valUintNotEmpty
			} else {
				fi.value = valUint
			}
		}
	case reflect.Uint8:
		if asString {
			if omitEmpty {
				fi.value = valUint8NotEmptyAsString
			} else {
				fi.value = valUint8AsString
			}
		} else {
			if omitEmpty {
				fi.value = valUint8NotEmpty
			} else {
				fi.value = valUint8
			}
		}
	case reflect.Uint16:
		if asString {
			if omitEmpty {
				fi.value = valUint16NotEmptyAsString
			} else {
				fi.value = valUint16AsString
			}
		} else {
			if omitEmpty {
				fi.value = valUint16NotEmpty
			} else {
				fi.value = valUint16
			}
		}
	case reflect.Uint32:
		if asString {
			if omitEmpty {
				fi.value = valUint32NotEmptyAsString
			} else {
				fi.value = valUint32AsString
			}
		} else {
			if omitEmpty {
				fi.value = valUint32NotEmpty
			} else {
				fi.value = valUint32
			}
		}
	case reflect.Uint64:
		if asString {
			if omitEmpty {
				fi.value = valUint64NotEmptyAsString
			} else {
				fi.value = valUint64AsString
			}
		} else {
			if omitEmpty {
				fi.value = valUint64NotEmpty
			} else {
				fi.value = valUint64
			}
		}
	case reflect.Float32:
		if asString {
			if omitEmpty {
				fi.value = valFloat32NotEmptyAsString
			} else {
				fi.value = valFloat32AsString
			}
		} else {
			if omitEmpty {
				fi.value = valFloat32NotEmpty
			} else {
				fi.value = valFloat32
			}
		}
	case reflect.Float64:
		if asString {
			if omitEmpty {
				fi.value = valFloat64NotEmptyAsString
			} else {
				fi.value = valFloat64AsString
			}
		} else {
			if omitEmpty {
				fi.value = valFloat64NotEmpty
			} else {
				fi.value = valFloat64
			}
		}
	case reflect.String:
		if omitEmpty {
			fi.value = valStringNotEmpty
		} else {
			fi.value = valString
		}
	case reflect.Struct:
		fi.value = valJustVal
	case reflect.Ptr:
		if omitEmpty {
			fi.value = valPtrNotEmpty
		} else {
			fi.value = valJustVal
		}
	case reflect.Interface:
		if omitEmpty {
			fi.value = valPtrNotEmpty
		} else {
			fi.value = valJustVal
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if omitEmpty {
			fi.value = valSliceNotEmpty
		} else {
			fi.value = valJustVal
		}
	}
	return &fi
}
