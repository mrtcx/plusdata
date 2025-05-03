// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package btree

import "github.com/mrtcx/plusdata/tree"

type element struct {
	btree *bTree
	node  *Node
	idx   int
}

func (e *element) Key() interface{} {
	return e.node.keys[e.idx]
}

func (e *element) Value() interface{} {
	return e.node.values[e.idx]
}

func (e *element) SetValue(value interface{}) {
	e.node.values[e.idx] = value
}

func (e *element) Next() tree.Element {
	return e.btree.Next(e.node.keys[e.idx])
}

func (e *element) Prev() tree.Element {
	return e.btree.Prev(e.node.keys[e.idx])
}

func (bp *bTree) Find(key interface{}) tree.Element {
	node, idx := bp.findNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{btree: bp, node: node, idx: idx}
}

func (bp *bTree) Left() tree.Element {
	if bp.root == nil {
		return nil
	}
	mleft := mostLeft(bp.root)
	if mleft == nil {
		return nil
	}
	return &element{btree: bp, node: mleft, idx: 0}
}

func (bp *bTree) Right() tree.Element {
	if bp.root == nil {
		return nil
	}
	mright := mostRight(bp.root)
	if mright == nil {
		return nil
	}
	return &element{btree: bp, node: mright, idx: len(mright.keys) - 1}
}

func (bp *bTree) Prev(key interface{}) tree.Element {
	node, idx := bp.findPrevNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{btree: bp, node: node, idx: idx}
}

func (bp *bTree) Next(key interface{}) tree.Element {
	node, idx := bp.findNextNodeIdx(key)
	if node == nil {
		return nil
	}
	return &element{btree: bp, node: node, idx: idx}
}
