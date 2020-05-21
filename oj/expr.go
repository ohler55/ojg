// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// Expr is a JSON path expression composed of fragments.
type Expr []Frag

const (
	fragIndexMask    = 0x0000ffff
	descentFlag      = 0x00010000
	descentChildFlag = 0x00020000
)

// String returns a string representation of the expression.
func (x Expr) String() string {
	return string(x.Append(nil))
}

// Append a string representation of the expression to a byte slice and return
// the expanded buffer.
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

// X creates an empty Expr.
func X() Expr {
	return Expr{}
}

// R creates an Expr with a Root fragment.
func R() Expr {
	return Expr{Root('$')}
}

// B creates an Expr with a Bracket fragment.
func B() Expr {
	return Expr{Bracket(' ')}
}

// D creates an Expr with a recursive Descent fragment.
func D() Expr {
	return Expr{Descent('.')}
}

// W creates an Expr with a Wildcard fragment.
func W() Expr {
	return Expr{Wildcard('*')}
}

// C creates an Expr with a Child fragment.
func C(key string) Expr {
	return Expr{Child(key)}
}

// N creates an Expr with an Nth fragment.
func N(n int) Expr {
	return Expr{Nth(n)}
}

// U creates an Expr with an Union fragment.
func U(keys []string, indexes []int) Expr {
	return Expr{&Union{Keys: keys, Indexes: indexes}}
}

// S creates an Expr with a Slice fragment.
func S(start, end, step int) Expr {
	return Expr{&Slice{Start: start, End: end, Step: step}}
}

// A appends an At fragment to the Expr.
func (x Expr) A() Expr {
	return append(x, At('@'))
}

// At appends an At fragment to the Expr.
func (x Expr) At() Expr {
	return append(x, At('@'))
}

// B appends a Bracket fragment to the Expr.
func (x Expr) B() Expr {
	return append(x, Bracket(' '))
}

// C appends a Child fragment to the Expr.
func (x Expr) C(key string) Expr {
	return append(x, Child(key))
}

// Child appends a Child fragment to the Expr.
func (x Expr) Child(key string) Expr {
	return append(x, Child(key))
}

// W appends a Wildcard fragment to the Expr.
func (x Expr) W() Expr {
	return append(x, Wildcard('*'))
}

// Wildcard appends a Wildcard fragment to the Expr.
func (x Expr) Wildcard() Expr {
	return append(x, Wildcard('*'))
}

// N appends an Nth fragment to the Expr.
func (x Expr) N(n int) Expr {
	return append(x, Nth(n))
}

// Nth appends an Nth fragment to the Expr.
func (x Expr) Nth(n int) Expr {
	return append(x, Nth(n))
}

// R appends a Root fragment to the Expr.
func (x Expr) R() Expr {
	return append(x, Root('$'))
}

// Root appends a Root fragment to the Expr.
func (x Expr) Root() Expr {
	return append(x, Root('$'))
}

// D appends a recursive Descent fragment to the Expr.
func (x Expr) D() Expr {
	return append(x, Descent('.'))
}

// Descent appends a recursive Descent fragment to the Expr.
func (x Expr) Descent() Expr {
	return append(x, Descent('.'))
}

// U appends a Union fragment to the Expr.
func (x Expr) U(keys []string, indexes []int) Expr {
	return append(x, &Union{Keys: keys, Indexes: indexes})
}

// Union appends a Union fragment to the Expr.
func (x Expr) Union(keys []string, indexes []int) Expr {
	return append(x, &Union{Keys: keys, Indexes: indexes})
}

// S appends a Slice fragment to the Expr.
func (x Expr) S(start, end, step int) Expr {
	return append(x, &Slice{Start: start, End: end, Step: step})
}

// Slice appends a Slice fragment to the Expr.
func (x Expr) Slice(start, end, step int) Expr {
	return append(x, &Slice{Start: start, End: end, Step: step})
}

