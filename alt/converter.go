// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt

import (
	"math"
	"strconv"
	"time"
)

type Converter struct {
	Int    []func(val int64) (interface{}, bool)
	Float  []func(val float64) (interface{}, bool)
	String []func(val string) (interface{}, bool)
	Map    []func(val map[string]interface{}) (interface{}, bool)
	Array  []func(val []interface{}) (interface{}, bool)
}

// if map or slice then orig is not changed
func (c *Converter) Convert(v interface{}) interface{} {
	v, _ = c.convert(v)
	return v
}

func (c *Converter) convert(v interface{}) (interface{}, bool) {
	switch tv := v.(type) {
	case int64:
		for _, fun := range c.Int {
			if cv, ok := fun(tv); ok {
				return cv, true
			}
		}
	case float64:
		for _, fun := range c.Float {
			if cv, ok := fun(tv); ok {
				return cv, true
			}
		}
	case string:
		for _, fun := range c.String {
			if cv, ok := fun(tv); ok {
				return cv, true
			}
		}
	case []interface{}:
		for _, fun := range c.Array {
			if cv, ok := fun(tv); ok {
				return cv, true
			}
		}
		for i, m := range tv {
			if cv, ok := c.convert(m); ok {
				tv[i] = cv
			}
		}
	case map[string]interface{}:
		for _, fun := range c.Map {
			if cv, ok := fun(tv); ok {
				return cv, true
			}
		}
		for k, m := range tv {
			if cv, ok := c.convert(m); ok {
				tv[k] = cv
			}
		}

	case int:
		return c.convert(int64(tv))
	case int8:
		return c.convert(int64(tv))
	case int16:
		return c.convert(int64(tv))
	case int32:
		return c.convert(int64(tv))
	case uint:
		return c.convert(int64(tv))
	case uint8:
		return c.convert(int64(tv))
	case uint16:
		return c.convert(int64(tv))
	case uint32:
		return c.convert(int64(tv))
	case uint64:
		return c.convert(int64(tv))
	case float32:
		// This small rounding makes the conversion from 32 bit to 64 bit
		// display nicer.
		f, i := math.Frexp(float64(tv))
		f = float64(int64(f*fracMax)) / fracMax
		return c.convert(math.Ldexp(f, i))
	}
	return v, false
}

func Convert(v interface{}, funcs ...interface{}) interface{} {
	c := Converter{}
	for _, fun := range funcs {
		switch tf := fun.(type) {
		case func(val int64) (interface{}, bool):
			c.Int = append(c.Int, tf)
		case func(val float64) (interface{}, bool):
			c.Float = append(c.Float, tf)
		case func(val string) (interface{}, bool):
			c.String = append(c.String, tf)
		case func(val map[string]interface{}) (interface{}, bool):
			c.Map = append(c.Map, tf)
		case func(val []interface{}) (interface{}, bool):
			c.Array = append(c.Array, tf)
		}
	}
	v, _ = c.convert(v)

	return v
}

var (
	TimeRFC3339Converter = Converter{
		String: []func(val string) (interface{}, bool){
			func(val string) (interface{}, bool) {
				if 20 <= len(val) && len(val) <= 37 {
					for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
						if t, err := time.ParseInLocation(layout, val, time.UTC); err == nil {
							return t, true
						}
					}
				}
				return val, false
			},
		},
	}

	TimeNanoConverter = Converter{
		Int: []func(val int64) (interface{}, bool){
			func(val int64) (interface{}, bool) {
				if 946684800000000000 <= val { // 2000-01-01
					return time.Unix(0, val), true
				}
				return val, false
			},
		},
	}

	MongoConverter = Converter{
		Map: []func(val map[string]interface{}) (interface{}, bool){
			func(val map[string]interface{}) (interface{}, bool) {
				if len(val) != 1 {
					return val, false
				}
				for k, v := range val {
					s, ok := v.(string)
					if !ok {
						break
					}
					switch k {
					case "$numberLong":
						if i, err := strconv.ParseInt(s, 10, 64); err == nil {
							return i, true
						}
					case "$date":
						if t, err := time.ParseInLocation("2006-01-02T15:04:05.999Z07:00", s, time.UTC); err == nil {
							return t, true
						}
					case "$numberDecimal":
						if f, err := strconv.ParseFloat(s, 64); err == nil {
							return f, true
						}
					case "$oid":
						return s, true
					}
				}
				return val, false
			},
		},
	}
)
