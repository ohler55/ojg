// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
)

var fnMap = map[string]Fn{}

type Fn struct {
	Name     string
	Eval     func(root map[string]interface{}, at interface{}, args ...interface{}) interface{}
	Args     []interface{}
	Desc     string
	Compile  func(*Fn)
	compiled bool
}

func Define(f *Fn) {
	if _, has := fnMap[f.Name]; has {
		panic(fmt.Errorf("%s already defined", f.Name))
	}
	fnMap[f.Name] = *f
}

func FnDocs() map[string]string {
	docs := map[string]string{}
	for k, f := range fnMap {
		docs[k] = f.Desc
	}
	return docs
}

func NewFn(name string) (fn *Fn) {
	if f, has := fnMap[name]; has {
		fn = &f
	}
	return
}

func (f *Fn) Simplify() interface{} {
	simple := make([]interface{}, 0, len(f.Args)+1)
	simple = append(simple, f.Name)
	for _, a := range f.Args {
		if sa, _ := a.(alt.Simplifier); sa != nil {
			simple = append(simple, sa.Simplify())
		} else {
			simple = append(simple, a)
		}
	}
	return simple
}

func (f *Fn) String() string {
	return sen.String(f)
}

func (f *Fn) appendArg(arg interface{}) {
	f.Args = append(f.Args, arg)
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
			switch ta[0].(type) {
			case jp.Root:
				val = ta.First(root)
			case jp.At:
				val = ta.First(at)
			default:
				val = ta.First(root)
			}
		}
	default:
		val = arg
	}
	return val
}
