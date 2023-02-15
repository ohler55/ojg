// Copyright (c) 2023, Peter Ohler, All rights reserved.

package alt

import (
	"strings"
	"time"
)

// Filter is a simple filter for matching against arbitrary date.
type Filter map[string]any

// NewFilter creates a new filter from the spec which whould be a map where
// the keys are simple paths of keys delimited by the dot ('.') character. An
// example is "top.child.grandchild". The matching will either match the key
// when the data is traversed directly or in the case of a slice the elements
// of the slice are also traversed. Generally a Filter is created and reused
// as there is some overhead in creating the Filter. An alternate format is a
// nested set of maps.
func NewFilter(spec map[string]any) Filter {
	f := Filter{}
	for k, v := range spec {
		path := strings.Split(k, ".")
		f2 := f
		for _, k2 := range path[:len(path)-1] {
			sub, _ := f2[k2].(Filter)
			if sub == nil {
				sub = Filter{}
				f2[k2] = sub
			}
			f2 = sub
		}
		// TBD if value is map create sub filter
		switch tv := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			f2[path[len(path)-1]], _ = asInt(tv)
		case float32, float64:
			f2[path[len(path)-1]], _ = asFloat(tv)
		default:
			f2[path[len(path)-1]] = v
		}
	}
	return f
}

// Match returns true if the target matches the Filter.
func (f Filter) Match(target any) bool {
	tm, ok := target.(map[string]any)
	if !ok {
		return false
	}
top:
	for k, fv := range f {
		switch tv := tm[k].(type) {
		case map[string]any:
			if sub, ok := fv.(Filter); ok {
				if !sub.Match(tv) {
					return false
				}
			}
		case []any:
			for _, v := range tv {
				if f.Match(v) {
					continue top
				}
			}
		case nil:
			if fv != nil {
				return false
			}
		case bool:
			if b, ok := fv.(bool); !ok || tv != b {
				return false
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			v, _ := asInt(tv)
			if i, ok := asInt(fv); !ok || v != i {
				return false
			}
		case float32, float64:
			v, _ := asFloat(tv)
			if ff, ok := asFloat(fv); !ok || v != ff {
				return false
			}
		case string:
			if fs, ok := fv.(string); !ok || fs != tv {
				return false
			}
		case time.Time:
			if ft, ok := fv.(time.Time); !ok || !ft.Equal(tv) {
				return false
			}
		default:
			// TBD reflect as map or slice
		}
	}
	return true
}

// Simplify returns a simplified representation of the Filter.
func (f Filter) Simplify() any {
	simple := map[string]any{}
	for k, v := range f {
		switch tv := v.(type) {
		case map[string]any:
			simple[k] = Filter(tv).Simplify()
		case Filter:
			simple[k] = tv.Simplify()
		default:
			simple[k] = tv
		}
	}
	return simple
}
