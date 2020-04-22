// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg

// should make callback when end of top
type Handler interface {
	ObjectHandler
	ArrayHandler
	NullHandler
	BoolHandler
	IntHandler
	FloatHandler
	KeyHandler
	ErrorHandler
	Caller
}

type ErrorHandler interface {
	Error(err error, line, col int64)
}

type ObjectHandler interface {
	ObjectStart()
	ObjectEnd()
}

type ArrayHandler interface {
	ArrayStart()
	ArrayEnd()
}

type NullHandler interface {
	Null()
}

type IntHandler interface {
	Int(value int64)
}

type FloatHandler interface {
	Float(value float64)
}

type BoolHandler interface {
	Bool(value bool)
}

type KeyHandler interface {
	KeyO(key string)
}

type Caller interface {
	Call()
}
