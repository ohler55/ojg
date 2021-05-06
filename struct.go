// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unsafe"
)

const (
	MaskByTag  = byte(0x10)
	MaskExact  = byte(0x08) // exact key vs lowwer case first letter
	MaskPretty = byte(0x04)
	MaskNested = byte(0x02)
	MaskSen    = byte(0x01)
	MaskSet    = byte(0x20)
	MaskIndex  = byte(0x1f)
)

// Struct holds reflect information about a struct.
type Struct struct {
	Type   reflect.Type
	Fields [32][]*Field
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

	for u := byte(0); u < MaskSet; u++ {
		if (MaskByTag&u) != 0 && (MaskExact&u) != 0 { // reuse previously built
			st.Fields[u] = st.Fields[u & ^MaskExact]
			continue
		}
		st.Fields[u] = buildFields(st.Type, u)
	}
	return
}

func buildFields(rt reflect.Type, u byte) (fa []*Field) {
	if (MaskByTag & u) != 0 {
		fa = buildTagFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0)
	} else if (MaskExact & u) != 0 {
		fa = buildExactFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0)
	} else {
		fa = buildLowFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 > strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildTagFields(rt reflect.Type, out, pretty, sen bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildTagFields(f.Type, out, pretty, sen) {
				fi.Index = append([]int{i}, fi.Index...)
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
			fa = append(fa, newField(f, key, omitEmpty, asString, pretty, sen))
		}
	}
	return
}

func buildExactFields(rt reflect.Type, out, pretty, sen bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildExactFields(f.Type, out, pretty, sen) {
				fi.Index = append([]int{i}, fi.Index...)
				fi.offset += f.Offset
				fa = append(fa, fi)
			}
		} else {
			fa = append(fa, newField(f, f.Name, false, false, pretty, sen))
		}
	}
	return
}

func buildLowFields(rt reflect.Type, out, pretty, sen bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildLowFields(f.Type, out, pretty, sen) {
				fi.Index = append([]int{i}, fi.Index...)
				fi.offset += f.Offset
				fa = append(fa, fi)
			}
		} else {
			if 3 < len(name) {
				if name[0] < 0x80 {
					name[0] = name[0] | 0x20
				}
			} else {
				name = bytes.ToLower(name)
			}
			fa = append(fa, newField(f, string(name), false, false, pretty, sen))
		}
	}
	return
}
