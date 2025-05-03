// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rbtree

import "github.com/mrtcx/plusdata/tree"

type element struct {
	rb   *rbTree
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
	return e.rb.Next(e.node.key)
}

func (e *element) Prev() tree.Element {
	return e.rb.Prev(e.node.key)
}

func (rb *rbTree) Find(key interface{}) tree.Element {
	node := rb.findNode(key)
	if node == nil {
		return nil
	}
	return &element{rb: rb, node: node}
}

func (rb *rbTree) Left() tree.Element {
	mleft := rb.leftNode()
	if mleft == nil {
		return nil
	}
	return &element{rb: rb, node: mleft}
}

func (rb *rbTree) Right() tree.Element {
	mright := rb.rightNode()
	if mright == nil {
		return nil
	}
	return &element{rb: rb, node: mright}
}

func (rb *rbTree) Prev(key interface{}) tree.Element {
	node := rb.findPrevNode(key)
	if node == nil {
		return nil
	}
	return &element{rb: rb, node: node}
}

func (rb *rbTree) Next(key interface{}) tree.Element {
	node := rb.findNextNode(key)
	if node == nil {
		return nil
	}
	return &element{rb: rb, node: node}
}
