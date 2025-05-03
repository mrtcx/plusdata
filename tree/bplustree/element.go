// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bplustree

import "github.com/mrtcx/plusdata/tree"

type element struct {
	node *Node
	idx  int
}

func (e *element) Key() interface{} {
	return e.node.keys[e.idx]
}

func (e *element) Value() interface{} {
	return e.node.valus[e.idx]
}

func (e *element) SetValue(value interface{}) {
	e.node.valus[e.idx] = value
}

func (e *element) Next() tree.Element {
	if e.idx == len(e.node.keys)-1 {
		if e.node.next == nil {
			return nil
		}
		return &element{node: e.node.next, idx: 0}
	}
	return &element{
		node: e.node,
		idx:  e.idx + 1,
	}
}

func (e *element) Prev() tree.Element {
	if e.idx == 0 {
		if e.node.prev == nil {
			return nil
		}
		return &element{
			node: e.node.prev,
			idx:  len(e.node.prev.keys) - 1,
		}
	}
	return &element{
		node: e.node,
		idx:  e.idx - 1,
	}
}

func (bp *bplusTree) Find(key interface{}) tree.Element {
	node, idx := bp.getNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{node: node, idx: idx}
}

func (bp *bplusTree) Left() tree.Element {
	mleft := bp.mostLeft()
	if mleft == nil {
		return nil
	}
	return &element{node: mleft, idx: 0}
}

func (bp *bplusTree) Right() tree.Element {
	mright := bp.mostRight()
	if mright == nil {
		return nil
	}
	return &element{node: mright, idx: len(mright.keys) - 1}
}

func (bp *bplusTree) Prev(key interface{}) tree.Element {
	node, idx := bp.prevNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{node: node, idx: idx}
}

func (bp *bplusTree) Next(key interface{}) tree.Element {
	node, idx := bp.nextNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{node: node, idx: idx}
}
