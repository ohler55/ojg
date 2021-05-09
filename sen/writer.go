// Copyright (c) 2020, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

const (
	spaces = "\n                                                                                                                                "
	tabs   = "\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
)

// Writer is a SEN writer that includes a reused buffer for reduced
// allocations for repeated encoding calls.
type Writer struct {
	ojg.Options
	buf []byte
	w   io.Writer
}

// SEN writes data, SEN encoded. On error, an empty string is returned.
func (wr *Writer) SEN(data interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
		}
	}()
	return wr.MustSEN(data)
}

// MustSEN writes data, SEN encoded. On error a panic is called with the error.
func (wr *Writer) MustSEN(data interface{}) string {
	if wr.InitSize <= 0 {
		wr.InitSize = 256
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	wr.buildSen(data, 0, false)

	return string(wr.buf)
}

// Write a SEN string for the data provided.
func (wr *Writer) Write(w io.Writer, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	wr.MustWrite(w, data)
	return
}

// MustWrite a SEN string for the data provided. If an error occurs panic is
// called with the error.
func (wr *Writer) MustWrite(w io.Writer, data interface{}) {
	wr.w = w
	if wr.InitSize <= 0 {
		wr.InitSize = 256
	}
	if wr.WriteLimit <= 0 {
		wr.WriteLimit = 1024
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	if wr.Color {
		wr.cbuildSen(data, 0) // TBD embedded
	} else {
		wr.buildSen(data, 0, false)
	}
	if 0 < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
	}
}

func (wr *Writer) buildSen(data interface{}, depth int, embedded bool) {
	switch td := data.(type) {
	case nil:
		wr.buf = append(wr.buf, []byte("null")...)

	case bool:
		if td {
			wr.buf = append(wr.buf, []byte("true")...)
		} else {
			wr.buf = append(wr.buf, []byte("false")...)
		}

	case int:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int8:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int16:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int32:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case int64:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(td, 10))...)
	case uint:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint8:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint16:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint32:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)
	case uint64:
		wr.buf = append(wr.buf, []byte(strconv.FormatInt(int64(td), 10))...)

	case float32:
		wr.buf = append(wr.buf, []byte(strconv.FormatFloat(float64(td), 'g', -1, 32))...)
	case float64:
		wr.buf = append(wr.buf, []byte(strconv.FormatFloat(td, 'g', -1, 64))...)

	case string:
		wr.buf = ojg.AppendSENString(wr.buf, td, !wr.HTMLUnsafe)

	case time.Time:
		wr.buf = wr.AppendTime(wr.buf, td, true)

	case []interface{}:
		wr.buildSimpleArray(td, depth)

	case map[string]interface{}:
		wr.buildSimpleObject(td, depth)

	default:
		if simp, _ := data.(alt.Simplifier); simp != nil {
			data = simp.Simplify()
			wr.buildSen(data, depth, false)
			return
		}
		if g, _ := data.(alt.Genericer); g != nil {
			wr.buildSen(g.Generic().Simplify(), depth, false)
			return
		}
		if 0 < len(wr.CreateKey) {
			ao := alt.Options{CreateKey: wr.CreateKey, OmitNil: wr.OmitNil, FullTypePath: wr.FullTypePath}
			wr.buildSen(alt.Decompose(data, &ao), depth, false)
			return
		} else {
			wr.buildSen(alt.Decompose(data, &alt.Options{OmitNil: wr.OmitNil}), depth, false)
			return
		}
	}
	if wr.w != nil && wr.WriteLimit < len(wr.buf) {
		if _, err := wr.w.Write(wr.buf); err != nil {
			panic(err)
		}
		wr.buf = wr.buf[:0]
	}
}

func (wr *Writer) buildSimpleArray(n []interface{}, depth int) {
	wr.buf = append(wr.buf, '[')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
		if wr.Tab {
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
			x := depth*wr.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*wr.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		for _, m := range n {
			wr.buf = append(wr.buf, []byte(cs)...)
			wr.buildSen(m, d2, false)
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		var prev interface{}
		for j, m := range n {
			if 0 < j {
				switch prev.(type) {
				case []interface{}, map[string]interface{}:
				default:
					wr.buf = append(wr.buf, ' ')
				}
			}
			prev = m
			wr.buildSen(m, depth, false)
		}
	}
	wr.buf = append(wr.buf, ']')
}

func (wr *Writer) buildSimpleObject(n map[string]interface{}, depth int) {
	wr.buf = append(wr.buf, '{')
	d2 := depth + 1
	var is string
	var cs string

	if wr.Tab || 0 < wr.Indent {
		if wr.Tab {
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
			x := depth*wr.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is = spaces[0:x]
			x = d2*wr.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs = spaces[0:x]
		}
		if wr.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				wr.buf = append(wr.buf, []byte(cs)...)
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buf = append(wr.buf, ' ')
				wr.buildSen(m, d2, false)
			}
		}
		wr.buf = append(wr.buf, []byte(is)...)
	} else {
		first := true
		if wr.Sort {
			keys := make([]string, 0, len(n))
			for k := range n {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				m := n[k]
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0, false)
			}
		} else {
			for k, m := range n {
				if m == nil && wr.OmitNil {
					continue
				}
				if first {
					first = false
				} else {
					wr.buf = append(wr.buf, ' ')
				}
				wr.buf = ojg.AppendSENString(wr.buf, k, !wr.HTMLUnsafe)
				wr.buf = append(wr.buf, ':')
				wr.buildSen(m, 0, false)
			}
		}
	}
	wr.buf = append(wr.buf, '}')
}
