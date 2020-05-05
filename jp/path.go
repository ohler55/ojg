// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"time"

	"github.com/ohler55/ojg/gd"
)

func Must(path string) Frag {
	// TBD
	return nil
}

func Parse(path string) (Frag, error) {
	// TBD
	return nil, nil
}

func Mustb(path []string) Frag {
	// TBD
	return nil
}

func Bracket(path []string) (Frag, error) {
	// TBD
	return nil, nil
}

func Get(n gd.Node, path interface{}) (result []gd.Node) {
	if x, _ := pathToExpr(path); x != nil {
		result = x.Get(n)
	}
	return
}

func First(n gd.Node, path interface{}) (result gd.Node) {
	if x, _ := pathToExpr(path); x != nil {
		result = x.First(n)
	}
	return
}

func Bool(n interface{}, path interface{}, defVal ...bool) (v bool) {
	// TBD
	return
}

func Int(n interface{}, path interface{}, defVal ...int64) (v int64) {
	// TBD
	return
}

func Float(n interface{}, path interface{}, defVal ...float64) (v float64) {
	// TBD
	return
}

func Time(n interface{}, path interface{}, defVal ...time.Time) (v time.Time) {
	// TBD
	return
}

func Array(n gd.Node, path interface{}) gd.Array {
	// TBD
	return nil
}

func Object(n gd.Node, path interface{}) gd.Object {
	// TBD
	return nil
}

// Set a child node value.
func Set(n, value gd.Node, path interface{}) error {
	// TBD
	return nil
}

func SetOne(n, value gd.Node, path interface{}) error {
	// TBD
	return nil
}

// Remove removes nodes returns then in an array.
func Remove(n gd.Node, path interface{}) []gd.Node {
	// TBD
	return nil
}

func RemoveOne(n gd.Node, path interface{}) gd.Node {
	// TBD
	return nil
}

func Sget(n interface{}, path interface{}) []gd.Node {
	// TBD
	return nil
}

func Sfirst(n interface{}, path interface{}) gd.Node {
	// TBD
	return nil
}

func Sarray(n interface{}, path interface{}) []interface{} {
	// TBD
	return nil
}

func Sobject(n interface{}, path interface{}) map[string]interface{} {
	// TBD
	return nil
}

func Sset(n, value interface{}, path interface{}) error {
	// TBD
	return nil
}

func SsetOne(n, value interface{}, path interface{}) error {
	// TBD
	return nil
}

func Sremove(n interface{}, path interface{}) []interface{} {
	// TBD
	return nil
}

func SremoveOne(n interface{}, path interface{}) interface{} {
	// TBD
	return nil
}

func pathToExpr(path interface{}) (x Expr, err error) {
	switch tp := path.(type) {
	case string:
		// TBD parse string
	case []string:
		// TBD parse each bracket
	case Expr:
		x = tp
	}
	return
}
