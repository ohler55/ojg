// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strings"
	"time"

	"github.com/ohler55/ojg/gd"
)

type Expr []Frag

func (x Expr) String() string {
	var b strings.Builder

	// TBD check for bracketed as first element

	return b.String()
}

func (x Expr) Get(n gd.Node) (result []gd.Node) {
	// TBD
	return
}

func (x Expr) First(n gd.Node) (result gd.Node) {
	// TBD
	return
}

func (x Expr) Bool(n gd.Node, defVal ...bool) (v bool) {
	if 0 < len(defVal) {
		v = defVal[0]
	}
	if n := x.First(n); n != nil {
		switch tn := n.(type) {
		case gd.Bool:
			v = bool(tn)
			//case bool:
			//v = tn
		}
	}
	return
}

func (x Expr) Int(n gd.Node, defVal ...int64) (v int64) {
	// TBD
	return
}

func (x Expr) Float(n gd.Node, defVal ...float64) (v float64) {
	// TBD
	return
}

func (x Expr) Time(n gd.Node, defVal ...time.Time) (v time.Time) {
	// TBD
	return
}

func (x Expr) Array(n gd.Node) gd.Array {
	// TBD
	return nil
}

func (x Expr) Object(n gd.Node) gd.Object {
	// TBD
	return nil
}

// Set a child node value.
func (x Expr) Set(n, value gd.Node) error {
	// TBD
	return nil
}

func (x Expr) SetOne(n, value gd.Node) error {
	// TBD
	return nil
}

// Remove removes nodes returns then in an array.
func (x Expr) Remove(n gd.Node) []gd.Node {
	// TBD
	return nil
}

func (x Expr) RemoveOne(n gd.Node) gd.Node {
	// TBD
	return nil
}

func (x Expr) Iget(n interface{}) []gd.Node {
	// TBD
	return nil
}

func (x Expr) Ifirst(n interface{}) gd.Node {
	// TBD
	return nil
}

func (x Expr) Ibool(n interface{}, defVal ...bool) (v bool) {
	// TBD
	return
}

func (x Expr) Iint(n interface{}, defVal ...int64) (v int64) {
	// TBD
	return
}

func (x Expr) Ifloat(n interface{}, defVal ...float64) (v float64) {
	// TBD
	return
}

func (x Expr) Itime(n interface{}, defVal time.Time) (v time.Time) {
	// TBD
	return
}

func (x Expr) Iarray(n interface{}) []interface{} {
	// TBD
	return nil
}

func (x Expr) Iobject(n interface{}) map[string]interface{} {
	// TBD
	return nil
}

func (x Expr) Iset(n, value interface{}) error {
	// TBD
	return nil
}

func (x Expr) IsetOne(n, value interface{}) error {
	// TBD
	return nil
}

func (x Expr) Iremove(n interface{}) []interface{} {
	// TBD
	return nil
}

func (x Expr) IremoveOne(n interface{}) interface{} {
	// TBD
	return nil
}

/*
	Get(n gd.Node) []gd.Node
	First(n gd.Node) gd.Node

	Bool(n gd.Node, defVal ...bool) bool
	Int(n gd.Node, defVal ...int64) int64
	Float(n gd.Node, defVal ...float64) float64
	Time(n gd.Node, defVal ...time.Time) time.Time
	Array(n gd.Node) gd.Array
	Object(n gd.Node) gd.Object

	// Set a child node value.
	Set(n, value gd.Node) error
	SetOne(n, value gd.Node) error

	// Remove removes nodes returns then in an array.
	Remove(n gd.Node) []gd.Node
	RemoveOne(n gd.Node) gd.Node

	Iget(n interface{}) []gd.Node
	Ifirst(n interface{}) gd.Node

	Ibool(n interface{}, defVal ...bool) bool
	Iint(n interface{}, defVal ...int64) int64
	Ifloat(n interface{}, defVal ...float64) float64
	Itime(n interface{}, defVal ...time.Time) time.Time
	Iarray(n interface{}) []interface{}
	Iobject(n interface{}) map[string]interface{}

	Iset(n, value interface{}) error
	IsetOne(n, value interface{}) error

	Iremove(n interface{}) []interface{}
	IremoveOne(n interface{}) interface{}
*/
