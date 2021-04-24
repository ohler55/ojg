// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

const (
	spaces = "\n                                                                                                                                "
	tabs   = "\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
	hex    = "0123456789abcdef"
)

// JSON returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an
// int as an indent or a *Options.
func JSON(data interface{}, args ...interface{}) string {
	o := &DefaultOptions

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			oi := *o
			oi.Indent = ta
			o = &oi
		case *Options:
			o = ta
		}
	}
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.buf = o.buf[:0]
		}
	}()
	o.buildJSON(data, 0, false)

	return string(o.buf)
}

// Marshal returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent or a *Options. An error will be returned if the Option.Strict
// flag is true and a value is encountered that can not be encoded other than
// by using the %v format of the fmt package.
func Marshal(data interface{}, args ...interface{}) (out []byte, err error) {
	o := &GoOptions
	o.KeyExact = true
	o.UseTags = true

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			oi := *o
			oi.Indent = ta
			o = &oi
		case *Options:
			o = ta
		}
	}
	o.strict = true
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.buf = o.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	o.buildJSON(data, 0, false)
	out = o.buf

	return
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	o := &DefaultOptions

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			oi := *o
			oi.Indent = ta
			o = &oi
		case *Options:
			o = ta
		}
	}
	o.w = w
	if o.InitSize == 0 {
		o.InitSize = 256
	}
	if o.WriteLimit == 0 {
		o.WriteLimit = 1024
	}
	if cap(o.buf) < o.InitSize {
		o.buf = make([]byte, 0, o.InitSize)
	} else {
		o.buf = o.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			o.buf = o.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if o.Color {
		o.cbuildJSON(data, 0) // TBD embedded
	} else {
		o.buildJSON(data, 0, false)
	}
	if err == nil && w != nil && 0 < len(o.buf) {
		_, err = o.w.Write(o.buf)
	}
	return
}

func (o *Options) buildJSON(data interface{}, depth int, embedded bool) {
	switch td := data.(type) {
	case nil:
		o.buf = append(o.buf, []byte("null")...)

	case bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}
	case gen.Bool:
		if td {
			o.buf = append(o.buf, []byte("true")...)
		} else {
			o.buf = append(o.buf, []byte("false")...)
		}

	case int:
		o.buf = strconv.AppendInt(o.buf, int64(td), 10)
	case int8:
		o.buf = strconv.AppendInt(o.buf, int64(td), 10)
	case int16:
		o.buf = strconv.AppendInt(o.buf, int64(td), 10)
	case int32:
		o.buf = strconv.AppendInt(o.buf, int64(td), 10)
	case int64:
		o.buf = strconv.AppendInt(o.buf, td, 10)
	case uint:
		o.buf = strconv.AppendUint(o.buf, uint64(td), 10)
	case uint8:
		o.buf = strconv.AppendUint(o.buf, uint64(td), 10)
	case uint16:
		o.buf = strconv.AppendUint(o.buf, uint64(td), 10)
	case uint32:
		o.buf = strconv.AppendUint(o.buf, uint64(td), 10)
	case uint64:
		o.buf = strconv.AppendUint(o.buf, td, 10)
	case gen.Int:
		o.buf = strconv.AppendInt(o.buf, int64(td), 10)

	case float32:
		o.buf = strconv.AppendFloat(o.buf, float64(td), 'g', -1, 32)
	case float64:
		o.buf = strconv.AppendFloat(o.buf, float64(td), 'g', -1, 64)
	case gen.Float:
		o.buf = strconv.AppendFloat(o.buf, float64(td), 'g', -1, 64)

	case string:
		o.buildString(td)
	case gen.String:
		o.buildString(string(td))

	case time.Time:
		o.buildTime(td)
	case gen.Time:
		o.buildTime(time.Time(td))

	case []interface{}:
		o.buildSimpleArray(td, depth)
	case gen.Array:
		o.buildArray(td, depth)
	case []gen.Node:
		o.buildArray(gen.Array(td), depth)

	case map[string]interface{}:
		o.buildSimpleObject(td, depth)
	case gen.Object:
		o.buildObject(td, depth)

	default:
		if g, _ := data.(alt.Genericer); g != nil {
			o.buildJSON(g.Generic(), depth, false)
			return
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			o.buildJSON(data, depth, false)
			return
		}
		if 0 < len(o.CreateKey) {
			ao := alt.Options{
				CreateKey:    o.CreateKey,
				OmitNil:      o.OmitNil,
				FullTypePath: o.FullTypePath,
				UseTags:      o.UseTags,
				KeyExact:     o.KeyExact,
				NestEmbed:    o.NestEmbed,
				BytesAs:      o.BytesAs,
			}
			o.buildJSON(alt.Decompose(data, &ao), depth, embedded)
			return
		}
		if !o.NoReflect {
			rv := reflect.ValueOf(data)
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			ao := alt.Options{
				CreateKey:    o.CreateKey,
				OmitNil:      o.OmitNil,
				FullTypePath: o.FullTypePath,
				UseTags:      o.UseTags,
				KeyExact:     o.KeyExact,
				NestEmbed:    o.NestEmbed,
				BytesAs:      o.BytesAs,
			}
			switch kind {
			case reflect.Struct:
				o.buildStruct(rv, &ao, depth, embedded)
			case reflect.Slice, reflect.Array:
				o.buildSlice(rv, &ao, depth)
			default:
				// Not much should get here except Map, Complex and un-decomposable
				// values.
				dec := alt.Decompose(data, &ao)
				o.buildJSON(dec, depth, false)
				return
			}
		} else if o.strict {
			panic(fmt.Errorf("%T can not be encoded as a JSON element", data))
		} else {
			o.buildString(fmt.Sprintf("%v", td))
		}
	}
	if o.w != nil && o.WriteLimit < len(o.buf) {
		if _, err := o.w.Write(o.buf); err != nil {
			panic(err)
		}
		o.buf = o.buf[:0]
	}
}

