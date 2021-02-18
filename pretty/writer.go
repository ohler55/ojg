// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/sen"
)

const (
	nullStr  = "null"
	trueStr  = "true"
	falseStr = "false"
	spaces   = "\n                                                                                                                                "
)

type writer struct {
	sen.Options
	edge  int
	lazy  bool
	depth int
}

func JSON(data interface{}, args ...interface{}) string {
	w := writer{
		Options: sen.DefaultOptions,
		edge:    80,
		lazy:    false,
		depth:   0,
	}
	w.Quote = true

	b, _ := w.encode(data, args...)

	return string(b)
}

func SEN(data interface{}, args ...interface{}) string {
	w := writer{
		Options: sen.DefaultOptions,
		edge:    80,
		lazy:    true,
		depth:   0,
	}
	w.Quote = false
	b, _ := w.encode(data, args...)

	return string(b)
}

func (w *writer) encode(data interface{}, args ...interface{}) (out []byte, err error) {
	for _, arg := range args {
		switch ta := arg.(type) {
		case int:
			w.edge = ta
		case *sen.Options:
			w.Options = *ta
		}
	}
	if w.InitSize == 0 {
		w.InitSize = 256
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
			}
		}
	}()
	tree := w.build(data)
	w.Buf = w.Buf[:0]
	w.Indent = 2
	if w.edge*3/8 < w.depth {
		w.Indent = 1
	}
	w.fill(tree, 0, false)
	out = w.Buf

	return
}

func (w *writer) build(data interface{}) (n *node) {
	w.depth++
	n = &node{}
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
	if w.Color {
		w.Buf = append(w.Buf, sen.Normal...)
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
		n.buf = append(append([]byte(w.NullColor), n.buf...), sen.Normal...)
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
		n.buf = append(append([]byte(w.BoolColor), n.buf...), sen.Normal...)
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
		n.buf = append(append([]byte(w.NumberColor), n.buf...), sen.Normal...)
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
		n.buf = append(append([]byte(w.NumberColor), n.buf...), sen.Normal...)
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
		n.buf = append(append([]byte(w.NumberColor), n.buf...), sen.Normal...)
	}
	return
}

func (w *writer) buildStringNode(v string) (n *node) {
	w.Buf = w.Buf[:0]
	w.BuildString(v)
	n = &node{size: len(w.Buf), kind: bytesNode}
	n.buf = make([]byte, len(w.Buf))
	copy(n.buf, w.Buf)
	if w.Color {
		n.buf = append(append([]byte(w.StringColor), n.buf...), sen.Normal...)
	}
	return
}

func (w *writer) buildTimeNode(v time.Time) (n *node) {
	w.Buf = w.Buf[:0]
	w.BuildTime(v)
	n = &node{size: len(w.Buf), kind: bytesNode}
	n.buf = make([]byte, len(w.Buf))
	copy(n.buf, w.Buf)
	if w.Color {
		// TBD could be more detailed or better, have a separate time color
		n.buf = append(append([]byte(w.StringColor), n.buf...), sen.Normal...)
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
	}
	return
}

func (w *writer) buildMapNode(v map[string]interface{}) (n *node) {
	n = &node{
		members: make([]*node, 0, len(v)),
		size:    2, // {}
		kind:    mapNode,
	}
	for k, m := range v {
		mn := w.build(m)
		n.members = append(n.members, mn)
		// build key
		w.Buf = w.Buf[:0]
		w.BuildString(k)
		mn.key = make([]byte, len(w.Buf))
		copy(mn.key, w.Buf)
		if 2 < n.size {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += len(mn.key) + 2 + mn.size // key, colon, space, value
		if w.Color {
			mn.key = append(append([]byte(w.KeyColor), mn.key...), sen.Normal...)
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
	for k, m := range v {
		mn := w.build(m)
		n.members = append(n.members, mn)
		// build key
		w.Buf = w.Buf[:0]
		w.BuildString(k)
		mn.key = make([]byte, len(w.Buf))
		copy(mn.key, w.Buf)
		if 2 < n.size {
			n.size++ // space
			if !w.lazy {
				n.size++ // comma
			}
		}
		n.size += len(mn.key) + 2 + mn.size // key, colon, space, value
		if w.Color {
			mn.key = append(append([]byte(w.KeyColor), mn.key...), sen.Normal...)
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
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
		}
		w.Buf = append(w.Buf, '[')
		d2 := depth + 1
		if flat || start+n.size < w.edge {
			for i, m := range n.members {
				if 0 < i {
					if !w.lazy {
						w.Buf = append(w.Buf, ',')
					}
					w.Buf = append(w.Buf, ' ')
				}
				w.fill(m, d2, true)
			}
		} else {
			x := depth*w.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is := spaces[0:x]
			x = d2*w.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs := spaces[0:x]
			for i, m := range n.members {
				if 0 < i && !w.lazy {
					w.Buf = append(w.Buf, ',')
				}
				w.Buf = append(w.Buf, []byte(cs)...)
				w.fill(m, d2, flat)
			}
			w.Buf = append(w.Buf, []byte(is)...)
		}
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
		}
		w.Buf = append(w.Buf, ']')
	case mapNode:
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
		}
		w.Buf = append(w.Buf, '{')
		d2 := depth + 1
		if flat || start+n.size < w.edge {
			for i, m := range n.members {
				if 0 < i {
					if !w.lazy {
						w.Buf = append(w.Buf, ',')
					}
					w.Buf = append(w.Buf, ' ')
				}
				if w.Color {
					w.Buf = append(w.Buf, w.KeyColor...)
				}
				w.Buf = append(w.Buf, m.key...)
				if w.Color {
					w.Buf = append(w.Buf, w.SyntaxColor...)
					w.Buf = append(w.Buf, ": "...)
				} else {
					w.Buf = append(w.Buf, ": "...)
				}
				w.fill(m, d2, true)
			}
		} else {
			x := depth*w.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			is := spaces[0:x]
			x = d2*w.Indent + 1
			if len(spaces) < x {
				x = len(spaces)
			}
			cs := spaces[0:x]
			for i, m := range n.members {
				if 0 < i && !w.lazy {
					w.Buf = append(w.Buf, ',')
				}
				w.Buf = append(w.Buf, []byte(cs)...)
				if w.Color {
					w.Buf = append(w.Buf, w.KeyColor...)
				}
				w.Buf = append(w.Buf, m.key...)
				if w.Color {
					w.Buf = append(w.Buf, w.SyntaxColor...)
					w.Buf = append(w.Buf, ": "...)
				} else {
					w.Buf = append(w.Buf, ": "...)
				}
				w.fill(m, d2, flat)
			}
			w.Buf = append(w.Buf, []byte(is)...)
		}
		if w.Color {
			w.Buf = append(w.Buf, w.SyntaxColor...)
		}
		w.Buf = append(w.Buf, '}')
	}
}
