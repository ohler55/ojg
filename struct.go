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
	// MaskByTag is the mask for byTag fields.
	MaskByTag = byte(0x10)
	// MaskExact is the mask for Exact fields.
	MaskExact = byte(0x08) // exact key vs lowwer case first letter
	// MaskPretty is the mask for Pretty fields.
	MaskPretty = byte(0x04)
	// MaskNested is the mask for Nested fields.
	MaskNested = byte(0x02)
	// MaskSen is the mask for Sen fields.
	MaskSen = byte(0x01)
	// MaskSet is the mask for Set fields.
	MaskSet = byte(0x20)
	// MaskIndex is the mask for an index that has been set up.
	MaskIndex = byte(0x1f)
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

// GetTypeStruct gets the struct information about the reflect type. This is
// used internally and is not expected to be used externally.
func GetTypeStruct(rt reflect.Type) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&rt))[1]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(rt, x, false)
}

// Non-locking version used in field creation.
func getTypeStruct(rt reflect.Type, embedded bool) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&rt))[1]
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(rt, x, embedded)
}

// GetStruct gets the struct information for the provided value. This is use
// internally and is not expected to be used externally.
func GetStruct(v interface{}) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(reflect.TypeOf(v), x, false)
}

func buildStruct(rt reflect.Type, x uintptr, embedded bool) (st *Struct) {
	st = &Struct{Type: rt}
	structMap[x] = st

	// TBD create value of type rt to use for addressable and same offset
	//
	//rv := reflect.New(rt)
	for u := byte(0); u < MaskSet; u++ {
		if (MaskByTag&u) != 0 && (MaskExact&u) != 0 { // reuse previously built
			st.Fields[u] = st.Fields[u & ^MaskExact]
			continue
		}
		st.Fields[u] = buildFields(st.Type, u, embedded)
	}
	return
}

func buildFields(rt reflect.Type, u byte, embedded bool) (fa []*Field) {
	if (MaskByTag & u) != 0 {
		fa = buildTagFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0, embedded)
	} else if (MaskExact & u) != 0 {
		fa = buildExactFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0, embedded)
	} else {
		fa = buildLowFields(rt, (MaskNested&u) != 0, (MaskPretty&u) != 0, (MaskSen&u) != 0, embedded)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 > strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildTagFields(rt reflect.Type, out, pretty, sen, embedded bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildTagFields(f.Type, out, pretty, sen, embedded) {
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
			// TBD check field lookup vs offset (maybe vs can address)
			//  change arg from embedded to canAddr or direct+
			fa = append(fa, newField(f, key, omitEmpty, asString, pretty, sen, embedded))
		}
	}
	return
}

func buildExactFields(rt reflect.Type, out, pretty, sen, embedded bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildExactFields(f.Type, out, pretty, sen, embedded) {
				fi.Index = append([]int{i}, fi.Index...)
				fi.offset += f.Offset
				fa = append(fa, fi)
			}
		} else {
			fa = append(fa, newField(f, f.Name, false, false, pretty, sen, embedded))
		}
	}
	return
}

func buildLowFields(rt reflect.Type, out, pretty, sen, embedded bool) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous && !out {
			for _, fi := range buildLowFields(f.Type, out, pretty, sen, embedded) {
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
			fa = append(fa, newField(f, string(name), false, false, pretty, sen, embedded))
		}
	}
	return
}