func (o *Options) buildString(s string) {
	o.buf = alt.AppendJSONString(o.buf, s, !o.HTMLUnsafe)
}

func (o *Options) buildTime(t time.Time) {
	if o.TimeMap {
		o.buf = append(o.buf, []byte(`{"`)...)
		o.buf = append(o.buf, o.CreateKey...)
		o.buf = append(o.buf, []byte(`":"`)...)
		if o.FullTypePath {
			o.buf = append(o.buf, []byte("time/Time")...)
		} else {
			o.buf = append(o.buf, []byte("Time")...)
		}
		o.buf = append(o.buf, []byte(`","value":`)...)
	} else if 0 < len(o.TimeWrap) {
		o.buf = append(o.buf, []byte(`{"`)...)
		o.buf = append(o.buf, []byte(o.TimeWrap)...)
		o.buf = append(o.buf, []byte(`":`)...)
	}
	switch o.TimeFormat {
	case "", "nano":
		o.buf = append(o.buf, []byte(strconv.FormatInt(t.UnixNano(), 10))...)
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			o.buf = append(o.buf, []byte(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))...)
		} else {
			o.buf = append(o.buf, []byte(fmt.Sprintf("%d.%09d", secs, -(nano-(secs*int64(time.Second)))))...)
		}
	default:
		o.buf = append(o.buf, '"')
		o.buf = append(o.buf, []byte(t.Format(o.TimeFormat))...)
		o.buf = append(o.buf, '"')
	}
	if 0 < len(o.TimeWrap) || o.TimeMap {
		o.buf = append(o.buf, '}')
	}
}

func (o *Options) buildArray(n gen.Array, depth int) {
	o.buf = append(o.buf, '[')
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		d2 := depth + 1
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buildJSON(m, d2, false)
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buildJSON(m, depth, false)
		}
	}
	o.buf = append(o.buf, ']')
}

func (o *Options) buildSimpleArray(n []interface{}, depth int) {
	o.buf = append(o.buf, '[')
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		d2 := depth + 1
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buildJSON(m, d2, false)
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j, m := range n {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buildJSON(m, depth, false)
		}
	}
	o.buf = append(o.buf, ']')
}

func (o *Options) buildObject(n gen.Object, depth int) {
	o.buf = append(o.buf, '{')
	first := true
	d2 := depth + 1
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildJSON(m, d2, false)
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildJSON(m, d2, false)
			}
		}
	}
	o.buf = append(o.buf, '}')
}

