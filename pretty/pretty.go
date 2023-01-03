// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"io"

	"github.com/ohler55/ojg"
)

// JSON encoded output. Arguments can be used to set the writer options. An
// int sets the width while a float64 is separated into a width as the integer
// portion of the float and the 10ths sets the maximum depth per line. A bool
// sets the align option and a *ojg.Options replaces the options portion of
// the writer.
func JSON(data any, args ...any) string {
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

// SEN encoded output. Arguments can be used to set the writer options. An int
// sets the width while a float64 is separated into a width as the integer
// portion of the float and the 10ths sets the maximum depth per line. A bool
// sets the align option and a *ojg.Options replaces the options portion of
// the writer.
func SEN(data any, args ...any) string {
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

// WriteJSON encoded output written to the provided io.Writer. Arguments can
// be used to set the writer options. An int sets the width while a float64 is
// separated into a width as the integer portion of the float and the 10ths
// sets the maximum depth per line. A bool sets the align option and a
// *ojg.Options replaces the options portion of the writer.
func WriteJSON(w io.Writer, data any, args ...any) (err error) {
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

// WriteSEN encoded output written to the provided io.Writer. Arguments can be
// used to set the writer options. An int sets the width while a float64 is
// separated into a width as the integer portion of the float and the 10ths
// sets the maximum depth per line. A bool sets the align option and a
// *ojg.Options replaces the options portion of the writer.
func WriteSEN(w io.Writer, data any, args ...any) (err error) {
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
