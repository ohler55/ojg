// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"io"

	"github.com/ohler55/ojg"
)

// JSON encoded output.
func JSON(data interface{}, args ...interface{}) string {
	w := Writer{
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      false,
	}
	w.config(args)
	b, _ := w.encode(data)

	return string(b)
}

// SEN encoded output.
func SEN(data interface{}, args ...interface{}) string {
	w := Writer{
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      true,
	}
	w.config(args)
	b, _ := w.encode(data)

	return string(b)
}

// JSON encoded output written to the provided io.Writer.
func WriteJSON(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := Writer{
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      false,
	}
	pw.w = w
	pw.config(args)
	_, err = pw.encode(data)

	return
}

// SEN encoded output written to the provided io.Writer.
func WriteSEN(w io.Writer, data interface{}, args ...interface{}) (err error) {
	pw := Writer{
		Options:  ojg.DefaultOptions,
		Width:    80,
		MaxDepth: 3,
		SEN:      true,
	}
	pw.w = w
	pw.config(args)
	_, err = pw.encode(data)

	return
}
