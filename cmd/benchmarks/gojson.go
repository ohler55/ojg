// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	var result any
	for n := 0; n < b.N; n++ {
		if err := json.Unmarshal(sample, &result); err != nil {
			panic(err)
		}
	}
}

func goUnmarshalPatient(b *testing.B) {
	sample, _ := ioutil.ReadFile(patFilename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var out Patient
		if err := json.Unmarshal(sample, &out); err != nil {
			panic(err)
		}
	}
}

func goUnmarshalCatalog(b *testing.B) {
	sample, _ := ioutil.ReadFile(catFilename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var out Catalog
		if err := json.Unmarshal(sample, &out); err != nil {
			panic(err)
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
			var data any
			if err := dec.Decode(&data); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				panic(err)
			}
		}
	}
}

func goDecode(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	for n := 0; n < b.N; n++ {
		dec := json.NewDecoder(bytes.NewReader(sample))
		for {
			_, err := dec.Token()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				panic(err)
			}
		}
	}
}

func goParseChan(b *testing.B) {
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
	for n := 0; n < b.N; n++ {
		var result any
		// The go json package does not have a chan based result handler so
		// fake it to set the baseline for others.
		if err := json.Unmarshal(sample, &result); err != nil {
			panic(err)
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

func goMarshalCatalog(b *testing.B) {
	sample, _ := ioutil.ReadFile(catFilename)
	var cat Catalog
	if err := json.Unmarshal(sample, &cat); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(&cat); err != nil {
			panic(err)
		}
	}
}

func goMarshalPatient(b *testing.B) {
	sample, _ := ioutil.ReadFile(patFilename)
	var patient Patient
	if err := json.Unmarshal(sample, &patient); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(&patient); err != nil {
			panic(err)
		}
	}
}

func marshalJSON(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.Marshal(data); err != nil {
			panic(err)
		}
	}
}

func marshalJSONIndent(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := json.MarshalIndent(data, "", "  "); err != nil {
			panic(err)
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
			panic(err)
		}
	}
}
