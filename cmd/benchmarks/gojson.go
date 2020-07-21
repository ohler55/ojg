// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func goParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var result interface{}
	for n := 0; n < b.N; n++ {
		if err := json.Unmarshal(sample, &result); err != nil {
			log.Fatal(err)
		}
	}
}

func goDecodeReader(b *testing.B) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := json.NewDecoder(f)
		for {
			var data interface{}
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func goParseChan(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	rc := make(chan interface{}, b.N)
	go func() {
		for {
			if v := <-rc; v == nil {
				break
			}
		}
	}()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var result interface{}
		if err := json.Unmarshal(sample, &result); err != nil {
			log.Fatal(err)
		}
		rc <- result
	}
	rc <- nil
}

func goValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !json.Valid(sample) {
			log.Fatal("JSON not valid")
		}
	}
}

func marshalJSON(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(data); err != nil {
			log.Fatal(err)
		}
	}
}

func marshalJSONIndent(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.MarshalIndent(data, "", "  "); err != nil {
			log.Fatal(err)
		}
	}
}

func jsonEncodeIndent(b *testing.B) {
	data := loadSample()
	var buf strings.Builder
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			log.Fatal(err)
		}
	}
}
