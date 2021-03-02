// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

import (
	"sort"
)

const (
	strNode   = 's'
	numNode   = 'n'
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

type col struct {
	key  interface{} // string or int
	size int
	subs []*col
}

func (c *col) Simplify() interface{} {
	simple := map[string]interface{}{"key": c.key, "size": c.size}
	if 0 < len(c.subs) {
		var subs []interface{}
		for _, sub := range c.subs {
			subs = append(subs, sub.Simplify())
		}
		simple["subs"] = subs
	}
	return simple
}

func (n *node) subKind() (kind byte) {
	for _, m := range n.members {
		if kind != m.kind {
			if kind != 0 {
				return 0
			}
			kind = m.kind
		}
	}
	return
}

func (n *node) genCols(lazy bool) *col {
	switch n.subKind() {
	case arrayNode:
		c := col{}
		for _, m := range n.members {
			m.updateArrayCol(&c, lazy)
		}
		return &c
	case mapNode:
		return nil // TBD
	default:
		return nil
	}
}

func (n *node) updateArrayCol(c *col, lazy bool) {
	for i, m := range n.members {
		var sub *col
		for _, s := range c.subs {
			if s.key == i {
				sub = s
			}
		}
		if sub == nil {
			sub = &col{key: i}
			c.subs = append(c.subs, sub)
		}
		switch m.kind {
		case strNode, numNode:
			if sub.size < m.size {
				sub.size = m.size
			}
		case arrayNode:
			m.updateArrayCol(sub, lazy)
		case mapNode:
			// TBD
		}
	}
	sort.Slice(c.subs, func(i, j int) bool {
		ki, _ := c.subs[i].key.(int)
		kj, _ := c.subs[j].key.(int)
		return ki < kj
	})
	c.size = 0
	for _, sub := range c.subs {
		c.size += sub.size
	}
	if lazy {
		c.size += len(c.subs) + 1
	} else {
		c.size += len(c.subs) * 2
	}
}
