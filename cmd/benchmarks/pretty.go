// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"log"
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
)

func prettyJSON(b *testing.B) {
	data := loadSample()
	opt := sen.Options{OmitNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = pretty.JSON(data, &opt)
	}
}

func prettySEN(b *testing.B) {
	data := loadSample()
	opt := sen.Options{OmitNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = pretty.SEN(data, &opt)
	}
}

func prettyWriteJSON(b *testing.B) {
	data := loadSample()
	var w noWriter
	b.ResetTimer()
	opt := oj.Options{OmitNil: true, Indent: 2}
	for n := 0; n < b.N; n++ {
		if err := pretty.WriteJSON(w, data, &opt); err != nil {
			log.Fatal(err)
		}
	}
}

func prettyWriteSEN(b *testing.B) {
	data := loadSample()
	var w noWriter
	b.ResetTimer()
	opt := oj.Options{OmitNil: true, Indent: 2}
	for n := 0; n < b.N; n++ {
		if err := pretty.WriteSEN(w, data, &opt); err != nil {
			log.Fatal(err)
		}
	}
}
