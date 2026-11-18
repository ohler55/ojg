// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unsafe"
)

const (
	maskByTag  = byte(0x01)
	maskExact  = byte(0x02) // exact key vs lowwer case first letter
	maskNested = byte(0x04)
	maskPretty = byte(0x08)
	maskMax    = byte(0x10)
)

type sinfo struct {
	rt     reflect.Type
	fields [16][]*finfo
}

var (
	structMut sync.Mutex
	// Keyed by the pointer to the type.
	structMap      = map[uintptr]*sinfo{}
	structEmptyMap = map[uintptr]*sinfo{}
)

// Non-locking version used in field creation.
func getTypeStruct(rt reflect.Type, embedded, omitEmpty bool) (st *sinfo) {
	x := (*[2]uintptr)(unsafe.Pointer(&rt))[1]
	if st = structMap[x]; st != nil {
		return
	}
	return buildStruct(rt, x, embedded, omitEmpty)
}

func getSinfo(v any, omitEmpty bool) (st *sinfo) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	sm := structMap
	if omitEmpty {
		sm = structEmptyMap
	}
	structMut.Lock()
	defer structMut.Unlock()
	if st = sm[x]; st != nil {
		return
	}
	return buildStruct(reflect.TypeOf(v), x, false, omitEmpty)
}

func buildStruct(rt reflect.Type, x uintptr, embedded, omitEmpty bool) (st *sinfo) {
	st = &sinfo{rt: rt}
	if omitEmpty {
		structEmptyMap[x] = st
	} else {
		structMap[x] = st
	}
	for u := byte(0); u < maskMax; u++ {
		if (maskByTag&u) != 0 && (maskExact&u) != 0 { // reuse previously built
			st.fields[u] = st.fields[u & ^maskExact]
			continue
		}
		st.fields[u] = buildFields(st.rt, u, embedded, omitEmpty)
	}
	return
}

func buildFields(rt reflect.Type, u byte, embedded, omitEmpty bool) (fa []*finfo) {
	switch {
	case (maskByTag & u) != 0:
		fa = buildTagFields(rt, (maskNested&u) != 0, (maskPretty&u) != 0, embedded, omitEmpty)
	case (maskExact & u) != 0:
		fa = buildExactFields(rt, (maskNested&u) != 0, (maskPretty&u) != 0, embedded, omitEmpty)
	default:
		fa = buildLowFields(rt, (maskNested&u) != 0, (maskPretty&u) != 0, embedded, omitEmpty)
	}
	fa = dominantFields(rt, fa, (maskByTag&u) != 0)
	sort.Slice(fa, func(i, j int) bool { return 0 > strings.Compare(fa[i].key, fa[j].key) })
	return
}

func buildTagFields(rt reflect.Type, out, pretty, embedded, omitEmpty bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		if skipField(&f, !out) {
			continue
		}
		if f.Anonymous && !out && isEmbeddedStruct(f.Type) {
			if f.Type.Kind() == reflect.Pointer {
				for _, fi := range buildTagFields(f.Type.Elem(), out, pretty, embedded, omitEmpty) {
					fi.index = append([]int{i}, fi.index...)
					fi.Append = skipNilEmbed(fi.iAppend)
					fi.iAppend = fi.Append
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildTagFields(f.Type, out, pretty, embedded, omitEmpty) {
					fi.index = append([]int{i}, fi.index...)
					fi.offset += f.Offset
					fa = append(fa, fi)
				}
			}
		} else {
			asString := false
			omit := omitEmpty
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
						omit = true
					case "string":
						asString = true
					}
				}
			}
			fa = append(fa, newFinfo(&f, key, omit, asString, pretty, embedded))
		}
	}
	return
}

func buildExactFields(rt reflect.Type, out, pretty, embedded, omitEmpty bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		if skipField(&f, !out) {
			continue
		}
		if f.Anonymous && !out && isEmbeddedStruct(f.Type) {
			if f.Type.Kind() == reflect.Pointer {
				for _, fi := range buildExactFields(f.Type.Elem(), out, pretty, embedded, omitEmpty) {
					fi.index = append([]int{i}, fi.index...)
					fi.Append = skipNilEmbed(fi.iAppend)
					fi.iAppend = fi.Append
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildExactFields(f.Type, out, pretty, embedded, omitEmpty) {
					fi.index = append([]int{i}, fi.index...)
					fi.offset += f.Offset
					fa = append(fa, fi)
				}
			}
		} else {
			fa = append(fa, newFinfo(&f, f.Name, omitEmpty, false, pretty, embedded))
		}
	}
	return
}

