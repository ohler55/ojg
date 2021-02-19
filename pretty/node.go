// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

const (
	bytesNode = 'b'
	arrayNode = 'a'
	mapNode   = 'm'
)

type node struct {
	key     []byte
	members []*node
	buf     []byte
	size    int
	depth   int
	kind    byte
}
