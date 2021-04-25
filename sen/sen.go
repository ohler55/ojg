// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"fmt"
	"io"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
)

type Options = ojg.Options

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
	wr := &Writer{
		Options: ojg.DefaultOptions,
	}
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			wr.Indent = ta
		case *ojg.Options:
			wr.Options = *ta
		}
	}
	if wr.InitSize == 0 {
		wr.InitSize = 256
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
		}
	}()
	wr.buildSen(data, 0)

	return string(wr.buf)
}

// Write a JSON string for the data provided. The data can be a simple type of
// nil, bool, int, floats, time.Time, []interface{}, or map[string]interface{}
// or a Node type, The args, if supplied can be an int as an indent or a
// *Options.
func Write(w io.Writer, data interface{}, args ...interface{}) (err error) {
	wr := &Writer{
		Options: ojg.DefaultOptions,
	}
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case int:
			wr.Indent = ta
		case *ojg.Options:
			wr.Options = *ta
		}
	}
	wr.w = w
	if wr.InitSize == 0 {
		wr.InitSize = 256
	}
	if wr.WriteLimit == 0 {
		wr.WriteLimit = 1024
	}
	if cap(wr.buf) < wr.InitSize {
		wr.buf = make([]byte, 0, wr.InitSize)
	} else {
		wr.buf = wr.buf[:0]
	}
	defer func() {
		if r := recover(); r != nil {
			wr.buf = wr.buf[:0]
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if wr.Color {
		wr.cbuildJSON(data, 0)
	} else {
		wr.buildSen(data, 0)
	}
	if w != nil && 0 < len(wr.buf) {
		_, err = wr.w.Write(wr.buf)
	}
	return
}
