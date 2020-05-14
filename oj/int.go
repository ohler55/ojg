// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"strconv"
	"time"

	"github.com/ohler55/ojg/gen"
)

// Int convert the value provided to an int64. If conversion is not possible
// such if the provided value is an array then the first option default value
// is returned or if not provide 0 is returned. If the type is not one of the
// int or uint types and there is a second optional default then that value is
// returned. This approach keeps the return as a single value and give the
// user the choice of how to indicate a bad value.
func Int(v interface{}, defaults ...int64) (i int64) {
	switch tv := v.(type) {
	case nil:
		if 1 < len(defaults) {
			i = defaults[1]
		}
	case bool:
		if 1 < len(defaults) {
			i = defaults[1]
		} else if tv {
			i = 1
		} else {
			i = 0
		}
	case int64:
		i = tv
	case int:
		i = int64(tv)
	case int8:
		i = int64(tv)
	case int16:
		i = int64(tv)
	case int32:
		i = int64(tv)
	case uint:
		i = int64(tv)
	case uint8:
		i = int64(tv)
	case uint16:
		i = int64(tv)
	case uint32:
		i = int64(tv)
	case uint64:
		i = int64(tv)
	case float32:
		i = int64(tv)
		if float32(i) != tv {
			if 1 < len(defaults) {
				i = defaults[1]
			}
		}
	case float64:
		i = int64(tv)
		if float64(i) != tv {
			if 1 < len(defaults) {
				i = defaults[1]
			}
		}
	case string:
		var err error
		if 1 < len(defaults) {
			i = defaults[1]
		} else if i, err = strconv.ParseInt(tv, 10, 64); err != nil {
			if f, err := strconv.ParseFloat(tv, 64); err == nil {
				i = int64(f)
				if float64(i) != f {
					if 0 < len(defaults) {
						i = defaults[0]
					}
				}
			} else if 0 < len(defaults) {
				i = defaults[0]
			}
		}

	case time.Time:
		if 1 < len(defaults) {
			i = defaults[1]
		} else {
			i = tv.UnixNano()
		}

	case gen.Bool:
		if 1 < len(defaults) {
			i = defaults[1]
		} else if tv {
			i = 1
		} else {
			i = 0
		}
	case gen.Int:
		i = int64(tv)
	case gen.Float:
		i = int64(tv)
		if float64(i) != float64(tv) {
			if 1 < len(defaults) {
				i = defaults[1]
			}
		}
	case gen.String:
		i = Int(string(tv), defaults...)
	case gen.Time:
		if 1 < len(defaults) {
			i = defaults[1]
		} else {
			i = time.Time(tv).UnixNano()
		}
	case gen.Big:
		return Int(string(tv), defaults...)

	default:
		if 0 < len(defaults) {
			i = defaults[0]
		}
	}
	return
}
