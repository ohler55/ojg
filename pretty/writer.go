// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

const (
	nullStr  = "null"
	trueStr  = "true"
	falseStr = "false"
	spaces   = "\n                                                                                                                                "
	hex      = "0123456789abcdef"
)

type writer struct {
	sen.Options
	edge     int
	maxDepth int
	lazy     bool
}

func JSON(data interface{}, args ...interface{}) string {
	w := writer{
		Options:  sen.DefaultOptions,
		edge:     80,
		maxDepth: 2,
		lazy:     false,
	}
	b, _ := w.encode(data, args...)

	return string(b)
}

func SEN(data interface{}, args ...interface{}) string {
	w := writer{
		Options:  sen.DefaultOptions,
		edge:     80,
		maxDepth: 2,
		lazy:     true,
	}
	b, _ := w.encode(data, args...)

	return string(b)
}

func WriteJSON(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := writer{
		Options:  sen.DefaultOptions,
		edge:     80,
		maxDepth: 2,
		lazy:     false,
	}
	pw.W = w
	_, err = pw.encode(data, args...)

	return

}

func WriteSEN(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := writer{
		Options:  sen.DefaultOptions,
		edge:     80,
		maxDepth: 2,
		lazy:     true,
	}
	pw.W = w
	_, err = pw.encode(data, args...)

	return
}

func (w *writer) encode(data interface{}, args ...interface{}) (out []byte, err error) {
	for _, arg := range args {
		switch ta := arg.(type) {
		case int:
			w.edge = ta
		case float64:
			if 0.0 < ta {
				if ta < 1.0 {
					w.maxDepth = int(math.Round(ta * 10.0))
				} else {
					w.edge = int(ta)
					w.maxDepth = int(math.Round((ta - float64(w.edge)) * 10.0))
				}
				if w.maxDepth == 0 { // use the default
					w.maxDepth = 2
				}
			}
		case *sen.Options:
			sw := w.W
			w.Options = *ta
			w.W = sw
		case *oj.Options:
			sw := w.W
			w.Options.Indent = ta.Indent
			w.Options.Tab = ta.Tab
			w.Options.Sort = ta.Sort
			w.Options.OmitNil = ta.OmitNil
			w.Options.InitSize = ta.InitSize
			w.Options.WriteLimit = ta.WriteLimit
			w.Options.TimeFormat = ta.TimeFormat
			w.Options.TimeWrap = ta.TimeWrap
			w.Options.CreateKey = ta.CreateKey
			w.Options.FullTypePath = ta.FullTypePath
			w.Options.Color = ta.Color
			w.Options.SyntaxColor = ta.SyntaxColor
			w.Options.KeyColor = ta.KeyColor
			w.Options.NullColor = ta.NullColor
			w.Options.BoolColor = ta.BoolColor
			w.Options.NumberColor = ta.NumberColor
			w.Options.StringColor = ta.StringColor
			w.Options.NoColor = ta.NoColor
			w.W = sw
		}
	}
	if w.InitSize == 0 {
		w.InitSize = 256
	}
	if w.WriteLimit == 0 {
		w.WriteLimit = 1024
	}
	if cap(w.Buf) < w.InitSize {
		w.Buf = make([]byte, 0, w.InitSize)
	} else {
		w.Buf = w.Buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
				out = []byte{}
				if w.Color && w.W != nil {
					_, err = w.W.Write([]byte(w.NoColor))
				}
			}
		}
	}()
	tree := w.build(data)
	w.Buf = w.Buf[:0]
	w.Indent = 2
	if w.edge*3/8 < tree.depth {
		w.Indent = 1
	}
	w.fill(tree, 0, false)
	if w.W != nil && 0 < len(w.Buf) {
		_, err = w.W.Write(w.Buf)
		w.Buf = w.Buf[:0]
	}
	out = w.Buf

	return
}

func (w *writer) build(data interface{}) (n *node) {
	switch td := data.(type) {
	case nil:
		n = w.buildNull()
	case bool:
		n = w.buildBool(td)
	case gen.Bool:
		n = w.buildBool(bool(td))
	case int:
		n = w.buildInt(int64(td))
	case int8:
		n = w.buildInt(int64(td))
	case int16:
		n = w.buildInt(int64(td))
	case int32:
		n = w.buildInt(int64(td))
	case int64:
		n = w.buildInt(int64(td))
	case uint:
		n = w.buildInt(int64(td))
	case uint8:
		n = w.buildInt(int64(td))
	case uint16:
		n = w.buildInt(int64(td))
	case uint32:
		n = w.buildInt(int64(td))
	case uint64:
		n = w.buildInt(int64(td))
	case gen.Int:
		n = w.buildInt(int64(td))
	case float32:
		n = w.buildFloat32(td)
	case float64:
		n = w.buildFloat64(td)
	case gen.Float:
		n = w.buildFloat64(float64(td))
	case string:
		n = w.buildStringNode(td)
	case gen.String:
		n = w.buildStringNode(string(td))
	case time.Time:
		n = w.buildTimeNode(td)
	case gen.Time:
		n = w.buildTimeNode(time.Time(td))
	case []interface{}:
		n = w.buildArrayNode(td)
	case gen.Array:
		n = w.buildGenArrayNode(td)
	case map[string]interface{}:
		n = w.buildMapNode(td)
	case gen.Object:
		n = w.buildGenMapNode(td)
	default:
		if g, _ := data.(alt.Genericer); g != nil {
			return w.build(g.Generic())
		}
		if simp, _ := data.(alt.Simplifier); simp != nil {
			return w.build(simp.Simplify())
		}
		if 0 < len(w.CreateKey) {
			ao := alt.Options{CreateKey: w.CreateKey, OmitNil: w.OmitNil, FullTypePath: w.FullTypePath}
			return w.build(alt.Decompose(data, &ao))
		} else {
			n = w.buildStringNode(fmt.Sprintf("%v", td))
		}
	}
	return
}