func (o *Options) buildSimpleObject(n map[string]interface{}, depth int) {
	o.buf = append(o.buf, '{')
	first := true
	d2 := depth + 1
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buf = append(o.buf, []byte(cs)...)
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buf = append(o.buf, ' ')
				o.buildJSON(m, d2, false)
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		if o.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildJSON(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && o.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					o.buf = append(o.buf, ',')
				}
				o.buildString(k)
				o.buf = append(o.buf, ':')
				o.buildJSON(m, d2, false)
			}
		}
	}
	o.buf = append(o.buf, '}')
}

func (o *Options) buildStruct(rv reflect.Value, opt *alt.Options, depth int, embedded bool) {
	dc := alt.LookupDecomposer(rv.Interface())
	var fields []*alt.Field
	d2 := depth + 1
	if opt.NestEmbed {
		if opt.UseTags {
			fields = dc.OutTag
		} else if opt.KeyExact {
			fields = dc.OutName
		} else {
			fields = dc.OutLow
		}
	} else {
		if opt.UseTags {
			fields = dc.ByTag
		} else if opt.KeyExact {
			fields = dc.ByName
		} else {
			fields = dc.ByLow
		}
	}
	o.buf = append(o.buf, '{')
	first := true
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for _, fi := range fields {
			v, omit := fi.Value(rv, opt.OmitNil, embedded)
			if omit || (opt.OmitNil && v == nil) {
				continue
			}
			if first {
				first = false
			} else {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			o.buildString(fi.Key)
			o.buf = append(o.buf, []byte{':', ' '}...)
			o.buildJSON(v, d2, true)
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		var v interface{}
		var has bool
		var wrote bool
		for _, fi := range fields {
			o.buf, v, wrote, has = fi.Append(o.buf, rv, opt.OmitNil, embedded)
			if wrote {
				first = false
				continue
			}
			if !has {
				continue
			}
			o.buildJSON(v, d2, true)
			o.buf = append(o.buf, ',')
			first = false

			/*
				v, omit := fi.Value(rv, opt.OmitNil, embedded)
				if omit || (opt.OmitNil && v == nil) {
					continue
				}
				o.buildString(fi.Key)
				o.buf = append(o.buf, ':')
				o.buildJSON(v, d2, true)
				first = false
				o.buf = append(o.buf, ',')
			*/
		}
		if !first {
			o.buf = o.buf[:len(o.buf)-1]
		}
		// TBD always add comma after but at end, if last byte is a comma, shorten buf

	}
	o.buf = append(o.buf, '}')
}

func (o *Options) buildSlice(rv reflect.Value, opt *alt.Options, depth int) {
	d2 := depth + 1
	end := rv.Len()
	o.buf = append(o.buf, '[')
	if o.Tab || 0 < o.Indent {
		var is string
		var cs string
		if o.Tab {
			x := depth + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			is = tabs[0:x]
			x = d2 + 1
			if len(tabs) < x {
				x = len(tabs)
			}
			cs = tabs[0:x]
		} else {
			x := depth*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*o.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for j := 0; j < end; j++ {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			o.buf = append(o.buf, []byte(cs)...)
			rm := rv.Index(j)
			switch rm.Kind() {
			case reflect.Struct:
				o.buildStruct(rm, opt, d2, false)
			case reflect.Slice, reflect.Array:
				o.buildSlice(rm, opt, d2)
			default:
				o.buildJSON(rm.Interface(), d2, false)
			}
		}
		o.buf = append(o.buf, []byte(is)...)
	} else {
		for j := 0; j < end; j++ {
			if 0 < j {
				o.buf = append(o.buf, ',')
			}
			rm := rv.Index(j)
			switch rm.Kind() {
			case reflect.Struct:
				o.buildStruct(rm, opt, d2, false)
			case reflect.Slice, reflect.Array:
				o.buildSlice(rm, opt, d2)
			default:
				o.buildJSON(rm.Interface(), d2, false)
			}
		}
	}
	o.buf = append(o.buf, ']')
}
