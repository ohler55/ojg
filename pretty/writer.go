// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"fmt"
	"io"
	"math"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

const (
	nullStr  = "null"
	trueStr  = "true"
	falseStr = "false"
	spaces   = "\n                                                                " +
		"                                                                "
	hex = "0123456789abcdef"
)

// Writer writes data in either JSON or SEN format using setting to determine
// the output.
type Writer struct {
	sen.Options

	// Width is the suggested maximum width. In some cases it may not be
	// possible to stay withing the specified width.
	Width int

	// MaxDepth is the maximum depth of an element on a single line.
	MaxDepth int

	// Align if true attempts to align elements of children in list.
	Align bool

	// SEN format if true otherwise JSON encoding.
	SEN bool
}

// JSON encoded output.
func JSON(data interface{}, args ...interface{}) string {
	w := Writer{
		Options:  sen.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      false,
	}
	w.config(args)
	b, _ := w.encode(data)

	return string(b)
}

// SEN encoded output.
func SEN(data interface{}, args ...interface{}) string {
	w := Writer{
		Options:  sen.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      true,
	}
	w.config(args)
	b, _ := w.encode(data)

	return string(b)
}

// JSON encoded output written to the provided io.Writer.
func WriteJSON(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := Writer{
		Options:  sen.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      false,
	}
	pw.W = w
	pw.config(args)
	_, err = pw.encode(data)

	return
}

// SEN encoded output written to the provided io.Writer.
func WriteSEN(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := Writer{
		Options:  sen.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      true,
	}
	pw.W = w
	pw.config(args)
	_, err = pw.encode(data)

	return
}

// Encode data. Any panics during encoding will cause an empty return but will
// not fail.
func (w *Writer) Encode(data interface{}) []byte {
	b, _ := w.encode(data)

	return b
}

// Marshal data. The same as Encode but a panics during encoding will result
// in an error return.
func (w *Writer) Marshal(data interface{}) ([]byte, error) {
	return w.encode(data)
}

// Write encoded data to the op.Writer.
func (w *Writer) Write(wr io.Writer, data interface{}) (err error) {
	w.W = wr
	_, err = w.encode(data)

	return
}

func (w *Writer) config(args []interface{}) {
	for _, arg := range args {
		switch ta := arg.(type) {
		case int:
			w.Width = ta
		case float64:
			if 0.0 < ta {
				if ta < 1.0 {
					w.MaxDepth = int(math.Round(ta * 10.0))
				} else {
					w.Width = int(ta)
					w.MaxDepth = int(math.Round((ta - float64(w.Width)) * 10.0))
				}
				if w.MaxDepth == 0 { // use the default
					w.MaxDepth = 2
				}
			}
		case bool:
			w.Align = ta
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
			w.Options.HTMLSafe = !ta.HTMLUnsafe
			w.W = sw
		}
	}
}

