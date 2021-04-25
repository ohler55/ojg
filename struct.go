// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg

import (
	"reflect"
	"sort"
	"strings"
	"sync"
	"unsafe"
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

func GetStruct(v interface{}) (st *Struct) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	rt := reflect.TypeOf(v)
	st = &Struct{Type: rt}
	st.ByTag = buildTagFields(st.Type)
	sort.Slice(st.ByTag, func(i, j int) bool { return 0 < strings.Compare(st.ByTag[i].Key, st.ByTag[j].Key) })
	st.ByName = buildNameFields(st.Type)
	sort.Slice(st.ByName, func(i, j int) bool { return 0 < strings.Compare(st.ByName[i].Key, st.ByName[j].Key) })
	st.ByLow = buildLowFields(st.Type)
	sort.Slice(st.ByLow, func(i, j int) bool { return 0 < strings.Compare(st.ByLow[i].Key, st.ByLow[j].Key) })

	st.OutTag = buildOutTagFields(st.Type)
	st.OutName = buildOutNameFields(st.Type)
	st.OutLow = buildOutLowFields(st.Type)

	structMap[x] = st

	return
}

func buildTagFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildTagFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			fi := Field{
				Type:   f.Type,
				Key:    f.Name,
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
				parts := strings.Split(tag, ",")
				switch parts[0] {
				case "":
					fi.Key = f.Name
				case "-":
					if 1 < len(parts) {
						fi.Key = "-"
					} else {
						continue
					}
				default:
					fi.Key = parts[0]
				}
				for _, p := range parts[1:] {
					switch p {
					case "omitempty":
						fi.setOmitEmpty()
					case "string":
						fi.asString = true
					}
				}
			}
			fi.setup()
			fa = append(fa, &fi)
		}
	}
	return
}

func buildNameFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildNameFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			fi := Field{
				Type:   f.Type,
				Key:    f.Name,
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			fi.setup()
			fa = append(fa, &fi)
		}
	}
	return
}

func buildLowFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildLowFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			name[0] = name[0] | 0x20
			fi := Field{
				Type:   f.Type,
				Key:    string(name),
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			fi.setup()
			fa = append(fa, &fi)
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
		fi := Field{
			Type:   f.Type,
			Key:    f.Name,
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
			parts := strings.Split(tag, ",")
			switch parts[0] {
			case "":
				// ok as is
			case "-":
				if 1 < len(parts) {
					fi.Key = "-"
				} else {
					continue
				}
			default:
				fi.Key = parts[0]
			}
			for _, p := range parts[1:] {
				switch p {
				case "omitempty":
					fi.setOmitEmpty()
				case "string":
					fi.asString = true
				}
			}
		}
		fi.setup()
		fa = append(fa, &fi)
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
		fi := Field{
			Type:   f.Type,
			Key:    f.Name,
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		fi.setup()
		fa = append(fa, &fi)
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
		fi := Field{
			Type:   f.Type,
			Key:    string(name),
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		fi.setup()
		fa = append(fa, &fi)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}
