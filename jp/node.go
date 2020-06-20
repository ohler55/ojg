// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import (
	"strconv"

	"github.com/ohler55/ojg/gen"
)

type index int

// String returns the key as a string.
func (i index) String() string {
	return strconv.Itoa(int(i))
}

// Alter converts the node into it's native type. Note this will modify
// Objects and Arrays in place making them no longer usable as the
// original type. Use with care!
func (i index) Alter() interface{} {
	return i
}

// Simplify makes a copy of the node but as simple types.
func (i index) Simplify() interface{} {
	return int64(i)
}

// Dup returns a deep duplicate of the node.
func (i index) Dup() gen.Node {
	return i
}

// Empty returns true if the node is empty.
func (i index) Empty() bool {
	return false
}

// GetNodes the elements of the data identified by the path.
func (x Expr) GetNodes(n gen.Node) (results []gen.Node) {
	if len(x) == 0 {
		return
	}
	var v gen.Node
	var prev gen.Node
	var has bool

	stack := make([]gen.Node, 0, 64)
	stack = append(stack, n)

	f := x[0]
	fi := index(0) // frag index
	stack = append(stack, index(fi))

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(index); up {
			stack = stack[:len(stack)-1]
			fi = index(ii) & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			if tv, ok := prev.(gen.Object); ok {
				if v, has = tv[string(tf)]; has {
					if fi == index(len(x))-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Nth:
			i := int(tf)
			if tv, ok := prev.(gen.Array); ok {
				if i < 0 {
					i = len(tv) + i
				}
				var v gen.Node
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == index(len(x))-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case gen.Object:
				if fi == index(len(x))-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if fi == index(len(x))-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(gen.Int)
			top := (index(di) & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (int64(di) & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == index(len(x))-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == index(len(x))-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, gen.Int(fi|descentChildFlag))
				}
			} else {
				if fi == index(len(x))-1 { // last one
					if top {
						results = append(results, prev)
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == index(len(x))-1 { // last one
				results = append(results, n)
			} else {
				stack = append(stack, n)
			}
		case At, Bracket:
			if fi == index(len(x))-1 { // last one
				results = append(results, prev)
			} else {
				stack = append(stack, prev)
			}
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					switch tv := prev.(type) {
					case gen.Object:
						if v, has = tv[string(tu)]; has {
							if fi == index(len(x))-1 { // last one
								results = append(results, v)
							} else {
								switch v.(type) {
								case gen.Object, gen.Array:
									stack = append(stack, v)
								}
							}
						}
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v gen.Node
						if 0 < i && i < len(tv) {
							v = tv[i]
						}
						if fi == index(len(x))-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case Slice:
			start := 0
			end := -1
			step := 1
			if 0 < len(tf) {
				start = tf[0]
			}
			if 1 < len(tf) {
				end = tf[1]
			}
			if 2 < len(tf) {
				step = tf[2]
			}
			if tv, ok := prev.(gen.Array); ok {
				if start < 0 {
					start = len(tv) + start
				}
				if end < 0 {
					end = len(tv) + end
				}
				if start < 0 || end < 0 || len(tv) <= start || len(tv) <= end || step == 0 {
					continue
				}
				if 0 < step {
					for i := start; i <= end; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				} else {
					for i := start; end <= i; i += step {
						v = tv[i]
						if int(fi) == len(x)-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					}
				}
			}
		case *Filter:
			this := len(stack)
			stack, _ = tf.Eval(stack, prev).([]gen.Node)
			if int(fi) == len(x)-1 { // last one
				for ; this < len(stack); this++ {
					results = append(results, stack[this])
				}
				if this < len(stack) {
					stack = stack[:this]
				}
			}
		}
		if int(fi) < len(x)-1 {
			if _, ok := stack[len(stack)-1].(index); !ok {
				fi++
				f = x[fi]
				stack = append(stack, index(fi))
			}
		}
	}
	// Free up anything still on the stack.
	stack = stack[0:cap(stack)]
	for i := cap(stack) - 1; 0 <= i; i-- {
		stack[i] = nil
	}
	return
}

func (x Expr) FirstNode(n gen.Node) (result gen.Node) {
	if len(x) == 0 {
		return nil
	}
	var v gen.Node
	var prev gen.Node
	var has bool

	stack := make([]gen.Node, 0, 64)
	defer func() {
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	}()
	stack = append(stack, n)
	f := x[0]
	fi := index(0) // frag index
	stack = append(stack, index(fi))

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(index); up {
			stack = stack[:len(stack)-1]
			fi = index(ii) & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			if tv, ok := prev.(gen.Object); ok {
				if v, has = tv[string(tf)]; has {
					if fi == index(len(x))-1 { // last one
						return v
					}
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Nth:
			i := int(tf)
			if tv, ok := prev.(gen.Array); ok {
				if i < 0 {
					i = len(tv) + i
				}
				if 0 <= i && i < len(tv) {
					v = tv[i]
				}
				if fi == index(len(x))-1 { // last one
					return v
				} else {
					switch v.(type) {
					case gen.Object, gen.Array:
						stack = append(stack, v)
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case gen.Object:
				if fi == index(len(x))-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			case gen.Array:
				if fi == index(len(x))-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(gen.Int)
			top := (int64(di) & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (int64(di) & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case gen.Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, gen.Int(di|descentFlag))
					if fi == index(len(x))-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				case gen.Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == index(len(x))-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, gen.Int(fi|descentChildFlag))
				}
			} else {
				if fi == index(len(x))-1 { // last one
					if top {
						return prev
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == index(len(x))-1 { // last one
				return n
			}
			stack = append(stack, n)
		case At, Bracket:
			if fi == index(len(x))-1 { // last one
				return prev
			}
			stack = append(stack, prev)
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					switch tv := prev.(type) {
					case gen.Object:
						if v, has = tv[string(tu)]; has {
							if fi == index(len(x))-1 { // last one
								return v
							}
							switch v.(type) {
							case gen.Object, gen.Array:
								stack = append(stack, v)
							}
						}
					default:
						// TBD try reflection
						continue
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case gen.Array:
						if i < 0 {
							i = len(tv) + i
						}
						if 0 <= i && i < len(tv) {
							v = tv[i]
						}
						if fi == index(len(x))-1 { // last one
							return v
						}
						switch v.(type) {
						case gen.Object, gen.Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Slice:
			start := 0
			if 0 < len(tf) {
				start = tf[0]
			}
			if tv, ok := prev.(gen.Array); ok {
				if start < 0 {
					start = len(tv) + start
				}
				if start < 0 || len(tv) <= start {
					continue
				}
				v := tv[start]
				if int(fi) == len(x)-1 { // last one
					return v
				}
				v = tv[start]
				if int(fi) == len(x)-1 { // last one
					return v
				}
				switch v.(type) {
				case gen.Object, gen.Array:
					stack = append(stack, v)
				}
			}
		case *Filter:
			this := len(stack)
			stack, _ := tf.Eval(stack, prev).([]gen.Node)
			if int(fi) == len(x)-1 { // last one
				if this < len(stack) {
					result := stack[this]
					stack = stack[:this]
					return result
				}
			}
		}
		if int(fi) < len(x)-1 {
			if _, ok := stack[len(stack)-1].(index); !ok {
				fi++
				f = x[fi]
				stack = append(stack, index(fi))
			}
		}
	}
	return nil
}
