// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unsafe"

	"github.com/ohler55/ojg"
)

const (
	maskByTag  = byte(0x01)
	maskExact  = byte(0x02) // exact key vs lowwer case first letter
	maskNested = byte(0x04)
	maskSet    = byte(0x08)
)

// sinfo holds reflect information about a struct.
type sinfo struct {
	rt     reflect.Type
	fields [8][]*finfo
}

var (
	structMut sync.Mutex
	// Keyed by the pointer to the type.
	structMap = map[uintptr]*sinfo{}
)

func (si *sinfo) getFields(o *ojg.Options) []*finfo {
	var index byte
	if o.NestEmbed {
		index |= maskNested
	}
	if o.UseTags {
		index |= maskByTag
	} else if o.KeyExact {
		index |= maskExact
	}
	return si.fields[index]
}

// getSinfo gets the struct information for the provided value. This is use
// internally and is not expected to be used externally.
func getSinfo(v interface{}) (st *sinfo) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	structMut.Lock()
	defer structMut.Unlock()
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(reflect.TypeOf(v), x)
}

func buildStruct(rt reflect.Type, x uintptr) (st *sinfo) {
	st = &sinfo{rt: rt}
	structMap[x] = st

	for u := byte(0); u < maskSet; u++ {
		if (maskByTag&u) != 0 && (maskExact&u) != 0 { // reuse previously built
			st.fields[u] = st.fields[u & ^maskExact]
			continue
		}
		st.fields[u] = buildFields(st.rt, u)
	}
	return
}

func buildFields(rt reflect.Type, u byte) (fa []*finfo) {
	switch {
	case (maskByTag & u) != 0:
		fa = buildTagFields(rt, (maskNested&u) == 0)
	case (maskExact & u) != 0:
		fa = buildExactFields(rt, (maskNested&u) == 0)
	default:
		fa = buildLowFields(rt, (maskNested&u) == 0)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 > strings.Compare(fa[i].key, fa[j].key) })
	return
}

func buildTagFields(rt reflect.Type, nested bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		var fx byte
		if f.Anonymous && nested {
			if f.Type.Kind() == reflect.Ptr {
				for _, fi := range buildTagFields(f.Type.Elem(), nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.value = fi.ivalue
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildTagFields(f.Type, nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.offset += f.Offset
					fa = append(fa, fi)
				}
			}
		} else {
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
						fx |= omitMask
					case "string":
						fx |= strMask
					}
				}
			}
			fa = append(fa, newFinfo(&f, key, fx))
		}
	}
	return
}

func buildExactFields(rt reflect.Type, nested bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		var fx byte
		if f.Anonymous && nested {
			if f.Type.Kind() == reflect.Ptr {
				for _, fi := range buildExactFields(f.Type.Elem(), nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.value = fi.ivalue
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildExactFields(f.Type, nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.offset += f.Offset
					fa = append(fa, fi)
				}
			}
		} else {
			fa = append(fa, newFinfo(&f, f.Name, fx))
		}
	}
	return
}

func buildLowFields(rt reflect.Type, nested bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		var fx byte
		if f.Anonymous && nested {
			if f.Type.Kind() == reflect.Ptr {
				for _, fi := range buildLowFields(f.Type.Elem(), nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.value = fi.ivalue
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildLowFields(f.Type, nested) {
					fi.index = append([]int{i}, fi.index...)
					fi.offset += f.Offset
					fa = append(fa, fi)
				}
			}
		} else {
			if 3 < len(name) {
				if name[0] < 0x80 {
					name[0] |= 0x20
				}
			} else {
				name = bytes.ToLower(name)
			}
			fa = append(fa, newFinfo(&f, string(name), fx))
		}
	}
	return
}