func (w *Writer) encode(data interface{}) (out []byte, err error) {
	if w.InitSize == 0 {
		w.InitSize = 256
	}
	if len(spaces)-1 < w.Width {
		w.Width = len(spaces) - 1
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
	if w.Width*3/8 < tree.depth {
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

func (w *Writer) fill(n *node, depth int, flat bool) {
	start := depth * w.Indent
	switch n.kind {
	case strNode, numNode:
		w.Buf = append(w.Buf, n.buf...)
	case arrayNode:
		var comma []byte
		if w.Color {
			if !w.SEN {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, '[')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			if !w.SEN {
				comma = append(comma, ',')
			}
			w.Buf = append(w.Buf, '[')
		}
		if !flat && start+n.size < w.Width && n.depth < w.MaxDepth {
			flat = true
		}
		d2 := depth + 1
		var cs []byte
		var is []byte

		if flat {
			cs = []byte{' '}
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
		if !w.Align || w.MaxDepth < n.depth || len(n.members) < 2 || w.checkAlign(n, start, comma, cs) {
			for i, m := range n.members {
				if 0 < i {
					w.Buf = append(w.Buf, comma...)
					w.Buf = append(w.Buf, []byte(cs)...)
				} else if !flat {
					w.Buf = append(w.Buf, []byte(cs)...)
				}
				w.fill(m, d2, flat)
			}
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
			if !w.SEN {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.Buf = append(w.Buf, w.SyntaxColor...)
			w.Buf = append(w.Buf, '{')
			w.Buf = append(w.Buf, w.NoColor...)
		} else {
			if !w.SEN {
				comma = append(comma, ',')
			}
			w.Buf = append(w.Buf, '{')
		}
		d2 := depth + 1
		var cs []byte
		var is []byte

		if !flat && start+n.size < w.Width && n.depth < w.MaxDepth {
			flat = true
		}
		if flat {
			cs = []byte{' '}
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

// Return true if not filled.
func (w *Writer) checkAlign(n *node, start int, comma, cs []byte) bool {
	c := n.genTables(w.SEN)
	if c == nil || w.Width < start+c.size {
		return true
	}
	for i, m := range n.members {
		if 0 < i {
			w.Buf = append(w.Buf, comma...)
		}
		w.Buf = append(w.Buf, []byte(cs)...)
		switch m.kind {
		case arrayNode:
			w.alignArray(m, c, comma, cs)
		case mapNode:
			w.alignMap(m, c, comma, cs)
		}
	}
	return false
}

func (w *Writer) alignArray(n *node, t *table, comma, cs []byte) {
	if w.Color {
		w.Buf = append(w.Buf, w.SyntaxColor...)
		w.Buf = append(w.Buf, '[')
		w.Buf = append(w.Buf, w.NoColor...)
	} else {
		w.Buf = append(w.Buf, '[')
	}
	for k, col := range t.columns {
		if len(n.members) <= k {
			break
		}
		if 0 < k {
			w.Buf = append(w.Buf, comma...)
			w.Buf = append(w.Buf, ' ')
		}
		m := n.members[k]
		cw := col.size
		switch m.kind {
		case strNode:
			w.Buf = append(w.Buf, m.buf...)
			if m.size < cw {
				w.Buf = append(w.Buf, spaces[1:cw-m.size+1]...)
			}
		case numNode:
			if m.size < cw {
				w.Buf = append(w.Buf, spaces[1:cw-m.size+1]...)
			}
			w.Buf = append(w.Buf, m.buf...)
		case arrayNode:
			w.alignArray(m, col, comma, []byte{' '})
		case mapNode:
			w.alignMap(m, col, comma, []byte{' '})
		}
	}
	if w.Color {
		w.Buf = append(w.Buf, w.SyntaxColor...)
		w.Buf = append(w.Buf, ']')
		w.Buf = append(w.Buf, w.NoColor...)
	} else {
		w.Buf = append(w.Buf, ']')
	}
}

func (w *Writer) alignMap(n *node, t *table, comma, cs []byte) {
	if w.Color {
		w.Buf = append(w.Buf, w.SyntaxColor...)
		w.Buf = append(w.Buf, '{')
		w.Buf = append(w.Buf, w.NoColor...)
	} else {
		w.Buf = append(w.Buf, '{')
	}
	prevExist := false
	for i, col := range t.columns {
		k, _ := col.key.(string)
		var m *node
		for _, mm := range n.members {
			if string(mm.key) == k {
				m = mm
				break
			}
		}
		if prevExist {
			w.Buf = append(w.Buf, comma...)
			w.Buf = append(w.Buf, ' ')
		}
		if m == nil {
			prevExist = false
			pad := len(k) + 2 + col.size
			if i < len(t.columns)-1 {
				if w.SEN {
					pad += 1
				} else {
					pad += 2
				}
			}
			w.Buf = append(w.Buf, spaces[1:pad+1]...)
		} else {
			prevExist = true
			w.Buf = append(w.Buf, k...)
			w.Buf = append(w.Buf, ':')
			w.Buf = append(w.Buf, ' ')
			cw := col.size
			switch m.kind {
			case strNode:
				w.Buf = append(w.Buf, m.buf...)
				if m.size < cw {
					w.Buf = append(w.Buf, spaces[1:cw-m.size+1]...)
				}
			case numNode:
				if m.size < cw {
					w.Buf = append(w.Buf, spaces[1:cw-m.size+1]...)
				}
				w.Buf = append(w.Buf, m.buf...)
			case arrayNode:
				w.alignArray(m, col, comma, []byte{' '})
			case mapNode:
				w.alignMap(m, col, comma, []byte{' '})
			}
		}
	}
	if w.Color {
		w.Buf = append(w.Buf, w.SyntaxColor...)
		w.Buf = append(w.Buf, '}')
		w.Buf = append(w.Buf, w.NoColor...)
	} else {
		w.Buf = append(w.Buf, '}')
	}
}
