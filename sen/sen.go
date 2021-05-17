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

	writerPool = sync.Pool{
		New: func() interface{} {
			return &Writer{Options: DefaultOptions, buf: make([]byte, 0, 1024)}
		},
	}
	parserPool = sync.Pool{
		New: func() interface{} {
			return &Parser{}
		},
	}
)

// Parse SEN into a simple type. Arguments are optional and can be a
// func(interface{}) bool for callbacks or a chan interface{} for chan based
// result delivery. The SEN parser will also Parse JSON.
//
// A func argument is the callback for the parser if processing multiple
// SENs. If no callback function is provided the processing is limited to
// only one SEN.
//
// A chan argument will be used to deliver parse results.
func Parse(buf []byte, args ...interface{}) (interface{}, error) {
	p, _ := parserPool.Get().(*Parser)
	p.cb = nil
	p.resultChan = nil
	p.OnlyOne = false
	p.Reuse = false
	defer parserPool.Put(p)
	return p.Parse(buf, args...)
}

// MustParse SEN into a simple type. Arguments are optional and can be a
// func(interface{}) bool for callbacks or a chan interface{} for chan based
// result delivery. The SEN parser will also Parse JSON. Panics on error.
//
// A func argument is the callback for the parser if processing multiple
// SENs. If no callback function is provided the processing is limited to
// only one SEN.
//
// A chan argument will be used to deliver parse results.
func MustParse(buf []byte, args ...interface{}) interface{} {
	p := parserPool.Get().(*Parser)
	p.cb = nil
	p.resultChan = nil
	p.OnlyOne = false
	p.Reuse = false
	defer parserPool.Put(p)
	val, err := p.Parse(buf, args...)
	if err != nil {
		panic(err)
	}
	return val
}

// ParseReader reads and parses SEN into a simple type. Arguments are optional
// and can be a func(interface{}) bool for callbacks or a chan interface{} for
// chan based result delivery. The SEN parser will also Parse JSON.
//
// A func argument is the callback for the parser if processing multiple
// SENs. If no callback function is provided the processing is limited to
// only one SEN.
//
// A chan argument will be used to deliver parse results.
func ParseReader(r io.Reader, args ...interface{}) (data interface{}, err error) {
	p, _ := parserPool.Get().(*Parser)
	p.cb = nil
	p.resultChan = nil
	p.OnlyOne = false
	p.Reuse = false
	defer parserPool.Put(p)
	return p.ParseReader(r, args...)
}

// MustParseReader reads and parses SEN into a simple type. Arguments are
// optional and can be a func(interface{}) bool for callbacks or a chan
// interface{} for chan based result delivery. The SEN parser will also Parse
// JSON. Panics on error.
//
// A func argument is the callback for the parser if processing multiple
// SENs. If no callback function is provided the processing is limited to
// only one SEN.
//
// A chan argument will be used to deliver parse results.
func MustParseReader(r io.Reader, args ...interface{}) (data interface{}) {
	p := parserPool.Get().(*Parser)
	p.cb = nil
	p.resultChan = nil
	p.OnlyOne = false
	p.Reuse = false
	defer parserPool.Put(p)
	var err error
	if data, err = p.ParseReader(r, args...); err != nil {
		panic(err)
	}
	return
}

// Unmarshal parses the provided JSON and stores the result in the value
// pointed to by vp.
func Unmarshal(data []byte, vp interface{}, recomposer ...*alt.Recomposer) (err error) {
	p := Parser{}
	var v interface{}
	if v, err = p.Parse(data); err == nil {
		if 0 < len(recomposer) {
			_, err = recomposer[0].Recompose(v, vp)
		} else {
			_, err = alt.Recompose(v, vp)
		}
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

// MustWrite SEN for the data provided. The data can be a simple type of nil,
// bool, int, floats, time.Time, []interface{}, or map[string]interface{} or a
// Node type, The args, if supplied can be an int as an indent, *ojg.Options,
// or a *Writer. Panics on error.
func MustWrite(w io.Writer, data interface{}, args ...interface{}) {
	if err := Write(w, data, args...); err != nil {
		panic(err)
	}
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
