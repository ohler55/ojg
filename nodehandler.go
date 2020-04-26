// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import "github.com/ohler55/ojg/gd"

type nodeHandler struct {
	cb func(gd.Node) bool
	// TBD stack of structs with key, node, and reusable node array
	foo gd.Node
}

func (h *nodeHandler) ObjectStart() {
	h.foo = gd.Object{}
	// TBD
}

func (h *nodeHandler) ObjectEnd() {
	// TBD
}

func (h *nodeHandler) ArrayStart() {
	h.foo = gd.Array{}
	// TBD
}

func (h *nodeHandler) ArrayEnd() {
	// TBD
}

func (h *nodeHandler) Null() {
	h.foo = nil
	// TBD
}

func (h *nodeHandler) Bool(value bool) {
	h.foo = gd.Bool(value)
	// TBD
}

func (h *nodeHandler) Int(value int64) {
	h.foo = gd.Int(value)
	// TBD
}

func (h *nodeHandler) Float(value float64) {
	h.foo = gd.Float(value)
	// TBD
}

func (h *nodeHandler) Str(s string) {
	h.foo = gd.String(s)
	// TBD
}

func (h *nodeHandler) Key(key string) {
	h.foo = gd.String(key)
	// TBD
}

func (h *nodeHandler) Call() bool {
	// TBD top of stack is node in cb
	return h.cb(h.foo)
}
