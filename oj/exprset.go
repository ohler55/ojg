// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

type delFlagType struct {
}

var delFlag = &delFlagType{}

// Del removes matching nodes.
func (x Expr) Del(n interface{}) {
	_ = x.Set(n, delFlag)
}

// Del removes at most one node.
func (x Expr) DelOne(n interface{}) {
	_ = x.SetOne(n, delFlag)
}

// Set all matching child node values. An error is returned if it is not
// possible. If the path to the child does not exist array and map elements
// are added.
func (x Expr) Set(n, value interface{}) error {
	// TBD
	return nil
}

// Set a child node value. An error is returned if it is not possible. If the
// path to the child does not exist array and map elements are added.
func (x Expr) SetOne(n, value interface{}) error {
	// TBD
	return nil
}
