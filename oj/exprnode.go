// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

// GetNodes the elements of the data identified by the path.
func (x Expr) GetNodes(n Node) (results []Node) {
	if len(x) == 0 {
		return
	}
	var v Node
	var prev Node

	stack := make([]Node, 0, 64)
	stack = append(stack, n)

	f := x[0]
	fi := int64(0) // frag index
	stack = append(stack, Int(fi))

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(Int); up {
			stack = stack[:len(stack)-1]
			fi = int64(ii) & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			if tv, ok := prev.(Object); ok {
				if v, has = tv[string(tf)]; has {
					if fi == int64(len(x))-1 { // last one
						results = append(results, v)
					} else {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Nth:
			i := int(tf)
			if tv, ok := prev.(Array); ok {
				if i < 0 {
					i = len(tv) + i
				}
				var v Node
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == int64(len(x))-1 { // last one
					results = append(results, v)
				} else {
					switch v.(type) {
					case Object, Array:
						stack = append(stack, v)
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case Object:
				if fi == int64(len(x))-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
				if fi == int64(len(x))-1 { // last one
					for _, v = range tv {
						results = append(results, v)
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(Int)
			top := (int64(di) & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (int64(di) & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == int64(len(x))-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == int64(len(x))-1 { // last one
						for _, v = range tv {
							results = append(results, v)
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, Int(fi|descentChildFlag))
				}
			} else {
				if fi == int64(len(x))-1 { // last one
					if top {
						results = append(results, prev)
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == int64(len(x))-1 { // last one
				results = append(results, n)
			} else {
				stack = append(stack, n)
			}
		case At, Bracket:
			if fi == int64(len(x))-1 { // last one
				results = append(results, prev)
			} else {
				stack = append(stack, prev)
			}
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					var has bool
					switch tv := prev.(type) {
					case Object:
						if v, has = tv[string(tu)]; has {
							if fi == int64(len(x))-1 { // last one
								results = append(results, v)
							} else {
								switch v.(type) {
								case Object, Array:
									stack = append(stack, v)
								}
							}
						}
					default:
						// TBD try reflection
						continue
					}
				case int64:
					i := int(tu)
					switch tv := prev.(type) {
					case Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v Node
						if 0 < i && i < len(tv) {
							v = tv[i]
						}
						if fi == int64(len(x))-1 { // last one
							results = append(results, v)
						} else {
							switch v.(type) {
							case Object, Array:
								stack = append(stack, v)
							}
						}
					default:
						// TBD reflection
						continue
					}
				}
			}
		}
		if fi < int64(len(x))-1 {
			if _, ok := stack[len(stack)-1].(Int); !ok {
				fi++
				f = x[fi]
				stack = append(stack, Int(fi))
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

func (x Expr) FirstNode(n Node) (result Node) {
	if len(x) == 0 {
		return nil
	}
	var v Node
	var prev Node

	stack := make([]Node, 0, 64)
	defer func() {
		stack = stack[0:cap(stack)]
		for i := len(stack) - 1; 0 <= i; i-- {
			stack[i] = nil
		}
	}()
	stack = append(stack, n)
	f := x[0]
	fi := int64(0) // frag index
	stack = append(stack, Int(fi))

	for 1 < len(stack) { // must have at least a data element and a fragment index
		prev = stack[len(stack)-2]
		if ii, up := prev.(Int); up {
			stack = stack[:len(stack)-1]
			fi = int64(ii) & fragIndexMask
			f = x[fi]
			continue
		}
		stack[len(stack)-2] = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch tf := f.(type) {
		case Child:
			var has bool
			if tv, ok := prev.(Object); ok {
				if v, has = tv[string(tf)]; has {
					if fi == int64(len(x))-1 { // last one
						return v
					}
					switch v.(type) {
					case Object, Array:
						stack = append(stack, v)
					}
				}
			}
		case Nth:
			i := int(tf)
			if tv, ok := prev.(Array); ok {
				if i < 0 {
					i = len(tv) + i
				}
				var v Node
				if 0 < i && i < len(tv) {
					v = tv[i]
				}
				if fi == int64(len(x))-1 { // last one
					return v
				} else {
					switch v.(type) {
					case Object, Array:
						stack = append(stack, v)
					}
				}
			}
		case Wildcard:
			switch tv := prev.(type) {
			case Object:
				if fi == int64(len(x))-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			case Array:
				if fi == int64(len(x))-1 { // last one
					for _, v = range tv {
						return v
					}
				} else {
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					}
				}
			}
		case Descent:
			di, _ := stack[len(stack)-1].(Int)
			top := (int64(di) & descentChildFlag) == 0
			// first pass expands, second continues evaluation
			if (int64(di) & descentFlag) == 0 {
				self := false
				switch tv := prev.(type) {
				case Object:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, Int(di|descentFlag))
					if fi == int64(len(x))-1 { // last one
						for _, v = range tv {
							return v
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				case Array:
					// Put prev back and slide fi.
					stack[len(stack)-1] = prev
					stack = append(stack, di|descentFlag)
					if fi == int64(len(x))-1 { // last one
						if 0 < len(tv) {
							return tv[0]
						}
					}
					for _, v = range tv {
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
							self = true
						}
					}
				}
				if self {
					stack = append(stack, Int(fi|descentChildFlag))
				}
			} else {
				if fi == int64(len(x))-1 { // last one
					if top {
						return prev
					}
				} else {
					stack = append(stack, prev)
				}
			}
		case Root:
			if fi == int64(len(x))-1 { // last one
				return n
			}
			stack = append(stack, n)
		case At, Bracket:
			if fi == int64(len(x))-1 { // last one
				return prev
			}
			stack = append(stack, prev)
		case Union:
			for _, u := range tf {
				switch tu := u.(type) {
				case string:
					var has bool
					switch tv := prev.(type) {
					case Object:
						if v, has = tv[string(tu)]; has {
							if fi == int64(len(x))-1 { // last one
								return v
							}
							switch v.(type) {
							case Object, Array:
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
					case Array:
						if i < 0 {
							i = len(tv) + i
						}
						var v Node
						if 0 < i && i < len(tv) {
							v = tv[i]
						}
						if fi == int64(len(x))-1 { // last one
							return v
						}
						switch v.(type) {
						case Object, Array:
							stack = append(stack, v)
						}
					default:
						// TBD reflection
						continue
					}
				}
			}
		}
		if fi < int64(len(x))-1 {
			if _, ok := stack[len(stack)-1].(Int); !ok {
				fi++
				f = x[fi]
				stack = append(stack, Int(fi))
			}
		}
	}
	return nil
}
