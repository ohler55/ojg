// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt

import (
	"fmt"

	"github.com/ohler55/ojg/gen"
)

var emptySlice = []interface{}{}

// Builder is a basic type builder. It uses a stack model to build where maps
// (objects) and slices (arrays) add pushed on the stack and closed with a
// pop.
type Builder struct {
	stack  []interface{}
	starts []int
}

// Reset the builder.
func (b *Builder) Reset() {
	if 0 < cap(b.stack) && 0 < len(b.stack) {
		b.stack = b.stack[:0]
		b.starts = b.starts[:0]
	} else {
		b.stack = make([]interface{}, 0, 64)
		b.starts = make([]int, 0, 16)
	}
}

// Object pushs a map[string]interface{} onto the stack. A key must be
// provided if the top of the stack is an object (map) and must not be
// provided if the op of the stack is an array or slice.
func (b *Builder) Object(key ...string) error {
	newObj := map[string]interface{}{}
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		if obj, _ := b.stack[len(b.stack)-1].(map[string]interface{}); obj != nil {
			obj[key[0]] = newObj
		}
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	}
	b.starts = append(b.starts, -1)
	b.stack = append(b.stack, newObj)

	return nil
}

// Array pushs a []interface{} onto the stack. A key must be provided if the
// top of the stack is an object (map) and must not be provided if the op of
// the stack is an array or slice.
func (b *Builder) Array(key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		b.stack = append(b.stack, gen.Key(key[0]))
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	}
	b.starts = append(b.starts, len(b.stack))
	b.stack = append(b.stack, emptySlice)

	return nil
}

// Value pushs a value onto the stack. A key must be provided if the top of
// the stack is an object (map) and must not be provided if the op of the
// stack is an array or slice.
func (b *Builder) Value(value interface{}, key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		if obj, _ := b.stack[len(b.stack)-1].(map[string]interface{}); obj != nil {
			obj[key[0]] = value
		}
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	} else {
		b.stack = append(b.stack, value)
	}
	return nil
}

// Pop the stack, closing an array or object.
func (b *Builder) Pop() {
	if 0 < len(b.starts) {
		start := b.starts[len(b.starts)-1]
		if 0 <= start { // array
			start++
			size := len(b.stack) - start
			a := make([]interface{}, size)
			copy(a, b.stack[start:len(b.stack)])
			b.stack = b.stack[:start]
			b.stack[start-1] = a
			if 2 < len(b.stack) {
				if k, ok := b.stack[len(b.stack)-2].(gen.Key); ok {
					if obj, _ := b.stack[len(b.stack)-3].(map[string]interface{}); obj != nil {
						obj[string(k)] = a
						b.stack = b.stack[:len(b.stack)-2]
					}
				}
			}
		}
		b.starts = b.starts[:len(b.starts)-1]
	}
}

// PopAll repeats Pop until all open arrays or objects are closed.
func (b *Builder) PopAll() {
	for 0 < len(b.starts) {
		b.Pop()
	}
}

// Result of the builder is returned. This is the first item pushed on to the
// stack.
func (b *Builder) Result() (result interface{}) {
	if 0 < len(b.stack) {
		result = b.stack[0]
	}
	return
}
