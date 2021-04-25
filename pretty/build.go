// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"sort"
	"strconv"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

func (w *Writer) build(data interface{}) (n *node) {
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
			return w.build(alt.Decompose(data, &alt.Options{OmitNil: w.OmitNil}))
		}
	}
	return
}

func (w *Writer) buildNull() *node {
	n := node{
		buf:  []byte(nullStr),
		size: 4,
		kind: strNode,
	}
	if w.Color {
		n.buf = append(append([]byte(w.NullColor), n.buf...), w.NoColor...)
	}
	return &n
}

func (w *Writer) buildBool(v bool) (n *node) {
	if v {
		n = &node{
			buf:  []byte(trueStr),
			size: 4,
			kind: strNode,
		}
	} else {
		n = &node{
			buf:  []byte(falseStr),
			size: 5,
			kind: strNode,
		}
	}
	if w.Color {
		n.buf = append(append([]byte(w.BoolColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildInt(v int64) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatInt(v, 10)),
		kind: numNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildFloat32(v float32) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatFloat(float64(v), 'g', -1, 32)),
		kind: numNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildFloat64(v float64) (n *node) {
	n = &node{
		buf:  []byte(strconv.FormatFloat(v, 'g', -1, 64)),
		kind: numNode,
	}
	n.size = len(n.buf)
	if w.Color {
		n.buf = append(append([]byte(w.NumberColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildStringNode(v string) (n *node) {
	w.buf = w.buf[:0]
	if w.SEN {
		w.buf = ojg.AppendSENString(w.buf, v)
	} else {
		w.buf = ojg.AppendJSONString(w.buf, v, !w.HTMLUnsafe)
	}
	n = &node{size: len(w.buf), kind: strNode}
	n.buf = make([]byte, len(w.buf))
	copy(n.buf, w.buf)
	if w.Color {
		n.buf = append(append([]byte(w.StringColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildTimeNode(v time.Time) (n *node) {
	w.buf = w.buf[:0]
	w.buf = w.BuildTime(w.buf, v)
	n = &node{size: len(w.buf), kind: strNode}
	n.buf = make([]byte, len(w.buf))
	copy(n.buf, w.buf)
	if w.Color {
		n.buf = append(append([]byte(w.TimeColor), n.buf...), w.NoColor...)
	}
	return
}

func (w *Writer) buildArrayNode(v []interface{}) (n *node) {
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
			if !w.SEN {
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

func (w *Writer) buildGenArrayNode(v gen.Array) (n *node) {
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
			if !w.SEN {
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

func (w *Writer) buildMapNode(v map[string]interface{}) (n *node) {
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
		w.buf = w.buf[:0]
		if w.SEN {
			w.buf = ojg.AppendSENString(w.buf, k)
		} else {
			w.buf = ojg.AppendJSONString(w.buf, k, !w.HTMLUnsafe)
		}
		mn.key = make([]byte, len(w.buf))
		copy(mn.key, w.buf)
		if 2 < n.size {
			n.size++ // space
			if !w.SEN {
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

func (w *Writer) buildGenMapNode(v gen.Object) (n *node) {
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
		w.buf = w.buf[:0]
		if w.SEN {
			w.buf = ojg.AppendSENString(w.buf, k)
		} else {
			w.buf = ojg.AppendJSONString(w.buf, k, !w.HTMLUnsafe)
		}
		mn.key = make([]byte, len(w.buf))
		copy(mn.key, w.buf)
		if 2 < n.size {
			n.size++ // space
			if !w.SEN {
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
