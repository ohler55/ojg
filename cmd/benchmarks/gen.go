// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ohler55/ojg/gen"
)

func genParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &gen.Parser{}
	for n := 0; n < b.N; n++ {
		if _, err := p.Parse(sample); err != nil {
			log.Fatal(err)
		}
	}
}

func genParseReader(b *testing.B) {
	var p gen.Parser
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
