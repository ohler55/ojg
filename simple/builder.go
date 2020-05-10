// Copyright (c) 2020, Peter Ohler, All rights reserved.

package simple

import "fmt"

var emptyArray = []interface{}{}

type sKey string

type Builder struct {
	stack  []interface{}
	starts []int
}

func (b *Builder) Reset() {
	if 0 < cap(b.stack) && 0 < len(b.stack) {
		b.stack = b.stack[:0]
		b.starts = b.starts[:0]
	} else {
		b.stack = make([]interface{}, 0, 64)
		b.starts = make([]int, 0, 16)
	}
}

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

func (b *Builder) Array(key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		b.stack = append(b.stack, sKey(key[0]))
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	}
	b.starts = append(b.starts, len(b.stack))
	b.stack = append(b.stack, emptyArray)

	return nil
}

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
				if k, ok := b.stack[len(b.stack)-2].(sKey); ok {
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

func (b *Builder) PopAll() {
	for 0 < len(b.starts) {
		b.Pop()
	}
}

func (b *Builder) Result() (result interface{}) {
	if 0 < len(b.stack) {
		result = b.stack[0]
	}
	return
}
