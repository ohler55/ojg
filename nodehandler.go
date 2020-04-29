// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

import (
	"github.com/ohler55/ojg/gd"
)

type keyStr string

type nodeHandler struct {
	cb          func(gd.Node) bool
	stack       []gd.Node
	arrayStarts []int
}

func newNodeHandler() *nodeHandler {
	h := nodeHandler{
		stack:       make([]gd.Node, 0, 64),
		arrayStarts: make([]int, 0, 16),
	}
	return &h
}

func (h *nodeHandler) ObjectStart() {
	top := len(h.stack) == 0
	n := gd.Object{}
	h.add(n)
	if !top {
		h.stack = append(h.stack, n)
	}
}

func (h *nodeHandler) ObjectEnd() {
	// does nothing
}

func (h *nodeHandler) ArrayStart() {
	h.arrayStarts = append(h.arrayStarts, len(h.stack))
}

func (h *nodeHandler) ArrayEnd() {
	start := h.arrayStarts[len(h.arrayStarts)-1]
	size := len(h.stack) - start
	n := gd.Array(make([]gd.Node, size))
	copy(n, h.stack[start:len(h.stack)])
	h.stack = h.stack[0:start]
	h.add(n)
}

func (h *nodeHandler) Null() {
	h.add(nil)
}

func (h *nodeHandler) Bool(value bool) {
	h.add(gd.Bool(value))
}

func (h *nodeHandler) Int(value int64) {
	h.add(gd.Int(value))
}

func (h *nodeHandler) Float(value float64) {
	h.add(gd.Float(value))
}

func (h *nodeHandler) Str(value string) {
	h.add(gd.String(value))
}

func (h *nodeHandler) Key(key string) {
	h.stack = append(h.stack, keyStr(key))
}

func (h *nodeHandler) Call() bool {
	return h.cb(h.stack[0])
}

func (h *nodeHandler) add(n gd.Node) {
	if 2 <= len(h.stack) {
		if k, ok := h.stack[len(h.stack)-1].(keyStr); ok {
			obj, _ := h.stack[len(h.stack)-2].(gd.Object)
			obj[string(k)] = n
			h.stack = h.stack[0 : len(h.stack)-1]

			return
		}
	}
	h.stack = append(h.stack, n)
}

// String returns the key as a string.
func (k keyStr) String() string {
	return string(k)
}

// Alter converts the node into it's native type. Note this will modify
// Objects and Arrays in place making them no longer usable as the
// original type. Use with care!
func (k keyStr) Alter() interface{} {
	return string(k)
}

// Simplify makes a copy of the node but as simple types.
func (k keyStr) Simplify() interface{} {
	return string(k)
}

// Dup returns a deep duplicate of the node.
func (k keyStr) Dup() gd.Node {
	return k
}

// Empty returns true if the node is empty.
func (k keyStr) Empty() bool {
	return false
}

// AsBool returns the Bool value of the node if possible. The ok return is
// true if successful.
func (k keyStr) AsBool() (v gd.Bool, ok bool) {
	return false, false
}

// AsInt returns the Int value of the node if possible. The ok return is
// true if successful.
func (k keyStr) AsInt() (v gd.Int, ok bool) {
	return 0, false
}

// AsFloat returns the Float value of the node if possible. The ok return
// is true if successful.
func (k keyStr) AsFloat() (v gd.Float, ok bool) {
	return 0.0, false
}
