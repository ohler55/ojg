// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
)

var fnMap = map[string]Fn{}

// Fn encapsulates the information about a formula function in the package.
type Fn struct {
	Name     string
	Eval     func(root map[string]interface{}, at interface{}, args ...interface{}) interface{}
	Args     []interface{}
	Desc     string
	Compile  func(*Fn)
	compiled bool
}

// Define a function for assembly use.
func Define(f *Fn) {
	if _, has := fnMap[f.Name]; has {
		panic(fmt.Errorf("%s already defined", f.Name))
	}
	fnMap[f.Name] = *f
}

// FnDocs returns the documentation for all function.
func FnDocs() map[string]string {
	docs := map[string]string{}
	for k, f := range fnMap {
		docs[k] = f.Desc
	}
	return docs
}

// NewFn create a new named function of the named behavior.
func NewFn(name string) (fn *Fn) {
	if f, has := fnMap[name]; has {
		fn = &f
	}
	return
}

// Simplify a function in to simple types that can be encodes as JSON or SEN.
func (f *Fn) Simplify() interface{} {
	simple := make([]interface{}, 0, len(f.Args)+1)
	simple = append(simple, f.Name)
	for _, a := range f.Args {
		switch ta := a.(type) {
		case alt.Simplifier:
			simple = append(simple, ta.Simplify())
		case fmt.Stringer:
			simple = append(simple, ta.String())
		default:
			simple = append(simple, a)
		}
	}
	return simple
}

// String return a string representation of the function.
func (f *Fn) String() string {
	return sen.String(f)
}

func (f *Fn) compile() {
	if f.Compile != nil {
		f.Compile(f)
	} else {
		for i, a := range f.Args {
			if list, _ := a.([]interface{}); 0 < len(list) {
				if name, _ := list[0].(string); 0 < len(name) {
					if af := NewFn(name); af != nil {
						af.Args = list[1:]
						af.compile()
						f.Args[i] = af
					}
				}
			} else if str, _ := a.(string); 0 < len(str) && (str[0] == '$' || str[0] == '@') {
				if x, err := jp.Parse([]byte(str)); err == nil {
					f.Args[i] = x
				}
			}
		}
	}
	f.compiled = true
}

func evalArg(root map[string]interface{}, at, arg interface{}) (val interface{}) {
	switch ta := arg.(type) {
	case *Fn:
		val = ta.Eval(root, at, ta.Args...)
	case jp.Expr:
		if 0 < len(ta) {
			if _, ok := ta[0].(jp.At); ok {
				val = ta.First(at)
			} else {
				val = ta.First(root)
			}
		}
	default:
		val = arg
	}
	return val
}
