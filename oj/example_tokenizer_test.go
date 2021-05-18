// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"bytes"
	"fmt"

	"github.com/ohler55/ojg/oj"
)

type Toker struct {
	buf []byte
}

func (h *Toker) Null() {
	h.buf = append(h.buf, "null "...)
}

func (h *Toker) Bool(v bool) {
	h.buf = append(h.buf, fmt.Sprintf("%t ", v)...)
}

func (h *Toker) Int(v int64) {
	h.buf = append(h.buf, fmt.Sprintf("%d ", v)...)
}

func (h *Toker) Float(v float64) {
	h.buf = append(h.buf, fmt.Sprintf("%g ", v)...)
}

func (h *Toker) Number(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *Toker) String(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *Toker) ObjectStart() {
	h.buf = append(h.buf, '{')
	h.buf = append(h.buf, ' ')
}

func (h *Toker) ObjectEnd() {
	h.buf = append(h.buf, '}')
	h.buf = append(h.buf, ' ')
}

func (h *Toker) Key(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s: ", v)...)
}

func (h *Toker) ArrayStart() {
	h.buf = append(h.buf, '[')
	h.buf = append(h.buf, ' ')
}

func (h *Toker) ArrayEnd() {
	h.buf = append(h.buf, ']')
	h.buf = append(h.buf, ' ')
}

func ExampleTokenizer_Parse() {
	toker := oj.Tokenizer{}
	h := Toker{}
	src := `[true,null,123,12.3]{"x":12345678901234567890}`
	if err := toker.Parse([]byte(src), &h); err != nil {
		panic(err)
	}
	fmt.Println(string(bytes.TrimSpace(h.buf)))

	// Output: [ true null 123 12.3 ] { x: 12345678901234567890 }
}
