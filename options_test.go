// Copyright (c) 2021, Peter Ohler, All rights reserved.

package ojg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

func TestOptionsAppendTimeMap(t *testing.T) {
	when := time.Date(2021, time.May, 21, 10, 11, 12, 123456789, time.UTC)
	var buf []byte
	o := ojg.Options{TimeMap: true, CreateKey: "^", FullTypePath: true, TimeFormat: "", TimeWrap: ""}

	buf = o.AppendTime(buf, when, false)
	tt.Equal(t, `{"^":"time/Time","value":1621591872123456789}`, string(buf))

	buf = buf[:0]
	buf = o.AppendTime(buf, when, true)
	tt.Equal(t, `{^:"time/Time" value:1621591872123456789}`, string(buf))

	buf = buf[:0]
	o.FullTypePath = false
	buf = o.AppendTime(buf, when, false)
	tt.Equal(t, `{"^":"Time","value":1621591872123456789}`, string(buf))

	buf = buf[:0]
	buf = o.AppendTime(buf, when, true)
	tt.Equal(t, `{^:Time value:1621591872123456789}`, string(buf))
}

func TestOptionsAppendTimeWrap(t *testing.T) {
	when := time.Date(2021, time.May, 21, 10, 11, 12, 123456789, time.UTC)
	var buf []byte
	o := ojg.Options{TimeFormat: time.RFC3339Nano, TimeWrap: "@"}

	buf = o.AppendTime(buf, when, false)
	tt.Equal(t, `{"@":"2021-05-21T10:11:12.123456789Z"}`, string(buf))

	buf = buf[:0]
	buf = o.AppendTime(buf, when, true)
	tt.Equal(t, `{@:"2021-05-21T10:11:12.123456789Z"}`, string(buf))

	o.TimeFormat = "second"
	buf = buf[:0]
	buf = o.AppendTime(buf, when, false)
	tt.Equal(t, `{"@":1621591872.123456789}`, string(buf))

	buf = buf[:0]
	buf = o.AppendTime(buf, when, true)
	tt.Equal(t, `{@:1621591872.123456789}`, string(buf))
}

func TestOptionsAppendTimeSecond(t *testing.T) {
	when := time.Date(2021, time.May, 21, 10, 11, 12, 123456789, time.UTC)
	var buf []byte
	o := ojg.Options{TimeFormat: "second"}

	buf = o.AppendTime(buf, when, false)
	tt.Equal(t, `1621591872.123456789`, string(buf))

	buf = buf[:0]
	when = time.Date(1954, time.May, 21, 10, 11, 12, 123456789, time.UTC)
	buf = o.AppendTime(buf, when, true)
	tt.Equal(t, `-492788927.876543211`, string(buf))
}

func TestOptionsDecomposeTime(t *testing.T) {
	when := time.Date(2021, time.May, 21, 10, 11, 12, 123456789, time.UTC)
	o := ojg.Options{TimeFormat: "time"}
	v := o.DecomposeTime(when)
	_, ok := v.(time.Time)
	tt.Equal(t, true, ok, fmt.Sprintf("%T", v))

	o.TimeFormat = "nano"
	v = o.DecomposeTime(when)
	i, _ := v.(int64)
	tt.Equal(t, int64(1621591872123456789), i)

	o.TimeFormat = "second"
	v = o.DecomposeTime(when)
	f, _ := v.(float64)
	tt.Equal(t, float64(1621591872123456789)/float64(time.Second), f)

	o.TimeFormat = time.RFC3339Nano
	v = o.DecomposeTime(when)
	s, _ := v.(string)
	tt.Equal(t, "2021-05-21T10:11:12.123456789Z", s)

	o.TimeMap = true
	o.CreateKey = "^"
	v = o.DecomposeTime(when)
	m, _ := v.(map[string]interface{})
	tt.Equal(t, map[string]interface{}{"^": "Time", "value": "2021-05-21T10:11:12.123456789Z"}, m)

	o.FullTypePath = true
	v = o.DecomposeTime(when)
	m, _ = v.(map[string]interface{})
	tt.Equal(t, map[string]interface{}{"^": "time/Time", "value": "2021-05-21T10:11:12.123456789Z"}, m)

	o.TimeMap = false
	o.TimeWrap = "@"
	v = o.DecomposeTime(when)
	m, _ = v.(map[string]interface{})
	tt.Equal(t, map[string]interface{}{"@": "2021-05-21T10:11:12.123456789Z"}, m)

}
