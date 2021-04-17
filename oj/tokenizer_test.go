// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

type testHandler struct {
	buf []byte
}

func (h *testHandler) Null() {
	h.buf = append(h.buf, "null "...)
}

func (h *testHandler) Bool(v bool) {
	h.buf = append(h.buf, fmt.Sprintf("%t ", v)...)
}

func (h *testHandler) Int(v int64) {
	h.buf = append(h.buf, fmt.Sprintf("%d ", v)...)
}

func (h *testHandler) Float(v float64) {
	h.buf = append(h.buf, fmt.Sprintf("%g ", v)...)
}

func (h *testHandler) Number(v gen.Big) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *testHandler) String(v string) {
	h.buf = append(h.buf, fmt.Sprintf("%s ", v)...)
}

func (h *testHandler) ObjectStart() {
	h.buf = append(h.buf, '{')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ObjectEnd() {
	h.buf = append(h.buf, '}')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ArrayStart() {
	h.buf = append(h.buf, '[')
	h.buf = append(h.buf, ' ')
}

func (h *testHandler) ArrayEnd() {
	h.buf = append(h.buf, ']')
	h.buf = append(h.buf, ' ')
}

func TestTokenizerParseBasic(t *testing.T) {
	toker := oj.Tokenizer{}
	h := testHandler{}
	src := `[true,null,123,12.3]{"x":12345678901234567890}`
	err := toker.Parse([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 12345678901234567890 } ", string(h.buf))

	h.buf = h.buf[:0]
	err = toker.Parse([]byte(src), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 12345678901234567890 } ", string(h.buf))
}

func TestTokenizerLoad(t *testing.T) {
	toker := oj.Tokenizer{}
	h := testHandler{}
	err := toker.Load(strings.NewReader(`[true,null,123,12.3]{"x":3}`), &h)
	tt.Nil(t, err)
	tt.Equal(t, "[ true null 123 12.3 ] { x 3 } ", string(h.buf))
}

func TestZeroHandler(t *testing.T) {
	h := oj.ZeroHandler{}
	src := `[true,null,123,12.3]{"x":12345678901234567890}`
	err := oj.TokenizeString(src, &h)
	tt.Nil(t, err)

	err = oj.Tokenize([]byte(src), &h)
	tt.Nil(t, err)

	err = oj.TokenizeLoad(strings.NewReader(src), &h)
	tt.Nil(t, err)
}
