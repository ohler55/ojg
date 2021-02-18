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
	kind    byte
}

// Simplify is use for debugging.
func (n *node) Simplify() interface{} {
	var members []interface{}
	for _, m := range n.members {
		members = append(members, m)
	}
	return map[string]interface{}{
		"key":     string(n.key),
		"buf":     string(n.buf),
		"size":    n.size,
		"kind":    n.kind,
		"members": members,
	}
}
