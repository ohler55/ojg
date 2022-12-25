// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/oj"
)

func ojParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &oj.Parser{}
	for n := 0; n < b.N; n++ {
		if _, err := p.Parse(sample); err != nil {
			panic(err)
		}
	}
}

func ojParseReuse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		if _, err := p.Parse(sample); err != nil {
			panic(err)
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
			panic(err)
		}
	}
}

func ojParseReaderReuse(b *testing.B) {
	p := oj.Parser{Reuse: true}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, err = p.ParseReader(f); err != nil {
			panic(err)
		}
	}
}

func ojTokenize(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	h := oj.ZeroHandler{}
	t := oj.Tokenizer{}
	for n := 0; n < b.N; n++ {
		if err := t.Parse(sample, &h); err != nil {
			panic(err)
		}
	}
}

func ojTokenizeLoad(b *testing.B) {
	t := oj.Tokenizer{}
	h := oj.ZeroHandler{}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if err := t.Load(f, &h); err != nil {
			panic(err)
		}
	}
}

func ojMarshalCatalog(b *testing.B) {
	sample, _ := ioutil.ReadFile(catFilename)
	var cat Catalog
	if err := json.Unmarshal(sample, &cat); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := oj.Marshal(&cat); err != nil {
			panic(err)
		}
	}
}

func ojMarshalPatient(b *testing.B) {
	sample, _ := ioutil.ReadFile(patFilename)
	var patient Patient
	if err := json.Unmarshal(sample, &patient); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := oj.Marshal(&patient); err != nil {
			panic(err)
		}
	}
}

func ojUnmarshalPatient(b *testing.B) {
	sample, _ := ioutil.ReadFile(patFilename)
	p := oj.Parser{Reuse: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var out Patient
		if err := p.Unmarshal(sample, &out); err != nil {
			panic(err)
		}
	}
}

func ojUnmarshalCatalog(b *testing.B) {
	sample, _ := ioutil.ReadFile(catFilename)
	p := oj.Parser{Reuse: true}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var out Catalog
		if err := p.Unmarshal(sample, &out); err != nil {
			panic(err)
		}
	}
}

func ojParseChan(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	rc := make(chan any, b.N)
	ready := make(chan bool)
	go func() {
		ready <- true
		for {
			if v := <-rc; v == nil {
				break
			}
		}
	}()
	<-ready
	b.ResetTimer()
	var p oj.Parser
	for n := 0; n < b.N; n++ {
		if _, err := p.Parse(sample, rc); err != nil {
			panic(err)
		}
	}
	rc <- nil
}

func ojValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var v oj.Validator
	for n := 0; n < b.N; n++ {
		if err := v.Validate(sample); err != nil {
			panic(err)
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
			panic(err)
		}
	}
}

func ojJSON(b *testing.B) {
	data := loadSample()
	wr := oj.Writer{Options: ojg.Options{OmitNil: true}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = wr.MustJSON(data)
	}
}

func ojJSONIndent(b *testing.B) {
	data := loadSample()
	wr := oj.Writer{Options: ojg.Options{OmitNil: true, Indent: 2}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = wr.MustJSON(data)
	}
}

// JSON indented and sorted
func ojJSONSort(b *testing.B) {
	data := loadSample()
	wr := oj.Writer{Options: ojg.Options{OmitNil: true, Indent: 2, Sort: true}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = wr.MustJSON(data)
	}
}

func ojWriteIndent(b *testing.B) {
	data := loadSample()
	var w noWriter
	wr := oj.Writer{Options: ojg.Options{OmitNil: true, Indent: 2}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wr.MustWrite(w, data)
	}
}
