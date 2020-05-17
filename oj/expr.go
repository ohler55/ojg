// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

type Expr []Frag

func (x Expr) String() string {
	return string(x.Append(nil))
}

func (x Expr) Append(buf []byte) []byte {
	bracket := false
	for i, frag := range x {
		if _, ok := frag.(Bracket); ok {
			bracket = true
			continue
		}
		buf = frag.Append(buf, bracket, i == 0)
	}
	return buf
}

func (x Expr) Get(n interface{}) (result []interface{}) {
	if 0 < len(x) {
		return x[0].get(n, n, x[1:])
	}
	return
}

func (x Expr) GetNodes(n Node) (result []Node) {
	// TBD
	return
}

func (x Expr) First(n interface{}) (result interface{}) {
	if 0 < len(x) {
		result, _ = x[0].first(n, n, x[1:])
	}
	return
}

func (x Expr) FirstNode(n Node) (result Node) {
	// TBD
	return
}

// Set a child node value.
func (x Expr) Set(n, value interface{}) error {
	// TBD
	return nil
}

func (x Expr) SetOne(n, value interface{}) error {
	// TBD
	return nil
}

// Del removes nodes returns them in an array.
func (x Expr) Del(n interface{}) {
	// TBD
}

// Del removes nodes returns them in an array.
func (x Expr) DelOne(n interface{}) {
	// TBD
}

func X() Expr {
	return Expr{}
}

func R() Expr {
	return Expr{Root('$')}
}

func B() Expr {
	return Expr{Bracket(' ')}
}

func (x Expr) B() Expr {
	return append(x, Bracket(' '))
}

func (x Expr) C(key string) Expr {
	return append(x, Child(key))
}

func (x Expr) Child(key string) Expr {
	return append(x, Child(key))
}

func (x Expr) W() Expr {
	return append(x, Wildcard('*'))
}

func (x Expr) Wildcard() Expr {
	return append(x, Wildcard('*'))
}

func (x Expr) R() Expr {
	return append(x, Root('$'))
}

func (x Expr) Root() Expr {
	return append(x, Root('$'))
}

func (x Expr) D() Expr {
	return append(x, Descent('.'))
}

func (x Expr) Descent() Expr {
	return append(x, Descent('.'))
}
