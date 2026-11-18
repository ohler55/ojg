// Copyright (c) 2024, Peter Ohler, All rights reserved.

package jp

import (
	"encoding/json"
	"sort"
)

// TargetRest is used by the MatchHandler to associate a Target and Rest of a
// match search.
type TargetRest struct {
	Target Expr
	// Rest is set when a Filter is included in the initializing target. Since
	// Filters can only be evaluated when there is data for the evaluation a
	// traget with a Filter is split with the pre-filter portion and the rest
	// starting with the filter.
	Rest Expr
}

// PathHandler is a TokenHandler compatible with both the oj.TokenHandler and
// the sen.TokenHandler. Fields are public to allow derived types to access
// those fields.
type MatchHandler struct {
	Targets []*TargetRest
	Path    Expr
	Stack   []any
	OnData  func(path Expr, data any)
}

// NewMatchHandler creates a new MatchHandler.
func NewMatchHandler(onData func(path Expr, data any), targets ...Expr) *MatchHandler {
	h := MatchHandler{
		Path:   R(),
		OnData: onData,
	}
	for _, target := range targets {
		tr := TargetRest{Target: target}
		for i, f := range target {
			if _, ok := f.(*Filter); ok {
				tr.Rest = target[i:]
				tr.Target = target[:i]
				break
			}
		}
		h.Targets = append(h.Targets, &tr)
	}
	return &h
}

// Null is called when a JSON null is encountered.
func (h *MatchHandler) Null() {
	h.AddValue(nil)
}

// Bool is called when a JSON true or false is encountered.
func (h *MatchHandler) Bool(v bool) {
	h.AddValue(v)
}

// Int is called when a JSON integer is encountered.
func (h *MatchHandler) Int(v int64) {
	h.AddValue(v)
}

// Float is called when a JSON decimal is encountered that fits into a
// float64.
func (h *MatchHandler) Float(v float64) {
	h.AddValue(v)
}

// Number is called when a JSON number is encountered that does not fit
// into an int64 or float64.
func (h *MatchHandler) Number(num string) {
	h.AddValue(json.Number(num))
}

// String is called when a JSON string is encountered.
func (h *MatchHandler) String(v string) {
	h.AddValue(v)
}

// ObjectStart is called when a JSON object start '{' is encountered.
func (h *MatchHandler) ObjectStart() {
	h.objArrayStart(map[string]any{}, Child(""))
}

// ObjectEnd is called when a JSON object end '}' is encountered.
func (h *MatchHandler) ObjectEnd() {
	h.objArrayEnd()
}

// Key is called when a JSON object key is encountered.
func (h *MatchHandler) Key(k string) {
	h.Path[len(h.Path)-1] = Child(k)
}

// ArrayStart is called when a JSON array start '[' is encountered.
func (h *MatchHandler) ArrayStart() {
	h.objArrayStart([]any{}, Nth(0))
}

// ArrayEnd is called when a JSON array end ']' is encountered.
func (h *MatchHandler) ArrayEnd() {
	h.objArrayEnd()
}

// AddValue is called when a leaf value is encountered.
func (h *MatchHandler) AddValue(v any) {
	if 0 < len(h.Stack) {
		switch ts := h.Stack[len(h.Stack)-1].(type) {
		case map[string]any:
			ts[string(h.Path[len(h.Path)-1].(Child))] = v
		case []any:
			h.Stack[len(h.Stack)-1] = append(ts, v)
		}
	} else if h.pathMatch(true) {
		h.OnData(h.Path, v)
	}
	h.incNth()
}

func (h *MatchHandler) objArrayStart(v any, frag Frag) {
	if 0 < len(h.Stack) {
		switch ts := h.Stack[len(h.Stack)-1].(type) {
		case map[string]any:
			ts[string(h.Path[len(h.Path)-1].(Child))] = v
		case []any:
			h.Stack[len(h.Stack)-1] = append(ts, v)
		}
		h.Stack = append(h.Stack, v)
	} else if h.pathMatch(false) {
		h.Stack = append(h.Stack, v)
	}
	h.Path = append(h.Path, frag)
}

func (h *MatchHandler) objArrayEnd() {
	h.Path = h.Path[:len(h.Path)-1]
	if 0 < len(h.Stack) {
		if len(h.Stack) == 1 {
			h.checkRest(h.Stack[0])
		}
		v := h.Stack[len(h.Stack)-1]
		h.Stack = h.Stack[:len(h.Stack)-1]
		if 0 < len(h.Stack) {
			switch ts := h.Stack[len(h.Stack)-1].(type) {
			case map[string]any:
				ts[string(h.Path[len(h.Path)-1].(Child))] = v
			case []any:
				ts[h.Path[len(h.Path)-1].(Nth)] = v
			}
		}
	}
	h.incNth()
}

func (h *MatchHandler) incNth() {
	if last := len(h.Path) - 1; 0 <= last {
		if nth, ok := h.Path[last].(Nth); ok {
			h.Path[last] = nth + 1
		}
	}
}

func (h *MatchHandler) checkRest(v any) {
	var tr *TargetRest
	for _, t := range h.Targets {
		if PathMatch(t.Target, h.Path) {
			tr = t
			break
		}
	}
	if tr == nil {
		return
	}
	if tr.Rest == nil {
		h.OnData(h.Path, v)
		return
	}
	locs := tr.Rest.Locate(v, 0)
	// Filters walk their candidates from the top of a stack, so Locate can
	// return them in reverse order. Callbacks follow document order.
	sort.Slice(locs, func(i, j int) bool {
		for k := 0; k < len(locs[i]) && k < len(locs[j]); k++ {
			switch a := locs[i][k].(type) {
			case Nth:
				if b, ok := locs[j][k].(Nth); ok && a != b {
					return int64(a) < int64(b)
				}
			case Child:
				if b, ok := locs[j][k].(Child); ok && a != b {
					return string(a) < string(b)
				}
			}
		}
		return len(locs[i]) < len(locs[j])
	})
	for _, loc := range locs {
		// Copy the path: h.Path is reused by the tokenizer and appending in
		// place would overwrite the path of an earlier callback.
		p := make(Expr, 0, len(h.Path)+len(loc))
		p = append(p, h.Path...)
		p = append(p, loc...)
		h.OnData(p, loc.First(v))
	}
}

func (h *MatchHandler) pathMatch(leaf bool) bool {
	for _, tr := range h.Targets {
		if PathMatch(tr.Target, h.Path) {
			if !leaf || tr.Rest == nil {
				return true
			}
		}
	}
	return false
}