func (w *writer) buildNull() *node {
	n := node{
		buf:  []byte(nullStr),
		size: 4,
		kind: bytesNode,
	}
	if w.Color {
		n.buf = append(append([]byte(w.NullColor), n.buf...), w.NoColor...)
	}
	return &n
}

func (w *writer) buildBool(v bool) (n *node) {
	if v {
		n = &node{
			buf:  []byte(trueStr),
			size: 4,
			kind: bytesNode,
		}
	} else {
		n = &node{
			buf:  []byte(falseStr),
			size: 5,
			kind: bytesNode,
		}
	}
	if w.Color {
		n.buf = append(append([]byte(w.BoolColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) buildInt(v int64) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatInt(v, 10)),
		kind: bytesNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) buildFloat32(v float32) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatFloat(float64(v), 'g', -1, 32)),
		kind: bytesNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) buildFloat64(v float64) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatFloat(v, 'g', -1, 64)),
		kind: bytesNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) buildStringNode(v string) (n *node) {
	w.Buf = w.Buf[:0]
	if w.lazy {
		w.BuildString(v)
	} else {
		w.BuildQuotedString(v)
	}
	n = &node{size: len(w.Buf), kind: bytesNode}
	n.buf = make([]byte, len(w.Buf))
	copy(n.buf, w.Buf)
	if w.Color {
		n.buf = append(append([]byte(w.StringColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) BuildQuotedString(s string) {
	w.Buf = append(w.Buf, '"')
	for _, r := range s {
		switch r {
		case '\\':
			w.Buf = append(w.Buf, []byte{'\\', '\\'}...)
		case '"':
			w.Buf = append(w.Buf, []byte{'\\', '"'}...)
		case '\b':
			w.Buf = append(w.Buf, []byte{'\\', 'b'}...)
		case '\f':
			w.Buf = append(w.Buf, []byte{'\\', 'f'}...)
		case '\n':
			w.Buf = append(w.Buf, []byte{'\\', 'n'}...)
		case '\r':
			w.Buf = append(w.Buf, []byte{'\\', 'r'}...)
		case '\t':
			w.Buf = append(w.Buf, []byte{'\\', 't'}...)
		case '&', '<', '>': // prefectly okay for JSON but commonly escaped
			w.Buf = append(w.Buf, []byte{'\\', 'u', '0', '0', hex[r>>4], hex[r&0x0f]}...)
		case '\u2028':
			w.Buf = append(w.Buf, []byte(`\u2028`)...)
		case '\u2029':
			w.Buf = append(w.Buf, []byte(`\u2029`)...)
		default:
			if r < ' ' {
				w.Buf = append(w.Buf, []byte{'\\', 'u', '0', '0', hex[(r>>4)&0x0f], hex[r&0x0f]}...)
			} else if r < 0x80 {
				w.Buf = append(w.Buf, byte(r))
			} else {
				if len(w.Utf) < utf8.UTFMax {
					w.Utf = make([]byte, utf8.UTFMax)
				} else {
					w.Utf = w.Utf[:cap(w.Utf)]
				}
				n := utf8.EncodeRune(w.Utf, r)
				w.Buf = append(w.Buf, w.Utf[:n]...)
			}
		}
	}
	w.Buf = append(w.Buf, '"')
}

func (w *writer) buildTimeNode(v time.Time) (n *node) {
	w.Buf = w.Buf[:0]
	w.BuildTime(v)
	n = &node{size: len(w.Buf), kind: bytesNode}
	n.buf = make([]byte, len(w.Buf))
	copy(n.buf, w.Buf)
	if w.Color {
		// TBD could be more detailed or better, have a separate time color
		n.buf = append(append([]byte(w.StringColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *writer) buildArrayNode(v []interface{}) (n *node) {
	n = &node{
		members: make([]*node, 0, len(v)),
		size:    2, // []
		kind:    arrayNode,
	}
	for i, m := range v {
		mn := w.build(m)
		n.members = append(n.members, mn)
		if 0 < i {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += mn.size
		if n.depth < mn.depth+1 {
			n.depth = mn.depth + 1
		}
	}
	return
}

func (w *writer) buildGenArrayNode(v gen.Array) (n *node) {
	n = &node{
		members: make([]*node, 0, len(v)),
		size:    2, // []
		kind:    arrayNode,
	}
	for i, m := range v {
		mn := w.build(m)
		n.members = append(n.members, mn)
		if 0 < i {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += mn.size
		if n.depth < mn.depth+1 {
			n.depth = mn.depth + 1
		}
	}
	return
}

func (w *writer) buildMapNode(v map[string]interface{}) (n *node) {
	n = &node{
		members: make([]*node, 0, len(v)),
		size:    2, // {}
		kind:    mapNode,
	}
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		mn := w.build(v[k])
		n.members = append(n.members, mn)
		// build key
		w.Buf = w.Buf[:0]
		if w.lazy {
			w.BuildString(k)
		} else {
			w.BuildQuotedString(k)
		}
		mn.key = make([]byte, len(w.Buf))
		copy(mn.key, w.Buf)
		if 2 < n.size {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += len(mn.key) + 2 + mn.size // key, colon, space, value
		if n.depth < mn.depth+1 {
			n.depth = mn.depth + 1
		}
		if w.Color {
			mn.key = append(append([]byte(w.KeyColor), mn.key...), w.NoColor...)
		}
	}
	return
}

func (w *writer) buildGenMapNode(v gen.Object) (n *node) {
	n = &node{
		members: make([]*node, 0, len(v)),
		size:    2, // {}
		kind:    mapNode,
	}
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		mn := w.build(v[k])
		n.members = append(n.members, mn)
		// build key
		w.Buf = w.Buf[:0]
		if w.lazy {
			w.BuildString(k)
		} else {
			w.BuildQuotedString(k)
		}
		mn.key = make([]byte, len(w.Buf))
		copy(mn.key, w.Buf)
		if 2 < n.size {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += len(mn.key) + 2 + mn.size // key, colon, space, value
		if n.depth < mn.depth+1 {
			n.depth = mn.depth + 1
		}
		if w.Color {
			mn.key = append(append([]byte(w.KeyColor), mn.key...), w.NoColor...)
		}
	}
	return
}

func (w *writer) fill(n *node, depth int, flat bool) {
	start := depth * w.Indent
	switch n.kind {
	case bytesNode:
		w.Buf = append(w.Buf, n.buf...)
	case arrayNode:
		var comma []byte
		if w.Color {
			if !w.lazy {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, '[')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			if !w.lazy {
				comma = append(comma, ',')
			}
			w.Buf = append(w.Buf, '[')
		}
		d2 := depth + 1
		var cs []byte
		var is []byte

		if flat || (start+n.size < w.edge && n.depth < w.maxDepth) {
			cs = []byte{' '}
			flat = true
		} else {
			x := d2*w.Indent + 1
			if len(spaces) < x {
				flat = true
			} else {
				cs = []byte(spaces[0:x])
				x = depth*w.Indent + 1
				is = []byte(spaces[0:x])
			}
		}
		for i, m := range n.members {
			if 0 < i {
				w.Buf = append(w.Buf, comma...)
				w.Buf = append(w.Buf, []byte(cs)...)
			} else if !flat {
				w.Buf = append(w.Buf, []byte(cs)...)
			}
			w.fill(m, d2, flat)
		}
		w.Buf = append(w.Buf, []byte(is)...)
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, ']')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			w.Buf = append(w.Buf, ']')
		}
	case mapNode:
		var comma []byte
		if w.Color {
			if !w.lazy {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, '{')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			if !w.lazy {
				comma = append(comma, ',')
			}
			w.Buf = append(w.Buf, '{')
		}
		d2 := depth + 1
		var cs []byte
		var is []byte
		if flat || (start+n.size < w.edge && n.depth < w.maxDepth) {
			cs = []byte{' '}
			flat = true
		} else {
			x := d2*w.Indent + 1
			if len(spaces) < x {
				flat = true
			} else {
				cs = []byte(spaces[0:x])
				x = depth*w.Indent + 1
				is = []byte(spaces[0:x])
			}
		}
		for i, m := range n.members {
			if 0 < i {
				w.Buf = append(w.Buf, comma...)
				w.Buf = append(w.Buf, []byte(cs)...)
			} else if !flat {
				w.Buf = append(w.Buf, []byte(cs)...)
			}
			w.Buf = append(w.Buf, m.key...)
			if w.Color {
				w.Buf = append(w.Buf, w.SyntaxColor...)
				w.Buf = append(w.Buf, ':')
				w.Buf = append(w.Buf, w.NoColor...)
				w.Buf = append(w.Buf, ' ')
			} else {
				w.Buf = append(w.Buf, ": "...)
			}
			w.fill(m, d2, flat)
		}
		w.Buf = append(w.Buf, []byte(is)...)
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, '}')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			w.Buf = append(w.Buf, '}')
		}
	}
	if w.W != nil && w.WriteLimit < len(w.Buf) {
		if _, err := w.W.Write(w.Buf); err != nil {
			panic(err)
		}
		w.Buf = w.Buf[:0]
	}
}