// The easy way to implement the Get is to have each fragment handle the
// getting using recursion. The overhead of a go function call is rather high
// though so instead a psuedo call stack is implemented here that grows and
// shrinks as the getting takes place. The fragment index if placed on the
// stack as well mostly for a small degree of simplicity in what a few people
// might find a complex approach to the solution. Its twice as fast as the
// recursive function call approach.

// Get the elements of the data identified by the path.
func (x Expr) Get(data interface{}) (results []interface{}) {
	if len(x) == 0 {
		return
	}
	var v interface{}
	var prev interface{}

	stack := make([]interface{}, 0, 64)
	stack = append(stack, data)

	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(int); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]interface{}:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							fi++
							f = x[fi]
							stack = append(stack, fi)
						}
					}
				}
			case Object:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							fi++
							f = x[fi]
							stack = append(stack, fi)
						}
					}
				}
			default:
				// TBD try reflection
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Array:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case []interface{}:
				if fi == len(x)-1 { // last one
					results = append(results, tv...)
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Object:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Array:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(int)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				} else if fi < len(x)-1 {
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			} else {
				if fi == len(x)-1 { // last one
					if top {
						results = append(results, prev)
					}
				} else {
					stack = append(stack, prev)
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Root:
			if fi == len(x)-1 { // last one
				results = append(results, data)
			} else {
				stack = append(stack, data)
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		case At, Bracket:
			if fi == len(x)-1 { // last one
				results = append(results, prev)
			} else {
				stack = append(stack, prev)
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		}
	}
	// Free up anything still on the stack.
	stack = stack[0:cap(stack)]
	for i := len(stack) - 1; 0 <= i; i-- {
		stack[i] = nil
	}
	return
}

// First element of the data identified by the path.
func (x Expr) First(data interface{}) interface{} {
	if len(x) == 0 {
		return nil
	}
	var v interface{}
	var prev interface{}

	stack := make([]interface{}, 0, 64)
	defer func() {
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	}()
	stack = append(stack, data)
	f := x[0]
	fi := 0 // frag index
	stack = append(stack, fi)

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(int); up {
			stack = stack[:len(stack)-1]
			fi = ii & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			switch tv := prev.(type) {
			case map[string]interface{}:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						return v
					}
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
						fi++
						f = x[fi]
						stack = append(stack, fi)
					}
				}
			case Object:
				if v, has = tv[string(tf)]; has {
					if fi == len(x)-1 { // last one
						return v
					}
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
						fi++
						f = x[fi]
						stack = append(stack, fi)
					}
				}
			default:
				// TBD try reflection
			}
		case Nth:
			i := int(tf)
			switch tv := prev.(type) {
			case []interface{}:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					return v
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Array:
				if i < 0 {
					i = len(tv) + i
				}
				var v interface{}
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == len(x)-1 { // last one
					return v
				} else {
					switch v.(type) {
					case map[string]interface{}, []interface{}, Object, Array:
						stack = append(stack, v)
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case map[string]interface{}:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case []interface{}:
				if fi == len(x)-1 { // last one
					if 0 < len(tv) {
						return tv[0]
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Object:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			case Array:
				if fi == len(x)-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
						}
					}
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(int)
			top := (di & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (di & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case map[string]interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case []interface{}:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == len(x)-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case map[string]interface{}, []interface{}, Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, fi|descentChildFlag)
				} else if fi < len(x)-1 {
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			} else {
				if fi == len(x)-1 { // last one
					if top {
						return prev
					}
				} else {
					stack = append(stack, prev)
					fi++
					f = x[fi]
					stack = append(stack, fi)
				}
			}
		case Root:
			if fi == len(x)-1 { // last one
				return data
			}
			stack = append(stack, data)
			fi++
			f = x[fi]
			stack = append(stack, fi)
		case At, Bracket:
			if fi == len(x)-1 { // last one
				return prev
			}
			stack = append(stack, prev)
			fi++
			f = x[fi]
			stack = append(stack, fi)
		}
	}
	return nil
}