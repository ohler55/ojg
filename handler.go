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
	KeyObjectStart(key string)
	ObjectStart()
	ObjectEnd()
}

type ArrayHandler interface {
	KeyArrayStart(string)
	ArrayStart()
	ArrayEnd()
}

type NullHandler interface {
	KeyNull(key string)
	Null()
}

type IntHandler interface {
	KeyInt(key string, value int64)
	Int(value int64)
}

type FloatHandler interface {
	KeyFloat(key string, value float64)
	Float(value float64)
}

type BoolHandler interface {
	KeyBool(key string, value bool)
	Bool(value bool)
}

type Caller interface {
	Call()
}
