// Copyright (c) 2021, Peter Ohler, All rights reserved.

package pretty

const (
	strNode   = 's'
	numNode   = 'n'
	arrayNode = 'a'
	mapNode   = 'm'
)

type table struct {
	size    int
	columns map[interface{}]int
}

type node struct {
	key     []byte
	members []*node
	buf     []byte
	size    int
	depth   int
	kind    byte
	table   *table
}

// TBD walk all below and create table with dotted keys like x.0.y

// only called for lists
func (n *node) genTable(lazy bool) {
	t := table{columns: map[interface{}]int{}}
	mkind := byte(0)
	for _, m := range n.members {
		switch m.kind {
		case byte(0):
			mkind = m.kind
		case arrayNode:
			if mkind != m.kind {
				if mkind != 0 {
					return
				}
				mkind = m.kind
			}
			for i, m2 := range m.members {
				w := m2.size
				if m2.table != nil {
					w = m2.table.size
				}
				if t.columns[i] < w {
					t.columns[i] = w
				}
			}
		case mapNode:
			if mkind != m.kind {
				if mkind != 0 {
					return
				}
				mkind = m.kind
			}
			for _, m2 := range m.members {
				w := m2.size
				if m2.table != nil {
					w = m2.table.size
				}
				if t.columns[string(m2.key)] < w {
					t.columns[string(m2.key)] = w
				}
			}
		default:
			return
		}
	}
	for _, w := range t.columns {
		t.size += w
	}
	n.table = &t
}
