// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"reflect"
	"strings"
)

type composer struct {
	fun     RecomposeFunc
	any     RecomposeAnyFunc
	short   string
	full    string
	rtype   reflect.Type
	indexes map[string]reflect.StructField
}

func indexType(rt reflect.Type) (im map[string]reflect.StructField) {
	if rt.NumField() == 0 {
		return
	}
	cands := map[string][]reflect.StructField{}
	collectIndexes(rt, nil, cands)
	im = map[string]reflect.StructField{}
	for k, fa := range cands {
		if f, ok := dominantStructField(fa); ok {
			im[k] = f
		}
	}
	return
}

// collectIndexes adds every field of rt, including the fields of embedded
// structs, to cands under the key the field is matched with.
func collectIndexes(rt reflect.Type, index []int, cands map[string][]reflect.StructField) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		if skipField(&f, true) {
			continue
		}
		f.Index = append(append([]int{}, index...), i)
		if f.Anonymous && isEmbeddedStruct(f.Type) {
			ft := f.Type
			if ft.Kind() == reflect.Pointer {
				ft = ft.Elem()
			}
			collectIndexes(ft, f.Index, cands)
			continue
		}
		k := f.Name
		if tag, _ := f.Tag.Lookup("json"); 0 < len(tag) {
			parts := strings.Split(tag, ",")
			switch parts[0] {
			case "":
				k = strings.ToLower(f.Name)
			case "-":
				if 1 < len(parts) {
					k = "-"
				} else {
					continue
				}
			default:
				k = parts[0]
			}
		}
		cands[k] = append(cands[k], f)
	}
}

// dominantStructField picks the field used for a key that more than one
// field maps to, such as a field of the outer struct and a field of an
// embedded struct with the same name. The rules are the same as for
// encoding/json. The least nested field wins. If more than one field has
// that depth then the only one with a json tag name wins. Otherwise no field
// is used for the key.
func dominantStructField(fa []reflect.StructField) (dom reflect.StructField, ok bool) {
	depth := len(fa[0].Index)
	for _, f := range fa[1:] {
		if len(f.Index) < depth {
			depth = len(f.Index)
		}
	}
	var top []reflect.StructField
	for _, f := range fa {
		if len(f.Index) == depth {
			top = append(top, f)
		}
	}
	if len(top) == 1 {
		return top[0], true
	}
	for _, f := range top {
		if tag, _, _ := strings.Cut(f.Tag.Get("json"), ","); 0 < len(tag) {
			if ok {
				return dom, false
			}
			dom, ok = f, true
		}
	}
	return
}
