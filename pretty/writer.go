// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"fmt"
	"io"
	"math"

	"github.com/ohler55/ojg"
)

const (
	nullStr  = "null"
	trueStr  = "true"
	falseStr = "false"
	spaces   = "\n                                                                " +
		"                                                                "
)

// Writer writes data in either JSON or SEN format using setting to determine
// the output.
type Writer struct {
	ojg.Options

	// Width is the suggested maximum width. In some cases it may not be
	// possible to stay withing the specified width.
	Width int

	// MaxDepth is the maximum depth of an element on a single line.
	MaxDepth int

	// Align if true attempts to align elements of children in list.
	Align bool

	// SEN format if true otherwise JSON encoding.
	SEN bool

	buf []byte
	w   io.Writer
}

// JSON encoded output.
func JSON(data interface{}, args ...interface{}) string {
	w := Writer{
		Options:  ojg.DefaultOptions,
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
		Options:  ojg.DefaultOptions,
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
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      false,
	}
	pw.w = w
	pw.config(args)
	_, err = pw.encode(data)

	return
}

// SEN encoded output written to the provided io.Writer.
func WriteSEN(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := Writer{
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      true,
	}
	pw.w = w
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
	w.w = wr
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
		case *ojg.Options:
			sw := w.w
			w.Options = *ta
			w.w = sw
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
	if cap(w.buf) < w.InitSize {
		w.buf = make([]byte, 0, w.InitSize)
	} else {
		w.buf = w.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
				out = []byte{}
				if w.Color && w.w != nil {
					_, err = w.w.Write([]byte(w.NoColor))
				}
			}
		}
	}()
	tree := w.build(data)
	w.buf = w.buf[:0]
	w.Indent = 2
	if w.Width*3/8 < tree.depth {
		w.Indent = 1
	}
	w.fill(tree, 0, false)
	if w.w != nil && 0 < len(w.buf) {
		_, err = w.w.Write(w.buf)
		w.buf = w.buf[:0]
	}
	out = w.buf

	return
}

func (w *Writer) fill(n *node, depth int, flat bool) {
	start := depth * w.Indent
	switch n.kind {
	case strNode, numNode:
		w.buf = append(w.buf, n.buf...)
	case arrayNode:
		var comma []byte
		if w.Color {
			if !w.SEN {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.buf = append(w.buf, w.SyntaxColor...)
			w.buf = append(w.buf, '[')
			w.buf = append(w.buf, w.NoColor...)
		} else {
			if !w.SEN {
				comma = append(comma, ',')
			}
			w.buf = append(w.buf, '[')
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
					w.buf = append(w.buf, comma...)
					w.buf = append(w.buf, []byte(cs)...)
				} else if !flat {
					w.buf = append(w.buf, []byte(cs)...)
				}
				w.fill(m, d2, flat)
			}
		}
		w.buf = append(w.buf, []byte(is)...)
		if w.Color {
			w.buf = append(w.buf, w.SyntaxColor...)
			w.buf = append(w.buf, ']')
			w.buf = append(w.buf, w.NoColor...)
		} else {
			w.buf = append(w.buf, ']')
		}
	case mapNode:
		var comma []byte
		if w.Color {
			if !w.SEN {
				comma = append(comma, w.SyntaxColor...)
				comma = append(comma, ',')
				comma = append(comma, w.NoColor...)
			}
			w.buf = append(w.buf, w.SyntaxColor...)
			w.buf = append(w.buf, '{')
			w.buf = append(w.buf, w.NoColor...)
		} else {
			if !w.SEN {
				comma = append(comma, ',')
			}
			w.buf = append(w.buf, '{')
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
				w.buf = append(w.buf, comma...)
				w.buf = append(w.buf, []byte(cs)...)
			} else if !flat {
				w.buf = append(w.buf, []byte(cs)...)
			}
			w.buf = append(w.buf, m.key...)
			if w.Color {
				w.buf = append(w.buf, w.SyntaxColor...)
				w.buf = append(w.buf, ':')
				w.buf = append(w.buf, w.NoColor...)
				w.buf = append(w.buf, ' ')
			} else {
				w.buf = append(w.buf, ": "...)
			}
			w.fill(m, d2, flat)
		}
		w.buf = append(w.buf, []byte(is)...)
		if w.Color {
			w.buf = append(w.buf, w.SyntaxColor...)
			w.buf = append(w.buf, '}')
			w.buf = append(w.buf, w.NoColor...)
		} else {
			w.buf = append(w.buf, '}')
		}
	}
	if w.w != nil && w.WriteLimit < len(w.buf) {
		if _, err := w.w.Write(w.buf); err != nil {
			panic(err)
		}
		w.buf = w.buf[:0]
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
			w.buf = append(w.buf, comma...)
		}
		w.buf = append(w.buf, []byte(cs)...)
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
		w.buf = append(w.buf, w.SyntaxColor...)
		w.buf = append(w.buf, '[')
		w.buf = append(w.buf, w.NoColor...)
	} else {
		w.buf = append(w.buf, '[')
	}
	for k, col := range t.columns {
		if len(n.members) <= k {
			break
		}
		if 0 < k {
			w.buf = append(w.buf, comma...)
			w.buf = append(w.buf, ' ')
		}
		m := n.members[k]
		cw := col.size
		switch m.kind {
		case strNode:
			w.buf = append(w.buf, m.buf...)
			if m.size < cw {
				w.buf = append(w.buf, spaces[1:cw-m.size+1]...)
			}
		case numNode:
			if m.size < cw {
				w.buf = append(w.buf, spaces[1:cw-m.size+1]...)
			}
			w.buf = append(w.buf, m.buf...)
		case arrayNode:
			w.alignArray(m, col, comma, []byte{' '})
		case mapNode:
			w.alignMap(m, col, comma, []byte{' '})
		}
	}
	if w.Color {
		w.buf = append(w.buf, w.SyntaxColor...)
		w.buf = append(w.buf, ']')
		w.buf = append(w.buf, w.NoColor...)
	} else {
		w.buf = append(w.buf, ']')
	}
}

func (w *Writer) alignMap(n *node, t *table, comma, cs []byte) {
	if w.Color {
		w.buf = append(w.buf, w.SyntaxColor...)
		w.buf = append(w.buf, '{')
		w.buf = append(w.buf, w.NoColor...)
	} else {
		w.buf = append(w.buf, '{')
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
			w.buf = append(w.buf, comma...)
			w.buf = append(w.buf, ' ')
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
			w.buf = append(w.buf, spaces[1:pad+1]...)
		} else {
			prevExist = true
			w.buf = append(w.buf, k...)
			w.buf = append(w.buf, ':')
			w.buf = append(w.buf, ' ')
			cw := col.size
			switch m.kind {
			case strNode:
				w.buf = append(w.buf, m.buf...)
				if m.size < cw {
					w.buf = append(w.buf, spaces[1:cw-m.size+1]...)
				}
			case numNode:
				if m.size < cw {
					w.buf = append(w.buf, spaces[1:cw-m.size+1]...)
				}
				w.buf = append(w.buf, m.buf...)
			case arrayNode:
				w.alignArray(m, col, comma, []byte{' '})
			case mapNode:
				w.alignMap(m, col, comma, []byte{' '})
			}
		}
	}
	if w.Color {
		w.buf = append(w.buf, w.SyntaxColor...)
		w.buf = append(w.buf, '}')
		w.buf = append(w.buf, w.NoColor...)
	} else {
		w.buf = append(w.buf, '}')
	}
}