func buildLowFields(rt reflect.Type, out, pretty, embedded, omitEmpty bool) (fa []*finfo) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		if skipField(&f, !out) {
			continue
		}
		name := []byte(f.Name)
		if f.Anonymous && !out && isEmbeddedStruct(f.Type) {
			if f.Type.Kind() == reflect.Pointer {
				for _, fi := range buildLowFields(f.Type.Elem(), out, pretty, embedded, omitEmpty) {
					fi.index = append([]int{i}, fi.index...)
					fi.Append = skipNilEmbed(fi.iAppend)
					fi.iAppend = fi.Append
					fa = append(fa, fi)
				}
			} else {
				for _, fi := range buildLowFields(f.Type, out, pretty, embedded, omitEmpty) {
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
			fa = append(fa, newFinfo(&f, string(name), omitEmpty, false, pretty, embedded))
		}
	}
	return
}

// isEmbeddedStruct returns true if the fields of an embedded field of type rt
// are written in place of the field itself. That is only the case for a
// struct or a pointer to a struct. Any other embedded type, such as
// `type Count int`, is written as a field named after the type, the same as
// encoding/json does.
// skipField reports whether a struct field is left out, the same as for
// encoding/json. An unexported field is left out unless it is an embedded
// struct, or a pointer to one, that is flattened. The exported fields of
// that struct are then used in its place.
func skipField(f *reflect.StructField, flatten bool) bool {
	return !f.IsExported() && !(flatten && f.Anonymous && isEmbeddedStruct(f.Type))
}

func isEmbeddedStruct(rt reflect.Type) bool {
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	return rt.Kind() == reflect.Struct
}

// skipNilEmbed wraps the append function of a field reached through an
// embedded pointer so that the field is skipped when that pointer is nil
// instead of panicking, the same as encoding/json does.
func skipNilEmbed(af appendFunc) appendFunc {
	return func(fi *finfo, buf []byte, rv reflect.Value, addr uintptr, safe bool) ([]byte, any, appendStatus) {
		if _, err := rv.FieldByIndexErr(fi.index); err != nil {
			return buf, nil, aSkip
		}
		return af(fi, buf, rv, addr, safe)
	}
}

// dominantFields removes the fields that are hidden by another field with the
// same key. That happens when a field of an embedded struct has the same key
// as a field of the outer struct or of another embedded struct. The rules are
// the same as for encoding/json. The field with the shortest index (the least
// nested) wins. If more than one field has that depth then a field with a
// json tag name wins when byTag is true and it is the only one. Otherwise all
// the fields with that key are dropped.
func dominantFields(rt reflect.Type, fa []*finfo, byTag bool) []*finfo {
	groups := make(map[string][]*finfo, len(fa))
	for _, fi := range fa {
		groups[fi.key] = append(groups[fi.key], fi)
	}
	if len(groups) == len(fa) {
		return fa
	}
	keep := fa[:0]
	for _, fi := range fa {
		if dominantField(rt, groups[fi.key], byTag) == fi {
			keep = append(keep, fi)
		}
	}
	return keep
}

func dominantField(rt reflect.Type, group []*finfo, byTag bool) (dom *finfo) {
	if len(group) == 1 {
		return group[0]
	}
	depth := len(group[0].index)
	for _, fi := range group[1:] {
		if len(fi.index) < depth {
			depth = len(fi.index)
		}
	}
	var top []*finfo
	for _, fi := range group {
		if len(fi.index) == depth {
			top = append(top, fi)
		}
	}
	if len(top) == 1 {
		return top[0]
	}
	if byTag {
		for _, fi := range top {
			if tag, _, _ := strings.Cut(rt.FieldByIndex(fi.index).Tag.Get("json"), ","); 0 < len(tag) {
				if dom != nil {
					return nil
				}
				dom = fi
			}
		}
	}
	return dom
}
