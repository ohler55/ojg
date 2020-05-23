// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

const (
	equal  = "=="
	length = "length"
)

var (
	//eq = op([]byte{'\x00', '=', '=', '='})
	eq = &op{prec: 0, code: '=', name: []byte("==")}
)

type op struct {
	name []byte
	prec byte
	code byte
}

// Script represents JSON Path script used in filters as well.
type Script []interface{}

type Filter Script

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (s Script) Append(buf []byte) []byte {

	// TBD

	return buf
}

func (f Filter) Eval(stack []interface{}, data interface{}) []interface{} {
	//fmt.Println("*** Eval")
	switch td := data.(type) {
	case []interface{}:
		estack := make([]interface{}, len(f))
		for _, v := range td {
			//fmt.Printf("*** checking %d: %v\n", foo, v)
			// Eval filter
			copy(estack, f)
			// resolve all expr members
			for i, ev := range estack {
			Normalize:
				switch x := ev.(type) {
				case Expr:
					if m, ok := v.(map[string]interface{}); ok && len(x) == 2 {
						if _, ok = x[0].(At); ok {
							var c Child
							if c, ok = x[1].(Child); ok {
								ev = m[string(c)]
								estack[i] = ev
								goto Normalize
							}
						}
					}
					ev = x.First(v)
					estack[i] = ev
					goto Normalize
				case int:
					estack[i] = int64(x)
				case int8:
					estack[i] = int64(x)
				case int16:
					estack[i] = int64(x)
				case int32:
					estack[i] = int64(x)
				case uint:
					estack[i] = int64(x)
				case uint8:
					estack[i] = int64(x)
				case uint16:
					estack[i] = int64(x)
				case uint32:
					estack[i] = int64(x)
				case uint64:
					estack[i] = int64(x)
				case float32:
					estack[i] = float64(x)
				default:
					//fmt.Printf("*** %T %v\n", x, x)
					// TBD normalize to simple types
				}
			}
			//fmt.Printf("*** 2. estack: %s\n", JSON(estack))
			for i := len(estack) - 1; 0 <= i; i-- {
				o, _ := estack[i].(*op)
				if o == nil {
					continue
				}
				switch o.code {
				case eq.code:
					if len(estack) <= i+2 {
						// TBD bad script
						fmt.Printf("******* bad script - %s\n", JSON(estack))
						return stack
					}
					//fmt.Printf("*** EQ %T %s %T %s\n", estack[i+1], JSON(estack[i+1]), estack[i+2], JSON(estack[i+2]))
					estack[i] = false
					switch left := estack[i+1].(type) {
					case int64:
						right, ok := estack[i+2].(int64)
						if ok && left == right {
							estack[i] = true
						}
					case int:
						right, ok := estack[i+2].(int)
						if ok && left == right {
							estack[i] = true
						}
					}
				}
			}
			//fmt.Printf("*** estack at end %s\n", JSON(estack))
			if b, _ := estack[0].(bool); b {
				stack = append(stack, v)
			}
		}
		for i, _ := range estack {
			estack[i] = nil
		}
	}
	return stack
}

func (s Script) Foo() Script {
	s = append(s, eq)
	s = append(s, A().C("a"))
	s = append(s, int64(52))
	return s
}

// TBD if list then walk and check each one
//  stack based again
//  [op, value, value]
//  [op, value, op, value, value, op, value, value]
// start from end and walk back to op
//  [&&, ==, @.foo, 3, >, 4, @.bar]
//  [&&, ==, @.foo, 3, true]
//  [&&, true, true]

// build like Expr
// Script(@.foo, ==, 3, &&, 4 < @.bar)
//  [@.foo, 3, ==]
//  [&&, @.foo, 3, ==]
//  [&&, @.foo, 3, ==, 4, @.bar, <]
// Script(@.foo, ==, 3, &&, (, 4 < @.bar, ) )

// bytes are ( and ), maybe + * - /

// should op be []byte{precedence,code,string} or a struct?

// startb out as public for ops then make private after parser is ready
