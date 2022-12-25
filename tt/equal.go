// Copyright (c) 2020, Peter Ohler, All rights reserved.

package tt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
)

// Equal return true if two values are equal and fails a test if not equal.
func Equal(t *testing.T, expect, actual any, args ...any) (eq bool) {
	eq = valuesEqual(expect, actual)
	if !eq {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\nexpect: (%T) %v\nactual: (%T) %v\n", expect, expect, actual, actual))
		finishFail(t, &b, args)
	}
	return
}

// NotEqual return true if two values are not equal and fails a test if equal.
func NotEqual(t *testing.T, expect, actual any, args ...any) (eq bool) {
	eq = valuesEqual(expect, actual)
	if eq {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\nexpect: (%T) %v to not equal\nactual: (%T) %v\n", expect, expect, actual, actual))
		finishFail(t, &b, args)
	}
	return
}

func valuesEqual(expect, actual any) (eq bool) {
	switch te := expect.(type) {
	case nil:
		eq = nil == actual
	case bool:
		switch ta := actual.(type) {
		case bool:
			eq = te == ta
		case gen.Bool:
			eq = te == bool(ta)
		default:
			eq = false
		}
	case gen.Bool:
		switch ta := actual.(type) {
		case gen.Bool:
			eq = bool(te) == bool(ta)
		default:
			eq = false
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, gen.Int:
		x, _ := asInt(expect)
		a, ok := asInt(actual)
		eq = x == a && ok
	case float32, float64, gen.Float:
		x, _ := asFloat(expect)
		a, ok := asFloat(actual)
		eq = x == a && ok
	case gen.Big:
		x, _ := actual.(gen.Big)
		eq = te == x
	case string:
		x, _ := asString(expect)
		a, ok := asString(actual)
		eq = x == a && ok
		if !eq && 2 < len(x) && x[0] == '/' && x[len(x)-1] == '/' {
			rx, err := regexp.Compile(x[1 : len(x)-1])
			if err != nil {
				fmt.Printf("*-*-* %q is not a valid regular expression. %s", x, err)
				return
			}
			eq = rx.MatchString(a)
		}
		/*
			if !eq {
					if !eq {
						tx, ta = colorizeStrings(tx, ta)
						expect = tx
						actual = ta
					}
			}
		*/
	case time.Time:
		tm, _ := actual.(time.Time)
		eq = tm == te
	case gen.Time:
		ta, _ := actual.(gen.Time)
		eq = ta == te

	case gen.String, json.Number:
		x, _ := asString(expect)
		a, ok := asString(actual)
		eq = x == a && ok
	case []any:
		switch ta := actual.(type) {
		case []any:
			eq = true
			for i := 0; i < len(te); i++ {
				if len(ta) <= i {
					eq = false
					break
				}
				if eq = valuesEqual(te[i], ta[i]); !eq {
					break
				}
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		case gen.Array:
			eq = valuesEqual(expect, ta.Simplify())
		default:
			eq = false
		}
	case gen.Array:
		switch ta := actual.(type) {
		case gen.Array:
			eq = true
			for i := 0; i < len(te); i++ {
				if len(ta) <= i {
					eq = false
					break
				}
				if eq = valuesEqual(te[i], ta[i]); !eq {
					break
				}
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		default:
			eq = false
		}
	case []gen.Node:
		switch ta := actual.(type) {
		case []gen.Node:
			eq = true
			for i := 0; i < len(te); i++ {
				if len(ta) <= i {
					eq = false
					break
				}
				if eq = valuesEqual(te[i], ta[i]); !eq {
					break
				}
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		default:
			eq = false
		}
	case map[string]any:
		switch ta := actual.(type) {
		case map[string]any:
			eq = true
			for k, ve := range te {
				va, has := ta[k]
				if !has {
					eq = false
					break
				}
				eq = valuesEqual(ve, va)
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		case gen.Object:
			eq = valuesEqual(expect, ta.Simplify())
		default:
			eq = false
		}
	case gen.Object:
		switch ta := actual.(type) {
		case gen.Object:
			eq = true
			for k, ve := range te {
				va, has := ta[k]
				if !has {
					eq = false
					break
				}
				eq = valuesEqual(ve, va)
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		default:
			eq = false
		}
	case alt.Path:
		if ta, ok := actual.(alt.Path); ok {
			eq = true
			for i := 0; i < len(te); i++ {
				if len(ta) <= i {
					eq = false
					break
				}
				if eq = valuesEqual(te[i], ta[i]); !eq {
					break
				}
			}
			if eq && len(te) != len(ta) {
				eq = false
			}
		}
	default:
		if actual != nil {
			eq = reflect.DeepEqual(expect, actual)
		}
	}
	return
}
