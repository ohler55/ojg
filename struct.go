// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"sort"
	"strings"
	"sync"
	"unsafe"
)

const (
	MaskByTag  = 0x10
	MaskExact  = 0x08 // exact key vs lowwer case first letter
	MaskPretty = 0x04
	MaskNested = 0x02
	MaskSen    = 0x01
	MaskSet    = 0x20
)

// Struct holds reflect information about a struct.
type Struct struct {
	Type    reflect.Type
	ByTag   []*Field
	ByName  []*Field
	ByLow   []*Field
	OutTag  []*Field
	OutName []*Field
	OutLow  []*Field
}

var (
	structMut sync.Mutex
	// Keyed by the pointer to the type.
	structMap = map[uintptr]*Struct{}
)

func GetTypeStruct(rt reflect.Type) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&rt))[1]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(rt, x)
}

// Non-locking version used in field creation.
func getTypeStruct(rt reflect.Type) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&rt))[1]
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(rt, x)
}

func GetStruct(v interface{}) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(reflect.TypeOf(v), x)
}

func buildStruct(rt reflect.Type, x uintptr) (st *Struct) {
	st = &Struct{Type: rt}
	structMap[x] = st

	st.ByTag = buildTagFields(st.Type, false)
	sort.Slice(st.ByTag, func(i, j int) bool { return 0 > strings.Compare(st.ByTag[i].Key, st.ByTag[j].Key) })
	st.ByName = buildNameFields(st.Type, false)
	sort.Slice(st.ByName, func(i, j int) bool { return 0 > strings.Compare(st.ByName[i].Key, st.ByName[j].Key) })
	st.ByLow = buildLowFields(st.Type, false)
	sort.Slice(st.ByLow, func(i, j int) bool { return 0 > strings.Compare(st.ByLow[i].Key, st.ByLow[j].Key) })

	st.OutTag = buildOutTagFields(st.Type)
	st.OutName = buildOutNameFields(st.Type)
	st.OutLow = buildOutLowFields(st.Type)

	return
}

func buildTagFields(rt reflect.Type, anon bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildTagFields(f.Type, true) {
				fi.index = append([]int{i}, fi.index...)
				fi.offset += f.Offset
				fa = append(fa, fi)
			}
		} else {
			omitEmpty := false
			asString := false
			key := f.Name
			if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
				parts := strings.Split(tag, ",")
				switch parts[0] {
				case "":
					key = f.Name
				case "-":
					if 1 < len(parts) {
						key = "-"
					} else {
						continue
					}
				default:
					key = parts[0]
				}
				for _, p := range parts[1:] {
					switch p {
					case "omitempty":
						omitEmpty = true
					case "string":
						asString = true
					}
				}
			}
			fa = append(fa, newField(f, key, omitEmpty, asString, anon))
		}
	}
	return
}

func buildNameFields(rt reflect.Type, anon bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildNameFields(f.Type, true) {
				fi.index = append([]int{i}, fi.index...)
				fa = append(fa, fi)
			}
		} else {
			fa = append(fa, newField(f, f.Name, false, false, anon))
		}
	}
	return
}

func buildLowFields(rt reflect.Type, anon bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildLowFields(f.Type, true) {
				fi.index = append([]int{i}, fi.index...)
				fa = append(fa, fi)
			}
		} else {
			name[0] = name[0] | 0x20
			fa = append(fa, newField(f, string(name), false, false, anon))
		}
	}
	return
}

func buildOutTagFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		omitEmpty := false
		asString := false
		key := f.Name
		if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
			parts := strings.Split(tag, ",")
			switch parts[0] {
			case "":
				// ok as is
			case "-":
				if 1 < len(parts) {
					key = "-"
				} else {
					continue
				}
			default:
				key = parts[0]
			}
			for _, p := range parts[1:] {
				switch p {
				case "omitempty":
					omitEmpty = true
				case "string":
					asString = true
				}
			}
		}
		fa = append(fa, newField(f, key, omitEmpty, asString, false))
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildOutNameFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		fa = append(fa, newField(f, f.Name, false, false, false))
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildOutLowFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		name[0] = name[0] | 0x20
		fa = append(fa, newField(f, string(name), false, false, false))
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}
