// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

// should make callback when end of top
type Handler interface {
	ObjectStart()
	ObjectEnd()
	ArrayStart()
	ArrayEnd()
	Null()
	Bool(value bool)
	Int(value int64)
	Float(value float64)
	Str(key string)
	Key(key string)
	Call() bool
}
