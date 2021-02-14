// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/ohler55/ojg/asm"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

var sopt = sen.Options{Sort: true}

func TestDefine(t *testing.T) {
	err := defineDup()
	tt.Nil(t, err)

	err = defineDup()
	tt.NotNil(t, err)
}

func defineDup() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()
	asm.Define(&asm.Fn{Name: "dup"})

	return
}

func TestFnDocs(t *testing.T) {
	docs := asm.FnDocs()
	tt.NotNil(t, docs)
}
