// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"strconv"
	"time"

	"github.com/ohler55/ojg/gen"
)

// Int convert the value provided to an int64 making every effort to complete
// the conversion. If the value can not be converted zero is returned.
func Int(v interface{}) (i int64) {
	i, _ = AsInt(v)
	return
}

// AsInt convert the value provided to an int64 making every effort to
// complete the conversion. If the value can not be converted zero is
// returned. A status code is returned indicating the conversion was an type
// match, a successful conversion (Ok), or the conversion was not possible
// (Fail).
func AsInt(v interface{}) (i int64, status Status) {
	status = Exact
	switch tv := v.(type) {
	case nil:
		i = 0
		status = Ok
	case bool:
		if tv {
			i = 1
		} else {
			i = 0
		}
		status = Ok
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
		status = Ok
	case float64:
		i = int64(tv)
		status = Ok
	case string:
		status = Fail
		if f, err := strconv.ParseFloat(tv, 64); err == nil {
			i = int64(f)
			if float64(i) == f {
				status = Ok
			}
		} else if i, err = strconv.ParseInt(tv, 10, 64); err == nil {
			status = Ok
		} else {
			status = Fail
		}
	case time.Time:
		i = tv.UnixNano()

	case gen.Bool:
		if tv {
			i = 1
		} else {
			i = 0
		}
		status = Ok
	case gen.Int:
		i = int64(tv)
	case gen.Float:
		i = int64(tv)
	case gen.String:
		status = Fail
		if f, err := strconv.ParseFloat(string(tv), 64); err == nil {
			i = int64(f)
			if float64(i) == f {
				status = Ok
			}
		} else if i, err = strconv.ParseInt(string(tv), 10, 64); err == nil {
			status = Ok
		} else {
			status = Fail
		}
	case gen.Time:
		i = time.Time(tv).UnixNano()
	case gen.Big:
		return AsInt(string(tv))

	default:
		status = Fail
	}
	return
}
