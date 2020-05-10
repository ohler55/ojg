// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"time"

	"github.com/ohler55/ojg/gen"
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

func Get(n gen.Node, path interface{}) (result []gen.Node) {
	if x, _ := pathToExpr(path); x != nil {
		result = x.Get(n)
	}
	return
}

func First(n gen.Node, path interface{}) (result gen.Node) {
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

func Array(n gen.Node, path interface{}) gen.Array {
	// TBD
	return nil
}

func Object(n gen.Node, path interface{}) gen.Object {
	// TBD
	return nil
}

// Set a child node value.
func Set(n, value gen.Node, path interface{}) error {
	// TBD
	return nil
}

func SetOne(n, value gen.Node, path interface{}) error {
	// TBD
	return nil
}

// Remove removes nodes returns then in an array.
func Remove(n gen.Node, path interface{}) []gen.Node {
	// TBD
	return nil
}

func RemoveOne(n gen.Node, path interface{}) gen.Node {
	// TBD
	return nil
}

func Sget(n interface{}, path interface{}) []gen.Node {
	// TBD
	return nil
}

func Sfirst(n interface{}, path interface{}) gen.Node {
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
