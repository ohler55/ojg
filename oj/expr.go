// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

type Expr []Frag

const (
	fragIndexMask = 0x0000ffff
	descentFlag   = 0x00010000
)

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

func (x Expr) GetNodes(n Node) (result []Node) {
	// TBD
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

func W() Expr {
	return Expr{Wildcard('*')}
}

func C(key string) Expr {
	return Expr{Child(key)}
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

// The easy way to implement the Get is to have each fragment handle the
// getting using recursion. The overhead of a go function call is rather high
// though so instead a psuedo call stack is implemented here that grows and
// shrinks as the getting takes place. The fragment index if placed on the
// stack as well mostly for a small degree of simplicity in what a few people
// might find a complex approach to the solution. Its twice as fast as the
// recursive function call approach.

// [map,a] - down
// [map,a,map,b] - down
//   append result
// [map,a,b] - down
// [map,a] - up
// []

// [map,*] - down
// [map,*,a-map,b-map,c-map,d-map,a] - down
// [map,*,a-map,b-map,c-map,d-map,a,map,b] - down
//   append result
// [map,*,a-map,b-map,c-map,d-map,a,b] - down
//   up
// [map,*,a-map,b-map,c-map,d-map,a] - up
//   remove d-map and a
//   append a if prev is not a frag
//   down
// [map,*,a-map,b-map,c-map,a] - down
// [map,*,a-map,b-map,c-map,a,c-map,b] - down
//   append result
//   up
// ...
// [map,*,a-map,a] - down
// [map,*,a-map,a,a-map,b] - down
//   append result
//   up
// [map,*,a-map,a] - up
//   remove a-map and a
//   prev is a frag so fi-- and remove 2 from stack
// [map,*] - up
//   remove 2
// []

// [map,*]
// [*,a-map,b-map,c-map,d-map,a]
// [*,a-map,b-map,c-map,a,d-a-map,b]
//   append d-a-map[b]
// [*,a-map,b-map,c-map,a,b]
// [*,a-map,b-map,c-map,a]
// [*,a-map,b-map,a,c-a-map,b]
// ...
// [*,a]
// [*,a,a-a-map,b]
//   append a-a-map[b]
// [*,a,b]
// [*,a]
// [*]
// []

// TBD on each iteration
// if [..., data, frag]
//   process data by frag
// if [..., frag, frag]
//   frag is finished so fi-- and pop

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
		stack[len(stack)-2] = fi
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
				}
			}
		case Descent:

			// TBD if index is unmasked then do next
			//   if masked then iterate with unmasked self

			// TBD like wildcard but put self at end of one pass on next for another pass
			//  how to avoid iterating more than one on maps?
			//  first next frag, how to indicate the second pass?
			//   maybe negative index in stack to indicate alt
		case Root:
			if fi == len(x)-1 { // last one
				results = append(results, data)
			} else {
				stack = append(stack, data)
				fi++
				f = x[fi]
				stack = append(stack, fi)
			}
		case Bracket:
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
		stack[len(stack)-2] = fi
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
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
					stack = append(stack, x[fi])
				}
			}
		case Descent:
			// TBD like wildcard but put self at end of one pass on next for another pass
			//  how to avoid iterating more than one on maps?
			//  first next frag, how to indicate the second pass?
			//   maybe negative index in stack to indicate alt
		case Root:
			if fi == len(x)-1 { // last one
				return data
			}
			stack = append(stack, data)
			fi++
			f = x[fi]
			stack = append(stack, fi)
		case Bracket:
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
