// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import (
	"fmt"
	"io"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

type Options = ojg.Options

var (
	DefaultOptions = ojg.DefaultOptions
	BrightOptions  = ojg.BrightOptions
	GoOptions      = ojg.GoOptions
	HTMLOptions    = ojg.HTMLOptions

	DefaultWriter = Writer{
		Options: ojg.DefaultOptions,
		buf:     make([]byte, 0, 1024),
	}
	GoWriter = Writer{
		Options: ojg.GoOptions,
		buf:     make([]byte, 0, 1024),
		strict:  true,
	}
)

// Parse JSON into a gen.Node. Arguments are optional and can be a bool
// or func(interface{}) bool.
//
// A bool indicates the NoComment parser attribute should be set to the bool
// value.
//
// A func argument is the callback for the parser if processing multiple
// JSONs. If no callback function is provided the processing is limited to
// only one JSON.
func Parse(b []byte, args ...interface{}) (n interface{}, err error) {
	p := Parser{}
	return p.Parse(b, args...)
}

// ParseString is similar to Parse except it takes a string
// argument to be parsed instead of a []byte.
func ParseString(s string, args ...interface{}) (n interface{}, err error) {
	p := Parser{}
	return p.Parse([]byte(s), args...)
}

// Load a JSON from a io.Reader into a simple type. An error is returned
// if not valid JSON.
func Load(r io.Reader, args ...interface{}) (interface{}, error) {
	p := Parser{}
	return p.ParseReader(r, args...)
}

// Validate a JSON string. An error is returned if not valid JSON.
func Validate(b []byte) error {
	v := Validator{}
	return v.Validate(b)
}

// ValidateString a JSON string. An error is returned if not valid JSON.
func ValidateString(s string) error {
	v := Validator{}
	return v.Validate([]byte(s))
}

// ValidateReader a JSON stream. An error is returned if not valid JSON.
func ValidateReader(r io.Reader) error {
	v := Validator{}
	return v.ValidateReader(r)
}

// Unmarshal parses the provided JSON and stores the result in the value
// pointed to by vp.
func Unmarshal(data []byte, vp interface{}, recomposer ...*alt.Recomposer) (err error) {
	p := Parser{}
	p.num.ForceFloat = true
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

// JSON returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an
// int as an indent or a *Options.
func JSON(data interface{}, args ...interface{}) string {
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
	return wr.JSON(data)
}

// Marshal returns a JSON string for the data provided. The data can be a
// simple type of nil, bool, int, floats, time.Time, []interface{}, or
// map[string]interface{} or a Node type, The args, if supplied can be an int
// as an indent or a *Options. An error will be returned if the Option.Strict
// flag is true and a value is encountered that can not be encoded other than
// by using the %v format of the fmt package.
func Marshal(data interface{}, args ...interface{}) (out []byte, err error) {
	wr := &GoWriter
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
			wr.strict = true
		}
	}
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	wr.MustJSON(data)
	out = wr.buf

	return
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
