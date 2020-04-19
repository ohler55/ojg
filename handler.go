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
	Null(key *string)
}

type IntHandler interface {
	Int(key *string, value int64)
}

type FloatHandler interface {
	Float(key *string, value float64)
}

type BoolHandler interface {
	Bool(key *string, value bool)
}

type Caller interface {
	Call()
}
