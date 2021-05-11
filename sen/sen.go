// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"io"

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
	// GoOptions are the options that match the go json.Marshal behavior.
	GoOptions = ojg.GoOptions
	// HTMLOptions are the options that can be used to encode as HTML JSON.
	HTMLOptions = ojg.HTMLOptions

	// DefaultWriter is the default writer. This is not concurrent
	// safe. Individual go routine writers should be used when writing
	// concurrently.
	DefaultWriter = Writer{
		Options: ojg.DefaultOptions,
		buf:     make([]byte, 0, 1024),
	}
)

// Parse a JSON string in to simple types. An error is returned if not valid JSON.
func Parse(buf []byte, args ...interface{}) (interface{}, error) {
	return DefaultParser.Parse(buf, args...)
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
// as an indent or a *Options.
func String(data interface{}, args ...interface{}) string {
	wr := &DefaultWriter
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			w2 := *wr
			wr = &w2
			wr.Indent = ta
			wr.findex = 0
		case *ojg.Options:
			w2 := *wr
			wr = &w2
			wr.Options = *ta
			wr.findex = 0
		case *Writer:
			wr = ta
		}
	}
	return wr.SEN(data)
}

// Bytes returns a SEN []byte for the data provided. The data can be a simple
// type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent or a *Options.
func Bytes(data interface{}, args ...interface{}) []byte {
	wr := &DefaultWriter
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			w2 := *wr
			wr = &w2
			wr.Indent = ta
			wr.findex = 0
		case *ojg.Options:
			w2 := *wr
			wr = &w2
			wr.Options = *ta
			wr.findex = 0
		case *Writer:
			wr = ta
		}
	}
	return wr.MustSEN(data)
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	wr := &DefaultWriter
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			w2 := *wr
			wr = &w2
			wr.Indent = ta
			wr.findex = 0
		case *ojg.Options:
			w2 := *wr
			wr = &w2
			wr.Options = *ta
			wr.findex = 0
		case *Writer:
			wr = ta
		}
	}
	return wr.Write(w, data)
}
