// Copyright (c) 2024, Peter Ohler, All rights reserved.

package jp

import "encoding/json"

// PathHandler is a TokenHandler compatible with both the oj.TokenHandler and
// the sen.TokenHandler. Fields are public to allow derived types to access
// those fields.
type PathDataHandler struct {
	Target  Expr
	Path    Expr
	StopLen int
	Data    any
	OnData  func(path Expr, data any)
}

// NewPathDataHandler creates a new PathDataHandler.
func NewPathDataHandler(target Expr, onData func(path Expr, data any)) *PathDataHandler {
	return &PathDataHandler{
		Target: target,
		OnData: onData,
	}
}

// Null is called when a JSON null is encountered.
func (h *PathDataHandler) Null() {
	h.AddValue(nil)
}

// Bool is called when a JSON true or false is encountered.
func (h *PathDataHandler) Bool(v bool) {
	h.AddValue(v)
}

// Int is called when a JSON integer is encountered.
func (h *PathDataHandler) Int(v int64) {
	h.AddValue(v)
}

// Float is called when a JSON decimal is encountered that fits into a
// float64.
func (h *PathDataHandler) Float(v float64) {
	h.AddValue(v)
}

// Number is called when a JSON number is encountered that does not fit
// into an int64 or float64.
func (h *PathDataHandler) Number(num string) {
	h.AddValue(json.Number(num))
}

// String is called when a JSON string is encountered.
func (h *PathDataHandler) String(v string) {
	h.AddValue(v)
}

// ObjectStart is called when a JSON object start '{' is encountered.
func (h *PathDataHandler) ObjectStart() {
	// TBD check target, if a match then set StopLen
	h.Path = append(h.Path, Child(""))
}

// ObjectEnd is called when a JSON object end '}' is encountered.
func (h *PathDataHandler) ObjectEnd() {
	h.Path = h.Path[:len(h.Path)-1]
	// TBD check StopLen and OnData if needed
}

// Key is called when a JSON object key is encountered.
func (h *PathDataHandler) Key(k string) {
	h.Path[len(h.Path)-1] = Child(k)
	// TBD set last in path to Child(k)
	// if StopLen matches len(path) then call CallOnData()
	// compare to target, if a match then set StopLen
}

// ArrayStart is called when a JSON array start '[' is encountered.
func (h *PathDataHandler) ArrayStart() {
	// TBD add Nth(-1) to path
}

// ArrayEnd is called when a JSON array end ']' is encountered.
func (h *PathDataHandler) ArrayEnd() {
	h.Path = h.Path[:len(h.Path)-1]
	// TBD check StopLen and OnData if needed
}

// AdValue is called when a leave value is encountered.
func (h *PathDataHandler) AddValue(v any) {
	last := len(h.Path) - 1
	if 0 <= last {
		if nth, ok := h.Path[last].(Nth); ok {
			h.Path[last] = nth + 1
		}
	}
	// TBD add to data if 0 < StopLen
	// TBD else if path matches target call OnData
	// TBD look at data, if n
}

// CallOnData calls the OnData function.
func (h *PathDataHandler) CallOnData() {
	h.OnData(h.Path, h.Data)
	h.Data = nil
	h.StopLen = 0
}
