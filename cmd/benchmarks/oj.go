// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ohler55/ojg/oj"
)

func ojParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &oj.Parser{}
	for n := 0; n < b.N; n++ {
		if _, err := p.Parse(sample); err != nil {
			log.Fatal(err)
		}
	}
}

func ojParseReader(b *testing.B) {
	var p oj.Parser
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, err = p.ParseReader(f); err != nil {
			log.Fatal(err)
		}
	}
}

func ojValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var v oj.Validator
	for n := 0; n < b.N; n++ {
		if err := v.Validate(sample); err != nil {
			log.Fatal(err)
		}
	}
}

func ojValidateReader(b *testing.B) {
	var v oj.Validator
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to read %s. %s\n", filename, err)
		return
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if err := v.ValidateReader(f); err != nil {
			log.Fatal(err)
		}
	}
}

func ojJSON(b *testing.B) {
	data := loadSample()
	opt := oj.Options{OmitNil: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojJSONIndent(b *testing.B) {
	data := loadSample()
	opt := oj.Options{OmitNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

// JSON indented and sorted
func ojJSONSort(b *testing.B) {
	data := loadSample()
	opt := oj.Options{OmitNil: true, Indent: 2, Sort: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}

func ojWriteIndent(b *testing.B) {
	data := loadSample()
	var w noWriter
	b.ResetTimer()
	opt := oj.Options{OmitNil: true, Indent: 2}
	for n := 0; n < b.N; n++ {
		if err := oj.Write(w, data, &opt); err != nil {
			log.Fatal(err)
		}
	}
}
