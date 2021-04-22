// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// Field hold information about a struct field.
type Field struct {
	Type      reflect.Type
	Key       string
	Kind      reflect.Kind
	Index     []int
	OmitEmpty bool
	String    bool
}

// Decomposer holds reflect information about a struct.
type Decomposer struct {
	Name    string
	Full    string
	Type    reflect.Type
	ByTag   []*Field
	ByName  []*Field
	ByLow   []*Field
	OutTag  []*Field
	OutName []*Field
	OutLow  []*Field
}

var (
	decompMut sync.Mutex
	decompMap = map[string]*Decomposer{}
)

func LookupDecomposer(rt reflect.Type) (dc *Decomposer) {
	decompMut.Lock()
	defer decompMut.Unlock()
	name := rt.Name()
	full := rt.PkgPath() + "/" + name
	if dc = decompMap[full]; dc != nil {
		return
	}
	dc = &Decomposer{Type: rt, Name: name, Full: full}
	dc.ByTag = buildTagFields(dc.Type)
	sort.Slice(dc.ByTag, func(i, j int) bool { return 0 < strings.Compare(dc.ByTag[i].Key, dc.ByTag[j].Key) })
	dc.ByName = buildNameFields(dc.Type)
	sort.Slice(dc.ByName, func(i, j int) bool { return 0 < strings.Compare(dc.ByName[i].Key, dc.ByName[j].Key) })
	dc.ByLow = buildLowFields(dc.Type)
	sort.Slice(dc.ByLow, func(i, j int) bool { return 0 < strings.Compare(dc.ByLow[i].Key, dc.ByLow[j].Key) })

	dc.OutTag = buildOutTagFields(dc.Type)
	dc.OutName = buildOutNameFields(dc.Type)
	dc.OutLow = buildOutLowFields(dc.Type)

	decompMap[full] = dc

	return
}

// GetDecomposer gets the Decomposer for a value type.
func GetDecomposer(v interface{}) (dc *Decomposer) {
	return LookupDecomposer(reflect.TypeOf(v))
}

func (fi *Field) Value(rv reflect.Value, opt *Options) (v interface{}, omit bool) {
	fv := rv.FieldByIndex(fi.Index)
	v = fv.Interface()
	if fi.OmitEmpty {
		switch tv := v.(type) {
		case nil:
			omit = true
		case bool:
			omit = !tv
		case string:
			omit = len(tv) == 0
		case []byte:
			omit = len(tv) == 0
		case time.Time:
			omit = tv.IsZero()
		case int:
			omit = tv == 0
		case int8:
			omit = tv == 0
		case int16:
			omit = tv == 0
		case int32:
			omit = tv == 0
		case int64:
			omit = tv == 0
		case uint:
			omit = tv == 0
		case uint8:
			omit = tv == 0
		case uint16:
			omit = tv == 0
		case uint32:
			omit = tv == 0
		case uint64:
			omit = tv == 0
		case float32:
			omit = tv == 0.0
		case float64:
			omit = tv == 0.0
		default:
			switch fi.Kind {
			case reflect.Map, reflect.Slice, reflect.Array:
				omit = reflect.ValueOf(v).IsZero()
			}
		}
	}
	if opt.OmitNil {
		switch v.(type) {
		case nil:
			omit = true
		default:
			switch fi.Kind {
			case reflect.Map, reflect.Slice, reflect.Array:
				omit = reflect.ValueOf(v).IsZero()
			}
		}
	}
	if fi.String {
		v = fmt.Sprintf("%v", v)
	}
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
				fi.Index = append([]int{i}, fi.Index...)
				fa = append(fa, fi)
			}
		} else {
			fi := Field{
				Type:  f.Type,
				Key:   f.Name,
				Index: f.Index,
				Kind:  f.Type.Kind(),
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
						fi.OmitEmpty = true
					case "string":
						fi.String = true
					}
				}
			}
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
				fi.Index = append([]int{i}, fi.Index...)
				fa = append(fa, fi)
			}
		} else {
			fa = append(fa, &Field{
				Type:  f.Type,
				Key:   f.Name,
				Index: f.Index,
				Kind:  f.Type.Kind(),
			})
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
				fi.Index = append([]int{i}, fi.Index...)
				fa = append(fa, fi)
			}
		} else {
			name[0] = name[0] | 0x20
			fa = append(fa, &Field{
				Type:  f.Type,
				Key:   string(name),
				Index: f.Index,
				Kind:  f.Type.Kind(),
			})
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
			Type:  f.Type,
			Key:   f.Name,
			Index: f.Index,
			Kind:  f.Type.Kind(),
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
					fi.OmitEmpty = true
				case "string":
					fi.String = true
				}
			}
		}
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
			Type:  f.Type,
			Key:   f.Name,
			Index: f.Index,
			Kind:  f.Type.Kind(),
		}
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
			Type:  f.Type,
			Key:   string(name),
			Index: f.Index,
			Kind:  f.Type.Kind(),
		}
		fa = append(fa, &fi)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}
