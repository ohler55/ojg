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

// Bool returns a bool for the value given. If the value give can not be cast
// to a bool the option default (defVal) is returned. If no defVal is given
// then false is returned.
func (x Expr) Bool(n interface{}, defVal ...bool) (v bool) {
	if 0 < len(defVal) {
		v = defVal[0]
	}
	// TBD first...
	switch tn := n.(type) {
	case gd.Bool:
		v = bool(tn)
	case bool:
		v = tn
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

func (x Expr) Sget(n interface{}) []interface{} {
	// TBD
	return nil
}

func (x Expr) Sfirst(n interface{}) interface{} {
	// TBD
	return nil
}

func (x Expr) Sarray(n interface{}) []interface{} {
	// TBD
	return nil
}

func (x Expr) Sobject(n interface{}) map[string]interface{} {
	// TBD
	return nil
}

func (x Expr) Sset(n, value interface{}) error {
	// TBD
	return nil
}

func (x Expr) SsetOne(n, value interface{}) error {
	// TBD
	return nil
}

func (x Expr) Sremove(n interface{}) []interface{} {
	// TBD
	return nil
}

func (x Expr) SremoveOne(n interface{}) interface{} {
	// TBD
	return nil
}
