// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"io"
	"sync"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

// Options is an alias for ojg.Options
type Options = ojg.Options

var (
	// DefaultOptions are the default options for the this package.
	DefaultOptions = ojg.DefaultOptions
	// BrightOptions are the bright color options.
	BrightOptions = ojg.BrightOptions
	// HTMLOptions are the options that can be used to encode as HTML JSON.
	HTMLOptions = ojg.HTMLOptions

	// DefaultWriter is the default writer. This is not concurrent
	// safe. Individual go routine writers should be used when writing
	// concurrently.
	DefaultWriter = Writer{
		Options: ojg.DefaultOptions,
		buf:     make([]byte, 0, 1024),
	}
	writerPool = sync.Pool{
		New: func() interface{} {
			return &Writer{Options: DefaultOptions, buf: make([]byte, 0, 1024)}
		},
	}
)

// Parse a SEN byte slice into simple types. An error is returned if not valid
// JSON.
func Parse(buf []byte, args ...interface{}) (interface{}, error) {
	return DefaultParser.Parse(buf, args...)
}

// MustParse a SEN byte slice into simple types. Panics on error.
func MustParse(buf []byte, args ...interface{}) interface{} {
	val, err := DefaultParser.Parse(buf, args...)
	if err != nil {
		panic(err)
	}
	return val
}

// ParseReader a JSON io.Reader. An error is returned if not valid JSON.
func ParseReader(r io.Reader, args ...interface{}) (data interface{}, err error) {
	return DefaultParser.ParseReader(r, args...)
}

// Unmarshal parses the provided JSON and stores the result in the value
// pointed to by vp.
func Unmarshal(data []byte, vp interface{}, recomposer ...alt.Recomposer) (err error) {
	p := Parser{}
	var v interface{}
	if v, err = p.Parse(data); err == nil {
		_, err = alt.Recompose(v, vp)
	}
	return
}

// String returns a SEN string for the data provided. The data can be a simple
// type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent, *ojg.Options, or a *Writer.
func String(data interface{}, args ...interface{}) string {
	var wr *Writer
	if 0 < len(args) {
		wr = pickWriter(args[0])
	}
	if wr == nil {
		wr, _ = writerPool.Get().(*Writer)
		defer writerPool.Put(wr)
	}
	return wr.SEN(data)
}

// Bytes returns a SEN []byte for the data provided. The data can be a simple
// type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent, *ojg.Options, or a *Writer.
func Bytes(data interface{}, args ...interface{}) []byte {
	var wr *Writer
	if 0 < len(args) {
		wr = pickWriter(args[0])
	}
	if wr == nil {
		wr, _ = writerPool.Get().(*Writer)
		defer writerPool.Put(wr)
	}
	return wr.MustSEN(data)
}

// Write SEN for the data provided. The data can be a simple type of nil,
// bool, int, floats, time.Time, []interface{}, or map[string]interface{} or a
// Node type, The args, if supplied can be an int as an indent, *ojg.Options,
// or a *Writer.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	var wr *Writer
	if 0 < len(args) {
		wr = pickWriter(args[0])
	}
	if wr == nil {
		wr, _ = writerPool.Get().(*Writer)
		defer writerPool.Put(wr)
	}
	return wr.Write(w, data)
}

func pickWriter(arg interface{}) (wr *Writer) {
	switch ta := arg.(type) {
	case int:
		wr = &Writer{
			Options: ojg.GoOptions,
			buf:     make([]byte, 0, 1024),
		}
		wr.Indent = ta
	case *ojg.Options:
		wr = &Writer{
			Options: *ta,
			buf:     make([]byte, 0, 1024),
		}
	case *Writer:
		wr = ta
	}
	return
}
