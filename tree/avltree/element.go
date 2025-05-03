// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package avltree

import "github.com/mrtcx/plusdata/tree"

type element struct {
	avl  *avlTree
	node *Node
}

func (e *element) Key() interface{} {
	return e.node.key
}

func (e *element) Value() interface{} {
	return e.node.value
}

func (e *element) SetValue(value interface{}) {
	e.node.value = value
}

func (e *element) Next() tree.Element {
	return e.avl.Next(e.node.key)
}

func (e *element) Prev() tree.Element {
	return e.avl.Prev(e.node.key)
}

func (avl *avlTree) Find(key interface{}) tree.Element {
	node := avl.findNode(key)
	if node == nil {
		return nil
	}
	return &element{avl: avl, node: node}
}

func (avl *avlTree) Left() tree.Element {
	mleft := avl.leftNode()
	if mleft == nil {
		return nil
	}
	return &element{avl: avl, node: mleft}
}

func (avl *avlTree) Right() tree.Element {
	mright := avl.rightNode()
	if mright == nil {
		return nil
	}
	return &element{avl: avl, node: mright}
}

func (avl *avlTree) Prev(key interface{}) tree.Element {
	node := avl.findPrevNode(key)
	if node == nil {
		return nil
	}
	return &element{avl: avl, node: node}
}

func (avl *avlTree) Next(key interface{}) tree.Element {
	node := avl.findNextNode(key)
	if node == nil {
		return nil
	}
	return &element{avl: avl, node: node}
}
