// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ohler55/ojg/gd"
)

const (
	spaces = "\n                                                                                                                                "

	hex = "0123456789abcdef"
)

type Options struct {
	Indent  int
	Sort    bool
	SkipNil bool

	// TimeFormat defines how time is encoded. Options are to use a time. layout
	// string format such as time.RFC3339Nano, "second" for a decimal
	// representation, "nano" for a an integer.
	TimeFormat string

	// TimeWrap if not empty encoded time as an object with a single member. For
	// example if set to "@" then and TimeFormat is RFC3339Nano then the encoded
	// time will look like '{"@":"2020-04-12T16:34:04.123456789Z"}'
	TimeWrap string

	b strings.Builder
}

type jstr struct {
	Options
	b strings.Builder
}

// String returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a gd.Node type, The args, if supplied can be an
// int as an indent or a *Options.
func String(data interface{}, args ...interface{}) string {
	var js jstr

	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			js.Indent = ta
		case *Options:
			js.Options = *ta
		}
	}
	js.buildJSON(data, 0)

	return js.b.String()
}

func (js *jstr) buildJSON(data interface{}, depth int) {
	switch td := data.(type) {
	case nil:
		js.b.Write([]byte("null"))

	case bool:
		if td {
			js.b.Write([]byte("true"))
		} else {
			js.b.Write([]byte("false"))
		}
	case gd.Bool:
		if td {
			js.b.Write([]byte("true"))
		} else {
			js.b.Write([]byte("false"))
		}

	case int:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case int8:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case int16:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case int32:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case int64:
		js.b.WriteString(strconv.FormatInt(td, 10))
	case uint:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case uint8:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case uint16:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case uint32:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case uint64:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))
	case gd.Int:
		js.b.WriteString(strconv.FormatInt(int64(td), 10))

	case float32:
		js.b.WriteString(strconv.FormatFloat(float64(td), 'g', -1, 64))
	case float64:
		js.b.WriteString(strconv.FormatFloat(td, 'g', -1, 64))
	case gd.Float:
		js.b.WriteString(strconv.FormatFloat(float64(td), 'g', -1, 64))

	case string:
		js.buildString(td)
	case gd.String:
		js.buildString(string(td))

	case time.Time:
		js.buildTime(td)
	case gd.Time:
		js.buildTime(time.Time(td))

	case []interface{}:
		js.buildSimpleArray(td, depth)
	case gd.Array:
		js.buildArray(td, depth)

	case map[string]interface{}:
		js.buildSimpleObject(td, depth)
	case gd.Object:
		js.buildObject(td, depth)

	default:
		js.buildString(fmt.Sprintf("%v", td))
	}
}

func (js *jstr) buildString(s string) {
	js.b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			js.b.Write([]byte{'\\', '\\'})
		case '"':
			js.b.Write([]byte{'\\', '"'})
		case '\b':
			js.b.Write([]byte{'\\', 'b'})
		case '\f':
			js.b.Write([]byte{'\\', 'f'})
		case '\n':
			js.b.Write([]byte{'\\', 'n'})
		case '\r':
			js.b.Write([]byte{'\\', 'r'})
		case '\t':
			js.b.Write([]byte{'\\', 't'})
		case '&', '<', '>': // prefectly okay for JSON but commonly escaped
			js.b.Write([]byte{'\\', 'u', '0', '0', hex[r>>4], hex[r&0x0f]})
		case '\u2028':
			js.b.Write([]byte(`\u2028`))
		case '\u2029':
			js.b.Write([]byte(`\u2029`))
		default:
			if r < ' ' {
				js.b.Write([]byte{'\\', 'u', hex[r>>12], hex[(r>>8)&0x0f], hex[(r>>4)&0x0f], hex[r&0x0f]})
			} else if r < 0x80 {
				js.b.WriteByte(byte(r))
			} else {
				js.b.WriteRune(r)
			}
		}
	}
	js.b.WriteByte('"')
}

