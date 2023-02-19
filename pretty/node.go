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
	skip    bool
}

type table struct {
	key     any // string or int
	size    int
	columns []*table
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

func (n *node) genTables(lazy bool) *table {
	switch n.subKind() {
	case arrayNode:
		t := table{}
		for _, m := range n.members {
			m.updateArrayTable(&t, lazy)
		}
		return &t
	case mapNode:
		t := table{}
		for _, m := range n.members {
			m.updateMapTable(&t, lazy)
		}
		return &t
	default:
		return nil
	}
}

func (n *node) updateArrayTable(t *table, lazy bool) {
	for i, m := range n.members {
		var col *table
		for _, s := range t.columns {
			if s.key == i {
				col = s
			}
		}
		if col == nil {
			col = &table{key: i}
			t.columns = append(t.columns, col)
		}
		switch m.kind {
		case strNode, numNode:
			if col.size < m.size {
				col.size = m.size
			}
		case arrayNode:
			m.updateArrayTable(col, lazy)
		case mapNode:
			m.updateMapTable(col, lazy)
		}
	}
	sort.Slice(t.columns, func(i, j int) bool {
		ki, _ := t.columns[i].key.(int)
		kj, _ := t.columns[j].key.(int)
		return ki < kj
	})
	t.size = 0
	for _, col := range t.columns {
		t.size += col.size
	}
	if lazy {
		t.size += len(t.columns) + 1
	} else {
		t.size += len(t.columns) * 2
	}
}

func (n *node) updateMapTable(t *table, lazy bool) {
	for _, m := range n.members {
		k := string(m.key)
		var col *table
		for _, s := range t.columns {
			if s.key == k {
				col = s
			}
		}
		if col == nil {
			col = &table{key: k}
			t.columns = append(t.columns, col)
		}
		switch m.kind {
		case strNode, numNode:
			if col.size < m.size {
				col.size = m.size
			}
		case arrayNode:
			m.updateArrayTable(col, lazy)
		case mapNode:
			m.updateMapTable(col, lazy)
		}
	}
	sort.Slice(t.columns, func(i, j int) bool {
		ki, _ := t.columns[i].key.(string)
		kj, _ := t.columns[j].key.(string)
		return ki < kj
	})
	t.size = 0
	for _, col := range t.columns {
		k, _ := col.key.(string)
		t.size += col.size + len(k)
	}
	if lazy {
		t.size += len(t.columns)*2 + 1
	} else {
		t.size += len(t.columns) * 4
	}
}
