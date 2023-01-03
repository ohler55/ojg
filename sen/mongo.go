// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen

import (
	"strconv"
	"time"
)

// AddMongoFuncs adds TokenFuncs for the common mongo Javascript functions
// that appear in the output from mongosh for some types. They functions
// included are:
//
//	ISODate(arg) returns time.Time when given either a RFC3339 string or milliseconds
//	ObjectId(arg) returns the arg as a string
//	NumberInt(arg)  returns the string argument as an int64 or if too large the original string
//	NumberLong(arg)  returns the string argument as an int64 or if too large the original string
//	NumberDecimal(arg)  returns the string argument as a float64 or if too large the original string
func (p *Parser) AddMongoFuncs() {
	if p.tokenFuncs == nil {
		p.tokenFuncs = map[string]TokenFunc{}
	}
	p.tokenFuncs["ISODate"] = isoDate
	p.tokenFuncs["ObjectId"] = objectID
	p.tokenFuncs["NumberInt"] = numberInt64
	p.tokenFuncs["NumberLong"] = numberInt64
	p.tokenFuncs["NumberDecimal"] = numberDecimal
}

func isoDate(args ...any) (t any) {
	if 0 < len(args) {
		switch ta := args[0].(type) {
		case string:
			t, _ = time.Parse(time.RFC3339Nano, ta)
		case int64:
			t = time.Unix(0, ta*1_000_000).UTC()
		}
	}
	return
}

func objectID(args ...any) (v any) {
	if 0 < len(args) {
		v = args[0]
	}
	return
}

func numberInt64(args ...any) (v any) {
	if 0 < len(args) {
		s, _ := args[0].(string)
		var err error
		if v, err = strconv.ParseInt(s, 10, 64); err != nil {
			v = args[0]
		}
	}
	return
}

func numberDecimal(args ...any) (v any) {
	if 0 < len(args) {
		s, _ := args[0].(string)
		var err error
		if v, err = strconv.ParseFloat(s, 64); err != nil {
			v = args[0]
		}
	}
	return
}