func (js *jstr) buildTime(t time.Time) {
	if 0 < len(js.TimeWrap) {
		js.b.WriteString(`{"`)
		js.b.WriteString(js.TimeWrap)
		js.b.WriteString(`":`)
	}
	switch js.TimeFormat {
	case "", "nano":
		js.b.WriteString(strconv.FormatInt(t.UnixNano(), 10))
	case "second":
		// Decimal format but float is not accurate enough so build the output
		// in two parts.
		nano := t.UnixNano()
		secs := nano / int64(time.Second)
		if 0 < nano {
			js.b.WriteString(fmt.Sprintf("%d.%09d", secs, nano-(secs*int64(time.Second))))
		} else {
			js.b.WriteString(fmt.Sprintf("%d.%09d", secs, -nano-(secs*int64(time.Second))))
		}
	default:
		js.b.WriteString(`"`)
		js.b.WriteString(t.Format(js.TimeFormat))
		js.b.WriteString(`"`)
	}
	if 0 < len(js.TimeWrap) {
		js.b.WriteString("}")
	}
}

func (js *jstr) buildArray(n gd.Array, depth int) {
	js.b.WriteByte('[')
	if 0 < js.Indent {
		x := depth*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		cs := spaces[0:x]

		for j, m := range n {
			if 0 < j {
				js.b.WriteByte(',')
			}
			js.b.WriteString(cs)
			if m == nil {
				js.b.WriteString("null")
			} else {
				js.buildJSON(m, d2)
			}
		}
		js.b.WriteString(is)
	} else {
		for j, m := range n {
			if 0 < j {
				js.b.WriteByte(',')
			}
			if m == nil {
				js.b.WriteString("null")
			} else {
				js.buildJSON(m, depth)
			}
		}
	}
	js.b.WriteByte(']')
}

func (js *jstr) buildSimpleArray(n []interface{}, depth int) {
	js.b.WriteByte('[')
	if 0 < js.Indent {
		x := depth*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		cs := spaces[0:x]

		for j, m := range n {
			if 0 < j {
				js.b.WriteByte(',')
			}
			js.b.WriteString(cs)
			if m == nil {
				js.b.WriteString("null")
			} else {
				js.buildJSON(m, d2)
			}
		}
		js.b.WriteString(is)
	} else {
		for j, m := range n {
			if 0 < j {
				js.b.WriteByte(',')
			}
			if m == nil {
				js.b.WriteString("null")
			} else {
				js.buildJSON(m, depth)
			}
		}
	}
	js.b.WriteByte(']')
}

func (js *jstr) buildObject(n gd.Object, depth int) {
	js.b.WriteByte('{')
	if 0 < js.Indent {
		x := depth*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		cs := spaces[0:x]
		if js.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && js.SkipNil {
					continue
				}
				if 0 < i {
					js.b.WriteByte(',')
				}
				js.b.WriteString(cs)
				js.buildString(k)
				js.b.WriteByte(':')
				if m := n[k]; m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, d2)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && js.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					js.b.WriteByte(',')
				}
				js.b.WriteString(cs)
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, d2)
				}
			}
		}
		js.b.WriteString(is)
	} else {
		if js.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && js.SkipNil {
					continue
				}
				if 0 < i {
					js.b.WriteByte(',')
				}
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, 0)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && js.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					js.b.WriteByte(',')
				}
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, 0)
				}
			}
		}
	}
	js.b.WriteByte('}')
}

func (js *jstr) buildSimpleObject(n map[string]interface{}, depth int) {
	js.b.WriteByte('{')
	if 0 < js.Indent {
		x := depth*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		is := spaces[0:x]
		d2 := depth + 1
		x = d2*js.Indent + 1
		if len(spaces) < x {
			x = depth*js.Indent + 1
		}
		cs := spaces[0:x]
		if js.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && js.SkipNil {
					continue
				}
				if 0 < i {
					js.b.WriteByte(',')
				}
				js.b.WriteString(cs)
				js.buildString(k)
				js.b.WriteByte(':')
				if m := n[k]; m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, d2)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && js.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					js.b.WriteByte(',')
				}
				js.b.WriteString(cs)
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, d2)
				}
			}
		}
		js.b.WriteString(is)
	} else {
		if js.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				m := n[k]
				if m == nil && js.SkipNil {
					continue
				}
				if 0 < i {
					js.b.WriteByte(',')
				}
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, 0)
				}
			}
		} else {
			first := true
			for k, m := range n {
				if m == nil && js.SkipNil {
					continue
				}
				if first {
					first = false
				} else {
					js.b.WriteByte(',')
				}
				js.buildString(k)
				js.b.WriteByte(':')
				if m == nil {
					js.b.WriteString("null")
				} else {
					js.buildJSON(m, 0)
				}
			}
		}
	}
	js.b.WriteByte('}')
}
