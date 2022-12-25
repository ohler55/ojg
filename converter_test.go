// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestConverterRFC3339(t *testing.T) {
	val := []any{
		"2021-03-05T10:11:12Z",
		"2021-03-05T10:11:12.123Z",
		"2021-03-05T10:11:12.123456789-05:00",
		"2021-03-05",
		"2021-03-05T10:11:12",                  // too short
		"2021-03-05T10:11:12.1234567890-05:00", // too long
		"2021-03-05 10:11:12.123Z",             // wrong format
	}
	v2, _ := ojg.TimeRFC3339Converter.Convert(val).([]any)
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		var ok bool
		if 4 <= i { // should not be converted
			_, ok = val[i].(string)
		} else {
			_, ok = val[i].(time.Time)
		}
		tt.Equal(t, true, ok, i, ":", val[i])
	}
}

func TestConverterNanoTime(t *testing.T) {
	val := []any{
		int64(946684800000000000),
		int64(946684800000000001),
		int64(1609804800000000000), // 2021-01-05
		uint64(946684800000000001),
		uint(946684800000000001),
		int(946684800000000001),

		int64(946684799999999999),
		int32(12345),
		int16(1234),
		int8(123),
		uint32(12345),
		uint16(1234),
		uint8(123),
		nil,
	}
	v2, _ := ojg.TimeNanoConverter.Convert(val).([]any)
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		_, ok := val[i].(time.Time)
		if 6 <= i { // should not be converted
			tt.Equal(t, false, ok, i, ":", val[i])
		} else {
			tt.Equal(t, true, ok, i, ":", val[i])
		}
	}
	vm := map[string]any{"x": int(946684800000000001)}
	_ = ojg.TimeNanoConverter.Convert(vm)
	_, ok := vm["x"].(time.Time)
	tt.Equal(t, true, ok)
}

func TestConverterFloat(t *testing.T) {
	fun := func(val float64) (any, bool) {
		if 946684800.0 <= val { // 2000-01-01
			secs := int64(val)
			return time.Unix(secs, int64(val*1000000000.0)-secs*1000000000), true
		}
		return val, false
	}
	val := []any{
		1609804800.000000000, // 2021-01-05
		float32(1609804800.0),
		123456789.123,
	}
	v2, _ := ojg.Convert(val, fun).([]any)
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		var ok bool
		if 2 <= i { // should not be converted
			_, ok = val[i].(float64)
		} else {
			_, ok = val[i].(time.Time)
		}
		tt.Equal(t, true, ok, i, ":", val[i])
	}
}

func TestConverterArray(t *testing.T) {
	fun := func(val []any) (any, bool) {
		if len(val) == 2 {
			if k, _ := val[0].(string); 0 < len(k) {
				return map[string]any{k: val[1]}, true
			}
		}
		return val, false
	}
	val := []any{
		[]any{"x", 2},
		[]any{1, 2},
		[]any{"x", 2, 3},
	}
	v2, _ := ojg.Convert(val, fun).([]any)
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		var ok bool
		if 1 <= i { // should not be converted
			_, ok = val[i].([]any)
		} else {
			_, ok = val[i].(map[string]any)
		}
		tt.Equal(t, true, ok, i, ":", val[i])
	}
}

func TestConverterMixed(t *testing.T) {
	val := []any{1, true, "ab"}
	v2, _ := ojg.Convert(val,
		func(val int64) (any, bool) { return val + 1, true },
		func(val string) (any, bool) { return val + "c", true },
		func(val map[string]any) (any, bool) { return true, true },
	).([]any)
	tt.Equal(t, []any{2, true, "abc"}, v2)
}

func TestConverterMongo(t *testing.T) {
	val := []any{
		map[string]any{"$oid": "507f191e810c19729de860ea"},
		map[string]any{"$date": "2021-03-05T11:22:33.123Z"},
		map[string]any{"$numberLong": "123456789"},
		map[string]any{"$numberDecimal": "123.456"},
		map[string]any{"$numberDecimal": "123.456", "x": 3},
		map[string]any{"$numberDecimal": 3},
	}
	v2, _ := ojg.MongoConverter.Convert(val).([]any)
	tt.Equal(t, 6, len(v2))
	tt.Equal(t, "507f191e810c19729de860ea", v2[0])
	tt.Equal(t, "time.Time 2021-03-05 11:22:33.123 +0000 UTC", fmt.Sprintf("%T %s", v2[1], v2[1]))
	tt.Equal(t, 123456789, v2[2])
	tt.Equal(t, 123.456, v2[3])
	tt.Equal(t, map[string]any{"$numberDecimal": "123.456", "x": 3}, v2[4])
	tt.Equal(t, map[string]any{"$numberDecimal": 3}, v2[5])
}
