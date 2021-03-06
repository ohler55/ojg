// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/tt"
)

func TestConverterRFC3339(t *testing.T) {
	val := []interface{}{
		"2021-03-05T10:11:12Z",
		"2021-03-05T10:11:12.123Z",
		"2021-03-05T10:11:12.123456789-05:00",
		"2021-03-05",
		"2021-03-05T10:11:12",                  // too short
		"2021-03-05T10:11:12.1234567890-05:00", // too long
		"2021-03-05 10:11:12.123Z",             // wrong format
	}
	v2, _ := alt.TimeRFC3339Converter.Convert(val).([]interface{})
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
	val := []interface{}{
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
	v2, _ := alt.TimeNanoConverter.Convert(val).([]interface{})
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		_, ok := val[i].(time.Time)
		if 6 <= i { // should not be converted
			tt.Equal(t, false, ok, i, ":", val[i])
		} else {
			tt.Equal(t, true, ok, i, ":", val[i])
		}
	}
	vm := map[string]interface{}{"x": int(946684800000000001)}
	_ = alt.TimeNanoConverter.Convert(vm)
	_, ok := vm["x"].(time.Time)
	tt.Equal(t, true, ok)
}

func TestConverterFloat(t *testing.T) {
	fun := func(val float64) (interface{}, bool) {
		if 946684800.0 <= val { // 2000-01-01
			secs := int64(val)
			return time.Unix(secs, int64(val*1000000000.0)-secs*1000000000), true
		}
		return val, false
	}
	val := []interface{}{
		1609804800.000000000, // 2021-01-05
		float32(1609804800.0),
		123456789.123,
	}
	v2, _ := alt.Convert(val, fun).([]interface{})
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
	fun := func(val []interface{}) (interface{}, bool) {
		if len(val) == 2 {
			if k, _ := val[0].(string); 0 < len(k) {
				return map[string]interface{}{k: val[1]}, true
			}
		}
		return val, false
	}
	val := []interface{}{
		[]interface{}{"x", 2},
		[]interface{}{1, 2},
		[]interface{}{"x", 2, 3},
	}
	v2, _ := alt.Convert(val, fun).([]interface{})
	for i := 0; i < len(val); i++ {
		tt.Equal(t, val[i], v2[i]) // verify they are the same
		var ok bool
		if 1 <= i { // should not be converted
			_, ok = val[i].([]interface{})
		} else {
			_, ok = val[i].(map[string]interface{})
		}
		tt.Equal(t, true, ok, i, ":", val[i])
	}
}

func TestConverterMixed(t *testing.T) {
	val := []interface{}{1, true, "ab"}
	v2, _ := alt.Convert(val,
		func(val int64) (interface{}, bool) { return val + 1, true },
		func(val string) (interface{}, bool) { return val + "c", true },
		func(val map[string]interface{}) (interface{}, bool) { return true, true },
	).([]interface{})
	tt.Equal(t, []interface{}{2, true, "abc"}, v2)
}

func TestConverterMongo(t *testing.T) {
	val := []interface{}{
		map[string]interface{}{"$oid": "507f191e810c19729de860ea"},
		map[string]interface{}{"$date": "2021-03-05T11:22:33.123Z"},
		map[string]interface{}{"$numberLong": "123456789"},
		map[string]interface{}{"$numberDecimal": "123.456"},
		map[string]interface{}{"$numberDecimal": "123.456", "x": 3},
		map[string]interface{}{"$numberDecimal": 3},
	}
	v2, _ := alt.MongoConverter.Convert(val).([]interface{})
	tt.Equal(t, 6, len(v2))
	tt.Equal(t, "507f191e810c19729de860ea", v2[0])
	tt.Equal(t, "time.Time 2021-03-05 11:22:33.123 +0000 UTC", fmt.Sprintf("%T %s", v2[1], v2[1]))
	tt.Equal(t, 123456789, v2[2])
	tt.Equal(t, 123.456, v2[3])
	tt.Equal(t, map[string]interface{}{"$numberDecimal": "123.456", "x": 3}, v2[4])
	tt.Equal(t, map[string]interface{}{"$numberDecimal": 3}, v2[5])
}
